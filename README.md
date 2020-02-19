# sofa

Simply connect to CouchDB database servers using Go.

## Documentation

View the [full documentation](https://pkg.go.dev/github.com/joelnb/sofa?tab=doc) which includes examples.

## Limitations

Large parts of the [CouchDB API](http://docs.couchdb.org/en/2.0.0/api/) are covered but not all functionality is currently implemented. Pull requests for any missing functionality would be welcomed!

## Contributing

Contributions of all sizes are welcomed. Simply make a pull request and I will be happy to discuss. If you don't have time to write the code please consider at least creating an issue so that I can ensure the issue gets sorted eventually.

### Running tests

The basic tests can be run using a simple `go test`. To run a more complete set of tests which access a real database you will need a temporary CouchDB instance. The simplest way to create this is using docker:

    docker run -d --name couchdb -p 5984:5984 couchdb:1

To run all the tests you will also need a version 2 server:

    docker run -d --name couchdb -p 5985:5984 couchdb:2

You can then set `SOFA_TEST_HOST_1` and `SOFA_TEST_HOST_2` to set the connection to each server:

    SOFA_TEST_HOST_1=http://127.0.0.1:5984 SOFA_TEST_HOST_2=http://127.0.0.1:5985 go test -v

If you have chosen to only start a single version then only include the appropriate environment variable to ensure tests for the other version are not run.
