:orphan:

.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; FAQ
  pair: Digital Rebar Provision; Troubleshooting

.. warning::
  This FAQ / Troubleshooting guide is being deprecated.  Please use the new
  :ref:`rs_knowledge_base` for the most up to date references.

.. All below cross-reference links were commented out, and have been replaced with the appropriate
   knowledge base cross-reference of the same value.

.. _rs_faq:

FAQ / Troubleshooting
~~~~~~~~~~~~~~~~~~~~~

The following section is designed to answer frequently asked questions and help troubleshoot Digital Rebar Provision installs.

Want ligher reading?  Checkout our :ref:`rs_fun`.

.. .. _rs_old_faq_drpcli:

Where can I get the DRPCLI?
---------------------------

If you need the DRPCLI, follow instructions for :ref:`rs_cli_download`.

.. .. _rs_faq_security:

Where can I learn more about Digital Rebar Security?
----------------------------------------------------

Please consult the dedicated :ref:`rs_security_faq`.

.. .. _rs_bind_error:

Bind Error
----------

Digital Rebar Provision will fail if it cannot attach to one of the required ports.

* Typical Error Message: "listen udp4 :67: bind: address already in use"
* Additional Information: The conflicted port will be included in the error between colons (e.g.: `:67:`)
* Workaround: If the conflicting service is not required, simply disable that service
* Resolution: Stop the offending service on the system.  Typical corrective actions are:

  * 67 - dhcp.  Correct with `sudo pkill dnsmasq`

See the port mapping list on start-up for a complete list.


.. .. _rs_tftp_error:

TFTP Error
----------

In the dr-provision logfiles you may occassionally see error messages relating to ``TFTP Aborted``.  These
errors are (typically) benign and expected behavior.  The TFTP protocol does not specify a mechanism to
obtain the size of a file to transfer for calculating completed transfer; without first requesting the file.
Digital Rebar Provision initiates the transfer request an then immediately aborts it.  This obtains the
file size for the next transfer to validate the file was served correctly.

Simply ignore these errors.  If you receive these errors and you believe you should be provisioning correctly,
check that you have correctly specified the default/unknown BootEnv, default Stage, and default Workflow
are set correctly.

error messages may appear similarly to:

  ::

    May 24 13:48:22 ubuntu dr-provision[7092]: dr-provision2018/05/24 20:48:22.006224 [280:13]static [error]: /home/travis/gopath/src/github.com/rackn/provision-server/v4/midlayer/tftp.go:82
    May 24 13:48:22 ubuntu dr-provision[7092]: [280:13]TFTP: lpxelinux.0: transfer error: sending block 0: code=0, error: TFTP Aborted


.. .. _rs_gen_cert:

Generate Certificate
--------------------

Sometimes the cert/key pair in the github tree is corrupt or not sufficient for the environment.  The following command can be used to rebuild a local cert/key pair.

  ::

    sudo openssl req -new -x509 -keyout server.key -out server.crt -days 365 -nodes

It may be necessary to install the openssl tools.


.. .. _rs_add_ssh:

Add SSH Keys to Authorized Keys
-------------------------------

VIDEO TUTORIAL: https://www.youtube.com/watch?v=StQql8Xn08c

NOTE: The UX offers a "Add Key" button on the top of the "Info & Preference" panel that will perform the following steps for you.

To have provisioned operating systems (including discovery/sledgehammer) add SSH keys, you should set the ``access-keys`` parameter with a hash of the desired keys.  This Param should be applied to the Machines you wish to update, either directly via adding the Param to the Machines, or by adding the Param to a Profile that is subsequently added to the Machines.  NOTE that the ``global`` Profile applies to all Machines, and you can add it to ``global`` should you desire to add the set of keys to ALL Machines being provisioned.

The below example adds *User1* and *User2* SSH keys to the profile *my-profile*.  Change appropriately for your enviornment.

  ::

    cat << END_KEYS > my-keys.json
    {
      "Params": {
        "access-keys": {
          "user1": "ssh-rsa user_1_key user1@krib",
          "user2": "ssh-rsa user_2_key user2@krib"
        }
      }
    }
    END_KEYS

    drpcli profiles update my-profile keys.json


.. .. _rs_docker_volume:

Example Docker Volume Usage
---------------------------

Digital Rebar Provision writes content in the Docker Container to the ``/provision/drp-data``
directory by default.  Most DRP Endpoint provisioning systems will want to have persistent
data across the container runtimes.  For this, you need to add a Docker Volume.  The below
example shows you how to use the local Docker host as the backing store for the volume. You
can also use any of the container based networked storage solutions to back your volume in.


1. Create a volume for the container

  ::

    export VOL="drp-data"

    # create a Docker volume
    docker volume create $VOL

2. Let's verify that the volume is currently empty

  ::

    docker volume inspect $VOL | jq '.[].Mountpoint'
    # outputs:
    # "/docker/volumes/drp-data/_data"

    # show the contents of the current (empty) volume
    ls -la $(docker volume inspect $VOL | jq -r '.[].Mountpoint')
    # total 0
    # drwxr-xr-x. 2 root root  40 Aug 21 00:41 .
    # dr-xr-x---. 1 root root 180 Aug 21 00:41 ..

