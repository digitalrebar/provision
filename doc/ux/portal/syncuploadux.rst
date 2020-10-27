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

* Preview - Provide a complete look into the Content Package including its Boot Envs, Params, Profiles, Stages, Tasks, and Templates. An option to see the raw JSON is also available
* Diff - Show the latest changes between the current and most recent version of the Content Package
* Download - Download the Content Package to latest (or other selected) version
* Remove - Remove the Content Package from the current DRP endpoint. Either by selecting multiple rows or using the version dropdown.

The top of the page has a set of blue buttons for additional information:

* Refresh - Update the current list of information on the page
* Install - Install selected content rows
* Remove - Uninstall selected content rows
* Upload - Upload content from a file
* Options - Change table display options
    * Installed Only - Only display installed content
    * Not Installed - Only display not installed content
    * Show Internal - Show DRP Internal content
    * Display Names - Show package display names instead of display IDs
    * Short Versions - Truncate version extension off of table versions
* Search - Search all content regardless of selected categories 
* Categories - Toggle categories for available content

Here's an overview of clickable things in the catalog table. All links that do not navigate internally within the UX will open in a new tab.

.. image:: images/catalog_clickables.png
  :width: 800
  :alt: Clickable objects in the catalog table

Boot ISOs
---------
This page shows all available Boot ISOs and Images for the DRP endpoint to use. 

The top of the page has a set of blue buttons for additional information:

* Refresh - Update the current list of ISOs
* Upload - Add a new ISO image to the DRP endpoint
* Delete - Remove an ISO image from the DRP endpoint 

Support Files
-------------
These files are located on the DRP webserver and are available via TFTP to start the PXE boot on new machines.  

The top of the page has a set of blue buttons for additional information:

* Refresh - Update the current list of files and folders
* Upload - Add a new file/folder to the DRP endpoint
* Delete - Remove a file/folder from the DRP endpoint 

Rows in the table with the folder icon can be clicked to preview its respective folder.

Rows in the table with the file icon can be clicked to download the respective file.

The folder labeled ".." will go up one directory.

The blue links in the "Root / path / folder / names" header can be clicked to traverse parent folders.
