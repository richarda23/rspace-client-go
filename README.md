README.md

rspace-client-go is an client for the RSpace API written in Go.

It is built and tested against RSpace version 1.69 - this is the minimum required version.


## Getting started

Set the following environment variables in your shell in order to be able to run tests:

   export RSPACE_API_KEY=<myApiKey> RSPACE_URL=<URL TO API>

e.g.

    export RSPACE_API_KEY=abcdefg RSPACE_URL=https://community.researchspace.com/api/v1

See [rspace.go](rspace/rspace.go) for package information

## Running tests

    go test

All tests are integration tests to be run against a real RSpace server. 

The test code is the best way to see how the library is used. 