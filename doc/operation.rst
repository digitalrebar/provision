.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Operations

.. _rs_operation:

Digital Rebar Provision Operations
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

This section will attempt to describe common operations and actions that can be done.  We will assume that ``drpcli`` is already somewhere in the path and have setup the environment variables to access Digital Rebar Provision.  See, :ref:`rs_cli`.

Some of these operations are in the :ref:`rs_ui`, but not all.  This will focus on CLI usage for now.  See the :ref:`rs_ui` section for UI usage.

.. note:: **drpcli** normally spits out JSON formatted objects or an array of objects.  To help with some of the command line functions, we use the **jq** tool.  This is available for both Linux and Darwin.  You can specify output format to be YAML, with the ``--format=yaml`` flag. 

Using the ``drpcli`` tool
+++++++++++++++++++++++++

The ``drpcli`` tool is designed to provide a single, easy to use compiled stand alone binary.  It very closely mimics the API call parameters and usage, and allows for easy "transition" between the use of the CLI binary, and the use of the API calls. 

We assume that you have the binary in your ``PATH`` somewhere (``/usr/local/bin/drpcli`` in "production" mode by default).  If you are using an "isolated" install mode, you will need to link the correct binary in to an apprpriate directory in your PATH.  For example:

  ::

    ln -s $HOME/bin/linux/amd64/dr-provision $HOME/bin/
    ln -s $HOME/bin/linux/amd64/drpcli $HOME/bin/

``drpcli`` is a self documenting binary file.  At any point you can can simply hit ``<Enter>`` to see the current contextual help output.  Some examples:

  ``drpcli <Enter>``
    Prints usage for the top-level resources that may be manipulated.

  ``drpcli templates <Enter>``
    Prints the resources that may be executed for the ``templates`` level.

  ``drpcli templates create <Enter>``
    Will display the commands associated with just ``templates/create/``

Note that the CLI syntax closely follows the API calls.  For example:

  ::

    drpcli templates create -< some_file.json

Would be equivalent to the HTTP REST API resource as follows:

  ::
    
    https://127.0.0.1:8092/api/v3/templates/ereate

.. note:: See the :ref:`rs_autocomplete` for BASH shell auto completion.


Preference Setting
++++++++++++++++++

By default Digital Rebar Provision (DRP) attempts to "do no harm".  This means that by default, any system that receives a DHCP lease from the DRP Endpoint, will by default, be set to ``local`` mode, which means boot from local disks.  You must explicitly change this behavior to enable provisioning activities of Machines. 

It is necessary to get or set the preferences for the system, to enable OS Installs (BootEnvs). 

  ::

    # Show the current preference settings
    drpcli prefs list

    # Or get a specific one
    drpcli prefs get unknownBootEnv


Set preference(s):

  ::

    # Set a preference
    drpcli prefs set unknownBootEnv discovery

    # or chain multiples as pairs
    drpcli prefs set unknownBootEnv discovery defaultBootEnv sledgehammer defaultStage discover

The system does validate values to make sure they are sane, so watch for errors.


BootEnv Operations
++++++++++++++++++

A :ref:`rs_model_bootenv` is the primary component for an Operating System installation definition.  A :ref:`rs_model_bootenv` is comprised of two primary pieces:

  #. A :ref:`rs_model_bootenv` JSON/YAML specification
  #. (usually) an ISO Image that installs that :ref:`rs_model_bootenv` 

