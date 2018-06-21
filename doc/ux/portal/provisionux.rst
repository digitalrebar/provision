.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; UX

.. _rs_provisionux:

Provision
=========
This section contains the setup information to provision a machine. The workflow process uses this information when configuring and deploying a new machine.

Boot Environments
-----------------
Configuring at least one Boot Environment is a critical first step in Digital Rebar Provision operation. The Digital Rebar CentOS based in-memory discovery image, Sledgehammer, will be installed on first use by default.

The UI will show a complete list of potential Boot Environments with the following information:

* Locked or Unlocked 
* Name - The name of the Boot Environment. 
* D/L 
* ISO - ISO or TAR Image involved with the named Boot Environment  
* Description - More information about the specific Boot Environment 

The top of the page includes a set of additional actions:

* Refresh - Refresh the list of available Boot Environments for the Endpoint to use
* Filter - Refine the list of Boot Environments based on these options: Available, Key, Name, OnlyUnknown, OSName, ReadOnly, Valid
* Add - Add a new Boot Environment 
* Clone - Clone a selected Boot Environment 
* Delete - Remove a selected Boot Environment

Templates
---------
Templates contain important instructions for the provisioning process, and are comprised of golang text/template strings. Once templates are rendered along with any assigned parameters, they are used by the BootEnv to boot the target machine. Templates may contain other templates, known as sub-templates.

The UI will show a complete list of Templates with the following information:

* Locked or Unlocked
* ID - Template name
* Preview - Description of what the Template does 

The top of the page includes a set of additional actions:

* Refresh - Refresh the list of available Templates for the Endpoint to use
* Filter - Refine the list of Templates based on these options: Available, ID, Key, ReadOnly, Valid
* Add - Add a new Template 
* Clone - Clone a selected Template 
* Delete - Remove a selected Template

Params
------
Parameters are passed to a template from a machine, and help to drive the template’s functions. They consist of key/value pairs that provide configuration to the renderer. Profiles allow params to be applied in bulk, or they can be attached to templates individually.

The UI will show a complete list of Params with the following information:

* Locked or Unlocked
* Name - Name of the Param
* Type - Value of the Param; e.g. Object, String, Array, etc
* Description - Additional information about the Param

The top of the page includes a set of additional actions:

* Refresh - Refresh the list of available Params for the Endpoint to use
* Filter - Refine the list of Params based on these options: Available, Key, Name, ReadOnly, Valid
* Add - Add a new Param        
* Clone - Clone a selected Param
* Delete - Remove a selected Param

Profiles
--------
Profiles provide a convenient way to apply sets of parameters to a machine. Multiple profiles can be assigned to one machine, and will be referenced in the order they are listed. Parameters can be linked to specific profiles through the profiles page, which can then be attached to machines through the machines UI.

The UI will show a complete list of Profiles with the following information:

* Locked or Unlocked
* Name - Name of the Profiles 
* Description - Additional information about the Profiles

The top of the page includes a set of additional actions:

* Refresh - Refresh the list of available Profiles for the Endpoint to use
* Filter - Refine the list of Params based on these options: Available, Key, Name, ReadOnly, Valid
* Add - Add a new Profile 
* Clone - Clone a selected Profile
* Delete - Remove a selected Profile
* Ansible - Present information about the Ansible Inventory Grid
