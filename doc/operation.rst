.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Operations

.. _rs_operation:

Digital Rebar Provision Operations
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

This section will attempt to describe common operations and actions that can be done.  We will assume that you have
drpcli somewhere in your path and have setup the environment variables to access Digital Rebar Provision.  See, :ref:`rs_cli`.

Some of these operations are in the :ref:`rs_ui`, but not all.  This will focus on CLI usage for now.  See the :ref:`rs_ui`
section for UI usage.

.. note:: **drpcli** normally spits out JSON formatted objects or an array of objects.  To help with some of the command line functions, we use the **jq** tool.  This is available for both Linux and Darwin.


Preference Setting
++++++++++++++++++

Usually, you need to get or set the preferences for your system.

  ::

    # Show the current preference settings
    drpcli prefs list

    # Or get a specific one
    drpcli prefs get unknownBootEnv


Set a preference

  ::

    # Set a preference
    drpcli prefs set unknownBootEnv discovery

    # or chain multiples as pairs
    drpcli prefs set unknownBootEnv discovery defaultBootEnv sledgehammer

The system does validate values to make sure they are sane, so watch for errors.


BootEnv Operations
++++++++++++++++++


Installing a "Canned" BootEnv
-----------------------------

Manipulating :ref:`rs_model_bootenv` and :ref:`rs_model_template` are handled by their own commands.  There are some
additional helpers especially when following the layout of the initial :ref:`rs_install`.

To install a provided :ref:`rs_model_bootenv`, do the following from your install location.

  ::

    cd assets
    drpcli bootenvs install bootenvs/ubuntu-16.04.yml

This is a CLI helper that is not in the API that will read the provided YAML :ref:`rs_model_bootenv` file,
upload the included or referenced :ref:`rs_model_template` files (from the *templates* peer directory), upload
the :ref:`rs_model_bootenv`, and check for an existing ISO in the ISO repository.  If an ISO is not present in
the already uploaded list, it will check a local isos directory for the file.  If that is not present and the
:ref:`rs_model_bootenv` contains a URL for the ISO, the ISO will attempt to be downloaded to the local isos
directory and then uploaded into Digital Rebar Provision.  Once upload, the ISO is "exploded" for access by
machines in the file server file system space.


Cloning a BootEnv
-----------------

Sometimes you have a :ref:`rs_model_bootenv` but want to make changes.  Now, these can be handled by :ref:`rs_model_template`
inclusion, but for now let's just focus on basic "cut and paste" style editing.

  ::

    drpcli bootenvs show ubuntu-16.04-install --format yaml > new-file.yaml
    # Edit the file
    #  change the Name field to something new. *MUST DO THIS*
    #  change the OS->Name field to something new if you don't want to sure the same iso directory.
    #  Edit other parameters as needed
    drpcli bootenvs create - < new-file.yaml

This is a shallow clone.  It will reuse the templates unless you explictly modify them.  You could use the *install*
command, but any new templates would need to be added to a *templates* directoy in the current directory.

Creating a BootEnv
------------------

You can additionally create an empty :ref:`rs_model_bootenv` by doing the following:

  ::

    drpcli bootenvs create emtpy_bootenv

This :ref:`rs_model_bootenv` will not be *Available*, but will allow for additional editing.

Editing a BootEnv
-----------------

Sometimes you just want to edit a :ref:`rs_model_bootenv`.  To do this, get the latest copy with the *show*
command.  Edit the file as needed.  Then using the *update* command, put the value back.  The *--format=yaml*
is optional, but I find YAML easier to edit.

  ::

    drpcli bootenvs show discovery --format=yaml > discovery.yaml
    # Edit the discovery.yaml as you want
    drpcli bootenvs update discovery - < discovery.yaml

Template Operations
+++++++++++++++++++

Cloning a Template
------------------

Sometimes you want to create a new template from an existing one.  To do this, do the following:

  ::

    drpcli templates show net_seed.tmpl | jq -r .Contents > new.tmpl
    # Edit the new.tmpl to be what you want
    drpcli templates upload new.tmpl as new_template

In this case, we are using **jq** to help us out.  **jq** is a JSON processing command line filter.  You send JSON in and you
get data back.  In this case, we are wanting the Contents of the template.  We save that to file, edit it, and upload it as a
new template, *new_template*.

You could also use the **create** subcommand of template, but often times **upload** is easier.

.. note:: Remember to add the new template to a :ref:`rs_model_bootenv` or another :ref:`rs_model_template` as an embedded template.


Updating a Template
-------------------

Sometimes you want to edit an existing template.  To do this, do the following:

  ::

    drpcli templates show net_seed.tmpl | jq -r .Contents > edit.tmpl
    # Edit the edit.tmpl to be what you want
    drpcli templates upload edit.tmpl as net_seed.tmpl

We use **jq** to get a copy of the current template, edit it, and use the upload command to replace the template.
If you aleady had a template, you could replace it with the upload command.

Profile Operations
++++++++++++++++++

Creating a Profile
------------------

