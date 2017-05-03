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
------------------

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

Editing a BootEnv
-----------------

Sometimes you just want to edit a :ref:`rs_model_bootenv`.  To do this, get the latest copy with the *show*
command.  Edit the file as needed.  Then using the *update* command, put the value back.  The *--format=yaml*
is optional, but I find YAML easier to edit.

  ::

    drpcli bootenvs show discovery --format=yaml > discovery.yaml
    # Edit the discovery.yaml as you want
    drpcli bootenvs update discovery - < discovery.yaml


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


