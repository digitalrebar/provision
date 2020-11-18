.. Copyright (c) 2019 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Setting Up Docker Contexts

.. _rs_contexts:

Docker Contexts
===============

About
-----

Contexts are a new top-level object that allow dr-provision to support the underlying idea that tasks for a machine can run somewhere besides on the machine itself. This document shows you how to get started with contexts by using our Docker (OCI) Execution Context plugin.

.. note:: Contexts are a new feature. RackN is actively working on an internal bootstrap deployment of the Context configuration environment.  Please check back for more details on that feature. In the meantime please use the below manual configuration steps to setup your environment to support contexts.

System Setup Information
------------------------

For the Docker (OCI) Execution Context plugin to work you need to make sure you have drp version 4.2.0 or newer. Your drp endpoint must also be running a recent version of docker or podman. This document only covers docker, but if podman is present on the system we will use it.
The following system information lists what was used to create this guide.

  ::

    DRP Endpoint OS: CentOS centos-release-7-7.1908.0.el7.centos.x86_64
    Digital Rebar: 4.2.0
    Docker:
    Client: Docker Engine - Community
     Version:           19.03.5
     API version:       1.40
     Go version:        go1.12.12
    Server: Docker Engine - Community
     Engine:
      Version:          19.03.5
      API version:      1.40 (minimum version 1.12)
      Go version:       go1.12.12
     containerd:
      Version:          1.2.10
      GitCommit:        b34a5c8af56e510852c35414db4c1f4fa6172339
     runc:
      Version:          1.0.0-rc8+dev
      GitCommit:        3e425f80a8c931f88e6d94a8c831b9d5aa481657
     docker-init:
      Version:          0.18.0
      GitCommit:        fec3683

.. note::
    * Older versions of docker will work just fine, but you could run into issues when building your container images.
    * Docker was setup using the following guide: https://docs.docker.com/install/linux/docker-ce/centos/
    * Podman can be used interchangably from Docker.

Prerequisites
-------------

For the Docker (OCI) Execution Context plugin requires that you have either Docker or Podman installed.  Operators may find that Podman provides better security and controls.

Getting Started
---------------

We need to verify that the DRP endpoint is properly configured to function
with the docker-context plugin. To do that make sure the user running
docker or podman and the user running drp are one in the same, or that the user
running drp can use docker. To make things easy in this example both users
will be root.

  ::

    root@li1147-65 ~]# ps waux|egrep 'docker|dr-pro'|grep -v grep
    root      1233  0.0 26.3 447484 267248 ?       Ssl  Dec12   3:11 /usr/local/bin/dr-provision
    root     13339  0.0  7.6 520844 77724 ?        Ssl  19:40   0:00 /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock


Next we need to install the `docker-context` plugin from the catalog, and verify we do not have any existing
contexts created.

  ::

    [root@li1147-65 ~]# drpcli catalog item install docker-context --version tip
    {
      "path": "docker-context",
      "size": 12269440
    }
    [root@li1147-65 ~]# drpcli contexts list Engine=docker-context
    []

Next we should create a docker image we can use.

  ::

    mkdir ~/drp-docker-example
    cd ~/drp-docker-example
    cp `which drpcli` .

    cat << EOF > dockerfile
    FROM alpine:latest
    RUN apk add bash jq
    COPY drpcli /usr/bin/drpcli
    RUN chmod 755 /usr/bin/drpcli
    ENTRYPOINT /usr/bin/drpcli machines processjobs
    EOF


Now we need to build the image:

  ::

    docker build --tag digitalrebar/runner .


Now to create the context our example workflow will use:

  ::

    [root@li1147-65 drp-docker-example]# drpcli contexts create '{"Name": "runner", "Description": "DRP Demo Context", "Engine": "docker-context", "Image": "digitalrebar/runner", "Meta": {"Title": "Example context"}}'
    {
      "Available": true,
      "Bundle": "",
      "Description": "DRP Demo Context",
      "Documentation": "",
      "Endpoint": "",
      "Engine": "docker-context",
      "Errors": [],
      "Image": "digitalrebar/runner",
      "Meta": {
        "Title": "Example context"
      },
      "Name": "runner",
      "ReadOnly": false,
      "Validated": true
    }


