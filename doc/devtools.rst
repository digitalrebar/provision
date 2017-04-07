.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. DigitalRebar Provision documentation under Digital Rebar master license
.. index::
  pair: DigitalRebar Provision; Developer Environment

.. _rs_dev_environment:

Developer Environment
~~~~~~~~~~~~~~~~~~~~~

This page is intended for people who are building DigitalRebar Provision from sources or contributing to the code base.  We maintain a inline documentation and test environment and expect contributors to participate in maintenance of those efforts.

.. note:: Prerequisites: we are using go version 1.8 or better.  These documents expect that you are able to install and update Golang.


.. _re_dev_quick:

Developer Quick Start
---------------------

To get started quickly, we've rolled all the installation steps into a script.  The script can be run directly from Github by copying the following lines:

  ::

    mkdir dr-provison-dev
    cd dr-provision-dev
    curl -fsSL https://raw.githubusercontent.com/digitalrebar/provision/master/tools/build.sh | bash


The end of the script includes environment configuration steps required to build the code.

If you want more details, consult the :ref:`rs_dev_install` section.

.. _rs_dev_build:

Building The Server
-------------------

Once you've got the code and assets once, you can repeat the process from the project root with the following steps:

  ::

    go generate server/assets.go
    go build -o dr-provision cmds/dr-provision.go


The generate step is only required when you are changing the embedded assets in the `embedded/assets/` directories.



Serving UI from File System
---------------------------

When working on the DigitalRebar Provision UI, you can skip the generate steps by using the `--dev-ui` flag.  Generally, this is started using `--dev-ui ./embedded/assets/ui`


.. _rs_testing:

Running the Tests
-----------------

DigitalRebar Provision uses the Golang test libraries and we work hard to maintain test coverage.

We use `tools/test.sh` in the provision root directory to test the entire code base.

You can test individual modules from their subdirectories by running `go test`

.. _rs_dev_install:

Developer Install Steps (manual)
--------------------------------

.. note:: Please review the `tools/build.sh` script also.  It may have been updated more recently than the documentation!

We are using `go-bindata <https://github.com/jteeuwen/go-bindata>`_ to embed binary assets in *dr-provision*  The following command 
will leave the *go-bindata* executable in *$GOPATH/bin*.

  ::

    go get -u github.com/jteeuwen/go-bindata/...


We are using `go-swagger <https://github.com/go-swagger/go-swagger>`_ to generate code from the API specification file.  The following
command will leave the *swagger* executable in *$GOPATH/bin*.

  ::

    go get -u github.com/go-swagger/go-swagger/cmd/swagger

We are using `swagger-ui <https://github.com/swagger-api/swagger-ui>`_ for a quick UI to inspect the API and drive the system.
It is customizable if we need, but we are running it straight up for now.

This basic dist files have been embedded into the dr-provision binary for the time being.  These are copied from the swagger-ui tree.


DigitalRebar Provision requires some basic files to provide a PXE environment.  This can be obtained by running the
*download-assets.sh* script.  This will populate the embedded/assets directory.

  ::
    ./tools/download-assets.sh


How to get Swagger-Ui
---------------------

DigiatlRebar Provision uses Swagger to generate interactive help for the API.  This is in the tree by default.  If you
need to update it, do the following:

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

* Rebuild the world (both generate and build)

Test Data
~~~~~~~~~

There is a test-data directory for local running.

Pulling pinned imports
----------------------

This must be done before building the client or the server.

* glide i

