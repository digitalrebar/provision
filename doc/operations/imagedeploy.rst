.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Platform documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Platform; CentOS Image Build & Deployment Guide

.. _rs_imagedeploy:

Building And Deploying CentOS Images
====================================


About
-----
Digital Rebar Platform can be used for image creation, as well as image based deployments. This means drp can
be used to both create, and then deploy "golden master" images. This guide will walk you through what is
required to first build and then deploy a CentOS 7 image using a few items from the global catalog. This guide
will start from an already installed endpoint. If you need assistance getting your endpoint installed
see :ref:`rs_quickstart`.

.. note:: Windows image creation is currently unsupported using this method. Please use a tool like Packer instead.


Expected System State
---------------------
The following lists the DRP version, content packs, plugins, and bootenvs with their versions expected to already be
installed before you begin following this guide.

* Digital Rebar Platform v4.2.x or higher
* drp-community-content v4.2.x or higher
* task-library v4.2.x or higher
* CentOS 7 bootenv iso uploaded (see :ref:`rs_quickstart` for more info)

You should be able to successfully provision a machine with CentOS 7 before proceeding.

Step By Step
------------

Log into the web interface by visiting ``https://<ip_address_of_your_endpoint>:8092/``

Next click the Catalog link (on left scroll down to bottom). Find the "Image Builder" and the "Image Deploy" items and
select the check boxes for each one. Click on the "Install" button on the top.

Now with the image builder and image deploy content installed an image needs to be created for use in the deploy process.
To build an image just run a machine through the image-builder-centos workflow. This workflow will use the centos-7-install
stage, it will also add ssh keys if you have the access-keys param defined, next it will install the drp agent, next it
creates an image from what was deployed, and prepares it to be managed by cloud-init. Finally the image is uploaded into
the drp endpoint before being marked as complete.

To be able to use this image we need to set some params for image-deploy. We need

* image-deploy/image-file
* image-deploy/image-os
* image-deploy/image-type

For the process we just went through the image-os will be "linux" and the image-type will be "tgz". For the value to use
for image-file we can get that from the output of the image-capture task from the image-builder-centos workflow that ran
earlier, or we can get it by looking at the web portals file browser. From the web ui click on the "Files" link on the
left near the bottom. Next on the right pane you should see an "images" folder. Click it. Next copy the name of the
tarball file. For me it is: centos-tarball-20200130-155447-877f710b40.tgz Your file will have different numbers in the name
due to time and date stamping. Now for the image-file value: ``image-deploy/image-file: files/images/centos-tarball-20200130-155447-877f710b40.tgz``

Now that we know the values we need to use for the parameters required for a CentOS-7 deployment we are almost ready to
deploy our image. We just need a workflow, and to set the params. You have several options when it comes to setting
the params. You can either create a profile add the params to the profile then apply the profile to the machine,
or you can apply them directly to the machine object it self, or you can do even more advanced things that are beyond
the scope of this document. The only thing left to do now is to make a basic workflow to deploy the image.

Click on the Workflows link on the left. Near the top click the Add button. Name your workflow something descriptive,
give it a good description, next add the stage "image-deploy", and "complete". That is the minimum required to test our
setup. Click save. Now you can apply this workflow to a machine that has had the above 3 params defined on it.


Customizing Storage Partitioning
--------------------------------

The Param ``curtin/partitions`` describes the disk layout, partitioning, and filesystem
structure that the Image is deployed to for images of type ``TGZ`` (RootFS Tarball).  For
*Raw* type images, the partitioning is for the most part "baked in" to the image itself.

Examples of :ref:`rs_imagedeploy_storage` can be referenced as a starting point for
creating custom partitioning layouts.


Troubleshooting
---------------
When doing CentOS image deployments using UEFI I get:

  ::

    grub2-install: error: /usr/lib/grub/x86_64-efi/modinfo.sh doesn't exist. Please specify --target or --directory.
    Traceback (most recent call last):
          File "/tmp/tmpudYdXr/target/curtin/curtin-hooks.py", line 386, in <module>
            main()
          File "/tmp/tmpudYdXr/target/curtin/curtin-hooks.py", line 365, in main
            grub2_install_efi(target)
          File "/tmp/tmpudYdXr/target/curtin/curtin-hooks.py", line 197, in grub2_install_efi
            '--recheck'])
          File "/usr/lib/python2.7/site-packages/curtin/util.py", line 640, in subp
            return subp(*args, **kwargs)
          File "/usr/lib/python2.7/site-packages/curtin/util.py", line 268, in subp
            return _subp(*args, **kwargs)
          File "/usr/lib/python2.7/site-packages/curtin/util.py", line 140, in _subp
            cmd=args)
        curtin.util.ProcessExecutionError: Unexpected error while running command.
        Command: ['unshare', '--fork', '--pid', '--', 'chroot', '/tmp/tmpudYdXr/target', u'grub2-install', u'--target=x86_64-efi', u'--efi-directory', u'/boot/efi', u'--recheck']
        Exit code: 1
        Reason: -
        Stdout: ''
        Stderr: ''
        finish: cmd-install/stage-curthooks/builtin/cmd-curthooks: FAIL: curtin command curthooks
        Unexpected error while running command.
        Command: ['/tmp/tmpudYdXr/target/curtin/curtin-hooks']
        Exit code: 1
        Reason: -
        Stdout: ''
        Stderr: ''

This means you are missing the ``grub2-efi-x64-modules`` from the image. Installing the missing package should correct this problem. 
