.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Architecture

.. _rs_architecture:


Architecture
~~~~~~~~~~~~

Digital Rebar Provision is intended to be a very simple service that can run with minimal overhead in nearly any environment.  For this reason, all the needed components are combined into the Golang binary server including the UI and Swagger UI assets.  The binary can be run as a user process or easily configured as a operating system service.

The service is designed to work with multiple backend data stores.  For stand alone operation, data is stored on the file system.  For Digital Rebar integration, data can be maintained in Consul.

The CLI is provided as a second executable so that it can be used remotely.

By design, there are minimal integrations between core services.  This allows the service to reduce complexity.  Beyond serving IPs and files, the primary action of the service is template expansion for boot environments (:ref:`rs_model_bootenv`).  The template expansion system allows subsitition properties to be set on a global, groups by profile, or per machine basis.

The architecture can be described in terms of the server and its data model.

.. _rs_server_features:

Key Features
============

Digital Rebar Provision is a new generation of data center automation designed for operators with a cloud-first approach. Data center provisioning is surprisingly complex because it’s caught between cutting edge hardware and arcane protocols embedded in firmware requirements that are still very much alive.

Swagger REST API & CLI
----------------------

Cloud-first means having a great, tested API. Years of provisioning experience went into this 3rd generation design and it shows. That includes a powerful API-driven DHCP.

Security & Authenticated API
----------------------------

Not an afterthought, we use both HTTPS and user authentication for the API. Our mix of basic and bearer token authentication recognizes that both users and automation will use the API. This brings a new level of security and control to data center provisioning.

Stand-alone multi-architecture Golang binary
--------------------------------------------

There are no dependencies or prerequisites, plus upgrades are drop in replacements. That allows users to experiment isolated on their laptop and then easily register it as a SystemD service.

Nested Template Expansion
-------------------------

In Digital Rebar Provision, Boot Environments are composed of reusable template snippets. These templates can incorporate global, profile or machine specific properties that enable users to set services, users, security or scripting extensions for their environment.
Configuration at Global, Group/Profile and Node level. Properties for templates can be managed in a wide range of ways that allows operators to manage large groups of servers in consistent ways.

Multi-mode (but optional) DHCP
------------------------------

Network IP allocation is a key component of any provisioning infrastructure; however, DHCP needs are highly site dependent. Digital Rebar Provision works as a multi-interface DHCP listener and can also resolve addresses from DHCP forwarders. It can even be disabled if the environment already has a DHCP service that can configure a the “next boot” provider.

Dynamic Provisioner templates for TFTP and HTTP
-----------------------------------------------

For security and scale, Digital Rebar Provision builds provisioning files dynamically based on the Boot Environment Template system. This means that critical system information is not written to disk and files do not have to be synchronized. Of course, when a file needs to be served it works too.

Node Discovery Bootstrapping
----------------------------

Digital Rebar’s long-standing discovery process is enabled in the Provisioner with the included discovery boot environment. That process includes an integrated secure token sequence so that new machines can self-register with the service via the API. This eliminates the need to pre-populate the Digital Rebar Provision system.

Multiple Seeding Operating Systems
----------------------------------

Digital Rebar Provision comes with a long list of Boot Environments and Templates including support for many Linux flavors, Windows, ESX and even LinuxKit. Our template design makes it easy to expand and update templates even on existing deployments.

Two-stage TFTP/HTTP Boot
------------------------

Our specialized Sledgehammer and Discovery images are designed for speed with optimized install cycles the improve boot speed by switching from PXE TFTP to IPXE HTTP in a two stage process. This ensures maximum hardware compatibility without creating excess network load.

.. _rs_server_architecture:

Server Architecture
===================

Digital Rebar Provision is provided by a single binary that contains tools and images needed to operate.
These are expanded on startup and made available by the file server services.

.. _rs_design_restriction:

Design Restrictions
-------------------

Since Digital Rebar Provision is part of the larger Digital Rebar system, it's scope is limited to handling DHCP and Provisioning actions.  Out of band management to control server flow or configure firmware plus other management features will be handled by other Digital Rebar services.

.. _rs_arch_services:

Services
--------

Provisioning requires handoffs between multiple services as described in the :ref:`rs_workflows` section.  Since several of services are standard protocols (DHCP, TFTP, HTTP), it may be difficult to change ports without breaking workflow.

