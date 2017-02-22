# rocket-skates

[![Build Status](https://travis-ci.org/rackn/rocket-skates.svg)](https://travis-ci.org/rackn/rocket-skates)
[![codecov](https://codecov.io/gh/rackn/rocket-skates/branch/master/graph/badge.svg)](https://codecov.io/gh/rackn/rocket-skates)
[![Go Report Card](https://goreportcard.com/badge/github.com/rackn/rocket-skates)](https://goreportcard.com/report/github.com/rackn/rocket-skates)
[![GoDoc](https://godoc.org/github.com/rackn/rocket-skates?status.svg)](https://godoc.org/github.com/rackn/rocket-skates)

## Tools

We are using [go-bindata](https://github.com/jteeuwen/go-bindata) to embed binary assets in *rocket-skates*  You can obtain it by running `go get -u github.com/jteeuwen/go-bindata/...`, which will leave the `go-bindata` executable in `$GOPATH/bin`

We are using [go-swagger](https://github.com/go-swagger/go-swagger) to generate code from the API specification file.
You will need a binary to generate the code.  This can be obtained from [here](https://github.com/go-swagger/go-swagger#static-binary). We are currently using version 0.8.0.

We are using [swagger-ui](https://github.com/swagger-api/swagger-ui) for a quick UI to inspect the API and drive the system.  It is customizable if we need, but we are running it straight up for now.

This basic dist files have been embedded into the rocket skates binary for the time being.  These are copied from the swagger-ui tree.

NOTES on how to get swagger-ui

* git clone https://github.com/swagger-api/swagger-ui
* cp -r swagger-ui/dist/\* embedded/assets/swagger-ui
* change in embedded/assets/swagger-ui/index.html:

```
@@ -38,7 +38,7 @@
       if (url && url.length > 1) {
         url = decodeURIComponent(url[1]);
       } else {
-        url = "http://petstore.swagger.io/v2/swagger.json";
+        url = "https://127.0.0.1:8092/swagger.json";
       }
 
       hljs.configure({
```

* Rebuild the world


## Test Data

There is a test-data directory for local running.

## Pulling pinned imports

This must be done before building the client or the server.

* glide i

## Building Server

* go generate server/main.go
* go build -o rocket-skates server/\*

## Running Server

* cd test-data
* sudo ../rocket-skates

NOTE: I need the sudo to bind the tftp port.  This is configurable, i.e.  *--tftp-port=30000*  

## STOP HERE!

## Generating the Client - DON'T DO THIS FOR NOW.

To generate the code:
* swagger generate client -P models.Principal -f swagger.json 

This generates the following directories:

* client
* cmd
* models
* restapi

The file that we edit for the client is:

* cmd/rocket-skates-client/main.go 

This is completely our own file.


## Building Client

* go build -o rscli cmd/rocket-skates-client/main.go

## Running Client

* ./rscli

This is a really stupid client but shows how to make the calls and get back structures.

