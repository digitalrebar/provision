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

To have provisioned operating systems (including discovery/sledgehammer) add keys, you should set the ``access_keys`` parameter with a hash of the desired keys.  This can be accomplished by editing the root access profile to add your key(s) and then update the profile via the CLI.

  ::

    vi assets/profiles/root-access.yaml
    drpcli profiles update root-access - < assets/profiles/root-access.yaml
    
NOTE: By default, these changes are targeted at the ``root-access`` profile and you will need to add that profile to selected machines for your keys to be injected.

If you want this parameter applied to all machines by default, then you should change ``root-access`` to ``global`` in the yaml file and command line.  

  ::

    cp assets/profiles/root-access.yaml assets/profiles/global.yaml
    # remember to change root-access to global!
    vi assets/profiles/global.yaml
    drpcli profiles update global - < assets/profiles/global.yaml

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