The figure below illustrates the three core Digital Rebar Provision services including protocols and default ports.  The services are:

#. Web - These services provide control for the other services

   #. API: REST endpoints with Swagger definition
   #. UI: User interface and Swagger helpers

#. DHCP: Address management includes numerous additional option fields used to tell systems how to interact with other data center services such as provisioning, DNS, NTP and routing.

#. Provision: sends files on request during provisioning process based on a template system:

   #. TFTP: very simple (but slow) protocol that's used by firmware boot processes because it is very low overhead.
   #. HTTP: faster file transfer protocol used by more advanced boot processes


.. figure::  ../images/core_services.png
   :alt: Core Digital Rebar Provision Services
   :target: https://docs.google.com/drawings/d/1SVGGwQZxopiVEYjIM3FXC92yG4DKCCejRBDNMsHmxKE/edit?usp=sharing


.. _rs_arch_ports:

Ports
-----

The table describes the ports that need to be available to run Digital Rebar Provision.  Firewall rules may need to be altered to enable these services.  The feature column indicates when the port is required.  For example, the DHCP server can be turned off and that port is no longer required.

========  =======   =====================
Ports     Feature   Usage
========  =======   =====================
67/udp    DHCP      DHCP Port
69/udp    PROV      TFTP Port
8091/tcp  PROV      HTTP-base File Server
8092/tcp  Always    DR Provision Mgmt
========  =======   =====================

.. _rs_data_architecture:

Data Architecture
=================

Digital Rebar Provision uses a fairly simple data model.  There are 4 main models for the provisioner server
and 4 models for the DHCP server.  Each model has a corresponding API in the :ref:`rs_api`.

The models define elements of the system.  The API provides basic CRUD (create, read, update, and delete) operations as well as
some additional actions for manipulating the state of the system.  The :ref:`rs_api` contains that definitions of the actual
structures and methods on those objects.  Additionally, the :ref:`rs_operation` will describe common actions to use and do with
these models and how to build them.  The :ref:`rs_cli` describes the Command Line manipulators for the model.

This section will describe its use and role in the system.

.. _rs_provisioner_models:

Provisioner Models
------------------

These models represent things that the provisioner server use and manipulate.


.. index::
  pair: Model; Machine

.. _rs_model_machine:

Machine
=======

The Machine object defines a machine that is being provisioned.  The Machine is represented by a unique **UUID**.  The UUID
is immutable after machine creation.  The machine's primary purpose is to map an incoming IP address to a :ref:`rs_model_bootenv`.
The :ref:`rs_model_bootenv` provides a set of rendered :ref:`rs_model_template` that will can be used to boot the machine.  The
machine provides parameters to the :ref:`rs_model_template`.  The Machine provides configuration to the renderer in the
form of parameters and fields.  Also, each :ref:`rs_model_machine` must have a :ref:`rs_model_bootenv` to boot from.
If a machine is created without a :ref:`rs_model_bootenv` specified, the system will assign the one specified by
the value of :ref:`rs_model_prefs` *defaultBootEnv*.

The **Name** field should contain the FQDN of the node.

The Machine object contains an **Error** field that represents errors encountered while operating on the machine.  In general,
these are errors pertaining to rendering the :ref:`rs_model_bootenv`.

The Machine parameters are defined as a special :ref:`rs_model_profile` on the Machine.  The profile stores a dictionary of
string keys to arbitrary objects.  These could be strings, booleans, numbers, arrays, or objects representing similarly
defined dictionaries.  The machine parameters are available to templates for expansion in them.

Additionally, the machine maintains an ordered list of profiles that are searched and then finally the **global profile**.  See :ref:`rs_model_profile` and :ref:`rs_model_template` for more information.

.. note:: When updating the Params part of the embedded Profile in the :ref:`rs_model_machine` object, using the **PUT** method will replace the Params map with the map from the input object.  The **PATCH** will merge the Params map in the input with the existing Params map in the current :ref:`rs_model_machine` object.  The **POST** method on the params subaction will replace the map with the input version.

.. index::
  pair: Model; Profile

.. _rs_model_profile:

Profile
=======

The Profile Object defines a set of key / value pairs (or parameters).  All of these may be manipulated by the :ref:`rs_api`.
The key space is a free form string and the value is an arbitrary data blob specified by JSON through
the :ref:`rs_api`.  The common parameters defined in :ref:`rs_model_template` can be set on these objects.
The system maintains a **global** profile for setting system wide parameters.  They are the lowest level of precedence.

