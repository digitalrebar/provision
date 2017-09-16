.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Upgrade

.. _rs_upgrade:

Upgrade
~~~~~~~

While not glamorous, existing code can be overwritten by a new install and restart.  That is about it.  Here are few more details.

Steps
=====

For isolated :ref:`rs_install`, update this way:

#. Stop dr-provision:
   ::

      killall dr-provision

#. Return to the install directory
#. Run the install again
   ::

     rm sha256sums
     # Remeber to use --drp-version to set a version other than stable
     # Curl/Bash from quickstart if desired, or this:
     tools/install.sh --isolated install

#. Restart dr-provision, as stated by the tools/install.sh output.

For non-isolated :ref:`rs_install`, update this way:

#. Stop dr-provision, using the system method of choice
   ::

     systemctl stop dr-provision

   or

   ::

     service dr-provision stop

#. Install new code - Use the same install technique as the first install.  :ref:`rs_install`
#. Start up dr-provision

  ::

    systemctl start dr-provision

  or

  ::

    service dr-provision start



Version to Version Notes
========================

In this section, notes about migrating from one release to another will be added.

v3.0.0 to v3.0.1
----------------
If parameters were added to machines or global, these will need to be manually re-added to the machine or 
global profile, respectively.  The machine's parameter setting cli is unchanged.  The global parameters will
need to be changed to a profiles call.

  ::
    
    drpcli parameters set fred greg

  to

  ::
    
    drpcli profiles set global fred greg


v3.0.1 to v3.0.2
----------------
There are changes to templates and bootenvs.  Upgrade will not update these automatically, because they may be in
use and working properly.  it is necessary to restart by removing the bootenvs and templates directory in
the data store directory (usually drp-data/digitalrebar or /var/lib/dr-provision/digitalrebar) and re-uploading
the bootenvs and templates (tools/discovery-load.sh).  Additionally, templates and bootenvs can be manually added and updated,
with drpcli.

v3.0.2 to v3.0.3
----------------
This is a quick turn release to address the issue with updating bootenvs.  This is a CLI code and docs only change.

v3.0.3 to v3.0.4
----------------
Nothing needs to be done.

v3.0.4 to v3.0.5
----------------
Nothing needs to be done.

v3.0.5 to v3.1.0
----------------

Subnet Enabled
~~~~~~~~~~~~~~

The subnet objects have an enabled flag that allows for subnets to be turned off without deleting them.  This value
defaults to false (off).  To enable existing subnets, you will need to do the following for each subnet in your system:

  ::

    drpcli subnets update subnet1 '{ "Enabled": true }'

Replace *subnet1* with the name of your subnet.

Local UI Removed
~~~~~~~~~~~~~~~~

The old UI has been removed and a redirect to the cloud-based UI is present instead.  The UI loads into the browswer
and then uses the API to access the Endpoint.  The DRP endpoint does not talk to the internet.  The browser acts as
a bridge for content transfers.  The only requirement is that the browser has access to the Endpoint and HTTPS-based
access to the internet.  The HTTPS-based access can be through a web proxy.

Install Script Changed
~~~~~~~~~~~~~~~~~~~~~~

There are minor changes to the install script for isolated mode.  Production mode installs are still done and updated
the same way.  For isolated, there are some new flags and options.  Please see the commands output for more details or 
check the updated :ref:`rs_quickstart`.

