package sofa

import (
	"crypto"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"

	"github.com/youmark/pkcs8"

	// Authentication with certificates can break if this is not included even though no methods
	// are called directly.
	// TODO: See if this can be controlled with tags.
	_ "crypto/sha512"
)

// Authenticator is an interface for anything which can supply authentication to a CouchDB server.
// The Authenticator is given access to every request made & also allowed to perform an initial setup
// on the connection.
type Authenticator interface {
	// Authenticate adds authentication to an existing http.Request.
	Authenticate(req *http.Request)

	// Client returns a client with the correct authentication setup to contact the CouchDB server.
	Client() (*http.Client, error)

	// Setup uses the provided connection to setup any authentication information which requires accessing
	// the CouchDB server.
	Setup(*Connection) error

	// Verify sets whether verififcation of server HTTPS certs will be performed by clients created
	// by this authenticator. The default for all authenticators should be to perform the verification
	// unless this method is called with the argument 'false' to specifically disable it.
	Verify(bool)
}

type nullAuthenticator struct {
	InsecureSkipVerify bool
}

func (a *nullAuthenticator) Authenticate(req *http.Request) {}

func (a *nullAuthenticator) Client() (*http.Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: a.InsecureSkipVerify,
		},
	}

	httpClient := &http.Client{
		Transport: tr,
	}

	return httpClient, nil
}

func (a *nullAuthenticator) Setup(con *Connection) error {
	return nil
}

func (a *nullAuthenticator) Verify(verify bool) {
	a.InsecureSkipVerify = !verify
}

// NullAuthenticator is an Authenticator which does no work - it implements the interface but
// does not supply any authentication information to the CouchDB server.
func NullAuthenticator() Authenticator {
	return &nullAuthenticator{}
}

type basicAuthenticator struct {
	Username           string
	Password           string
	InsecureSkipVerify bool
}

func (a *basicAuthenticator) Authenticate(req *http.Request) {
	// Basic auth headers must be set for every individual request
	req.SetBasicAuth(a.Username, a.Password)
}

func (a *basicAuthenticator) Client() (*http.Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: a.InsecureSkipVerify,
		},
	}

	httpClient := &http.Client{
		Transport: tr,
	}

	return httpClient, nil
}

func (a *basicAuthenticator) Setup(con *Connection) error {
	return nil
}

func (a *basicAuthenticator) Verify(verify bool) {
	a.InsecureSkipVerify = !verify
}

// BasicAuthenticator returns an implementation of the Authenticator interface which does HTTP basic
// authentication. If you are not using SSL then this will result in credentials being sent in plain
// text.
func BasicAuthenticator(user, pass string) Authenticator {
	return &basicAuthenticator{
		Username: user,
		Password: pass,
	}
}

type clientCertAuthenticator struct {
	CertPath           string
	KeyPath            string
	CaPath             string
	Password           string
	InsecureSkipVerify bool
}

// ClientCertAuthenticator provides an Authenticator which uses a client SSL certificate
// to authenticate to the couchdb server
func ClientCertAuthenticator(certPath, keyPath, caPath string) (Authenticator, error) {
	return &clientCertAuthenticator{
		CertPath: certPath,
		KeyPath:  keyPath,
		CaPath:   caPath,
	}, nil
}

// ClientCertAuthenticatorPassword provides an Authenticator which uses a client SSL certificate
// to authenticate to the couchdb server. This version allows the user to specify the password
// `the key is encrypted with.
func ClientCertAuthenticatorPassword(certPath, keyPath, caPath, password string) (Authenticator, error) {
	return &clientCertAuthenticator{
		CertPath: certPath,
		KeyPath:  keyPath,
		CaPath:   caPath,
		Password: password,
	}, nil
}

func (c *clientCertAuthenticator) Authenticate(req *http.Request) {}

