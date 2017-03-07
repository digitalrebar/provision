.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Rocket Skates documentation under Digital Rebar master license
.. index::
  pair: Rocket Skates; Developer Environment

.. _rs_dev_environment:

Developer Environment
~~~~~~~~~~~~~~~~~~~~~

This page is intended for people who are building Rocket Skates from sources or contributing to the code base.  We maintain a inline documentation and test environment and expect contributors to participate in maintenance of those efforts.

.. note:: Prerequisites: we are using go version 1.8 or better.  These documents expect that you are able to install and update Golang.


.. _re_dev_quick:

Developer Quick Start
---------------------

To get started quickly, we've rolled all the installation steps into a script.  The script can be run directly from Github by copying the following lines:

  ::

    mkdir rocket-skates-dev
    cd rocket-skates-dev
    curl -fsSL https://raw.githubusercontent.com/digitalrebar/rocket-skates/master/tools/build.sh | bash


The end of the script includes environment configuration steps required to build the code.

If you want more details, consult the :ref:`rs_dev_install` section.


.. _rs_dev_build:

Building The Server
-------------------

Once you've got the code built, you can repeat the process from the project root with the following steps:

  ::

    go generate server/main.go
    go build -o rocket-skates server/\*


The generate step is only required when you are changing the embedded assets in the `embedded/assets/` directories.

.. _rs_dev_install:

Developer Install Steps
-----------------------

.. note:: Please review the `tools/build.sh` script also.  It may have been updated more recently than the documentation!

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

