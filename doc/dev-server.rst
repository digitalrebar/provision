.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. DigitalRebar Provision documentation under Digital Rebar master license
.. index::
  pair: DigitalRebar Provision; Developer Environment

.. _rs_dev_environment:

Developer Environment
~~~~~~~~~~~~~~~~~~~~~

This page is intended for people who are building DigitalRebar Provision from sources or contributing to the code base.  We maintain inline documentation and test environment and expect contributors to participate in maintenance of those efforts.

.. note:: Prerequisites: we are using go version 1.8 or better.  These documents expect that you are able to install and update Golang.


.. _re_dev_quick:

Developer Quick Start
---------------------

To get started quickly, we've rolled all the installation steps into a script.  The script can be run directly from Github by copying the following lines:

  ::

    curl -fsSL https://raw.githubusercontent.com/digitalrebar/provision/master/tools/build.sh | bash

The script will use your current **GOPATH** variable for placing the code.  If **GOPATH** isn't set,
it will be set to *$HOME/go*.

Once the script is complete, you can change directory to the source area and continue development.

  ::

    export GOPATH=${GOPATH:-$HOME/go}
    cd "$GOPATH/src/github.com/digitalrebar/provision"


If you want more details on how to run the result, consult the :ref:`rs_install` section.  The **install.sh** script
can be used to install from the source directory after a build.

.. _rs_dev_build:

Building The Server
-------------------

Once you've got the code and assets once, you can repeat the process from the project root with the following command:

  ::

    tools/build.sh


Another reason to use the *tools/build.sh* script is that it will inject version information into the built binaries to make
it easier to track what version is deployed.  Developer builds and production builds will be identifiable through the *version*
command on both the server and the cli.

Serving UI from File System
---------------------------

When working on the DigitalRebar Provision UI, you can skip the generate steps by using the `--dev-ui` flag.  Generally, this is started using `--dev-ui ./embedded/assets/ui`


.. _rs_testing:

Running the Tests
-----------------

DigitalRebar Provision uses the Golang test libraries and we work hard to maintain test coverage.

We use `tools/test.sh` in the provision root directory to test the entire code base.

You can test individual modules from their subdirectories by running `go test`

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

* Rebuild the world (**tools/build.sh**)

Packaging the Code
------------------

Once the code is built, the code can be package for storage in Github or for use by the **install.sh** script.

Runing the **tools/package.sh** script will generate a **dr-provision.zip** and **dr-provision.sha256** file.  These files
can be used with the :ref:`rs_install` process.