func (c *clientCertAuthenticator) Client() (*http.Client, error) {
	var cert tls.Certificate
	var err error

	if c.Password == "" {
		cert, err = tls.LoadX509KeyPair(c.CertPath, c.KeyPath)
	} else {
		certBytes, readErr := os.ReadFile(c.CertPath)
		if readErr != nil {
			return nil, readErr
		}

		keyBytes, readErr := os.ReadFile(c.KeyPath)
		if readErr != nil {
			return nil, readErr
		}

		pemBlock, _ := pem.Decode(keyBytes)
		if pemBlock == nil {
			return nil, errors.New("expecting a PEM block in encrypted private key file")
		}

		pwBytes := []byte(c.Password)
		//nolint // DecryptPEMBlock may be insecure (as the warning says) but use can be required to integrate with existing systems
		decBytes, decErr := x509.DecryptPEMBlock(pemBlock, pwBytes)
		if decErr != nil {
			var myPkey crypto.PrivateKey
			rsaKey, rsaPKLoadErr := pkcs8.ParsePKCS8PrivateKeyRSA(pemBlock.Bytes, pwBytes)
			if rsaPKLoadErr != nil {
				ecdsaKey, ecdsaPKLoadErr := pkcs8.ParsePKCS8PrivateKeyECDSA(pemBlock.Bytes, pwBytes)

				if ecdsaPKLoadErr != nil {
					return nil, ecdsaPKLoadErr
				} else {
					myPkey = ecdsaKey
					// fmt.Println("ECDSA")
				}

				// return nil, pkErr
			} else {
				myPkey = rsaKey
				// fmt.Println("RSA")
			}

			cert = tls.Certificate{}
			certPEMBlock := certBytes

			for {
				var certDERBlock *pem.Block
				certDERBlock, certPEMBlock = pem.Decode(certPEMBlock)

				if certDERBlock == nil {
					break
				}

				if certDERBlock.Type == "CERTIFICATE" {
					cert.Certificate = append(cert.Certificate, certDERBlock.Bytes)
				}
			}

			cert.PrivateKey = myPkey

			// return nil, decErr
		} else {
			pemBlock.Bytes = decBytes
			pemBlock.Headers = nil

			cert, err = tls.X509KeyPair(certBytes, pem.EncodeToMemory(pemBlock))
		}
	}

	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: c.InsecureSkipVerify,
	}

	if c.CaPath != "" {
		caCert, err := os.ReadFile(c.CaPath)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig.RootCAs = caCertPool
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return &http.Client{
		Transport: transport,
	}, nil
}

func (c *clientCertAuthenticator) Setup(con *Connection) error {
	return nil
}

func (c *clientCertAuthenticator) Verify(verify bool) {
	c.InsecureSkipVerify = !verify
}

type cookieAuthenticator struct {
	InsecureSkipVerify bool
}

// CookieAuthenticator returns an implementation of the Authenticator interface which supports
// authentication
func CookieAuthenticator() Authenticator {
	return &cookieAuthenticator{}
}

func (a *cookieAuthenticator) Authenticate(req *http.Request) {}

func (a *cookieAuthenticator) Client() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &http.Client{Jar: jar}, nil
}

func (a *cookieAuthenticator) Setup(con *Connection) error {
	return nil
}

func (a *cookieAuthenticator) Verify(verify bool) {
	a.InsecureSkipVerify = !verify
}

type proxyAuthenticator struct {
	Username           string
	Roles              string
	Token              string
	InsecureSkipVerify bool
}

// ProxyAuthenticator returns an implementation of the Authenticator interface which supports
// the proxy authentication method described in the CouchDB documentation. This should not be
// used against a production server as the proxy would be expected to set the headers in that
// case.
func ProxyAuthenticator(username string, roles []string, secret string) Authenticator {
	var token = ""

	if secret != "" {
		mac := hmac.New(sha1.New, []byte(secret))
		_, err := io.WriteString(mac, username)
		if err != nil {
			panic(err)
		}

		token = fmt.Sprintf("%x", mac.Sum(nil))
	}

	return &proxyAuthenticator{
		Username: username,
		Roles:    strings.Join(roles, ","),
		Token:    token,
	}
}

func (a *proxyAuthenticator) Authenticate(req *http.Request) {
	req.Header.Set("X-Auth-CouchDB-UserName", a.Username)
	req.Header.Set("X-Auth-CouchDB-Roles", a.Roles)
	if a.Token != "" {
		req.Header.Set("X-Auth-CouchDB-Token", a.Token)
	}
}

func (a *proxyAuthenticator) Client() (*http.Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: a.InsecureSkipVerify,
		},
	}

	httpClient := &http.Client{
		Transport: tr,
	}

	return httpClient, nil
}

func (a *proxyAuthenticator) Setup(con *Connection) error {
	return nil
}

func (a *proxyAuthenticator) Verify(verify bool) {
	a.InsecureSkipVerify = !verify
}
