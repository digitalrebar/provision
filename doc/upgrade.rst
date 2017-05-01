.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Upgrade

.. _rs_upgrade:

Upgrade
~~~~~~~

While not glamorous, you can install over the existing code and restart.  That is about it.  Here are few more details.

Steps
=====

For isolated :ref:`rs_install`, update this way:

#. Stop dr-provision:
   ::

      killall dr-provision

#. Return to your install directory
#. Run the install again
   ::

     rm sha256sums
     # Remeber to use --drp-version is you want something other than stable
     # Curl/Bash from quickstart if you truly believe, or this:
     tools/install.sh --isolated install

#. Restart dr-provision, as stated by the tools/install.sh output.

For non-isolated :ref:`rs_install`, update this way:

#. Stop dr-provision, using your system method of choice
   ::

     systemctl stop dr-provision

   or

   ::

     service dr-provision stop

#. Install new code - How ever you installed before, do it again.  :ref:`rs_install`
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
If parameters were added to machines or global, these will need to be manually readded to the machine or 
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
use and working for you.  You will need to start over by removing the bootenvs and templates directory in
your data store directory (usually drp-data/digitalrebar or /var/lib/dr-provision/digitalrebar) and re-uploading
the bootenvs and templates (tools/discovery-load.sh).  Or you can manually add and update, templates and bootenvs
with drpcli.


