.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; FAQ
  pair: Digital Rebar Provision; Troubleshooting

.. _rs_faq:

FAQ / Troubleshooting
~~~~~~~~~~~~~~~~~~~~~

The following section is designed to answer frequently asked questions and help troubleshoot Digital Rebar Provision installs.

Want ligher reading?  Checkout our :ref:`rs_fun`.

.. _rs_bind_error:

Bind Error
----------

Digital Rebar Provision will fail if it cannot attach to one of the required ports.

* Typical Error Message: "listen udp4 :67: bind: address already in use"
* Additional Information: The conflicted port will be included in the error between colons (e.g.: `:67:`)
* Workaround: If the conflicting service is not required, simply disable that service
* Resolution: Stop the offending service on the system.  Typical corrective actions are:

  * 67 - dhcp.  Correct with `sudo pkill dnsmasq`

See the port mapping list on start-up for a complete list.

.. _rs_gen_cert:

Generate Certificate
--------------------

Sometimes the cert/key pair in the github tree is corrupt or not sufficient for the environment.  The following command can be used to rebuild a local cert/key pair.

  ::

    sudo openssl req -new -x509 -keyout server.key -out server.crt -days 365 -nodes

It may be necessary to install the openssl tools.

.. _rs_add_ssh:

Add SSH Keys to Authorized Keys
-------------------------------

VIDEO TUTORIAL: https://www.youtube.com/watch?v=StQql8Xn08c 

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

    drpcli profiles update my-profile -< keys.json
    
.. _rs_access_ssh_root_mode:

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


.. _rs_autocomplete:

Turn on autocomplete for the CLI
--------------------------------

The DRP CLI has built in support to generate autocomplete (tab completion) capabilities for the BASH shell.  To enable, you must generate the autocomplete script file, and add it to your system.  This can also be added to your global shell ``rc`` files to enable autocompletion every time you log in.  NOTE that most Linux distros do this slightly differently.  Select the method that works for your distro.  

You must specify a filename as an argument to the DRP CLI autocomplete command.  The filename will be created with the autocomplete script.  If you are writing to system areas, you need ``root`` access (eg via `sudo`).  

For Debian/Ubuntu and RHEL/CentOS distros:
  ::
  
    sudo drpcli autocomplete /etc/bash_completion.d/drpcli

For Mac OSX (Darwin):
  ::

    sudo drpcli autocomplete /usr/local/etc/bash_completion.d/drpcli

Once the autocomplete file has been created, either log out and log back in, or ``source`` the created file to enable autocomplete in the current shell session (example for Linux distros, adjust accordingly):
  ::

    source /etc/bash_completion.d/drpcli
    
.. _rs_more_debug:

Turn Up the Debug
-----------------

To get additional debug from dr-provision, set debug preferences to increase the logging.  See :ref:`rs_model_prefs`.

.. _rs_vboxnet:

Missing VBoxNet Network
-----------------------

Virtual Box does not add host only networks until a VM is attempting to use them.  If you are using the interfaces API (or UX wizard) to find available networks and ``vboxnet0`` does not appear then start your VM and recreate the address.

Virtual Box may also fail to allocate an IP to the host network due to incomplete configuration.  In this case, ``ip addr`` will show the network but no IPv4 address has been allocated; consequently, Digital Rebar will not report this as a working interface. 

.. _rs_debug_sledgehammer:

Debug Sledgehammer
------------------

If the sledgehammer discovery image should fail to launch Runner jobs successfully, or other issues arise with the start up sequences, you can debug start up via the systemd logging.  Log in to the console of the Machine in question (or if SSH is running and you have ``access-keys`` setup, you can SSH in), and run the following command to output logging:
  ::

      journalctl -u sledgehammer


.. _rs_convert_to_production_mode:

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
#. make sure you're start up scripts are in place for your production mode (eg: ``/etc/systemd/system/dr-provision.service``)
#. start the new production version with  ``systemctl start dr-provision.service``
#. verify everything is running fine
#. delete the ``drp-data`` directory (suggest retaining the backup copy for later just in case)

.. note::  WARNING:  If you install a new version of the Digital Rebar Provision service, you must verify that there are no Contents differences between the two versions.  Should the ``dr-provision`` service fail to start up; it's entirely likely that there may be some content changes that need to be addressed in the JSON/YAML files prior to the new version being started.  See the :ref:`rs_upgrade` notes for any version-to-version specific documentation.


.. _rs_kickseed:

Custom Kickstart and Preseeds
-----------------------------

Starting with ``drp-community-content`` version 1.5.0 and newer, you can now define a custom Kickstart or Preseed (aka *kickseed*) to override the defaults in the selected BootEnv.  You simply need to only define a single Param with the name of the Kickstart or Preseed you wish to override the default value. 
  ::
    
    export UUID="f6ca7bb6-d74f-4bc1-8544-f3df500fb15e"
    drpcli machines set $UUID param kickseed to "my_kickstart.cfg"

Of course, you can apply a Param to a Profile, and apply that Profile to a group of Machines if desired. 

.. note:: The Digital Rebar default kickstart and preseeds have Digital Rebar specific interactions that may be necessary to replicate.  Please review the default kickstart and preseeds for patterns ane examples you may need to re-use. 


.. _rs_jq_examples:

JQ Usage Examples
-----------------

JQ Raw Mode
===========

Raw JSON output is usefull when passing the results of one ``jq`` command in to another for scripted interaction.  Be sure to specify "Raw" mode in this case - to prevent colorization and extraneous quotes being wrapped around Key/Value data output.
  ::

      <some command> | jq -r ... 

.. _rs_filter_gohai:

Filter Out gohai-inventory
==========================

The ``gohai-inventory`` module is extremely useful for providing Machine classification information for use by other stages or tasks.  However, it is very long and causes a lot of content to be output to the console when listing Machine information.  Using a simple ``jq`` filter, you can delete the ``gohai-inventory`` content from the output display. 

Note that since the Param name is ``gohai-inventory``, we have to provide some quoting of the Param name, since the dash (``-``) has special meaning in JSON parsing.  
  ::

    drpcli machines list | jq 'del(.[].Params."gohai-inventory")'

Subsequently, if you are listing an individual Machine, then you can also filter it's ``gohai-inventory`` output as well, with:
  ::

    drpcli machines show <UUID> | jq 'del(.Params."gohai-inventory")'

List BootEnv Names
==================

Get list of bootenvs available in the installed content, by name:
  ::

    drpcli bootenvs list | jq '.[].Name'


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


Extract Specific Key From Output
================================

``jq`` can also pull out only specific Keys from the JSON input.  Here is an example to get ISO File name for a bootenv:
  ::

    drpcli contents show os-discovery | jq '.sections.bootenvs.discovery.OS.IsoFile'


Display Job Logs for Specific Machine
=====================================

The Job Logs provide a lot of information about the provisioning process of your DRP Endpoint.  However, you often only want to see Job Logs for a specific Machine to evaluate provisioning status.  To get specific Jobs from the job list - based on Machine UUID, do:
  ::

    export UUID=`abcd-efgh-ijkl-mnop-qrps"
    drpcli jobs list | jq ".[] | select(.Machine==\"$UUID\")"

Change a Machines Name
======================
If you wish to update/change a Machine Name, you can do: 
  ::

    export UUID="abcd-efgh-ijkl-mnop-qrst"
    drpcli machines update $UUID '{ "Name": "foobar" }'

Or, you can use the ``set`` method:
  ::

    export UUID="abcd-efgh-ijkl-mnop-qrst"
    drpcli machines set $UUID param Name to foobar