The profiles are free form dictionaries and default empty.  Any key/value pair can be added and referenced.

Other profiles may be created to group parameters together to apply to sets of machines.  The machine's profile
list allows the administrator to specify an ordered set of profiles that apply to that machine as well.
Additionally, the system maintains a special
profile for each machine to store custom parameters specific to that machine.  This profile is embedded in the :ref:`rs_model_machine` object.

When the system needs to render a template parameter, the machine's specific profile is checked, then the order
list of profiles stored in the Machine Object are checked, and finally the **global** profile is checked.  The
key and its value are used if found in template rendering.

.. note:: When updating the Params part of the :ref:`rs_model_profile`, using the **PUT** method will replace the Params map with the map from the input object.  The **PATCH** method will merge the Params map in the input with the existing Params map in the current :ref:`rs_model_profile` object.  The **POST** method on the params subaction will replace the map with the input version.


.. index::
  pair: Model; BootEnv

.. _rs_model_bootenv:

BootEnv
=======

The BootEnv object defines an environment to boot a machine.  It has two main components an OS information section and a templates
list.  The OS information section defines what makes up the installation base for this bootenv.  It defines the install ISO, a
URL to get the ISO, and SHA256 checksum to validate the image.  These are used to provide the basic install image, kernel, and
base packages for the bootenv.

The other primary section is a set of templates that represent files in the file server's file space that can served via HTTP or
TFTP.  The templates can be in-line in the BootEnv object or reference a :ref:`rs_model_template`.  The templates are specified as
a list of paths in the filesystem and either an ID of a :ref:`rs_model_template` or inline content.  The path field of the
template information can use the same template expansion that is used in the template.  See :ref:`rs_model_template` for more
information.

Additionally, the BootEnv defines required and optional parameters.  The required parameters validated at render time to be
present or an error is generated.  These parameters can be met by the parameters on the machine, the profiles in machine's profiles list,
or from the global :ref:`rs_model_profile`.

BootEnvs can be marked **OnlyUnknown**.  This tells the rest of the system that this BootEnv is not for specific machines.  It is a
general BootEnv.  For example, *discovery* and *ignore* are **OnlyUnknown**.  *discovery* is used to discover unknown machines and
add them to Digital Rebar Provision.  *ignore* is a special bootenv that tells machines to boot their local disk.  These BootEnvs
populate the pxelinux.0, ipxe, and elilo default fallthrough files.  These are different than their counterpart BootEnvs,
*sledgehammer* and *local* which are machine specific BootEnvs that populate configuration files that are specific to a single
machine.  A machine boots *local*; an unknown machine boots *ignore*.  There can only be one **OnlyUnknown** BootEnv active
at a time.  This is specified by the :ref:`rs_model_prefs` *unknownBootEnv*.

.. index::
  pair: Model; Template

.. _rs_model_template:

Template
========

The Template object defines a templated content that can be referenced by its ID.  The content of the template (or
in-line template in a :ref:`rs_model_bootenv`) is a `golang text/template <https://golang.org/pkg/text/template/#hdr-Actions>`_ string.
The template has a set of special expansions.  The normal expansion syntax is:

  ::

    {{ .Machine.Name }}

This would expand to the machine's **Name** field.  There are helpers for the parameter spaces, the :ref:`rs_model_bootenv` object,
and some miscellaneous functions.  Additionally, the normal `golang text/template <https://golang.org/pkg/text/template/#hdr-Actions>`_
functions are available as well.  Things like **range**, **len**, and comparators are available as well.  **template** inclusion is supported by the following syntax:

  ::

    {{ template "ID of Template" }}
    {{ template .Param.MyFavoriteTemplate }}


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
.Env.OS.Family                 An optional string from the BootEnv that is used to represent the OS Family.  Ubuntu preseed uses this to determine debian vs ubuntu as an example.
.Env.OS.Version                An optional string from the BootEnv that is used to represent the OS Version.  Ubuntu preseed uses this to determine what version of ubuntu is being installed.
.Env.JoinInitrds <proto>       A comma separated string of all the initrd files specified in the BootEnv reference through the specified proto (**tftp** or **http**)
.BootParams                    This renders the **BootParam** field of :ref:`rs_model_bootenv` at that spot.  Template expansion applies to that field as well.
.ProvisionerAddress            An IP address that is on the provisioner that is the most direct access to the machine.
.ProvisionerURL                An HTTP URL to access the base file server root
.ApiURL                        An HTTPS URL to access the Digital Rebar Provision API
.GenerateToken                 This generates limited use access token for the machine to either update itself if it exists or create a new machine.  The token's validity is limited in time by global preferences.  See :ref:`rs_model_prefs`.
.ParseURL <segment> <url>      Parse the specified URL and return the segment requested.
.ParamExists <key>             Returns true if the specified key is a valid parameter available for this rendering.
.Param <key>                   Returns the structure for the specified key for this rendering.
template <string> .            Includes the template specified by the string.  String can be a variable and note that template does NOT have a dot (.) in front.
============================== =================================================================================================================================================================================================

