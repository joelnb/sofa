package sofa

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"

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
}

type nullAuthenticator struct{}

func (t *nullAuthenticator) Authenticate(req *http.Request) {}

func (t *nullAuthenticator) Client() (*http.Client, error) {
	return &http.Client{}, nil
}

func (t *nullAuthenticator) Setup(con *Connection) error {
	return nil
}

// NullAuthenticator is an Authenticator which does no work - it implements the interface but
// does not supply any authentication information to the CouchDB server.
func NullAuthenticator() Authenticator {
	return &nullAuthenticator{}
}

type basicAuthenticator struct {
	Username string
	Password string
}

func (a *basicAuthenticator) Authenticate(req *http.Request) {
	// Basic auth headers must be set for every individual request
	req.SetBasicAuth(a.Username, a.Password)
}

func (a *basicAuthenticator) Client() (*http.Client, error) {
	return &http.Client{}, nil
}

func (a *basicAuthenticator) Setup(con *Connection) error {
	return nil
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
	CertPath string
	KeyPath  string
	CaPath   string
	Password string
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
		keyBytes, err := ioutil.ReadFile(c.KeyPath)
		if err != nil {
			return nil, err
		}

		pemBlock, _ := pem.Decode(keyBytes)
		if pemBlock == nil {
			return nil, errors.New("expecting a PEM block in encrypted private key file")
		}

		decBytes, err := x509.DecryptPEMBlock(pemBlock, []byte(c.Password))
		if err != nil {
			return nil, err
		}

		certBytes, err := ioutil.ReadFile(c.CertPath)
		if err != nil {
			return nil, err
		}

		cert, err = tls.X509KeyPair(certBytes, decBytes)
	}

	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	if c.CaPath != "" {
		caCert, err := ioutil.ReadFile(c.CaPath)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig.RootCAs = caCertPool
	}

	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return &http.Client{
		Transport: transport,
	}, nil
}

func (c *clientCertAuthenticator) Setup(con *Connection) error {
	return nil
}

type cookieAuthenticator struct{}

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

type proxyAuthenticator struct {
	Username string
	Roles    string
	Token    string
}

// ProxyAuthenticator returns an implementation of the Authenticator interface which supports
// the proxy authentication method described in the CouchDB documentation. This should not be
// used against a production server as the proxy would be expected to set the headers in that
// case.
func ProxyAuthenticator(username string, roles []string, secret string) Authenticator {
	var token = ""

	if secret != "" {
		mac := hmac.New(sha1.New, []byte(secret))
		io.WriteString(mac, username)
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
	return &http.Client{}, nil
}

func (a *proxyAuthenticator) Setup(con *Connection) error {
	return nil
}