In the command above we created a context and named it "runner" then we told the system it would be using the "docker-context" context engine.
Another important bit of info in our command above was the image, which if you notice matches the image we created above using docker.

Finally we need to create a special type of machine called the context runner machine. Pay close attention to the payload being sent with our machine create.

  ::

    [root@li1147-65 drp-docker-example]# docker ps
    CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
    [root@li1147-65 drp-docker-example]# drpcli machines create  '{"Name": "context-runner", "Meta": {"BaseContext": "runner"}}'
    {
      "Address": "",
      "Arch": "amd64",
      "Available": true,
      "BootEnv": "local",
      "Bundle": "",
      "Context": "runner",
      "CurrentJob": "",
      "CurrentTask": -1,
      "Description": "",
      "Endpoint": "",
      "Errors": [],
      "HardwareAddrs": [],
      "Locked": false,
      "Meta": {
        "BaseContext": "runner",
        "feature-flags": "change-stage-v2"
      },
      "Name": "context-runner",
      "OS": "",
      "Params": {},
      "Partial": false,
      "Profiles": [],
      "ReadOnly": false,
      "Runnable": true,
      "Secret": "Ox0Kd7pza_MBNUZQ",
      "Stage": "none",
      "Tasks": [],
      "Uuid": "70264cbb-8db8-4643-bf26-07d8239d7e38",
      "Validated": true,
      "Workflow": ""
    }
    [root@li1147-65 drp-docker-example]# docker ps
    CONTAINER ID        IMAGE                 COMMAND                  CREATED             STATUS              PORTS               NAMES
    7fbff8d9e630        digitalrebar/runner   "/bin/sh -c '/usr/bi…"   3 seconds ago       Up 1 second                             drp-runner-70264cbb-8db8-4643-bf26-07d8239d7e38

First we checked to make sure no containers were running yet using docker ps. Specifically we wanted to make sure no instances of "digitalrebar/runner" were being used yet.
Next we created the machine object using drpcli and a bit of json to craft the object.

.. note:: Notice the **Meta: BaseContext: runner**

The name of this context must match the context name we created when we created the context above.
Finally we ran docker ps again and found that now we have a container running. This should mean success. We can now verify it all works by adding a workflow, stages and tasks that utilize this new context.

To test our setup we will use an example content pack that we will need to clone using git.

  ::

    git clone https://github.com/digitalrebar/colordemo
    cd colordemo
    git checkout context-demo
    cd content
    drpcli contents bundle ../context-demo-1.0.yaml
    drpcli contents create ../context-demo-1.0.yaml

This downloads the colordemo repo from github, next we switch from the master branch to the *context-demo* branch. Next we bundle and upload our content pack which adds a workflow to our endpoint called `context-demo`.

Now lets create another machine object, only this time we will run our empty machine through our new context workflow we added above.

  ::

    drpcli machines create example
    drpcli machines update Name:example '{"Workflow":"context-demo"}'
    drpcli machines update Name:example '{"Context":"runner"}'

We have created an empty machine object, next we updated the workflow and set it to our new workflow, finally we set the context to the **runner** context which causes our workflow to trigger and run.

At this point you should be able to navigate to the jobs page of the web portal, look for the task `context-demo-example` and see it has run for about 10 seconds, and if you click into the job you will see the output from the task.


Container Image Locations
-------------------------

The Docker-Context plugin does not use `docker pull` or related approaches to retrieve images because this assumes external connectivity.  Instead, the Docker-Context relies on the needed images being stored as single artifact archives in the DRP server's `/files/contexts/docker-contexts` path.

Here is the process the plugin uses when starting a Context container:

* When a Context is requested, the Docker-Context plugin will check the system for an image tag that matches the `Context.Image` value.
* If that image is not installed, the plugin will download it from the DRP server's `/files/contexts/docker-contexts` path and register it as an image with the correct tag.
* Start the container as a Context.

If the tagged container image already exists on the DRP endpoint, then no download is attempted.  This allows operators to also pre-stage images for testing.

Additional Resources
--------------------

Here are some additional resources you may find useful:

    * MeetUp Intro To Contexts: https://youtu.be/4UGozDUGxy4
    * Terraform Contexts Tutorial: https://youtu.be/_e9F_QAAMYg