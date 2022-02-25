#! /usr/bin/env bash

set -euo pipefail

echo "Cleaning any leftover containers"
docker rm -f sofa_couchdb{1,2,3}

echo "Starting test containers"
docker run -d --name sofa_couchdb1 -p 5984:5984 couchdb:1
docker run -d --name sofa_couchdb2 -p 5985:5984 couchdb:2
docker run -d --name sofa_couchdb3 -p 5986:5984 -e COUCHDB_USER=admin -e COUCHDB_PASSWORD=adm1nP4rty couchdb:3