**GenerateToken** is very special.  This generates either a *known token* or an *unknown token* for use by the template to update objects
in Digital Rebar Provision.  The tokens are valid for a limited time as defined by the **knownTokenTimeout** and **unknownTokenTimeout**
:ref:`rs_model_prefs` respectively.  The tokens are also restricted to the function the can perform.  The *known token* is limited to only
reading and updating the specific machine the template is being rendered for.  If a machine is not present during the render, an
*unknown token* is generated that has the ability to query and create machines.  These are used by the install process to indicate that
the install is finished and that the *local* BootEnv should be used for the next boot and during the discovery process to create
the newly discovered machine.

.. note::
  **.Machine.Path** is particularly useful for ensuring that templates are expanded into a unique file space for
  each machine.  An example of this is per machine kickstart files.  These can be seen in the `assets/bootenvs/ubuntu-16.04.yml <https://github.com/digitalrebar/provision/blob/master/assets/bootenvs/ubuntu-16.04.yml>`_.

With regard to the **.Param** and **.ParamExists** functions, these return the parameter or existence of
the parameter specified by the *key* input.  The parameters are examined from most specific to global.  This means
that the Machine object's profile is checked first, then the list of :ref:`rs_model_profile` associated with the machine,
and finally the global :ref:`rs_model_profile`.  The parameters are stored in a :ref:`rs_model_profile`.

The default :ref:`rs_model_template` and :ref:`rs_model_bootenv` use the following optional (unless marked with an \*)
parameters.

=================================  ================  =================================================================================================================================
Parameter                          Type              Description
=================================  ================  =================================================================================================================================
ntp_servers                        Array of string   The format is an array of IP addresses in dotted quad format.
proxy-servers                      Array of objects  See below, :ref:`rs_arch_proxy_server` as well as some kickstart templates.
operating-system-disk              String            A string to use as the default install drive.  /dev/sda or sda depending upon kickstart or preseed.
access_keys                        Map of strings    The key is the name of the public key.  The value is the public key.  All keys are placed in the .authorized_keys file of root.
provisioner-default-password-hash  String            The password hash for the initial default password, **RocketSkates**
provisioner-default-user           String            The initial user to create for ubuntu/debian installs
dns-domain                         String            DNS Domain to use for this system's install
\*operating-system-license-key     String            Windows Only
\*operating-system-install-flavor  String            Windows Only
=================================  ================  =================================================================================================================================

For some examples of this in use, see :ref:`rs_operation` as well as the example profiles in the assets
:ref:`rs_install` directory.


Sub-templates
_____________

A :ref:`rs_model_template` may contain other templates as described above.  The system comes with some pre-existing
sub-templates to make kickstart and preseed generation easier.  The following templates are available had
have some parameters that drive them.  The required parameters can be applied through profiles or the
:ref:`rs_model_machine` profile.  The templates contain comments with how to use and parameters to set.

.. index::
  pair: SubTemplate; Update DRP BootEnv

Update Digital Rebar Provisioner BootEnv
++++++++++++++++++++++++++++++++++++++++

This sub-template updates the :ref:`rs_model_machine` object's BootEnv to the parameter, **next_boot_env**.  If
**next_boot_env** is not defined, the BootEnv will be set to *local*.  This template uses the **GenerateToken**
function to securely update Digital Rebar Provision.  To use, add the following to the post install section of
the kickstart or net-post-install.sh template.

  ::

    {{ template "update-drp-local.tmpl" . }}

An example :ref:`rs_model_profile` that sets the next BootEnv would be:

  ::

    Name: post-install-bootenv
    Params:
      next_boot_env: cores-live


