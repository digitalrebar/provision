.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Release v3.x
  pair: Digital Rebar Provision; Release Notes


.. _rs_release_v3x:

Digital Rebar version 3.x
-------------------------

Historical notes about v3 releases.

v3.6.0 to v3.7.0
~~~~~~~~~~~~~~~~

The plugin system has been updated to a new version.  All plugins have been updated to
use the new version.  After updating to *v3.7.0*, all plugins must be updated to function.
The system will start after update, but the plugin-providers will not load until they are
udpated.  Use the RackN UX to get the updates for the plugins.

The Task subsystem has been updated to default to `sane-exit-codes`.  This is a change from
the default of `original-exit-codes`.  This was done to address the need of task authors to
match some basic assumptions about exit codes.  *1* should be a fail and not reboot your box.
