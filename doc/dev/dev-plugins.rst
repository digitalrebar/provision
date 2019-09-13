.. Copyright (c) 2019 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Developer Environment

.. _rs_dev_plugins:

Developing Plugins
~~~~~~~~~~~~~~~~~~

This page is intended for people who are building Digital Rebar Provision plugins.

.. note:: Prerequisites: go version 1.12 or better.  These documents assumes the operator/devloper has the ability to both install and update Golang.

.. _re_dev_plugins_quick:

Plugin Build Quickstart
-----------------------

The follow is an example quickstart on how to setup Golang, get the source code, and build a plugin.
In this document, we are showing how to do this on a stock CentOS 7 system. You can compile Golang
binaries on any platform, but setup and usage on different platforms my vary.


With a freshly installed CentOS 7 system, get, install, and setup Golang:

  ::

    sudo yum -y install git
    mkdir $HOME/go
    wget https://dl.google.com/go/go1.12.7.linux-amd64.tar.gz
    tar -xzvf go1.12.7.linux-amd64.tar.gz
    sudo mv go /usr/local/
    export GOROOT=/usr/local/go
    export GOPATH=$HOME/go
    export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

Verify that Golang is installed and working for you:

  ::

    go version

Use ``go get ...`` syntax to get the source code pieces:

  ::

    go get github.com/digitalrebar/provision
    go get github.com/digitalrebar/provision-plugins

Now, lets move to the plugins directory, and compile a plugin:

  ::

    cd $HOME/go/src/github.com/digitalrebar/provision-plugins
    tools/build-one.sh cmds/ipmi

You should now have the IPMI plugin compiled for Linux amd64 (64bit) version:

  ::

      [shane@fuji provision-plugins]$ ls -l bin/linux/amd64/ipmi
      -rwxr-xr-x. 1 shane shane 12311392 Sep 12 19:37 bin/linux/amd64/ipmi
      [shane@fuji provision-plugins]$ file bin/linux/amd64/ipmi
      bin/linux/amd64/ipmi: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, stripped


.. _rs_dev_plugins_crosscompile:

Crosscompile Plugin for Different Platform
==========================================

Crosscompiling with Golang is very easy, and it is controlled with a few basic
environment variables:

  * ``GOOS`` - defines the Operating System target to build for
  * ``GOARCH`` - defines the platform architecture to build for

An example to force building Linux amd64 binaries (if you are building on a Mac OS X
system, for example):

  ::

    GOOS=linux GOARCH=amd64 tools/build-one.sh cmds/ipmi

If your build system is the same Architecture as your target build binaries platform, simply
drop the ``GOARCH`` variable.

Mac OS X platforms use ``darwin`` as the ``GOOS`` variable setting if you wanted to compile
for a Mac OS X system.