.. index::
  pair: SubTemplate; Web Proxy

.. _rs_arch_proxy_server:

Web Proxy
+++++++++

This sub-template sets up the environment variables and conditionally the apt repo to use a web proxy.  The
sub-template uses the **proxy-servers** parameter.  The place the template in the post-install section of
the kickstart or the net-post-install.sh script.

  ::

    {{ template "web-proxy.tmpl" . }}


An example :ref:`rs_model_profile` that sets proxies would look like this yaml.

  ::

    Name: proxy-config
    Params:
      proxy-servers:
        - url: http://1.1.1.1:3128
          address: 1.1.1.1
        - url: http://1.1.1.2:3128
          address: 1.1.1.2

.. index::
  pair: SubTemplate; Local Repos

Local Repos
+++++++++++

It is possible to use the exploded ISOs as repositories for post-installation work.  This can be helpful
when missing internet connectivity.  To cause the local repos to replace the public repos, set the *local_repo*
parameter to *true*.  This will force them to be changed.  There is one for ubuntu/debian-based systems,
**ubuntu-drp-only-repos.tmpl** and one for centos/redhat-based systems, **centos-drp-only-repos.tmpl**.
The place the template in the post-install section of the kickstart or the net-post-install.sh script.

  ::

    {{ template "ubuntu-drp-only-repos.tmpl" . }}
    {{ template "centos-drp-only-repos.tmpl" . }}


An example :ref:`rs_model_profile` that sets proxies would look like this yaml.

  ::

    Name: local-repos
    Params:
      local_repo: true

.. index::
  pair: SubTemplate; Set Hostname

.. _rs_st_set_hostname:

Set Hostname
++++++++++++

To set the hostname on the post-installed system, include this template.  It will work for ubuntu and centos-based
systems.  The place the template in the post-install section of the kickstart or the net-post-install.sh script.
The template uses the :ref:`rs_model_machine` built in parameters.

  ::

    {{ template "set-hostname.tmpl" . }}

.. index::
  pair: SubTemplate; Remote Root Access

.. _rs_st_remote_root_access:

Remote Root Access
++++++++++++++++++

This templates installs an authorized_keys file in the root user's home directory.  Multiple keys may be provided.
The template also sets the **/etc/ssh/sshd_config** entry *PermitRootLogin*.  The default setting is
*without-password* (keyed access only), but other values are available, *no*, *yes*, *forced-commands-only*.

  ::

    {{ template "root-remote-access.tmpl" . }}

An example :ref:`rs_model_profile` that sets the keys and *PermitRootLogin* would look like this yaml.

  ::

    Name: root-access
    Params:
      access_keys:
        key1:  ssh-rsa abasbaksl;gksj;glasgjasyyp
        key2:  ssh-rsa moreblablabalkhjlkasjg
      access_ssh_root_mode: yes

.. index::
  pair: SubTemplate; Digital Rebar Integration

Digital Rebar Integration
+++++++++++++++++++++++++

This template will join the newly installed node into Digital Rebar.  This template requires the
use of the :ref:`rs_st_remote_root_access` and :ref:`rs_st_set_hostname` subtemplates as well.  To use, include
these in the kickstart post install section or the net-post-install.sh script.  The **join-to-dr.tmpl** requires
setting the *join_dr* parameter to *true* and credentials to access Digital Rebar.  Digital Rebar's Endpoint is
specified with the *CommandURL* parameter, e.g. https://70.2.3.5.  The username and password used to access
Digital Rebar is specified with *rebar-machine_key*.  This should be the machine key in the rebar-access role
in the system deployment.  It is necessary to make sure that the rebar root access key is added to the **access_keys**
parameter.  To get these last two values, see the commands below.


  ::

    {{ template "set-hostname.tmpl" . }}
    {{ template "root-remote-access.tmpl" . }}
    {{ template "join-to-dr.tmpl" . }}


An example :ref:`rs_model_profile`.

  ::

    # Contains parameters for join-to-dr.tmpl and root-remote-access.tmpl
    Name: dr-int
    Params:
      access_keys:
        key1:  ssh-rsa abasbaksl;gksj;glasgjasyyp
      dr_join: true
      CommandURL: https://70.2.3.5
      rebar-machine_key: machine_install:109asdga;hkljhjha3aksljdga

