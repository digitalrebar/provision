.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Upgrade

.. _rs_upgrade:

Upgrade
~~~~~~~

While not glamorous, existing code can be overwritten by a new install and 
restart.  That is about it.  Here are a few more details.

We recommend that you backup your existing install as a safey measure.


Backup
======

It's always a good policy to backup any important data, configuration, and 
content information that may be related to an application before an upgrade.  
We strongly encourage you to backup your content prior to doing any upgrade activity.

Isolated Install
----------------

For "isolated" modes (eg. originally installed with something like
``install.sh | bash -s -- install --isolated``) , perform the following tasks:

#. log in to your Provision server as the user you performed the original install as
#. copy the drp-data directory to a backup location:
   ::

     D=`date +%Y%m%d-%H%M%S`
     cp -r drp-data drp-data.backup.$D


Production Install
------------------

For "production" install modes (no ``--isolated`` flag provided to ``install.sh``), perform the following tasks

#. log in to  your Digital Rebar Provision server
#. for DRP version 3.0.5 or older:
   ::

     D=`date +%Y%m%d-%H%M%S`
     mkdir backups
     sudo cp -r /var/lib/dr-provision backups/dr-provision.backup.$D
     sudo cp -r /var/lib/tftpboot backups/dr-provision.backup.$D

#. for DRP version 3.1.0 or newer:
   ::

     D=`date +%Y%m%d-%H%M%S`
     mkdir backups
     sudo cp -r /var/lib/dr-provision backups/dr-provision.backup.$D

Upgrade Steps
=============

The basic steps are the same for both Isolated and Production install modes:

  1. stop the existing service
  2. run the installer with the "upgrade" flag
  3. restart the service

Isolated Install
----------------

For isolated :ref:`rs_install`, update this way:

#. Stop dr-provision:
   ::

      killall dr-provision

#. Return to the install directory
#. Run the install again
   ::

     # Remember to use --drp-version to set a version other than stable if desired
     # Curl/Bash from quickstart if desired, or this:
     tools/install.sh upgrade --isolated 

#. Restart *dr-provision*, as stated by the ``tools/install.sh`` output.

Production Install
------------------

For non-isolated (aka "production mode") :ref:`rs_install`, update this way:

#. Stop dr-provision, using the system method of choice
   ::

     systemctl stop dr-provision

   or

   ::

     service dr-provision stop

#. Install new code - Use the same install technique as the first install, but change ``install`` to ``upgrade`` option.  :ref:`rs_install`
   ::

     tools/install.sh upgrade --isolated

#. Start up dr-provision

  ::

    systemctl start dr-provision

  or

  ::

    service dr-provision start



Version to Version Notes
========================

In this section, notes about migrating from one release to another will be added.

Release Notes for each version can be found at:  https://github.com/digitalrebar/provision/releases

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
`Release Notes for v3.1.0 <https://github.com/digitalrebar/provision/releases/tag/v3.1.0>`_

The v3.1.0 ``install.sh`` script now supports an ``--upgrade`` flag.  Depending on your installation method (eg ``isolated`` or ``production`` mode), the behavior of the flag will alter the installation process slightly.  Please ensure you `Backup`_ your content and configurations first just in case.

For ``isolated`` mode:

  ::

    install.sh --upgrade --isolated install


.. note:: You must be in the same directory path that you performed the initial install from for the upgrade to be successful.


For ``production`` mode:

The ``production`` mode update process will move around several directories and consolidate them to a single location.  In previous versions (v3.0.5 and older), the following two default directories were used in ``production`` mode:

  ::

    /var/lib/dr-provision - Digital Rebar Provision configurations and information
    /var/lib/tftpboot - TFTP boot root directory for serving content when TFTPD service enabled

In DRP v3.1.0 and newer, the content will be moved by the ``--upgrade`` function as follows:

  ::

    /var/lib/dr-provision/digitalrebar - old "dr-provision" directory
    /var/lib/dr-provision/tftpboot - old "tftpboot" directory


.. note:: Digital Rebar Provision version 3.1.0 introduced a new behavior to the ``subnets`` definitions.  ``subnets`` may now be ``enabled`` or ``disableed`` to selectively turn on/off provisioning for a given subnet.  By default, a subnet witll be disabled.  After an upgrade, you MUST enable the subnet for it to function again. See `Subnet Enabled`_ for additional details.


Subnet Enabled
++++++++++++++

Starting in v3.1.0, subnet objects have an enabled flag that allows for subnets to be turned off without deleting them.  This value defaults to false (off).  To enable existing subnets, you will need to do the following for each subnet in your system:

  ::

    drpcli subnets update subnet1 '{ "Enabled": true }'

Replace *subnet1* with the name of your subnet.  You may obtain a list of configured subnets with:

  ::

    drpcli subnets list | jq -r '.[].Name'


v3.1.0 to v3.2.0
----------------

`Release Notes for v3.2.0 <https://github.com/digitalrebar/provision/releases/tag/v3.2.0>`_

