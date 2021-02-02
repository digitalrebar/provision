.. Copyright (c) 2019 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; ESXi Getting Started Guide

.. _rs_esxi_gsg:


ESXi Getting Started Guide
==========================

About
-----
In this document we walk you through the minimal required steps to be able to provision machines with ESXi.


System Setup Information
------------------------

This document will assume that you have followed our :ref:`rs_quickstart` guide, and have installed `drp-community-content` as well as the `task-library`, and you are able to successfully discover machines with sledgehammer.


Getting Started
---------------
The following commands can be run from any machine that has access to your DRP endpoint. This can be your local workstation or the drp endpoint itself. To simplify matters I will be running the commands as root on my drp endpoint.

  ::

    drpcli catalog item install vmware
    drpcli bootenvs list|jq ''.[].Name|grep esxi

In the above steps we install the latest stable version of the vmware plugin. Next we list the bootenvs and grep for esxi. The output from that command shows the available bootenv that we have available
for you to install. Look for the version of esxi you want to install. Just like with the CentOS and Ubuntu examples from the quickstart guide, you will need to upload the iso the bootenv depends on. Unlike
our CentOS and Ubuntu bootenvs with ESXi you have to manually download the ISO file, we can not do that for you. To assist you we include links in the documentation field of the bootenv.

  ::

    drpcli bootenvs show esxi_700-15843807_dell-install|jq -r ''.Documentation
    Provides VMware BootEnv for ESXi 700-15843807 for dell
    For more details, and to download ISO see:

      - https://my.vmware.com/group/vmware/details?downloadGroup=OEM-ESXI70GA-DELLEMC&productId=974

    NOTE: The ISO filename and sha256sum must match this BootEnv exactly.


In the above example I wanted to find the ESXi 7.0 Dell specific ISO. After logging into the vmware download site and saving the iso to my drp end point in the /root directory I need to upload it into the system.

  ::

    drpcli isos upload VMware-VMvisor-Installer-7.0.0-15843807.x86_64-DellEMC_Customized-A00.iso

Similarly to upload the ESXi 7.0 standard release from vmware you would download that iso and run

  ::

    drpcli isos upload VMware-VMvisor-Installer-7.0.0-15843807.x86_64.iso

.. note:: The ISO filename and sha256sum must match the BootEnv exactly.

Once the commands are complete its safe to remove the iso file from your current working directory

  ::

    rm /root/VMware-VMvisor-Installer-7.0.0-15843807.x86_64.iso
    rm /root/VMware-VMvisor-Installer-7.0.0-15843807.x86_64-DellEMC_Customized-A00.iso

With these steps completed its time to move to creating a basic profile that we can apply to a machine before we apply an install workflow to it.

Creating A Basic Profile
------------------------

An example profile is provided by the vmware plugin we installed above. It is not required to make one, but in this example we will show you some options
you might want to use especially while troubleshooting. We will use the default profile as a starting place to make our own.

  ::

    drpcli profiles show esxi-profile > esxi-profile-clone.json
    $EDITOR esxi-profile-clone.json

.. note:: The default profile provided as an example in the vmware plugin WILL NOT function without being customized to fit the needs of your environment.

The changes you need to make in this file are to the "Name" and then we will be removing some of the params leaving behind a minimal profile. Below is what my file looks like now

  ::

      {
          "Available": true,
          "Bundle": "",
          "Description": "(clone me) Sample profile settings for ESXi kickstart install.",
          "Documentation": "Sets some basic Param values that are useful for dev/test deployments\nof VMware vSphere ESXi hypervisors.  Generally speaking these aren't\ngood to set for production systems.\n\nThis profile is intended to be cloned and applied to a Machine(s) for\nsubsequent use.  You can then remove/modify the values appropriate to\nyour use case, after you nave cloned it.\n",
          "Endpoint": "",
          "Errors": [],
          "Meta": {
            "color": "blue",
            "icon": "world",
            "title": "RackN Content"
          },
          "Name": "esxi-profile-clone",
          "Params": {
            "esxi/disk-install-options": "--firstdisk --overwritevmfs",
            "esxi/serial-console": "gdbPort=none logPort=none tty2Port=com1",
            "esxi/shell-local": true,
            "esxi/shell-remote": true,
          },
          "Partial": false,
          "Profiles": [],
          "ReadOnly": false,
          "Validated": true
      }

This profile will enable SSH to assist in troubleshooting should you need it. Next we need to upload this new profile to the endpoint

  ::

    drpcli profiles create esxi-profile-clone.json

The output of this command if successful will be the contents
of the esxi-profile-clone.json file printed to stdout. With this
final step complete we can now apply the new profile to a machine
we have waiting in discovery, then start the esxi install workflow.

  ::

    drpcli machines update Name:esxi-7-test '{"Profiles": ["esxi-profile-clone"], "Workflow": "esxi-install"}'

In this final command we apply the new profile to an existing machine named `esxi-7-test` that was in DRP and had already been discovered and was in `sledgehammer-wait`

Additional Resources
--------------------

This is the most minimal example of how to get started using the vmware plugin. For a comprehensive document which covers available Params, Stages, and more please see: :ref:`rs_cp_vmware`
