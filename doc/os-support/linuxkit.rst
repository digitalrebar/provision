.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Operating Support; LinuxKit

.. _rs_os_linuxkit:

LinuxKit
~~~~~~~~

`LinuxKit <https://github.com/linuxkit/linuxkit>`_ is one of Docker's latest creations to provide an immutable simple OS that
can run containers and other container platforms.  The goal is to provide a simple to update, secure, and maintain OS that
can facilitate container deployments.

It turns out that Digital Rebar Provision can easily deploy LinuxKit images.  The example assets provided with Digital Rebar
Provision contains some examples of `LinuxKit <https://github.com/linuxkit/linuxkit>`_ deployments.

There are three currently available.  Though they are really simple and easily cloned.

* lk-sshd - the example sshd
* lk-k8s-master - the k8s project master image
* lk-k8s-node - the k8s project node image

Here are the following steps that work on an Ubuntu 16.04 desktop with KVM already setup.  The basic overview is:

* Get and start dr-provision
* Configure dr-provision to handle KVM network
* Get and build LinuxKit
* Make some ISOs
* Install BootEnvs
* Run Nodes

Let's hit each of the steps.  For a simple play with it trial, create a directory for both Digital Rebar Provison
and LinuxKit.

Also, there is a `video <https://youtu.be/kITojfeYaPQ>`_ of these steps.

Get, Start, and Configure Digital Rebar Provision
-------------------------------------------------

There are already many pages for this, :ref:`rs_quickstart` and :ref:`rs_install`.  The main thing is to use the
**tip** version at the moment.  You should probably include the *discovery* and *sledgehammer* images.  The KVM
system used (the Digital Rebar test tool, kvm-slave) always PXE boots and machines can be easily added just by
starting them.

For a simple trail, use the install process with the *--isolated* flag.

Something like:

  ::

    mkdir -p lk-dr-trial/dr-provision
    cd lk-dr-trial/dr-provision
    curl -fsSL https://raw.githubusercontent.com/digitalrebar/provision/master/tools/install.sh | bash -s -- --isolated --rs-version=tip install
    # Follow the instructions at the end of the script


Get and Build LinuxKit
----------------------

First, create a directory for everything, clone LinuxKit, and build it.  This assumes that you go installed.

  ::

    cd lk-dr-trial
    git clone https://github.com/linuxkit/linuxkit.git
    cd linuxkit
    make

After a few minutes, you should have a bin directory with **moby** ready to go.

Make some ISOs
--------------

The provided :ref:`rs_model_bootenv` deploying the sshd example and k8s project.  To build these, we need to do a couple of things.

* Edit the *examples/sshd.yml*
  * Replace the "#your ssh key here" line with the contents of your SSH public key.  e.g. ~/.ssh/id_rsa.pub
* Run the moby build command

  ::

    # edit files
    bin/moby build examples/sshd.yml

This will generate an ISO, *sshd.iso*.  Copy this file into the assets/isos directory (creating it if it doesn't exist) in your
dr-provision install directory.

Additionally, we can build the Kubernetes images.  We still need to edit a couple of files.

* cd projects/kubernetes
* Edit the *k8s-master.yml*
  * Replace the "#your ssh key here" line with the contents of your SSH public key.  e.g. ~/.ssh/id_rsa.pub
  * At the end, make sure the *outputs:* section includes " - format: iso-bios".  Append it to the end of the file.
* Edit the *k8s-node.yml*
  * Replace the "#your ssh key here" line with the contents of your SSH public key.  e.g. ~/.ssh/id_rsa.pub
  * At the end, make sure the *outputs:* section includes " - format: iso-bios".  Append it to the end of the file.
* Run the make command

  ::

     cd projects/kubernetes
     # edit files
     make

This will generate two ISO images, *kube-master.iso* and *kube-node.iso*.  Copy these files into the assets/iso directory in your
dr-provision install directory.

Install BootEnvs
----------------

At this point, we can add the :ref:`rs_model_bootenv` to Digital Rebar Provision.

* Change to your Digital Rebar Provision directory and then to the assets directory.
* Run the following

  ::

    cd lk-dr-trial/dr-provision/assets
    export RS_KEY=rocketskates:r0cketsk8ts # or whatever you have it set to.
    ../drpcli bootenvs install bootenvs/lk-sshd.yml
    ../drpcli bootenvs install bootenvs/lk-k8s-master.yml
    ../drpcli bootenvs install bootenvs/lk-k8s-node.yml

This will make all three :ref:`rs_model_bootenv` available for new nodes.

Run Nodes
---------

At this point, you can boot some nodes and run them.  You can have pre-existing nodes or discovered nodes.  This will
use discovered nodes.

First, we start some nodes.  I used my kvm-slave tool that starts KVM on my Digital Rebar Provison network. .e.g. tools/kvm-slave
Anything that PXEs and you can three will work.

Once they are discovered, you will see something like this from **drpcli machines list**

  ::

    [
      {
        "Address": "192.168.124.21",
        "BootEnv": "sledgehammer",
        "Errors": null,
        "Name": "d52-54-54-07-00-00.example.com",
        "Uuid": "4cc8678e-cdc0-48ee-b898-799103840d7f"
      },
      {
        "Address": "192.168.124.23",
        "BootEnv": "sledgehammer",
        "Errors": null,
        "Name": "d52-54-55-00-00-00.example.com",
        "Uuid": "c22a3db3-dba8-4138-8375-7a546c8097e8"
      },
      {
        "Address": "192.168.124.22",
        "BootEnv": "sledgehammer",
        "Errors": null,
        "Name": "d52-54-54-7d-00-00.example.com",
        "Uuid": "d8d5b78a-976b-41c6-a968-31c73ba2b8a4"
      }
    ]

At this point, you should change the BootEnv field to the environment of choice.

  ::

    cd lk-dr-trial/dr-provision
    ./drpcli machines bootenv "4cc8678e-cdc0-48ee-b898-799103840d7f" lk-sshd
    ./drpcli machines bootenv "d8d5b78a-976b-41c6-a968-31c73ba2b8a4" lk-k8s-master
    ./drpcli machines bootenv "c22a3db3-dba8-4138-8375-7a546c8097e8" lk-k8s-node

At this point, you should reboot those kvm instances (close the KVM console window or kill the qemu process).  Once the systems
boot up, you should be able to ssh into them from the account your ssh key is from (as root).

And that is all for the sshd image.

For Kubernetes, you have to do a few more steps. In this example, 192.168.124.22 is the master.  We need to SSH into its kubelet
container and start kubeadm.  Something like this:

  ::

    ssh root@192.168.124.22
    nsenter --mount --target 1 runc exec --tty kubelet sh
    kubeadm-init.sh

This will run for a while and start up the master.  It will output a line that looks like this:

  ::

    kubeadm join --token bb38c6.117e66eabbbce07d 192.168.65.22:6443

This will need to run on each k8s-node.  We will need to SSH into the kubelet on the k8s node.  Something like this:

  ::

    ssh root@192.168.124.23
    nsenter --mount --target 1 runc exec --tty kubelet sh
    kubeadm join --token bb38c6.117e66eabbbce07d 192.168.65.22:6443

We wait for a while and if the KVM instances have internet access, then kubernetes will be up.  The default access for this cluster
is through the kubelet container though others are probably configurable.

  ::

    ssh root@192.168.124.22
    nsenter --mount --target 1 runc exec --tty kubelet sh
    kubectl get nodes


There are ssh helper scripts in the *linuxkit/projects/kubernetes* directory, but I found them to not always work with the latest
k8s containers.