Sometimes you want to create a :ref:`rs_model_profile`.  You can create an empty profile by doing the following:

  ::

    drpcli profiles create '{ "Name": "myprofile" }'

    or

    drpcli profiles create myprofile

If you just send a string, the system will attempt to use that as the Name of the profile.

Additionally, JSON can be provided to fill in some default values.

  ::

    drpcli profiles create '{ "Name": "myprofile", "Params": { "string_param1": "string", "map_parm1": { "key1": "value", "key2": "value2" } } }'


Deleting a Profile
------------------

Sometimes you want to delete a :ref:`rs_model_profile`.  You can use the destroy command in the profile CLI,
but the :ref:`rs_model_profile` must not be in use.  Use the following:

  ::

    drpcli profiles destroy myprofile


Altering an Existing Profile (including global)
-----------------------------------------------

Somtimes you want to update an existing :ref:`rs_model_profile`, including **global**.  You can *set*
parameter values by doing the following:

  ::

    drpcli profiles set myprofile param crazycat to true
    # These last two will show the value or the whole profile.
    drpcli profiles get myprofile param crazycat
    drpcli profiles show myprofile

.. note:: Setting a parameter's value to **null** will clear it from the structure.

Alternatively, you can also use the update command and send raw JSON similar to create.

  ::

    drpcli profiles update myprofile '{ "Params": { "string_param1": "string", "map_parm1": { "key1": "value", "key2": "value2" }, "crazycat": null } }'

Update is an additive operation by default.  So, to remove items, **null** must be passed as
the value of the key you wish to remove.

Machine Operations
++++++++++++++++++

Creating a Machine
------------------

Sometimes you want to create a :ref:`rs_model_machine`.  You know the IP address the machine is going to boot as and you just want to
create the machine and assign a :ref:`rs_model_bootenv`.  To do this, do the following:

  ::

    drpcli machine create '{ "Name": "greg.rackn.com", "Address": "1.1.1.1" }'

This would create the :ref:`rs_model_machine` named *greg.rackn.com* with an expected IP Address of *1.1.1.1*.  *dr-provision*
will create the machine, create a UUID for the node, and assign the :ref:`rs_model_bootenv` based upon the *defaultBootEnv*
:ref:`rs_model_prefs`.

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

Sometimes you want to add or remove a :ref:`rs_model_profile` to a :ref:`rs_model_machine`.  To add a profile, do the following:

  ::

    drpcli machines addprofile "dff3a693-76a7-49ce-baaa-773cbb6d5092" myprofile


To remove a profile, do the following:

  ::

    drpcli machines removeprofile "dff3a693-76a7-49ce-baaa-773cbb6d5092" myprofile

The :ref:`rs_model_machine` update command can also be used to modify the list of :ref:`rs_model_profile`.


Changing BootEnv on a Machine
-----------------------------

Sometimes you want to change the :ref:`rs_model_bootenv` associated with a :ref:`rs_model_machine`.  To do this, do the following:

  ::

    drpcli machines bootenv drpcli "dff3a693-76a7-49ce-baaa-773cbb6d5092" mybootenv

.. note:: The :ref:`rs_model_bootenv` *MUST* exists or the command will fail.


DHCP Operations
+++++++++++++++

Creating a Reservation
----------------------

Sometimes you want to create a :ref:`rs_model_reservation`.  This would be to make sure that a specific MAC Address received
a specific IP Adress.  Here is an example command.

  ::

     drpcli reservations create '{ "Addr": "1.1.1.1", "Token": "08:00:27:33:77:de", "Strategy": "MAC" }'

You can additionally add DHCP options or the Next Boot server.

  ::

     drpcli reservations create '{ "Addr": "1.1.1.5", "Token": "08:01:27:33:77:de", "Strategy": "MAC", "NextServer": "1.1.1.2", "Options": [ { "Code": 44, "Value": "1.1.1.1" } ] }'

Remember to add an option 1 (netmask) if you are not using a subnet to fill in default options.

User Operations
+++++++++++++++

Creating a User
---------------

Sometimes you want to create a :ref:`rs_model_user`.  By default, the user will be created without
a valid password.  The user will only be able to access the system through granted tokens.

To create a user, do the following:

  ::

    drpcli users create fred

.. note:: This :ref:`rs_model_user` will *NOT* be able to access the system without additional admin action.


Granting a User Token
---------------------

Sometimes as an administrator, you would like to grant a limited use and scope access token to a user.  To
grant a token, do the following:

  ::

    drpcli users token fred

This will create a token that is valid for 1 hour and can do anything.  Additionally, the CLI can take
additional parameters that alter the token's scope (model), actions, and key.

  ::

    drpcli users token fred ttl 600 scope users action password specfic fred

This will create a token that is valid for 10 minutes and can only execute the password API call on the
:ref:`rs_model_user` object named *fred*.

To use the token in with the CLI, use the -T option.

  ::

    drpcli -T <token> bootenvs list


Deleting a User
---------------

Sometimes you want to remove a reset from the system. To remove a user, do the following:

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


