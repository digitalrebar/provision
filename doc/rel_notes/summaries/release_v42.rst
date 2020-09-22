.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Release v4.2
  pair: Digital Rebar Provision; Release Notes


.. _rs_release_v42:

Digital Rebar version 4.2 [Winter 2019]
---------------------------------------

Release Date: December 31, 2019

Release Themes: Elevated Permissions and Access from Endpoint

Digital Rebar v4.2 significantly expanded the ability to operate in diverse operational environments with the introduction of Contexts and other Endpoint improvements.  Taken together, they expand and streamline the ability for workflows to securely modify and interact with a broad range of support systems.

Along with improvements to the core platform, RackN included numerous usability, bug fixes, and content extensions.  Some of the notable ones include enhanced elevated permissions for tasks, filtering enhancements that were tied into the UX and integrating commonly used tools into Digital Rebar components.

Note: since the v4.3 release included major architectural work and required more development time; consequently, v4.2 needed more patches than most RackN release cycles.

See :ref:`rs_release_summaries` for a complete list of all releases.

.. _rs_release_v42_otheritems:

Items of Note
~~~~~~~~~~~~~

* JQ integrated into DRP CLI - eliminates the need for installing JQ as a stand alone requirement.
* BSDTAR integrated into DRP Server - eliminates the need for installing bsdtar on the server
* Human/table formats for DRPCLI - makes it easier to read DRP output
* DRP CLI can install itself as a persistent agent - from script contexts makes it easier to run DRP runner as a service
* Plugins use Websocket events - migrate for legacy event model to improve stability.
* UX Save in Header - makes it easier to have multiple edits on the same item
* UX diff view - allow operator to compare new vs installed content (fixes regression)
* Updated license retrieval and install process - streamline process so that multiple users do not need RackN logins if a license is installed
* UX bulk actions for Endpoints - useful for multi-site management


Beta Items
~~~~~~~~~~

Bootstrap Agent Integrated into DRP Server - allows for on DRP server operations without use of Dangerzone context (v4.2)