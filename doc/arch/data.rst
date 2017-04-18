.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. DigitalRebar Provision documentation under Digital Rebar master license
.. index::
  pair: DigitalRebar Provision; Data Architecture

.. _rs_data_architecture:

Data Architecture
=================

DigitalRebar Provision uses a fairly simple data model.  There are 4 main models for the provisioner server
and 4 models for the DHCP server.  Each model has a cooresponding API in the :ref:`rs_api`.

Provisioner Models

These models represent things that the provisioner server use and manipulate.

* :ref:`rs_model_machine`
* :ref:`rs_model_bootenv`
* :ref:`rs_model_template`
* :ref:`rs_model_param`

DHCP Models

These models represent things that the DHCP server use and manipulate.

* :ref:`rs_model_subnet`
* :ref:`rs_model_reservation`
* :ref:`rs_model_lease`
* :ref:`rs_model_interface`

Additional Models

These models control additional parts and actions of the system.

* :ref:`rs_model_user`
* :ref:`rs_model_prefs`

Special Objects

These are not objects in the system but represent files and directories in the server space.

* :ref:`rs_model_file`
* :ref:`rs_model_iso`


.. _rs_models:

Models
------

The models define elements of the system.  The API provides basic CRUD (create, read, update, and delete) operations as well as
some additional actions for manipulating the state of the system.  The :ref:`rs_api` contains that definitions of the actual
structures and methods on those objects.  Additionally, the :ref:`rs_operation` will describe common actions to use and do with
these models and how to build them.  The :ref:`rs_cli` describes the Command Line manipulators for the model.

This section will describe its use and role in the system.


.. index::
  pair: Model; Machine

.. _rs_model_machine:

Machine
~~~~~~~

The Machine object defines a machine that is being provisioned.  The Machine is represented by a unique **UUID**.  The UUID 
is immutable after machine creation.  The machine's primary purpose is to map an incoming IP address to a :ref:`rs_model_bootenv`.
The :ref:`rs_model_bootenv` provides a set of rendered :ref:`rs_model_template` that will can be used to boot the machine.  The 
machine provides parameters to the :ref:`rs_model_template`.  The Machine provides configuration to the renderer in the 
form of parameters and fields.  The **Name** field should contain the FQDN of the node.

The Machine object contains an **Error** field that represents errors encountered while operating on the machine.  In general,
these are errors pertaining to rendering the :ref:`rs_model_bootenv`.

The Machine parameters are defined as a field on the Machine that is presented as a dicitionary of string keys to arbritary objects.
These could be strings, bools, numbers, arrays, or objects represented similarly defined dictionaries.  The machine parameters
are available to templates for expansion in them.

.. index::
  pair: Model; BootEnv

.. _rs_model_bootenv:

BootEnv
~~~~~~~

The BootEnv object defines an environment to boot a machine.  It has two main components an OS information section and a templates
list.  The OS information section defines what makes up the installation base for this bootenv.  It defines the install ISO, a
URL to get the ISO, and SHA256 checksum to validate the image.  These are used to provide the basic install image, kernel, and
base packages for the bootenv.

The other primary section is a set of templates that present files in the file server's name space that can served via HTTP or 
TFTP.  The templates can be in-line in the BootEnv object or reference a :ref:`rs_model_template`.  The templates are specified as
a list of paths in the filesystem and either an ID of a :ref:`rs_model_template` or inline content.  The path field of the 
template information can use the same template expansion that is used in the template.  See :ref:`rs_model_template` for more
information.

Additionally, the BootEnv defines required and optional parameters.  The required parameters validated at render time to be
present or an error is generated.  These parameters can be met by the parameters on the machine or from the global :ref:`rs_model_param`
space.


.. index::
  pair: Model; Template

.. _rs_model_template:

Template
~~~~~~~~

The Template object defines a templated content that can be referenced by its ID.  The content of the template (or 
in-line template in a :ref:`rs_model_bootenv`) is a `golang text/template <https://golang.org/pkg/text/template/#hdr-Actions>`_ string.
The template has a set of special expansions.  The normal expansion syntax is:

  ::

    {{ .Machine.Name }}

This would expand to the machine's **Name** field.  There are helpers for the parameter spaces, the :ref:`rs_model_bootenv` object,
and some miscellaneous functions.  Additionally, the normal `golang text/template <https://golang.org/pkg/text/template/#hdr-Actions>`_
functions are available as well.  Things like **range**, **len**, and comparators are available as well.  Currently, **template** inclusion
is not supported.

