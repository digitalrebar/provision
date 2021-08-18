.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Air Gapped Install

.. _rs_airgap:

Air Gap Install Instructions
============================


About
-----

In an air gapped install the DRP endpoint will not have access to the internet. In order to have a functional DRP endpoint you will
need to pre-fetch some packages from the internet. This document will cover what the minimal requirements will be to get up and running
in the air gapped environment. We will cover what is needed and where to get it. In this document it will be assumed you will be using
RHEL/CentOS 7 and starting the install of DRP after you have completed at least a "Minimal install" of the operating system. The
commands being run to pre-fetch content are being run from an Ubuntu workstation machine with internet access, a web browser, and wget,
but any workstation can be used as long as it has internet access and a way to download and save files from the internet.


Additional Operating System Packages Required
---------------------------------------------

The following packages will need to be installed on your air gapped endpoint:

* bsdtar
* libarchive (required by bsdtar)

It is assumed that these packages are already installed on the air gapped DRP endpoint at the time of the DRP install.


RackN & DRP Requirements
------------------------
The following is a list of DRP and RackN requirements needed for a successful air gapped install.

* DRP
* DRP Installer
* RackN Catalog
* RackN Portal (URL Provided by RackN support) *Optional*
* DRP Content
* RackN License
* Sledgehammer
* CentOS (*Optional*)

Step By Step
------------

The following commands should be run on the internet accessible workstation. The file versions will change so please do
not blindly copy & paste these commands.

  ::

    mkdir ~/air_gap
    cd ~/air_gap
    wget -O - https://rebar-catalog.s3-us-west-2.amazonaws.com/rackn-catalog.json | gunzip -c > rackn-catalog.json
    wget $(jq -r '.sections.catalog_items."drp-stable".Source' rackn-catalog.json)
    wget -O install.sh get.rebar.digital/stable
    chmod +x install.sh
    mkdir task-library
    cd task-library
    wget $(jq -r '.sections.catalog_items."task-library-stable".Source' ../rackn-catalog.json)
    cd ..
    mkdir drp-community-content
    cd drp-community-content
    wget $(jq -r '.sections.catalog_items."drp-community-content-stable".Source' ../rackn-catalog.json)
    mkdir ../sledgehammer
    cd ../sledgehammer
    grep -i sledgehammer ../drp-community-content/*.json | grep -i isourl | awk '{print $2}' | sort -u
    # This should yield 2 urls.
    # Download the tar file from each and save in this sledgehammer dir.
    wget (url1)
    wget (url2)
    mkdir ../centos
    cd ../centos
    # Note that this should be the most recent version
    # from the drp-community-content
    wget http://mirror.rackspace.com/CentOS/7/isos/x86_64/CentOS-7-x86_64-Minimal-1908.iso
    cd ..
    mkdir ux
    cd ux
    wget (url provided by RackN support)


The last thing required is the RackN License. At this time you need to either generate the license from an existing
internet accessible DRP endpoint, or by contacting RackN support to have a license generated for you. Once you have the
license file save it in the `~/air_gap` dir as `rackn-license.json`

Next you need to move the contents of `~/air_gap` to the air gap system. *The next set of commands are run on the air gapped system.*

  ::

    ./install.sh --zip-file=v4.1.3.zip --systemd=true --drp-home-dir=/srv/dr-provision --local-ui=true --no-content=true install
    cd ~/air_gap
    drpcli contents upload rackn-license.json
    cd drp-community-content
    drpcli contents upload v4.1.2.json
    cd ../task-library
    drpcli contents upload v4.1.2.json

    systemctl stop dr-provision

    cd /srv/dr-provision/tftpboot/isos
    mv ~/air_gap/sledgehammer/*.tar .
    mv ~/air_gap/centos/*.iso .

    cd /srv/dr-provision/tftpboot/files/ux
    bsdtar -xf ~/air_gap/ux/file.zip

    systemctl start dr-provision

Now to verify the portal is working. On a machine that can access the air gapped endpoint open a web browser and
visit `https://<YOUR IP>:8092/`. By default a self signed cert will need to be accepted (Note that you can provide
your own certs during deployment.) Log in using the default user and password. Once logged in we will verify the
portal is functioning properly by doing some final customizations. Set the default preferences for workflow.

  ::

    Click "Info & Preferences" on the left hand navigation
    On the right side of screen set "Default Workflow" to "discover-base"
    Set "Default Stage" to "discover"
    Set "Default BootEnv" to "sledgehammer"
    Set "Unknown BootEnv" to "discovery"
    Click save icon on the right (shaped like a floppy disk)


.. note::

    These tasks can be completed using the cli or api directly. We are using the portal here to test functionality of
    our self-hosted portal and for ease of configuration.

..

Next add and ssh key to the global profile.

  ::

    Click "Profiles" on the left hand navigation
    Click on "global"
    Click "Edit"
    Add "access-keys"
    Edit the "value" of "access-key" and place your ssh pub key in this value
    Click "Save"


Next add `package-repositories` to the global profile.

  ::

    Click "Edit" on the "global" profile
    Add "package-repositories"
    Edit the value of "package-repositories"
    Example:
    [
      {
        "arch": "x86_64",
        "installSource": true,
        "os": [
            "centos-7",
            "centos-7-install"
        ],
        "tag": "centos-7",
        "url": "https://10.0.0.10:8091/centos-7/install"
      }
    ]
    Click "Save"

Next you need to configure a subnet.

  ::

    Click "Subnets" on the left hand navigation
    Click "Add"
    Create Subnet From Interface
    Click "Use Interface"
    If defaults are acceptable scroll down and click "Add"

If you plan on using the RackN UX in air gapped mode, you will have to download the :ref:`rs_cp_ux_views` plugin if it is not already installed. You must also have a license with the "enable-airgap" feature enabled.

  ::

    Click "Catalog" on the left hand navigation
    Search for "ux-views"
    Click the green download button
    After the plugin finishes installing, refresh the webpage (CTRL+R)
    Click "UX Config" on the left hand navigation
    Click the "airgap" row under "Core" in the "UX" section
    Toggle the "airgap" config switch and click the blue save button.


Next its time to power on machine to deploy. Power it on, make sure its setup to PXE boot in the bios as its default
boot device. Once the machine has been discovered you should see it show up on the "Machines" view in the portal. You
may need to hit refresh. Once the machine shows up you can provision it by setting the workflow of the machine to
"centos-base". That will cause the machine to reboot which as long as the machine is set to PXE boot the CentOS 7 install
will begin. Once it has completed you can log in to verify everything worked correctly.

  ::

    ssh root@my_test_machine
    cat /etc/yum.repos.d/*

This should show that the only repo configured matches what you defined above in the "package-repositories" parameter.