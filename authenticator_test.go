package sofa

import (
	"os"
	"testing"
)

func fixtureTest(t *testing.T) {
	src, err := os.Stat("test-fixtures")
	if os.IsNotExist(err) || src.Mode().IsRegular() {
		t.Skip("test-fixtures does not exist or is not a directory")
	}
}

func TestClientCertAuthenticator(t *testing.T) {
	fixtureTest(t)

	auth, err := ClientCertAuthenticator(
		"test-fixtures/passwordless.crt",
		"test-fixtures/passwordless.key",
		"test-fixtures/ca.crt")

	if err != nil {
		t.Fatal(err)
	}

	_, err = auth.Client()
	if err != nil {
		t.Fatal(err)
	}
}

func TestClientCertAuthenticatorPassword(t *testing.T) {
	fixtureTest(t)

	auth, err := ClientCertAuthenticatorPassword(
		"test-fixtures/password.crt",
		"test-fixtures/password.key",
		"test-fixtures/ca.crt",
		"Th3P455w0rd")

	if err != nil {
		t.Fatal(err)
	}

	_, err = auth.Client()
	if err != nil {
		t.Fatal(err)
	}
}

func TestClientCertAuthenticatorPasswordPKCS8(t *testing.T) {
	fixtureTest(t)

	auth, err := ClientCertAuthenticatorPassword(
		"test-fixtures/password-pkcs8.crt",
		"test-fixtures/password-pkcs8.key",
		"test-fixtures/ca.crt",
		"Th3P455w0rd")

	if err != nil {
		t.Fatal(err)
	}

	_, err = auth.Client()
	if err != nil {
		t.Fatal(err)
	}
}

func TestClientCertAuthenticatorPasswordPKCS8_ECDSA(t *testing.T) {
	fixtureTest(t)

	auth, err := ClientCertAuthenticatorPassword(
		"test-fixtures/password-pkcs8-ecdsa.crt",
		"test-fixtures/password-pkcs8-ecdsa.key",
		"test-fixtures/ca.crt",
		"Th3P455w0rd")

	if err != nil {
		t.Fatal(err)
	}

	_, err = auth.Client()
	if err != nil {
		t.Fatal(err)
	}
}