3. Launch DRP, using our newly created volume:

  ::

    # now run DRP with our volume mapped to /provision/drp-data:
    docker run --volume $VOL:/provision/drp-data --name drp -itd --net host digitalrebar/provision:stable

4. Verify that DRP extracted the assets on the host in the mounted volume location:

  ::

    # when DRP starts up, it extracts and builds the default assets
    # in the writable backing store (directory structure):
    ls $(docker volume inspect drp-data | jq -r '.[].Mountpoint')
    # outputs:
    # digitalrebar  job-logs  plugins  replace  saas-content  secrets  tftpboot  ux


.. .. _rs_access_ssh_root_mode:

Set SSH Root Mode
-----------------

The Param ``access-ssh-root-mode`` defines the login policy for the *root* user.  The default vaule is ``without-password`` which means the remote SSH *root* user must access must be performed with SSH keys (see :ref:`rs_add_ssh`).  Possible values are:

========================  ==========================================================
value                     definition
========================  ==========================================================
``without-password``      require SSH public keys for root login, no forced commands
``yes``                   allow SSH *root* user login with password
``no``                    do not allow SSH *root* user login at all
``forced-commands-only``  only allow forced commands to run via remote login
========================  ==========================================================


.. .. _rs_default_password:

What are the default passwords?
-------------------------------

When using the community BootEnvs for installation, the password is set to a variant of ``RocketSkates``.  See :ref:`rs_configuring_default` for complete details.

For all bootenvs (sledgehammer, centos, ubuntu, etc.) the default pattern does NOT allow login via Password.  See :ref:`rs_add_ssh` for manaing SSH Authorized Keys login details.

We *strongly* recommend changing this default or, better, using SSH ``without-password`` options as per :ref:`rs_access_ssh_root_mode` details.

A quick reference table for passwords:

========================  ============  ============
use                       user          password
========================  ============  ============
``drp endpoint auth``     rocketskates  r0cketsk8ts
``sledgehammer``          root          rebar1
``most bootenvs`` (*)     root          RocketSkates
``debian`` / ``ubuntu``   rocketskates  RocketSkates
``cloud-init`` images     <varies> (*)  RocketSkates
``VMware ESXi``           root          <generated>
========================  ============  ============

For more information

.. note:: ``(*)`` "most bootenvs" and ``cloud-init`` images refers to CentOS, Ubuntu, CoreOS, ESXi, etc.  Generally speaking, this is the default "installed" credentials.  Note that each distro has it's own rules about ``root`` versus installed default user accounts.  DRP follows most vendors "patterns" with regards to ``root`` -vs- unprivileged user creation, with the username changed to "rocketskates".  Some vendor specific notes are below.

For ``debian / ubuntu`` bootenvs, the default user (``rocketskates``, can be changed by setting ``provisioner-default-user`` Param), has ``sudo`` privileges.

For Images with ``cloud-init`` pieces, there often is an injected ``centos`` user for CentOS, ``ubuntu`` for Ubuntu, etc. user.  This is controlled by the ``cloud-init`` configurations of the image build process.

For more about VMware, see :ref:`vmware_esxi_passwords`


.. .. _rs_drpclirc:

Using the ``.drpclirc`` File
----------------------------

If you need the DRPCLI, follow instructions for :ref:`rs_cli_download`.

In addition to the environment variables (eg ``RS_ENDPOINT``, ``RS_KEY``, etc) and setting explicit ``drpcli`` values via option flags (eg ``--enpdoint``, ``-E``, etc), you can now use a home _RC_ style configuration file to set these values.

To do so, create a file ``$HOME/.drpclirc`` with the following possible values and format:

====================== ============================================================================
value                  notes
====================== ============================================================================
``RS_ENDPOINT``        sets the endpoint API location (default: https://127.0.0.1:8092)
``RS_USERNAME``        sets username to auth to the Endpoint (default: "rocketskates")
``RS_PASSWORD``        sets the password for the auth (default: "r0cketsk8ts")
``RS_KEY``             sets user:pass pair for authentication (default: "rocketskates:r0cketsk8ts")
``RS_TOKEN``           a precreated Token (which may have a specific use scope)
``RS_FORMAT``          command line output format to use (json,yaml,text,table)
``RS_PRINT_FIELDS``    comma separate list of fields to show in output "table" or "text" format
``RS_NO_HEADER``       remove the header fields from "table" or "text" format output
``RS_TRUNCATE_LENGTH`` limits the length of fields displayed for "table" or "text" output formats
====================== ============================================================================

Example:
  ::

    RS_ENDPOINT=https://10.10.10.10.8092
    RS_PASSWORD=super_secure_secret_password_don't_share_with_anyone

Please note that you can not use Shell style ``export`` in front of the variable,
and do NOT surround the value with double or single quotes.

.. note:: The RS_FORMAT, RS_PRINT_FIELDS, RS_NO_HEADER, and RS_TRUNCATE_LENGTH variables are only valid for ``drpcli`` *v4.2.0-beta2.0* or newer versions.


.. .. _rs_human_formatters:

Using Table/Text Output Formatters
----------------------------------

As of ``v4.2.0-beta2.0``, the ``drpcli`` client supports additional output formats of *table* and *text* type.

Examples:
  ::

    drpcli --format table ...
    # or
    drpcli --format text ...

These output formats can be configured by setting environment shell variables,
or use of the .drpclirc (see: :ref:`rs_drpclirc`) file for setting default usage.

The shell environment variables are as follows.

====================== ============================================================================
value                  notes
====================== ============================================================================
``RS_FORMAT``          command line output format to use (json,yaml,text,table)
``RS_PRINT_FIELDS``    comma separate list of fields to show in output "table" or "text" format
``RS_NO_HEADER``       remove the header fields from "table" or "text" format output
``RS_TRUNCATE_LENGTH`` limits the length of fields displayed for "table" or "text" output formats
====================== ============================================================================

Examples of setting environment variables:
  ::

    export RS_FORMAT=table
    export RS_PRINT_FIELDS=Name,Uuid,Workflow,Stage,BootEnv
    export RS_NO_HEADER=true
    export RS_TRUNCATE_LENGTH=30

Examples of ``drpcli`` usage:
  ::

    drpcli subnets list --format table
    drpcli machines list --format table --print-fields Name,Uuid,Workflow,Stage,BootEnv --no-header
    drpcli machines list --format table --print-fields Name,Uuid,Workflow,Stage,BootEnv
    drpcli profiles list --format table --truncate-length 30
    drpcli machines params Name:jane --format=table --truncate-length=120
    drpcli extended -l endpoints list --format table --truncate-length 30

CLI FAQ:
--------

The CLI has a dedicated FAQ section.  Please see :ref:`rs_cli_faq`.

Topics include:
  * :ref:`rs_autocomplete`
  * :ref:`rs_cli_faq_zip`
  * :ref:`rs_download_rackn_content`


.. .. _rs_more_debug:

Turn Up the Debug
-----------------

To get additional debug from dr-provision, set debug preferences to increase the logging.  See :ref:`rs_model_prefs`.


.. .. _rs_vboxnet:

Missing VBoxNet Network
-----------------------

Virtual Box does not add host only networks until a VM is attempting to use them.  If you are using the interfaces API (or UX wizard) to find available networks and ``vboxnet0`` does not appear then start your VM and recreate the address.

Virtual Box may also fail to allocate an IP to the host network due to incomplete configuration.  In this case, ``ip addr`` will show the network but no IPv4 address has been allocated; consequently, Digital Rebar will not report this as a working interface.


.. .. _rs_vbox_no_boot:

VirtualBox "no bootable medium" on second boot
----------------------------------------------

VirtualBox PXE firmware does not handle PXE chaining effectively.  This happens because DRP treats known and unknown machines differently so the first boot gets more different boot instructions.

The workaround is to use DHCP option 67 to supply the correct boot file.  Setting DHCP option 67 to `lpxelinux.0` bypasses the chainloader after the machine has registered.

See also :ref:`rs_uefi_boot_option`


.. .. _rs_debug_sledgehammer:

Debug Sledgehammer
------------------

If the sledgehammer discovery image should fail to launch Runner jobs successfully, or other issues arise with the start up sequences, you can debug start up via the systemd logging.  Log in to the console of the Machine in question (or if SSH is running and you have ``access-keys`` setup, you can SSH in), and run the following command to output logging:
  ::

      journalctl -u sledgehammer


.. .. _rs_convert_to_production_mode:

Convert Isolated Install to Production Mode
-------------------------------------------

There currently is no officually supported *migration* tool to move from an ``Isolated`` to ``Production`` install mode.  However, any existing customizations, Machines, Leases, Reservations, Contents, etc. can be moved over from the Isolated install directory structure to a Production install directory, and you should be able to retain your Isolated mode environment.

All customized content is stored in subdirectories as follows:

  Isolated: in ``drp-data/`` in the Current Working Directory the installation was performed in
  Production:  in ``/var/lib/dr-provision``

The contents and structure of these locations is the same.  Follow the below procedure to safely move from Isolated to Production mode.

#. backup your current ``drp-data`` directory (eg ``tar -czvf /root/drp-isolated-backup.tgz drp-data/``)
#. ``pkill dr-provision`` service
#. perform fresh install on same host, without the ``--isolated`` flag
#. follow the start up scripts setup - BUT do NOT start the ``dr-provision`` service at this point
#.  copy the ``drp-data/*`` directories recursively to ``/var/lib/dr-provision`` (eg: ``unalias cp; cp -ra drp-data/* /var/lib/dr-provision/``)
#. make sure your startup scripts are in place for your production mode (eg: ``/etc/systemd/system/dr-provision.service``)
#. start the new production version with  ``systemctl start dr-provision.service``
#. verify everything is running fine
#. delete the ``drp-data`` directory (suggest retaining the backup copy for later just in case)

.. note::  WARNING:  If you install a new version of the Digital Rebar Provision service, you must verify that there are no Contents differences between the two versions.  Should the ``dr-provision`` service fail to start up; it's entirely likely that there may be some content changes that need to be addressed in the JSON/YAML files prior to the new version being started.  See the :ref:`rs_upgrade` notes for any version-to-version specific documentation.


.. .. _rs_customize_production_mode:

Customize Production Mode
-------------------------

You can use systemd drop configuration to alter dr-provision start up options.

To use, figure out the environment variable to set by checking the help of dr-provision.  e.g. dr-provision -h

You will need to create the drop-in directory if it doesn't exist.

* mkdir -p /etc/systemd/system/dr-provision.service.d

Then you will need to create a drop-in service file.  For example, to name your system, you would use this file, drpid.conf:

  ::

     [Service]
     Environment=RS_DRP_ID=mydrpserver

Then reload and restart the service.

* sudo systemctl daemon-reload && sudo systemctl restart dr-provision

This will work with multiple files and multiple variables.


.. .. _rs_kickseed:

Custom Kickstart and Preseeds
-----------------------------

Starting with ``drp-community-content`` version 1.5.0 and newer, you can now define a custom Kickstart or Preseed (aka *kickseed*) to override the defaults in the selected BootEnv.  You simply need to only define a single Param (``select-kickseed``) with the name of the Kickstart or Preseed you wish to override the default value.
  ::

    export UUID="f6ca7bb6-d74f-4bc1-8544-f3df500fb15e"
    drpcli machines set $UUID param select-kickseed to "my_kickstart.cfg"

Of course, you can apply a Param to a Profile, and apply that Profile to a group of Machines if desired.

.. note:: The Digital Rebar default kickstart and preseeds have Digital Rebar specific interactions that may be necessary to replicate.  Please review the default kickstart and preseeds for patterns and examples you may need to re-use.   We HIGHLY recommend you start with a `clone` operation of an existing Kickstart/Preseed file; and making appropriate modifications from that as a baseline.


.. .. _rs_kexec:

Can I eliminate reboots with kexec?
-----------------------------------

Yes.  Setting the `kexec-ok` param to `true` on the `global` or machine specific profile allows
BootEnvs that are kexec enabled to skip rebooting when changing to that BootEnv.  For example,
Sledgehammer enables kexec and can be started without a reboot from Linux environments.

This is a Linux specific feature.  For more about kexec: https://wiki.archlinux.org/index.php/Kexec

.. .. _rs_plugin_providers_license:

Import plugin failed pool: define failed
----------------------------------------

If you are using the DRPCLI to upload a licensed RackN plugin, the endpoint will reject the upload with a defined failed error.

Install the license content pack and try again.  If you've saved the `rackn-license.json` file then you can use the DRPCLI to upload it via `drpcli contents upload rackn-license.json`.


.. .. _rs_update_content_command_line:

Update Community Content via Command Line
-----------------------------------------

Here's a brief example of how to upgrade the Community Content installed in a DRP Endpoint using the command line.  See :ref:`rs_download_rackn_content` for additional steps with RackN content.

Perform the following steps to obtain new content.

View our currently installed Content version:
  ::

    $ drpcli contents show drp-community-content | jq .meta.Version
      "v1.4.0-0-ec1a3fa94e41a2d6a83fe8e6c9c0e99c5a039f79"

Get and upload our new version (in this example, explicitly set version to ``v1.5.0``.  However, you may also specify ``stable``, or ``tip``, and do not require specific version numbers for those.
  ::

    drpcli contents upload catalog:drp-community-content-v1.5.0

or
  ::
      drpcli catalog item install drp-community-content --version v1.5.0

Now verify that our installed content matches the new vesion we expected ...
  ::

    $ drpcli contents show drp-community-content | jq .meta.Version
      "v1.5.0-0-13f1aff688b53d5dfdab9a1a0c1098bd3c6dc76c"

Next check if sledgehammer needs to be updated
  ::

    drpcli bootenvs uploadiso sledgehammer

This command will update the sledgehammer isos if required. If your output from the above command is
  ::

    BootEnv sledgehammer already has all required ISO files

 Then your sledgehammer is all up to date.


.. .. _rs_reboot_faq:

Rebooting inside a Tasks, Stages and Workflows
----------------------------------------------

The Runner Task execution system supports many ways to cause a system reboot that allow for the task being marked as either complete or incomplete (so it can resume).  This can be very important for tasks that require a reboot mid-task.

These options are handled by using script helpers or sending specialized ``exit`` codes.  Please see :ref:`rs_workflow_reboot` for comprehensive documentation.


.. .. _rs_reboot_wo_ipmi:

Rebooting without IPMI plugins (without a Task)
-----------------------------------------------

The Runner will automatically reboot the system if the BootEnv changes during a Workflow.  You can force this behavior by changing the BootEnv to `local` on the machine manually then starting a Workflow with a different BootEnv like `discover`.  This will cause the runner to reboot the machine.


Steps:
  #. Clear the Machine Workflow
  #. Set the Machine BootEnv to `local`
  #. Update
  #. Set the Workflow to a workflow with a different BootEnv.
  #. Update and watch machine reboot


.. .. _rs_nested_templates:

Nested Templates (or "Sub-templates")
-------------------------------------

The Golang templating language does not provide a call-out to include another template.  However, at RackN, we've added the ability to include *nested templates* (sometimes referred to as *sub-templates*).  In any content piece that is valid to use the templating capabilities, simply use the following Template construct to refer to another     template.  The template referred to will be expanded inline in the calling template.  The nested template example below calls the template named (oddly enough) *nested.     tmpl*.
  ::

    {{template "nested.tmpl" .}}

    # or alternatively:

    {{$templateName := (printf "part-seed-%s.tmpl" (.Param "part-scheme")) -}}
    {{.CallTemplate $templateName .}}

The ``template`` construct is a text string that refers to a given template name which exists already.

The ``CallTemplate`` construct can be a variable or expression that evaluates to a string.


.. .. _rs_sprig:

How Can I manipulate values during Golang Template rendering?
-------------------------------------------------------------

The Digital Rebar Provision integrates most of the `Sprig Function Library <http://masterminds.github.io/sprig/>`_ in the Golang Template rendering operations.  That means that you may include their string, math and flow functions into your pipelines.

For example: `{{.Param "noCamelCase/hashiCorp" | snakecase }}` or `{{.Param "cool/tech" | regexMatch "([DRP]*)"}}`

Please consult the Sprig website for a full list of functions.

Note: Digital Rebar Provision blocks functions that could be used to operate on the endpoint outside of DRP template rendering for security reasons.


.. .. _rs_double_brace:

How Can I render Double Curly Braces `{{` and `}}` during Golang Template rendering?
------------------------------------------------------------------------------------

Moved to :ref:`rs_kb_00025`

.. .. _rs_change_machine_name:

Change a Machines Name
----------------------
If you wish to update/change a Machine Name, you can do:
  ::

    export UUID="abcd-efgh-ijkl-mnop-qrst"
    drpcli machines update $UUID '{ "Name": "foobar" }'

.. note:: Note that you can NOT use the ``drpcli machines set ...`` construct as it only sets Param values.  The Machines name is a Field, not a Parameter.  This will NOT work: ``drpcli machines set $UUID param Name to foobar``.


.. .. _rs_reservation_set_hostname:

Set `hostname` in a DHCP Reservation
------------------------------------

If you create a DHCP Reservation for a system (or convert an active Lease to Reservation), you can also set the Hostname for the Machine.  If you are pre-creating Reservations, this will allow you to have a pre-set hostname when the Machine first comes up.  Additionally, if you create/destroy your machine objects, but would like a hostname to persist with the Machine Reservation when the machine returns, you can do this.

.. note:: The UX version (at least as of v1.2.1 and older) does not support setting DHCP options to the Reservation.  You will have to perform these actions using either the CLI or API.  The CLI method is outlined below.

This procedure assumes you have a Reservation created already, and we are going to update the existing Reservation.  You can combine this procedure with creating a new Reservation, but only if you perform the operation via the CLI or API.

  ::

    # show the current Reservation:
    drpcli reservations show 192.168.8.100

    # create a Hostname specification in the DHCP Options section of the reservation:
    drpcli reservations update 192.168.8.100 '{ "Options": [ { "Code": 12, "Value": "pxe-client-8-100" } ] }'

In the above exmaple, we are assuming our DHCP Reservation is for a Reservation identified by the IP Address ``192.168.8.100``, and that we are setting the hostname (DHCP Option 12) to ``pxe-client-8-100``.


.. .. _rs_uefi_boot_option:

UEFI Boot Support - Option 67
-----------------------------
Starting with v3.7.1 and newer, a DHCP Subnet specification will try to automatically determine the correct values for the ``next-server`` and *DHCP Option 67* values.  In most cases, you shouldn't need to change this or set these fields.  Older versions of DRP may need the ``next-boot`` and/or the *DHCP Option 67* values set to work correctly.  This is especially true of Virtualbox environments prior to v3.7.1.  You will need to force the *DHCP Option 67* to ``lpxelinux.0``.

The DHCP service in Digital Rebar Provision can support fairly complex boot file service.  You can use advanced logic to ensure you send the right PXE boot file to a client, based on Legacy BIOS boot mode, or UEFI boot mode.  Note that UEFI boot mode can vary dramatically in implementations, and some (sadly; extensive) testing may be necessary to get it to work for your system.  We have several reports of field deployments with various UEFI implementations working with the new v3.7.0 and newer "magic" Option 67 values.

Here is an example of an advanced Option 67 parameter for a DHCP Subnet specification:

  ::

    {{if (eq (index . 77) "iPXE") }}default.ipxe{{else if (eq (index . 93) "0")}}ipxe.pxe{{else}}ipxe.efi{{end}}

If you run in to issues with UEFI boot support - please do NOT hesitate to contact us on the `Slack Channel <https://www.rackn.com/support/slack>`_ as we may have updated info to help you with UEFI boot support.

An example of adding this to your Subnet specification might look something like:
  ::

    # assumes your subnet name is "eth1" - change it to match your Subnet name:
    # you may need to delete the existing value if there is one, first, by doing:
    # drpcli subnets set eth1 option 67 to null # The setting to null is not needed with v3.7.1 and beyond.
    drpcli subnets set eth1 option 67 to '{{if (eq (index . 77) "iPXE") }}default.ipxe{{else if (eq (index . 93) "0")}}ipxe.pxe{{else}}ipxe.efi{{end}}'


.. note:: You should not have to add option 67 unless you are meeting a specific need.  Test without it first!


.. .. _rs_lpxelinux_no_such_file:

lpxelinux.0 error: no such file or directory
--------------------------------------------

After TFTPing lpxelinux.0, logs (or network packet traces) may show an error similar to:
  ::

    477    0.378296662    10.10.20.76    10.10.31.96    TFTP    159    Error Code, Code:
    File not found, Message: open /var/lib/dr-provision/tftpboot/pxelinux.cfg/16089a59-9abd-48c2-850a-2ac3bc134935: no such file or directory``

This is expected behavior that is standard PXE *waterfall* searching for a valid filename to boot from.  For full reference, please see the `syslinux <http://www.syslinux.org/>`_ reference documentation, at:

    http://www.syslinux.org/wiki/index.php?title=PXELINUX#Configuration

The expected behavior is for a client to attempt to download files in the following order:

    #. client id (DRP does not use this option, which is what generates the error)
    #. mac address (in the form of ``01-88-99-aa-bb-cc-dd``)
    #. ip  address in uppercase Hexadecimal format, stepping through IP, subnet, and classful boundaries
    #. fall back to the default defined file

Due to this behavior, filenames will be specified that do not exist, and the error message related to that probe request is a normal message.  This is NOT an indicator that provisioning is broken in your environment.


.. .. _rs_different_pxelinux_version:

Change Pxelinux Versions
------------------------

DRP ships with two versions of PXELinux, 6.03 and 3.86.  The default operation is to use 6.03 as lpxelinux.0 with
all the supporting files present in the tftpboot root directory.  This does not always work for all environments.
It is sometimes useful to change this.  In general, DRP attempts to serve iPXE based bootloaders through the
default DHCP operations.  Again, this is not always possible.

The 3.86 version is a single file shipped as esxi.0.

There are couple of ways to change the operation.

First, the file, esxi.0, can be used by changing the bootfile option in DHCP server.  For DRP, this can be at
the subnet or reservation level.

Second, the lpxelinux.0 file can be replaced.  To do this safely, a couple of steps need to be done.

#. In the tftpboot directory, copy lpxelinux.0 to lpxelinux.0.bak.
#. In the replace direcotry, copy esxi.0 to lpxelinux.0.  The replace directory is usually a peer to the tftpboot
   directory.
#. In the tftpboot directory, copy esxi.0 to lpxelinux.0.

The middle step keeps DRP from overwriting your changes on startup.


.. .. _rs_render_kickstart_preseed:

Render a Kickstart or Preseed
-----------------------------

Kickstart and Preseed files only created by request and are not stored on a filesystem that is viewable.  They are dynamically generated on the fly, and served from the virtual Filesystem space of the Digital Rebar HTTP server (on port 8091 by default).  However, it is possible to render a kickstart or preseed to evaluate how it is going to operate, or troubleshoot issues with your config files.

When a machine is in provisioning status, you can view the dynamically generated preseed or kickstart from the TFTP server (or via the HTTP gateway).  Provisioning status means the Machine has been placed into an installable BootEnv via a Stage.  If (for exaxmple) placed in to ``centos-8-install`` Stage, the ``compute.ks`` can be rendered for the machine.  Or, if placed in to ``ubuntu-18.04-install`` Stage, the ``seed`` can be rendered for the machine.

Get the Machine ID, then use the following constructed URL:
  ::

    MID="7f65279a-7e5c-4e69-af40-dd01af4c5667"
    DRP="10.10.10.10"
    TYPE="seed"   # seed for ubuntu, or compute.ks for centos

    http://${DRP}:8091/machines/${MID}/${TYPE}


Example URLs:

  ubuntu/debian:
    http://10.10.10.10:8091/machines/7f65279a-7e5c-4e69-af40-dd01af4c5667/seed

  centos/redhat:
    http://10.10.10.10:8091/machines/7f65279a-7e5c-4e69-af40-dd01af4c5667/compute.ks

.. note:: A simple trick ... you can create a non-existent Machine, and place that machine in different BootEnvs to render provisioning files for testing purposes.  For example, put the non-existent Machine in the ``centos-8-install`` Stage, then render the ``compute.ks`` kickstart URL above.

. _rs_render_does_not_explode_iso:

BootEnv Does Not "Explode ISO" after upload
-------------------------------------------

Problem: New or cloned BootEnv does not explode the uploaded ISO components in Digital Rebar after upload.

Possible Causes:

1. Uploaded ISOs must match shasums in the BootEnv
1. Install BootEnvs must end in ``-install``. See :ref:`rs_model_bootenv`
1. BootEnvs intended for network booting, _must_ include netbootable components.

.. .. _rs_ubuntu_local_repo:

Booting Ubuntu Without External Access
---------------------------------------

Default Ubuntu ISOs will attempt to check internet repositories, this can cause problems during provisioning if your environment does not have outbound access.

To workaround this problem, you need to supply a DNS and gateway for your subnet.  There are several ways to do this:

1. Internal to Digital Rebar: Define Options 3 (Gateway) and 6 (DNS) for your machines' Subnet.
2. External to Digital Rebar: Adding ``default_route=true`` to the boot parameters and include a DNS server on the local subnet in DHCP.


.. .. _rs_wget_timeout:

Network Unreachable from Wget / Second Stage Timeout
----------------------------------------------------

Throwing a ‘network unreachable’ error from `wget` when trying to fetch second stage initramfs; however, by the time you get dropped into a root console, eth0 has an IP address and can connect to the server fine.  May also see a baremetal PXE boot initial PXE boot works but then it's getting kicked to a shell before it can download root.squashfs.

Troubleshooting: You can manually grab the file with ``wget`` after it bails, so communications are working fine. It just appears it's not waiting long enough for DHCP and then fails to get the file before it gets an IP.

Note: You can set these changes the global profile so it will apply everywhere.  It shouldn’t hurt functioning systems (they will escape the loop early) and might fix this system.

Solution 1: Do you run your switches with Portfast? or spanning tree delays?

You add these to your kernel-console parameter to alter the retry and wait times.
  * `provisioner.portdelay=<Number of seconds>` - seconds to wait before bring up link
  * `provisioner.postportdelay=<Number of seconds>` - seconds to wait after bringing up link before dhcp
  * `provisioner.wgetretrycount=<Number of retries before failure>` - wget of squashfs occurs once a second for 10 times by default.

Solution 2: Is something is really “slower” than sledgehammer expects?

You could try setting `provisioner.wgetretrycount=60`.  `kernel-console` is a parameter that lets you changing the kernel parameters passed to bootenvs.
Sometimes it is used to tweak the kernel console that the kernel is using, but it can be used for other values as well.


.. .. _rs_sledgehammer_no_ip:

Sledgehammer Boots Without IP
-----------------------------

It is possible to have sledgehammer boot and hang without an IP. This can happen in environments where DRP is NOT the DHCP server and the subnet has
restricted number IPs with more HOSTs than IPs.  The ipxe and kernel boot components do DHCP releases as they boot.  This releases IPs back to pools.
In some DHCP servers, the address is immediately available for consumption.  DRP will delay returning the address to the pool for a few seconds to
prevent this problem.  If this happens as the linux user space of sledgehammer takes over operation, the user space DHCP server may not get an address
because the pool is empty.  This will make the machine appear as hung or not responsive.

This sometimes resolves itself as IP addresses become available.  Additional fixes including rebooting the machine, increasing DHCP scope, or using DRP as DHCP server.


.. .. _rs_no_matching_subnet:

No matching Subnet (MacOS DHCP)
-------------------------------

Problem: DHCP will not respond when running DRP from a Mac.  Log provides "No Matching Subnet" warning.

Cause: This is likely caused by not configuring the "MAC DARWIN" route correctly as per :ref:`rs_quickstart`.

Solution: Make sure that the address on the MAC should be outside the range.  Then set the ip, add the route, and then (re)start DRP.  Make sure all the broadcast routes are deleted first using `sudo route delete 255.255.255.255` multiple times.

.. .. _rs_kubernetes_dashboard:

Kubernetes Dashboard
--------------------

For :ref:`rs_krib`, the ``admin.conf`` files is saved into the ``krib/cluster-admin-conf`` profile parameter and can be downloaded after installation is complete.  Using this file ``kubectl --kubeconfig=admin.conf`` allows autheticated access to the cluster.  Please see the KRIB documentation for more details.

For other deployments such as Ansible Kubespray or the Kubeadm deployments of Kubernetes are all maintained by the respective Kubernetes communities.  Digital Rebar simply implements a basic version of those configurations.  Access to the Kubernetes Dashboard is often changing, and being updated by the community.  Please check with the respective communities about how to correctly access the Dashboard.

Some things to note in general:

  * Access is restricted; as well it should
  * You must configure/enable access to the Dashboard
  * Our implmentations usually have a mechanism configured, but this changes over time

Some things that have worked in the past:

  * ``kubectl proxy`` - enabled Proxy access to the Kubernetes Master to get to the Dashboard
  * try stopping the Proxy container, and running ``kubectl proxy --address 0.0.0.0 --accept-hosts '.*'``
     * carefully consider this implication - you are enable access from all hosts !!!
  * any other solutions, please let us know... we'll add them here


.. .. _rs_expand_templates:

Expand Templates from Failed Job
--------------------------------

If you have a task/template that has failed, once it's been run by the Job system, you can collect the rendered template.  The rendered template will be in JSON format, so it may be hard to parse.

  ::

    # set Endpoint and User/Pass appropriately for your environment
    export RS_ENDPOINT="https://127.0.0.1:8092"
    export RS_KEY="rocketskates:r0cketsk8ts"

    # get your Job ID from the failed job, and set accordingly:
    JOBID="abcdefghijklmnopqrstuvwxyz"
    curl -k -u $RS_KEY $RS_ENDPOINT/api/v3/jobs/$JOBID/actions > $JOBID.json

    # optional - if you have the remarshal tools installed:
    json2yaml $JOBID.json > $JOBID.yaml


RBAC - Limit Users to Just Poweron and Poweroff IPMI Controls
-------------------------------------------------------------

The Role Base Access and Controls subsystem allows an operator to construct user account permissions to limit the scope that a user can impact the Digital Rebar Provision system.  Below is an example of how to create a *Claim* that assigns the ``Role`` named ``prod-role`` that limits to only allow IPMI ``poweron`` and ``poweroff`` actions.  These permissions are applied to the _specific_ set of _scope_ *Machines*:

  ::

    drpcli roles update prod-role '"Claims": [{"action": "action:poweron, action:poweroff", "scope": "machines", "specific": "*"}]'

Now simply assign this Role to the given users you wish to limit their permissions on.


.. .. _rs_unblockRunnable_panic:


unblockRunnable Panic
---------------------

In some cases, DRP can panic with a message that contains unblockRunnable.  This error DRP protecting itself from an unhandled deadlock in the database system.
DRP will restart cleanly when restarted.  If run under a service watch system (e.g. systemd), the system will restart and continue.

Please gather the log failure and enter a new issue at `Digital Rebar Github <https://github.com/digitalrebar/provision>`_.


.. .. _rs_manager_system_time:

Manager and System Time
-----------------------

The Multi-site Manager system requires all DRP Endpoints that are being managed to have consistent and accurate
system clock date and time information.  Generally speaking, all Endpoints should have NTP services running,
and all RTC clocks set to UTC.  The Authentication Tokens and Secrets used for the token system will by design
fail if the clocks between two cooperating DRP Endpoint differ more than a few minutes.  This is an intentional
security measure.

If you encounter any of the following errors on "upstream" DRP Managers, this is often the system clocks being
out of sync:

  Machine Objects may show the following:
  ::

    (403) system:
    Invalid token: No valid key specified


  Plugins may show the following:
  ::

    Unable to create event stream: Bad Request

  or even Golang stack traces in some (eg IPMI plugin):
  ::

    Panic recovered: invalid WriteHeader code 0
    Stack trace:
    goroutine 33 [running]:
    runtime/debug.Stack(0x991080, 0xc000289510, 0x1)
	  /home/travis/.gimme/versions/go1.12.7.linux.amd64/src/runtime/debug/stack.go:24 +0x9d
    ...snip...

In addition, Machine objects may show additional failed validation error messages in the Machine details pages.

To correct the problem, install and verify all DRP Endpoints system clocks are in sync with NTP services.


.. .. _rs_jq_examples:

JQ Usage Examples
-----------------

JQ Raw Mode
===========

Raw JSON output is usefull when passing the results of one ``jq`` command in to another for scripted interaction.  Be sure to specify "Raw" mode in this case - to prevent colorization and extraneous quotes being wrapped around Key/Value data output.
  ::

      <some command> | jq -r ...


.. .. _rs_jq_filter_gohai:

Filter Out gohai-inventory
==========================

The ``gohai-inventory`` module is extremely useful for providing Machine classification information for use by other stages or tasks.  However, it is very long and causes a lot of content to be output to the console when listing Machine information.  Using a simple ``jq`` filter, you can delete the ``gohai-inventory`` content from the output display.

Note that since the Param name is ``gohai-inventory``, we have to provide some quoting of the Param name, since the dash (``-``) has special meaning in JSON parsing.
  ::

    drpcli machines list | jq 'del(.[].Params."gohai-inventory")'

Subsequently, if you are listing an individual Machine, then you can also filter it's ``gohai-inventory`` output as well, with:
  ::

    drpcli machines show <UUID> | jq 'del(.Params."gohai-inventory")'


.. .. _rs_jq_list_bootenvs:

List BootEnv Names
==================

Get list of bootenvs available in the installed content, by name:
  ::

    drpcli bootenvs list | jq '.[].Name'


.. .. _rs_jq_reformat_output:

Reformat Output With Specific Keys
==================================

Get list of machines, output "Name:Uuid" pairs from the the JSON output:
  ::

    drpcli machines list | jq -r '.[] | "\(.Name):\(.Uuid)"'

Output is printed as follows:
  ::

    machine1:05abe5dc-637a-4952-a1be-5ec85ba00686
    machine2:0d8b7684-9d0e-4c3e-9f89-eded02357521

You can modify the output separator (colon in this example) to suit your needs.


.. .. _rs_jq_extract_keys:

Extract Specific Key From Output
================================

``jq`` can also pull out only specific Keys from the JSON input.  Here is an example to get ISO File name for a bootenv:
  ::

    drpcli contents show os-discovery | jq '.sections.bootenvs.discovery.OS.IsoFile'


.. .. _rs_jq_display_job_logs:

Display Job Logs for Specific Machine
=====================================

The Job Logs provide a lot of information about the provisioning process of your DRP Endpoint.  However, you often only want to see Job Logs for a specific Machine to evaluate provisioning status.  To get specific Jobs from the job list - based on Machine UUID, do:
  ::

    export UUID=`abcd-efgh-ijkl-mnop-qrps"
    drpcli jobs list | jq ".[] | select(.Machine==\"$UUID\")"


.. .. _rs_jq_list_machines_with_profile:

List Machines with a Given Profile Added to Them
================================================

Starting sometime after v3.9.0 the API will allow you to filter Machines that have a given ``Profile`` applied to them.  If you don't have this version, you can use ``jq`` to list all Machines with a specified ``Profile`` by using the following construct:
  ::

    # set the PROFILE variable to the name you want to match
    export PROFILE=foobar
    drpcli machines list | jq -r ".[] | select(.Profiles[] == \"$PROFILE\") | \"\(.Name)\""

In this case, we simply list the output of the Machines ``Name``.  You can change the final ``\(.Name)`` to any valid JSON key(s) on the Machine Object.