The JSON/YAML specification will contain a set of definitions for the ISO image.  The default distributed :ref:`rs_model_bootenv` specs use the public mirror repos for the ISO images.  You can create a customer :ref:`rs_model_bootenv` with a pointer to your own hosted ISO images.  An example looks something like:

  ::

    root@demo:~$ drpcli bootenvs show ubuntu-16.04-install
    {
    "Available": true,
    "Name": "ubuntu-16.04-install",
    "OS": {
      "Family": "ubuntu",
      "IsoFile": "ubuntu-16.04.3-server-amd64.iso",
      "IsoSha256": "a06cd926f5855d4f21fb4bc9978a35312f815fbda0d0ef7fdc846861f4fc4600",
      "IsoUrl": "http://mirrors.kernel.org/ubuntu-releases/16.04/ubuntu-16.04.3-server-amd64.iso",
      "Name": "ubuntu-16.04",
    <...snip...>

This stanza shows the Ubuntu 16.04 :ref:`rs_model_bootenv` along with the associated Mirror HTTP location the ISO will be installed from.


Installing a "Canned" BootEnv
-----------------------------

Manipulating :ref:`rs_model_bootenv` and :ref:`rs_model_template` are handled by their own commands.  There are some additional helpers especially when following the layout of the initial :ref:`rs_install`.

To install a provided :ref:`rs_model_bootenv`, do the following from the install location.

  ::

    drpcli bootenvs uploadiso ubuntu-16.04-install

This is a CLI helper that is not in the API that will read the provided YAML :ref:`rs_model_bootenv` file,
upload the included or referenced :ref:`rs_model_template` files (from the *templates* peer directory), upload
the :ref:`rs_model_bootenv`, and check for an existing ISO in the ISO repository.  If an ISO is not present in
the already uploaded list, it will check a local isos directory for the file.  If that is not present and the
:ref:`rs_model_bootenv` contains a URL for the ISO, the ISO will attempt to be downloaded to the local isos
directory and then uploaded into Digital Rebar Provision.  Once upload, the ISO is "exploded" for access by
machines in the file server file system space.

Listing Installed BootEnvs
--------------------------

A list of all existing :ref:`rs_model_bootenv` installed on the DRP Endpoint can be obtained with the *list* command.  However, you usually do not wish to see all of the JSON values, and a simple ``jq`` filter can help output just the keys you are interested in, as follows:

  ::

    drpcli bootenvs list | jq -r '.[].Name'

    Outputs:
    centos-7-install
    centos-7.4.1708-install
    debian-8-install
    debian-9-install
    discovery
    ignore
    local
    sledgehammer
    ubuntu-16.04-install

Cloning a BootEnv
-----------------

Sometimes there is a :ref:`rs_model_bootenv` but it is necessary to make changes.  These can be handled by :ref:`rs_model_template`
inclusion, but for now let's just focus on basic "cut and paste" style editing.

  ::

    drpcli bootenvs show ubuntu-16.04-install --format yaml > new-file.yaml
    # Edit the file
    #  change the Name field to something new. *MUST DO THIS*
    #  change the OS->Name field to something new to avoid sharing an iso directory.
    #  Edit other parameters as needed
    drpcli bootenvs create - < new-file.yaml

This is a shallow clone.  It will reuse the templates unless they are explicitly modified.  It is possible to use the *install*
command, but any new templates would need to be added to a *templates* directory in the current directory.

Creating a BootEnv
------------------

It might be necessary to create an empty :ref:`rs_model_bootenv` by doing the following:

  ::

    drpcli bootenvs create emtpy_bootenv

This :ref:`rs_model_bootenv` will not be *Available*, but will allow for additional editing.

Editing a BootEnv
-----------------

It might be necessary to edit a :ref:`rs_model_bootenv`.  To do this, get the latest copy with the *show*
command.  Edit the file as needed.  Then using the *update* command, put the value back.  The *--format=yaml*
is optional, but I find YAML easier to edit.

  ::

    drpcli bootenvs show discovery --format=yaml > discovery.yaml
    # Edit the discovery.yaml as needed
    drpcli bootenvs update discovery - < discovery.yaml

Subnet Operations
+++++++++++++++++

Subnet definitions provide the necessary information for DHCP IP Address lease assignments, and allows Machines to be enrolled/discovered by a DRP Endpoint.  For any Layer 2 subnet/network that you wish to install Machines from, you must also specify a Subnet definition for.  In some environments, a Subnet definition may not be needed to allow Machines to be discovered. 

Cloning a Subnet
----------------

It might be necessary to create a new subnet from an existing one.  To do this, do the following:

  ::

    drpcli subnets show eth0 | jq -r > new_subnet.json
    # edit the new_subnet.json file with the new information
    drpcli subnets create -< new_subnet.json 

Creating a new Subnet
---------------------

A new subnet can be created from a JSON specification.  It is necessary to use all of the following JSON keys to successfully create a new Subnet

  ::

    echo '
    {
      "Name": "local_subnet",
      "Subnet": "10.10.16.10/24",
      "ActiveStart": "10.10.16.100",
      "ActiveEnd": "10.10.16.254",
      "NextServer": "10.10.16.10",
      "ActiveLeaseTime": 60,
      "Available": true,
      "Enabled": true,
      "Proxy": false,
      "ReadOnly": false,
      "ReservedLeaseTime": 7200,
      "Strategy": "MAC",
      "Validated": true,
      "OnlyReservations": false,
      "Pickers": [ "hint", "nextFree", "mostExpired" ],
      "Options": [
        { "Code": 1, "Value": "255.255.255.0", "Description": "Netmask" },
        { "Code": 3, "Value": "10.10.16.1", "Description": "Default Gateway" },
        { "Code": 6, "Value": "8.8.8.8", "Description": "DNS Servers" },
        { "Code": 15, "Value": "example.com", "Description": "Domain Name" },
        { "Code": 28, "Value": "10.10.16.255", "Description": "Broadcast Address" },
        { "Code": 67, "Value": "lpxelinux.0", "Description": "Boot file name" }
      ]
    } ' > /tmp/local_subnet.json 

    drpcli subnets create -< /tmp/local_subnet.json

Note that the "Description" is purely cosmetic and not used - however, it can be safely specified as it'll be ignored (it's added here for the readers reference).  You must provide the minimum DHCP Options as specified above.  You can find a complete set of DHCP Options at: 

  https://www.iana.org/assignments/bootp-dhcp-parameters/bootp-dhcp-parameters.xhtml

For complete documentation and information you can find the DHCP Options officially documented in `RFC2132 <https://tools.ietf.org/html/rfc2132>`_ 

Updating a Subnet
-----------------

From time to time, you may need to modify an existing Subnet definition.  Depending on your changes, you have a couple of options. 

Set the NTP Server pool via DHCP Option 42 for subnet "local_subnet":
  ::

    drpcli subnets set local_subnet option 42 to "0.pool.ntp.org"

Set the DHCP IP assignment from the following pick list for subnet "local_subnet".  See :ref:`rs_model_pickers` for a detailed description of the available Picker types:
  ::

    drpcli subnets pickers local_subnet hint,nextFree,mostExpired

Set the nextserver for PXE operation for subnet "local_subnet":
  ::

    drpcli subnets  nextserver  local_subnet 10.16.167.10

Set the subnet DHCP range of IP addresses for subnet "local_subnet":
  ::

    drpcli subnets range local_subnet 192.168.45.100 192.168.45.255

Set Active lease to 60 mins, and reserved lease to 7200 mins for subnet "local_subnet":
  ::

    drpcli subnets leasetimes local_subnet 60 7200

Update a subnet to set it to disabled (do not discover, and do not provision on this subnet, for subnet "local_subnet":
  ::

    drpcli subnets update local_subnet '{ "Enabled": false }'

Update a subnet with the contents of the specified JSON file, for subnet "local_subnet":
  ::

    drpcli subnets update local_subnet -< update-local_subnet.json 

Deleting a Subnet
-----------------

To remove a Subnet and subsequently cease PXE provisioning operations for that Subnet, perform the following:

  ::

    drpcli subnets destroy local_subnet 

List and Show Subnets
---------------------

Viewing configuration for all subnets can be done with the ``list`` command as follows:
  ::

    drpcli subnets list

To ``show`` an individual subnet, you will need the subnet name.  To show just the subnet names, you can use ``jq`` to filter the output, as follows:
  ::

    drpcli subnets list | jq '.[].Name'

Once you have determined which subnet you'd like to show specific information for, you can do so with the following command:
  ::

    # show the YAML formatted output for 'local_subnet' subnet
    drpcli subnets show local_subnet --format=yaml

Template Operations
+++++++++++++++++++

Templates are reusable blocks of code, that are dynamically expanded when used.  This allows for very sophisticated and complex operations.  It also allows for carefully crafted Templates to be re-usable across a broad set of use cases.

Cloning a Template
------------------

It might be necessary to create a new template from an existing one.  To do this, do the following:

  ::

    drpcli templates show net_seed.tmpl | jq -r .Contents > new.tmpl
    # Edit the new.tmpl to be what is required
    drpcli templates upload new.tmpl as new_template

In this case, we are using ``jq`` to help us out.  ``jq`` is a JSON processing command line filter.  JSON can be used to retrieve the required data.  In this case, we are wanting the Contents of the template.  We save that to file, edit it, and upload it as a new template, *new_template*.

It is possible to use the **create** subcommand of template, but often times **upload** is easier.

.. note:: Remember to add the new template to a :ref:`rs_model_bootenv` or another :ref:`rs_model_template` as an embedded template.


Updating a Template
-------------------

It might be necessary to edit an existing template.  To do this, do the following:

  ::

    drpcli templates show net_seed.tmpl | jq -r .Contents > edit.tmpl
    # Edit the edit.tmpl to be what is desired
    drpcli templates upload edit.tmpl as net_seed.tmpl

We use ``jq`` to get a copy of the current template, edit it, and use the upload command to replace the template.
If there already is a template present, then it can be replaced with the upload command.

Param Operations
++++++++++++++++

:ref:`rs_model_param` are simply key/value pairs.  However, DRP provides a strong typing model to enforce a specific type to a given Param.  This insures that Param values are valid elements as designed by the operator.

Creating a Param
----------------

It might be necessary to create a new :ref:`rs_model_param`, an empty Param may be created by doing the following:

  ::

    drpcli params create '{ "Name": "fluffy" }'

    or

    drpcli params create fluffy


The system will attempt to use any sent string as the Name of the Param.  To be complete, it is required to also speciy the Type that param must be:

  ::

    drpcli params create '{ "Description": "DNS domainname", "Name": "domainname", "Schema": { "type": "string" } }'

In this example, the type ``string`` was defined for the param.

Deleting a Param
----------------

It might be necessary to delete a :ref:`rs_model_param`. 

  ::

    drpcli params destroy fluffy


.. note:: The destroy operation will fail if the param is in use.

Editing a Param
---------------

It might be necessary to update a Param.  An example to add a ``type`` of ``string`` to our ``fluffy`` param above would be:

  ::

    drpcli params update fluffy '{ "Schema": { "type": "string" } }'


Profile Operations
++++++++++++++++++

:ref:`rs_model_profile` are simply collections of :ref:`rs_model_param` - they conveniently group multiple :ref:`rs_model_param` for easy consumption by other elements of the provisioning service.

Creating a Profile
------------------

It might be necessary to create a :ref:`rs_model_profile`. An empty profile can be created by doing the following:

  ::

    drpcli profiles create '{ "Name": "myprofile" }'

    or

    drpcli profiles create myprofile

The system will attempt to use any sent string as the Name of the profile.

Additionally, JSON can be provided to fill in some default values.

  ::

    drpcli profiles create '{ "Name": "myprofile", "Params": { "string_param1": "string", "map_parm1": { "key1": "value", "key2": "value2" } } }'

Alternatively, you can create profiles from an existing file containing JSON, as follows:

  ::

    echo '{ "Name": "myprofile", "Params": { "string_param1": "string", "map_parm1": { "key1": "value", "key2": "value2" } } }' > my_profile.json
    drpcli profiles create -< my_profile.json


Deleting a Profile
------------------

It might be necessary to delete a :ref:`rs_model_profile`.  It is possible to use the destroy command in the profile CLI,
but the :ref:`rs_model_profile` must not be in use.  Use the following:

  ::

    drpcli profiles destroy myprofile


Altering an Existing Profile (including the ``global`` profile)
---------------------------------------------------------------

It might be necessary to update an existing :ref:`rs_model_profile`, including **global**.  parameter values can be *set* by doing the following:

  ::

    drpcli profiles set myprofile param crazycat to true
    # These last two will show the value or the whole profile.
    drpcli profiles get myprofile param crazycat
    drpcli profiles show myprofile

.. note:: Setting a parameter's value to **null** will clear it from the structure.

Alternatively, the update command can be used to send raw JSON similar to create.

  ::

    drpcli profiles update myprofile '{ "Params": { "string_param1": "string", "map_parm1": { "key1": "value", "key2": "value2" }, "crazycat": null } }'

Update is an additive operation by default.  So, to remove items, **null** must be passed as
the value of the key to be removed.

Machine Operations
++++++++++++++++++

A :ref:`rs_model_machine` is typically a physical bare metal server, as DRP is intended to operate on bare metal infrastructure.  However, it can represent a Virtual Machine instance and provision it equally.  DRP does not provide *control plane* activities for virtualized environments (eg *VM Create*, etc. operations).

Creating a Machine
------------------

It might be necessary to create a :ref:`rs_model_machine`.  Given the IP that the machine will boot as all that is required is to create the machine and assign a :ref:`rs_model_bootenv`.  To do this, do the following:

  ::

    drpcli machine create '{ "Name": "greg.rackn.com", "Address": "1.1.1.1" }'

This would create the :ref:`rs_model_machine` named *greg.rackn.com* with an expected IP Address of *1.1.1.1*.  *dr-provision* will create the machine, create a UUID for the node, and assign the :ref:`rs_model_bootenv` based upon the *defaultBootEnv* :ref:`rs_model_prefs`.

  ::

    drpcli machine create '{ "Name": "greg.rackn.com", "Address": "1.1.1.1", "BootEnv": "ubuntu-16.04-install" }'

This would do the same thing as above, but would create the :ref:`rs_model_machine` with the *ubuntu-16.04-install*
:ref:`rs_model_bootenv`.

.. note:: The :ref:`rs_model_bootenv` MUST exist or the create will fail.

To create an empty :ref:`rs_model_machine`, do the following:

  ::

    drpcli machine create jill.rackn.com

This will create an empty :ref:`rs_model_machine` named *jill.rackn.com*.

.. note:: The *defaultBootEnv* :ref:`rs_model_bootenv` MUST exist or the create will fail.


Adding or Removing a Profile to a Machine
-----------------------------------------

It might be necessary to add or remove a :ref:`rs_model_profile` to or from a :ref:`rs_model_machine`.  To add a profile, do the following:

  ::

    drpcli machines addprofile "dff3a693-76a7-49ce-baaa-773cbb6d5092" myprofile


To remove a profile, do the following:

  ::

    drpcli machines removeprofile "dff3a693-76a7-49ce-baaa-773cbb6d5092" myprofile

The :ref:`rs_model_machine` update command can also be used to modify the list of :ref:`rs_model_profile`.


Changing BootEnv on a Machine
-----------------------------

It might be necessary to change the :ref:`rs_model_bootenv` associated with a :ref:`rs_model_machine`.  To do this, do the following:

  ::

    drpcli machines bootenv drpcli "dff3a693-76a7-49ce-baaa-773cbb6d5092" mybootenv

.. note:: The :ref:`rs_model_bootenv` *MUST* exists or the command will fail.


DHCP Operations
+++++++++++++++

.. _rs_create_reservation:

Creating a Reservation
----------------------

It might be necessary to create a :ref:`rs_model_reservation`.  This would be to make sure that a specific MAC Address received
a specific IP Address.  Here is an example command.

  ::

     drpcli reservations create '{ "Addr": "1.1.1.1", "Token": "08:00:27:33:77:de", "Strategy": "MAC" }'

Additionally, it is possible to add DHCP options or the Next Boot server.

  ::

     drpcli reservations create '{ "Addr": "1.1.1.5", "Token": "08:01:27:33:77:de", "Strategy": "MAC", "NextServer": "1.1.1.2", "Options": [ { "Code": 44, "Value": "1.1.1.1" } ] }'

Remember to add an option 1 (netmask) if a subnet is not being used to fill the default options.

.. _rs_advanced_workflow:

Advanced Workflow
+++++++++++++++++

Placeholder for Advanced Workflow overview.


.. _rs_stages:

Stages
------

Placeholder for Stages information.


.. _rs_stagemaps:

Stage Maps
----------

Placeholder for Stage Map information.


.. _rs_tasks:

Tasks
-----

Placeholder for Tasks information.


.. _rs_jobs:

Jobs
----

Placeholder for Jobs information.


User Operations
+++++++++++++++

Creating a User
---------------

It might be necessary to create a :ref:`rs_model_user`.  By default, the user will be created without
a valid password.  The user will only be able to access the system through granted tokens.

To create a user, do the following:

  ::

    drpcli users create fred

.. note:: This :ref:`rs_model_user` will *NOT* be able to access the system without additional admin action.


.. _rs_grant_token:

Granting a User Token
---------------------

Sometimes as an administrator, it may be necessary to grant a limited use and scope access token to a user.  To
grant a token, do the following:

  ::

    drpcli users token fred

This will create a token that is valid for 1 hour and can do anything.  Additionally, the CLI can take
additional parameters that alter the token's scope (model), actions, and key.

  ::

    drpcli users token fred ttl 600 scope users action password specific fred

This will create a token that is valid for 10 minutes and can only execute the password API call on the
:ref:`rs_model_user` object named *fred*.

To use the token in with the CLI, use the -T option.

  ::

    drpcli -T <token> bootenvs list


Deleting a User
---------------

It might be necessary to remove a reset from the system. To remove a user, do the following:

  ::

    drpcli users destroy fred


Revoking a User's Password
--------------------------

To clear the password from a :ref:`rs_model_user`, do the following:

  ::

    drpcli users update fred '{ "PasswordHash": "12" }'

This basically creates an invalid hash which matches no passwords.  Issued tokens will still continue to
function until their times expire.

Secure User Creation Pattern
----------------------------

A secure pattern would be the following:

* Admin creates a new account

  ::

    drpcli users create fred

* Admin creates a token for that account that only can set the password and sends that token to new user.

  ::

    drpcli users token fred scope users action password ttl 3600

* New user uses token to set their password

  ::

    drpcli -T <token> users password fred mypassword


