.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Rocket Skates documentation under Digital Rebar master license

Developer Environment
~~~~~~~~~~~~~~~~~~~~~

We are using `go-bindata <https://github.com/jteeuwen/go-bindata>`_ to embed binary assets in *rocket-skates*  You can obtain it by running `go get -u github.com/jteeuwen/go-bindata/...`, which will leave the `go-bindata` executable in `$GOPATH/bin`

We are using `go-swagger <https://github.com/go-swagger/go-swagger>`_ to generate code from the API specification file.  You can obtain it by running `go get -u github.com/go-swagger/go-swagger/cmd/swagger`, which will leave the `swagger` executabe in `$GOPATH/bin`

We are using `swagger-ui <https://github.com/swagger-api/swagger-ui>`_ for a quick UI to inspect the API and drive the system.  It is customizable if we need, but we are running it straight up for now.

This basic dist files have been embedded into the rocket skates binary for the time being.  These are copied from the swagger-ui tree.

*NOTES* on how to get swagger-ui

* git clone https://github.com/swagger-api/swagger-ui
* cp -r swagger-ui/dist/\* embedded/assets/swagger-ui
* change in embedded/assets/swagger-ui/index.html:

  ::

    @@ -38,7 +38,7 @@
           if (url && url.length > 1) {
             url = decodeURIComponent(url[1]);
           } else {
    -        url = "http://petstore.swagger.io/v2/swagger.json";
    +        url = "https://127.0.0.1:8092/swagger.json";
           }
     
           hljs.configure({

* Rebuild the world

Test Data
~~~~~~~~~

There is a test-data directory for local running.

Pulling pinned imports
----------------------

This must be done before building the client or the server.

* glide i


Building Server
---------------

* go generate server/main.go
* go build -o rocket-skates server/\*

