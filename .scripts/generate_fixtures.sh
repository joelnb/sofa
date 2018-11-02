#! /bin/bash

set -eu

DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

FIXTURE_DIR="$(realpath "${DIR}/../test-fixtures")"
FIXTURE_PASSWORD="Th3P455w0rd"

rm -rfv "${FIXTURE_DIR}"
mkdir -pv "${FIXTURE_DIR}"

pushd "${FIXTURE_DIR}"

# Create the CA certificate
openssl req -new -x509 -extensions v3_ca -keyout ca.key -out ca.crt -days 1825 -nodes -subj "/C=US/ST=Denial/L=Springfield/O=Sofa/CN=www.example.com"

# Create a key with no password & generate a certificate from it
openssl genrsa -out passwordless.key 2048
openssl req -new -key passwordless.key -out passwordless.csr -subj "/C=US/ST=Denial/L=Springfield/O=Sofa/CN=passwordless"
openssl x509 -req -days 1825 -in passwordless.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out passwordless.crt

# Create a key with a password & generate a certificate from it
openssl genrsa -des3 -passout "pass:${FIXTURE_PASSWORD}" -out password.key 2048
openssl req -new -passin "pass:${FIXTURE_PASSWORD}" -key password.key -out password.csr -subj "/C=US/ST=Denial/L=Springfield/O=Sofa/CN=password"
openssl x509 -req -days 1825 -in password.csr -CA ca.crt -CAkey ca.key -set_serial 02 -out password.crt

openssl genpkey -out password-pkcs8.key -des3 -algorithm RSA -pass "pass:${FIXTURE_PASSWORD}" -pkeyopt rsa_keygen_bits:2048
openssl req -new -passin "pass:${FIXTURE_PASSWORD}" -key password-pkcs8.key -out password-pkcs8.csr -subj "/C=US/ST=Denial/L=Springfield/O=Sofa/CN=password-pkcs8"
openssl x509 -req -days 1825 -in password-pkcs8.csr -CA ca.crt -CAkey ca.key -set_serial 02 -out password-pkcs8.crt

openssl ecparam -name prime256v1 -genkey | openssl pkcs8 -topk8 -v2 des3 -passout "pass:${FIXTURE_PASSWORD}" -out password-pkcs8-ecdsa.key
openssl req -new -passin "pass:${FIXTURE_PASSWORD}" -key password-pkcs8-ecdsa.key -out password-pkcs8-ecdsa.csr -subj "/C=US/ST=Denial/L=Springfield/O=Sofa/CN=password-pkcs8-ecdsa"
openssl x509 -req -days 1825 -in password-pkcs8-ecdsa.csr -CA ca.crt -CAkey ca.key -set_serial 02 -out password-pkcs8-ecdsa.crt

popd
