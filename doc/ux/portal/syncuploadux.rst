.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; UX

.. _rs_syncuploadux:

Synchronize and Upload
======================
This section contains information on the available content packages, ISOs, and Plugins available in the current DRP as well as tools not yet installed on this specific endpoint. 

Content Packages
----------------
Two separate tables are shown containing the following information:

* Endpoint Content - Manage all the installed content packages on the DRP endpoint : An Upgrade button indicates if the content packages have been updated and a new version needs to be downloaded
* Community & User Name Content - Additional content packages available to the user but not on the DRP endpoint  

Content Package Information

* Preview - Provide a complete look into the Content Packag including its Boot Envs, Params, Profiles, Stages, Tasks, and Templates. An option to see the raw JSON is also available
* Diff - Show the latest changes between the current and most recent version of the Content Package
* Upgrade - Upgrade the Content Package to its latest version
* Remove - Remove the Content Package from the current DRP endpoint

The top of the page has a set of blue boxes for additional information:

* Catalog - Gives a complete list of all possible content packages for the DRP endpoint
* Refresh - Update the current list of information on the page
* Show All - Show all options installed and not installed for the DRP endpoint 

Boot ISOs
---------
This page shows all available Boot ISOs and Images for the DRP endpoint to use. 

The top of the page has a set of blue boxes for additional information:

* Refresh - Update the current list of ISOs
* Upload - Add a new ISO image to the DRP endpoint
* Delete - Remove an ISO image from the DRP endpoint 

Plugin Providers
----------------
This section contains two separate tables showing Plugins available for the DRP endpoint.

* Endpoint Plugin Providers - Manage all the installed Plugins on the DRP endpoint.
* Organization Plugin Providers - Additional Plugins my organization can run however are not yet installed on my DRP endpoint

Plugin Information

* Plugin Name
* Upgrade - Get the latest code for a specific Plugin on the DRP endpoint
* Remove - Remove the Plugin from the DRP endpoint which will have it appear on the Organization Plugin Providers list
* Transfer - Move the Plugin to the DRP endpoint from the Organization Plugin Providers list

The top of the page has a set of blue boxes for additional information:

* Catalog - Gives a complete list of all possible Plugins for the DRP endpoint
* Upload - Add a new Plugin to the DRP endpoint 
* Refresh - Update the current list of information on the page

Support Files
-------------
These files are located on the DRP webserver and are available via TFTP to start the PXE boot on new machines.  