To get the values for the ssh key and the *rebar-machine_key*, check the *rebar-access* role's attributes or run the
following commands.

.. note:: DR Integration - commands to run on admin node to get values.

  * rebar-machine_key: docker exec -it compose_rebar_api_1 cat /etc/rebar.install.key
  * rebar root access key: docker exec -it compose_rebar_api_1 cat /home/rebar/.ssh/id_rsa.pub


.. _rs_dhcp_models:

DHCP Models
-----------

These models represent things that the DHCP server use and manipulate.

.. index::
  pair: Model; Subnet

.. _rs_model_subnet:

Subnet
======

The Subnet Object defines the configuration of a single subnet for the DHCP server to process.  Multiple subnets are allowed.  The Subnet
can represent a local subnet attached to a local interface (Broadcast Subnet) to the Digital Rebar Provision server or a subnet that is
being forwarded or relayed (Relayed Subnet) to the Digital Rebar Provision server.

The subnet is uniquely identified by its **Name**.  The subnet defines a CIDR-based range with a specific subrange to hand out for
nodes that do NOT have explicit reservations (**ActiveStart** thru **ActiveEnd**).  The subnet also defines the **NextServer** in
the PXE chain.  This is usually an IP associated with Digital Rebar Provision, but if the provisioner is disabled, this can be
any next hop server.  The lease times for both reserved and unreserved clients as specified here (**ReservedLeaseTime** and **ActiveLeaseTime**).
The subnet can also me marked as only working for explicitly reserved nodes (**ReservedOnly**).

The subnet also allows for the specification of DHCP options to be sent to clients.  These can be overridden by :ref:`rs_model_reservation`
specific options.  Some common options are:

========  ====  =================================
Type      #     Description
========  ====  =================================
IP        3     Default Gateway
IP        6     DNS Server
IP        15    Domain Name
String    67    Next Boot File - e.g. lpxelinux.0
========  ====  =================================

golang template expansion also works in these fields.  This can be used to make custom request-based reply options.

For example, this value in the Next Boot File option (67) will return a file based upon what type of machine is booting.  If
the machine supports, iPXE then an iPXE boot image is sent, if the system is marked for legacy bios, then lpxelinux.0 is returned,
otherwise return a 64-bit UEFI network boot loader:

  ::

    {{if (eq (index . 77) "iPXE") }}default.ipxe{{else if (eq (index . 93) "0")}}lpxelinux.0{{else}}bootx64.efi{{end}}


The data element for the template expansion as represented by the '.' above is a map of strings indexed by an integer.  The
integer is the option number from the DHCP request's incoming options.  The IP addresses and other data fields are converted to
a string form (dotted quads or base 10 numerals).

The final elements of a subnet are the **Strategy** and **Pickers** options.  These are described in the :ref:`rs_api` JSON description.
They define how a node should be identified (**Strategy**) and the algorithm for picking addresses (**Pickers**).  The strategy can
only be set to **MAC** currently.  This will use the MAC address of the node as its DHCP identifier.  Others may show up in time.

**Pickers** defines an ordered list of methods to determine the address to hand out.  Currently, this will default to the list:
*hint*, *nextFree*, and *mostExpired*.  The following options are available for the list.

* hint - Use what was provided in the DHCP Offer/Request
* nextFree - Within the subnet's pool of Active IPs, choose the next free making sure to loop over all addresses before reuse.
* mostExpired - If no free address is available, use the most expired address first.
* none - Do NOT hand out anything


.. index::
  pair: Model; Reservation

.. _rs_model_reservation:

Reservation
===========

The Reservation Object defines a mapping between a token and an IP address.  The token is defined by the assigned strategy.  Similar
to :ref:`rs_model_subnet`, the only current strategy is **MAC**.  This will use the MAC address of the incoming requests as the
identity token.  The reservation allows for the optional specification of specific options and a next server that override or
augment the options defined in a subnet.  Because the reservation is an explicit binding of the token to an IP address, the
address can be handed out without the definition of a subnet.  This requires that the reservation have the Netmask Option (Option 1)
specified.  In general, it is a good idea to define a subnet that will cover the reservation with default options and parameters, but
it is not required.

.. index::
  pair: Model; Lease

.. _rs_model_lease:

Lease
=====