The following table lists the current set of expansion custom functions:

============================== =================================================================================================================================================================================================
Expansion                      Description
============================== =================================================================================================================================================================================================
.Machine.Name                  The FQDN of the Machine in the Machine object stored in the **Name** field
.Machine.ShortName             The Name part of the FDQN of the Machine object stored in the **Name** field
.Machine.UUID                  The Machine's **UUID** field
.Machine.Path                  A path to a custom machine unique space in the file server name space.
.Machine.Address               The **Address** field of the Machine
.Machine.HexAddress            The **Address** field of the Machine in Hex format (useful for elilo config files
.Machine.URL                   A HTTP URL that references the Machine's specific unique filesystem space.
.Env.PathFor <proto> <file>    This references the boot environment and builds a string that presents a either a tftp or http specifier into exploded ISO space for that file.  *Proto* is **tftp** or **http**.  The *file* is a relative path inside the ISO.
.Env.InstallURL                An HTTP URL to the base ISO install directory.
.Env.JoinInitrds <proto>       A comma separated string of all the initrd files specified in the BootEnv reference through the specified proto (**tftp** or **http**)
.ProvisionerAddress            An IP address that is on the provisioner that is the most direct access to the machine.
.ProvisionerURL                An HTTP URL to access the base file server root
.ApiURL                        An HTTPS URL to access the DigitalRebar Provision API
.GenerateToken                 This generates limited use access token for the machine to either update itself if it exists or create a new machine.  The token's validity is limited in time by global preferences.  See :ref:`rs_model_prefs`.
.ParseURL <segment> <url>      Parse the specified URL and return the segment requested.
.ParamExists <key>             Returns true if the specified key is a valid parameter available for this rendering.
.Param <key>                   Returns the structure for the specified key for this rendering.
============================== =================================================================================================================================================================================================

.. note::
  **.Machine.Path** is particularly useful for ensure that templates are expanded into a unique file space for
  each machine.  An example of this is per machine kickstart files.  These can be seen in the **assets/bootenvs/ubuntu-16.04.yml**.

With regard to the **.Param** and **.ParamExists** functions, these functions return the parameter or existence of
the parameter specified by the *key* input.  The parameters are examined from most specific to global.  This means
that the Machine object's parameters are checked first, then the global :ref:`rs_model_param`.  The parameters on machines
and the global space are free form dictionaries and default empty.  Any key/value pair can be added and referenced.

The default :ref:`rs_model_template` and :ref:`rs_model_bootenv` use the following optional (unless marked with an \*)
parameters.

=================================  ================  =================================================================================================================================
Parameter                          Type              Description
=================================  ================  =================================================================================================================================
ntp_servers                        Array of objects  lookup format
proxy-servers                      Array of objects  lookup format
operating-system-disk              String            A string to use as the default install drive.  /dev/sda or sda depending upon kickstart or preseed.
access_keys                        Map of strings    The key is the name of the public key.  The value is the public key.  All keys are placed in the .authorized_keys file of root.
provisioner-default-password-hash  String            The password hash for the initial default password, **RocketSkates**
provisioner-default-user           String            The initial user to create for ubuntu/debian installs
dns-domain                         String            DNS Domain to use for this system's install
\*operating-system-license-key     String            Windows Only
\*operating-system-install-flavor  String            Windows Only
=================================  ================  =================================================================================================================================

For some examples of this in use, see :ref:`rs_operation`.

.. index::
  pair: Model; Param

.. _rs_model_param:

Param
~~~~~


.. index::
  pair: Model; Subnet

.. _rs_model_subnet:

Subnet
~~~~~~


.. index::
  pair: Model; Reservation

.. _rs_model_reservation:

Reservation
~~~~~~~~~~~


.. index::
  pair: Model; Lease

.. _rs_model_lease:

Lease
~~~~~


.. index::
  pair: Model; Interface

.. _rs_model_interface:

Interface
~~~~~~~~~



.. index::
  pair: Model; User

.. _rs_model_user:

User
~~~~


.. index::
  pair: Model; Prefs

.. _rs_model_prefs:

Prefs
~~~~~


.. index::
  pair: Model; Files

.. _rs_model_file:

Files
~~~~~


.. index::
  pair: Model; Isos

.. _rs_model_iso:

Isos
~~~~