There are fairly significant updates to the DRP Contents structure and layout in v3.2.0.  If you are upgrading to v3.2.0 you must remove any Digital Rebar and RackN content that you have installed in your Provisioning endpoint.  The following outline will help you understand the necessary steps.  If you have any issues with the upgrade process, please drop by the Slack #community channel for additional help.

Please read the steps through carefully, and make note of the current contents/plugins you currently have installed.  You will have to re-add these elements again.  You absolutely should backup your existing install prior to this upgrade.

  1. Overview

    Overiew of the update steps necessary, you should do in the following order.

    1. Update DRP to "stable" (v3.2.0)
    #. Remove Old Content
    #. Add Content back that was removed
    #. Update plugins
    #. Fix up things

  2. Updating DRP Endpoint

    If you are running isolated, do this (remove ``--isolated`` if you are not using isolated mode):

    ::

      curl -fsSL get.rebar.digital/stable | bash -s -- upgrade --isolated

    This will force the update of the local binaries to v3.2.0 stable.  Make sure you stop DRP process (``sudo killall dr-provision``, or ``sudo systemctl stop dr-provision.service``).

    Verify that your ``/etc/systemd/services/dr-provision`` start up file is still correct for your environment, if running a production install type.

    Restart DRP (follow ``--isolated`` mode start steps if in isolated mode; or ``sudo systemctl start dr-provision.service``)

    If in ``--isolated`` mode, don’t forget to copy ``drpcli`` and/or ``dr-provision`` binaries to where you prefer to keep them (eg ``$HOME/bin`` or ``/usr/local/bin``, etc... .

  3. Remove old content

    With the rework of content, you need to remove the following content packages if they were previously installed.

    ::

      os-linux
      os-discovery
      drp-community-content (if you are really behind, Digital Rebar Community Content).
      ipmi
      packet
      virtualbox

  4. Put the content back

    Install the new v3.2.0 content packs.  Note that the names have changed, and the mix of "ce-" and non-Community Content names has gone away.  For example; what originally was ``drp-community-content`` which included things like ``ce-sledgehammer`` is now moved to just ``sledgehammer``.  The RackN registered content of ``os-linux`` and ``os-discovery`` have now been folded in to the below content packs.

    ::

      drp-community-content - it is a must just get it.
      task-library - New RackN library of services for doing interesting things.
      drp-community-contrib - this is old or experimental things like centos6 or SL6.

  5. Update the plugins

    If you have any plugins installed, update them now.

    To facilitate version tracking, plugins provide their own content as a injected content from the plugin.  When the plugin is added, it will also add a content layer that will show up in the content packages section.

    Previously, a ``plugin-provider`` was installed separately from a Content of the same name.

  6. Fix things up

    This is mainly if you were using the Community Content version of things (``drp-community-content``, and BootEnvs with a prefix of ``ce-``).  The BootEnvs names change, by removing the prefix of "ce-" from the name.

    Make sure all the bootenvs are up to date and available.  This is a task you should always do after updating content.  If the BootEnv is marked with an "X" in the UX, or ``"Available": false`` from the CLI/API, you'll need to reload the ISO for the BootEnv.

    Then go to *Info & Preferences* and make sure your default stage and bootenvs are still valid.

    - This is where ``ce-sledgehammer`` becomes ``sledgehammer`` and ``ce-discovery`` becomes ``discovery``
    - The same with ``ce-ubuntu-16.04-install`` becomes ``ubuntu-16.04-install``.
    - The same with ``ce-centos-7.4.1708-install`` becomes ``centos-7-install``.

  Example pseudo-script to make changes:

    Please carefully read through this script and make sure it correlates to your installed content.  It is provided only as an example, and will absolutely require (possibly just minor) modifications for your environment.

    YOU MUST MODIFY THE *RACK_AUTH* variable appropriately for the download authentication to work correctly.

    ::

      # see all contents
      drpcli contents list

      # list JUST the names of the contents - note what you have installed,
      # you may need to re-install it below
      drpcli contents list | jq -r '.[].meta.Name' | egrep -v "BackingStore|BasicStore"

      # list which plugins you have installed - note it, you may need to install
      # it below
      drpcli plugin_providers list | jq '.[].Name'

      # go to RackN UX - log in, go to Hamburger menu (upper left, 3 horizontal lines)
      # go to Organization - User Profile - copy your UUID for Unique User Identity
      export RACKN_AUTH="?username=<UUID_Unique_User_Identity>"
      export CATALOG="https://qww9e4paf1.execute-api.us-west-2.amazonaws.com/main/catalog"

      # get raw output of just the content packs
      for CONTENT in `drpcli contents list | jq -r '.[].meta.Name' | egrep -v "BackingStore|BasicStore"`
      do
        echo "remove content:   $CONTENT"
        drpcli contents destroy $CONTENT
      done

      # install content
      for CONTENT in drp-community-content task-library drp-community-contrib
      do
        echo "install content:  $CONTENT"
        curl -s $CATALOG/content/${CONTENT}${RACKN_AUTH} -o $CONTENT.json
        drpcli contents create -< $CONTENT.json
      done

      # change "plug1", "plug2", etc... to the plugin provider names you need
      # examples:  "slack", "packet-ipmi", "ipmi"
      for PLUGIN in plug1 plug2 plug3
      do
        echo "install plugin:  $PLUGIN"
        curl -s $CATALOG/plugin/${PLUGIN}${RACKN_AUTH} -o $PLUGIN.json
        drpcli contents create -< $PLUGIN.json
      done

      # Ensure the Stage, Default, and Unknown BootEnv are set to valid values
      # adjust these as appropriate
      drpcli prefs set defaultStage discover defaultBootEnv sledgehammer unknownBootEnv discovery

    Again - make sure you modify things appropriately in the above scriptlet. YOU MUST MODIFY THE *RACK_AUTH* variable appropriately for the download authentication to work correctly.

v3.2.0 to v3.3.0
----------------

`Release Notes for v3.3.0 <https://github.com/digitalrebar/provision/releases/tag/v3.3.0>`_

No aditional steps required.

v3.3.0 to v3.4.0
----------------

`Release Notes for v3.4.0 <https://github.com/digitalrebar/provision/releases/tag/v3.4.0>`_

Content Changes
+++++++++++++++

Prior to restart Digital Rebar Provision endpoint - you may need to fix the Machines JSON entries for the ``Meta`` field.  It used to be an optional field, but is now required.  If your ``Meta`` field is set to ``null``, or non-existent, DRP will not startup correctly.  You will receive the following error message on start:
  ::

    dr-provision2018/01/07 15:14:01.275082 Extracting Default Assets
    panic: assignment to entry in nil map

To correct the problem, you will need to edit your JSON configuration files for your Machines. You can find your Machines spec files in ``/var/lib/dr-provision/digitalrebar/machines`` if you are running in *production* mode install.  If you are running in *isolated* mode, you will need to locate your ``drp-data`` directory which is in the base directory where you performed the install at; the machines directory will be ``drp-data/digitalrebar/machines``.

There may be two ``Meta`` tags.  You do NOT need to modify the ``Meta`` tag that is located in the *Params* section.

Change the first ``Meta`` tag as follows:
  ::

      # from:
      "Meta":null,

      # to something like:
      "Meta":{"feature-flags":"change-stage-v2"},

It is entirely possible that the ``Meta`` field is completely missing.  If so - inject the full ``Meta`` field as specified above.

``drpcli`` changes
++++++++++++++++++

Please see the `Release Notes <https://github.com/digitalrebar/provision/releases/tag/v3.4.0>`_ for information related to the ``drpcli`` command line changes.  The most notable changes that may impact your use (eg in existing scripts) of the tool:

#. Plugin upload method changed:

  ::

    # prior to v3.4.0
    drpcli plugin_providers upload $PLUGIN as $PLUG_NAME

    # v3.4.0 and newer version method:
    drpcli plugin_providers upload $PLUG_NAME from $PLUGIN

2. Many commands now have new *helper* capabilities.  See each command outputs relevant help statement.


v3.4.0 to v3.5.0
----------------

`Release Notes for v3.5.0 <https://github.com/digitalrebar/provision/releases/tag/v3.5.0>`_

No additional changes necessary.

v3.5.0 to v3.6.0
----------------

`Release Notes for v3.6.0 <https://github.com/digitalrebar/provision/releases/tag/v3.6.0>`_

No additional changes necessary.

v3.6.0 to v3.7.0
----------------

`Release Notes for v3.7.0 <https://github.com/digitalrebar/provision/releases/tag/v3.7.0>`_

The plugin system has been updated to a new version.  All plugins have been updated to
use the new version.  After updating to *v3.7.0*, all plugins must be updated to function.
The system will start after update, but the plugin-providers will not load until they are
udpated.  Use the RackN UX to get the updates for the plugins.

The Task subsystem has been updated to default to `sane-exit-codes`.  This is a change from
the default of `original-exit-codes`.  This was done to address the need of task authors to
match some basic assumptions about exit codes.  *1* should be a fail and not reboot your box.

Additionally, the default UX redirect has changed to the `stable portal <https://portal.rackn.io>`_.
This will result in more stable UX experience.

v3.7.0 to v3.8.0
----------------

`Release Notes for v3.8.0 <https://github.com/digitalrebar/provision/releases/tag/v3.8.0>`_

No additional changes necessary.

v3.8.0 to v3.9.0
----------------

`Release Notes for v3.9.0 <https://github.com/digitalrebar/provision/releases/tag/v3.9.0>`_

No additional changes necessary.

Local UI Removed
~~~~~~~~~~~~~~~~

The old UI has been removed and a redirect to the RackN Portal UI is present instead.  The UI loads into the browswer and then uses the API to access the Endpoint.  The DRP endpoint does not talk to the internet.  The browser acts as a bridge for content transfers.  The only requirement is that the browser has access to the Endpoint and HTTPS-based access to the internet.  The HTTPS-based access can be through a web proxy.

Install Script Changed
~~~~~~~~~~~~~~~~~~~~~~

There are minor changes to the install script for isolated mode.  Production mode installs are still done and updated the same way.  For isolated, there are some new flags and options.  Please see the commands output for more details or check the updated :ref:`rs_quickstart`.

For current ``install.sh`` script usage information, please run:

  ::

    install.sh --help


For complete details.