The Lease Object defines the ephemeral mapping of a token, as defined by the reservation's or subnet's strategy, and an IP address assigned
by the reservation or pulled form the subnet's pool.  The lease contains the Strategy used for the token and the expiration time.  The
contents of the lease are immutable with the exception of the expiration time.

.. index::
  pair: Model; Interface

.. _rs_model_interface:

Interface
=========

The Interface Object is a read-only object that is used to identify local interfaces and their addresses on the
Digital Rebar Provision server.  This is useful for determining what subnets to create and with what address ranges.
The :ref:`rs_ui_subnets` part of the :ref:`rs_ui` uses this to populate possible subnets to create.


.. _rs_additional_models:

Additional Models
-----------------

These models control additional parts and actions of the system.

.. index::
  pair: Model; User

.. _rs_model_user:

User
====

The User Object controls access to the system.  The User object contains a name and a password hash for validating access.  Additionally,
the User :ref:`rs_api` can be used to generate time-based, function restricted tokens for use in :ref:`rs_api` calls.  The
:ref:`rs_model_template` provides a helper function to generate these for restricted machine access in the discovery and post-install
process.

The User Object is usually created with an unset password.  This allows for the User have no access but still access the system
through constructed tokens.  The :ref:`rs_cli` has commands to set the password for a user.

More on access tokens, user creation,  and an control in :ref:`rs_operation`.


.. index::
  pair: Model; Prefs

.. _rs_model_prefs:

Prefs
=====

Most configuration is handle through the :ref:`rs_model_profile` system, but there are a few modifiable
options that can be changed over time in the server (outside of command line flags).  These are preferences.  The preferences are
key value pairs where both the key and the value are strings.  The use internally may be an integer, but the specification through
the :ref:`rs_api` is by string.

=================== ======= ==================================================================================================================================================================================
Pref                Type    Description
=================== ======= ==================================================================================================================================================================================
defaultBootEnv      string  This is a valid :ref:`rs_model_bootenv` the is assign to a :ref:`rs_model_machine` if the machine does not have a bootenv specified.  The default is **sledgehammer**.
unknownBootEnv      string  This is the :ref:`rs_model_bootenv` used when a boot request is serviced by an unknown machine.  The BootEnv must have **OnlyUnknown** set to true.  The default is **ignore**.
unknownTokenTimeout integer The amount of time in seconds that the token generated by **GenerateToken** is valid for unknown machines.  The default is 600 seconds.
knownTokenTimeout   integer The amount of time in seconds that the token generated by **GenerateToken** is valid for known machines.  The default is 3600 seconds.
debugRenderer       integer The debug level of the renderer system.  0 = off, 1 = info, 2 = debug
debugDhcp           integer The debug level of the DHCP system.  0 = off, 1 = info, 2 = debug
debugBootEnv        integer The debug level of the BootEnv system.  0 = off, 1 = info, 2 = debug
=================== ======= ==================================================================================================================================================================================

.. _rs_special_objects:

Special Objects
---------------

These are not objects in the system but represent files and directories in the server space.

.. index::
  pair: Model; Files

.. _rs_model_file:

Files
=====

File server has a managed filesystem space.  The :ref:`rs_api` defines methods to upload, destroy, and get these files outside of the
normal TFTP and HTTP path.  The TFTP and HTTP access paths are read-only.  The only way to modify this space is through the :ref:`rs_api`
or direct filesystem access underneath Digital Rebar Provision.  The filesystem space defaults to */var/lib/tftpboot*, but can be overridden
by the command line flag *--file-root*, e.g. *--file-root=`pwd`/drp-data* when using *--isolated* on install.  These directories can be
directly manipulated by administrators for faster loading times.

This space is also used by the :ref:`rs_model_bootenv` import process when "exploding" an ISO for use by :ref:`rs_model_machine`.

.. note:: Templates are **NOT** rendered to the file system.  They are in-memory generated on the fly content.

.. index::
  pair: Model; Isos

.. _rs_model_iso:

Isos
====

The ISO directory in the file server space is managed specially by the ISO :ref:`rs_api`.  The API handles upload and destroy
functionality.  The API also handles notification of the :ref:`rs_model_bootenv` system to "explode" ISOs that are needed by :ref:`rs_model_bootenv` and marking
the :ref:`rs_model_bootenv` as available.

ISOs can be directly placed into the **isos** directory in the file root, but the using :ref:`rs_model_bootenv` needs to be modified or deleted and
re-added to force the ISO to be exploded for use.


