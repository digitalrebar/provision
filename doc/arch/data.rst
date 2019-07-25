.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Data Architecture

.. _rs_data_architecture:

Data Architecture
=================

Digital Rebar Provision uses a fairly simple data model.  There are 4
main models for the provisioner server and 4 models for the DHCP
server.  Each model has a corresponding API in the :ref:`rs_api`.

The models define elements of the system.  The API provides basic CRUD
(create, read, update, and delete) operations as well as some
additional actions for manipulating the state of the system.  The
:ref:`rs_api` contains that definitions of the actual structures and
methods on those objects.  Additionally, the :ref:`rs_operation` will
describe common actions to use and do with these models and how to
build them.  The :ref:`rs_cli` describes the Command Line manipulators
for the model.

This section will describe its use and role in the system.

.. _rs_provisioner_models:

Provisioner Models
------------------

These models represent things that the provisioner server use and manipulate.


.. index::
  pair: Model; Machine

.. _rs_model_machine:

Machine
~~~~~~~

The Machine object defines a machine that is being provisioned.  The
Machine is represented by a unique **UUID**.  The UUID is immutable
after machine creation.  The machine's primary purpose is to map an
incoming IP address to a :ref:`rs_model_bootenv`.  The
:ref:`rs_model_bootenv` provides a set of rendered
:ref:`rs_model_template` that will can be used to boot the machine.
The machine provides parameters to the :ref:`rs_model_template`.  The
Machine provides configuration to the renderer in the form of
parameters and fields.  Also, each :ref:`rs_model_machine` must have a
:ref:`rs_model_bootenv` to boot from.  If a machine is created without
a :ref:`rs_model_bootenv` specified, the system will assign the one
specified by the value of :ref:`rs_model_prefs` *defaultBootEnv*.

The **Name** field should contain the FQDN of the node.

The Machine object contains an **Error** field that represents errors
encountered while operating on the machine.  In general, these are
errors pertaining to rendering the :ref:`rs_model_bootenv`.

The Machine parameters are defined as a special
:ref:`rs_model_profile` on the Machine.  The profile stores a
dictionary of string keys to arbitrary objects.  These could be
strings, booleans, numbers, arrays, or objects representing similarly
defined dictionaries.  The machine parameters are available to
templates for expansion in them.

Additionally, the machine maintains an ordered list of profiles that
are searched and then finally the **global profile**.  See
:ref:`rs_model_profile` and :ref:`rs_model_template` for more
information.

.. note:: When updating the Params part of the embedded Profile in the
          :ref:`rs_model_machine` object, using the **PUT** method
          will replace the Params map with the map from the input
          object.  The **PATCH** will merge the Params map in the
          input with the existing Params map in the current
          :ref:`rs_model_machine` object.  The **POST** method on the
          params subaction will replace the map with the input
          version.

.. index::
  pair: Model; Param

.. _rs_model_param:

Param
~~~~~

The Param Object is the lowest level building block.  It is a simple
key / value pair.  Each Param is a bounded type parameter, and type
definition is enforced.  The following types of parameters may be
used:

========================== ========================================================================
type                       description
========================== ========================================================================
integer                    A numerical value (eg "12" or "-3444")
boolean                    True or False ('true' or 'false')
string                     Textual string (eg "this is a string!")
array                      A series of elements of the same type
map                        a higher-order function that applies a given function to each element of a list, returning a list of results in the same order
========================== ========================================================================

.. index::
  pair: Model; Profile

.. _rs_model_profile:

Profile
~~~~~~~

The Profile Object defines a set of key / value pairs (or parameters).
All of these may be manipulated by the :ref:`rs_api`.  The key space
is a free form string and the value is an arbitrary data blob
specified by JSON through the :ref:`rs_api`.  The common parameters
defined in :ref:`rs_model_template` can be set on these objects.  The
system maintains a **global** profile for setting system wide
parameters.  They are the lowest level of precedence.

The profiles are free form dictionaries and default empty.  Any
key/value pair can be added and referenced.

Other profiles may be created to group parameters together to apply to
sets of machines.  The machine's profile list allows the administrator
to specify an ordered set of profiles that apply to that machine as
well.  Additionally, the system maintains a special profile for each
machine to store custom parameters specific to that machine.  This
profile is embedded in the :ref:`rs_model_machine` object.

When the system needs to render a template parameter, the machine's
specific profile is checked, then the order list of profiles stored in
the Machine Object are checked, and finally the **global** profile is
checked.  The key and its value are used if found in template
rendering.

.. note:: When updating the Params part of the
          :ref:`rs_model_profile`, using the **PUT** method will
          replace the Params map with the map from the input object.
          The **PATCH** method will merge the Params map in the input
          with the existing Params map in the current
          :ref:`rs_model_profile` object.  The **POST** method on the
          params subaction will replace the map with the input
          version.


.. index::
  pair: Model; BootEnv

.. _rs_model_bootenv:

BootEnv
~~~~~~~

The BootEnv object defines an environment to boot a machine.  It has
two main components an OS information section and a templates list.
The OS information section defines what makes up the installation base
for this bootenv.  It defines the install ISO, a URL to get the ISO,
and SHA256 checksum to validate the image.  These are used to provide
the basic install image, kernel, and base packages for the bootenv.

The other primary section is a set of templates that represent files
in the file server's file space that can served via HTTP or TFTP.  The
templates can be in-line in the BootEnv object or reference a
:ref:`rs_model_template`.  The templates are specified as a list of
paths in the filesystem and either an ID of a :ref:`rs_model_template`
or inline content.  The path field of the template information can use
the same template expansion that is used in the template.  See
:ref:`rs_model_template` for more information.

Additionally, the BootEnv defines required and optional parameters.
The required parameters validated at render time to be present or an
error is generated.  These parameters can be met by the parameters on
the machine, the profiles in machine's profiles list, or from the
global :ref:`rs_model_profile`.

BootEnvs can be marked **OnlyUnknown**.  This tells the rest of the
system that this BootEnv is not for specific machines.  It is a
general BootEnv.  For example, *discovery* and *ignore* are
**OnlyUnknown**.  *discovery* is used to discover unknown machines and
add them to Digital Rebar Provision.  *ignore* is a special bootenv
that tells machines to boot their local disk.  These BootEnvs populate
the pxelinux.0, ipxe, and elilo default fallthrough files.  These are
different than their counterpart BootEnvs, *sledgehammer* and *local*
which are machine specific BootEnvs that populate configuration files
that are specific to a single machine.  A machine boots *local*; an
unknown machine boots *ignore*.  There can only be one **OnlyUnknown**
BootEnv active at a time.  This is specified by the
:ref:`rs_model_prefs` *unknownBootEnv*.

.. index::
  pair: Model; Template

.. _rs_model_template:

Template
~~~~~~~~

The Template object defines a templated content that can be referenced
by its ID.  The content of the template (or in-line template in a
:ref:`rs_model_bootenv`) is a `golang text/template
<https://golang.org/pkg/text/template/#hdr-Actions>`_ string.  The
template has a set of special expansions.  The normal expansion syntax
is:

  ::

    {{ .Machine.Name }}

This would expand to the machine's **Name** field.  There are helpers
for the parameter spaces, the :ref:`rs_model_bootenv` object, and some
miscellaneous functions.  Additionally, the normal `golang
text/template <https://golang.org/pkg/text/template/#hdr-Actions>`_
functions are available as well.  Things like **range**, **len**, and
comparators are available as well.  **template** inclusion is
supported by the following syntax:

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
.ParseURL <segment> <url>      Parse the specified URL and return the segment requested. Supported segments can be one of *scheme* (eg "https"), *host* (eg "drp.example.com:8092"), or *path* (eg "/api/v3/machines").  *host* does not separate name and port.
.ParamExists <key>             Returns true if the specified key is a valid parameter available for this rendering.
.Param <key>                   Returns the structure for the specified key for this rendering.
.Repos <tag>, <tag>,...        Returns Repos (as defined by the package-repositories param currently in scope) with the matching tags.
.MachineRepos                  Returns all Repos that have the **OS** of the Machine defined in their os section.
.InstallRepos                  Returns exactly one Repo from the list chosen by MachineRepos that has the installSource bit set, and at most one Repo from the MachineRepos that has the securitySource bit set.
template <string> .            Includes the template specified by the string.  String can be a variable and note that template does NOT have a dot (.) in front.
============================== =================================================================================================================================================================================================

**GenerateToken** is very special.  This generates either a *known
token* or an *unknown token* for use by the template to update objects
in Digital Rebar Provision.  The tokens are valid for a limited time
as defined by the **knownTokenTimeout** and **unknownTokenTimeout**
:ref:`rs_model_prefs` respectively.  The tokens are also restricted to
the function the can perform.  The *known token* is limited to only
reading and updating the specific machine the template is being
rendered for.  If a machine is not present during the render, an
*unknown token* is generated that has the ability to query and create
machines.  These are used by the install process to indicate that the
install is finished and that the *local* BootEnv should be used for
the next boot and during the discovery process to create the newly
discovered machine.

.. note:: **.Machine.Path** is particularly useful for ensuring that
  templates are expanded into a unique file space for each machine.
  An example of this is per machine kickstart files.  These can be
  seen in the `assets/bootenvs/ubuntu-16.04.yml
  <https://github.com/digitalrebar/provision/blob/master/assets/bootenvs/ubuntu-16.04.yml>`_.

With regard to the **.Param** and **.ParamExists** functions, these
return the parameter or existence of the parameter specified by the
*key* input.  The parameters are examined from most specific to
global.  This means that the Machine object's profile is checked
first, then the list of :ref:`rs_model_profile` associated with the
machine, and finally the global :ref:`rs_model_profile`.  The
parameters are stored in a :ref:`rs_model_profile`.

The default :ref:`rs_model_template` and :ref:`rs_model_bootenv` use
the following optional (unless marked with an \*) parameters.

=================================  ================  =================================================================================================================================
Parameter                          Type              Description
=================================  ================  =================================================================================================================================
ntp_servers                        Array of string   The format is an array of IP addresses in dotted quad format.
proxy-servers                      Array of objects  See below, :ref:`rs_arch_proxy_server` as well as some kickstart templates.
operating-system-disk              String            A string to use as the default install drive.  /dev/sda or sda depending upon kickstart or preseed.
access-keys                        Map of strings    The key is the name of the public key.  The value is the public key.  All keys are placed in the .authorized_keys file of root.
provisioner-default-password-hash  String            The password hash for the initial default password, **RocketSkates**
provisioner-default-user           String            The initial user to create for ubuntu/debian installs
dns-domain                         String            DNS Domain to use for this system's install
\*operating-system-license-key     String            Windows Only
\*operating-system-install-flavor  String            Windows Only
=================================  ================  =================================================================================================================================

For some examples of this in use, see :ref:`rs_operation` as well as
the example profiles in the assets :ref:`rs_install` directory.


Sub-templates
_____________

A :ref:`rs_model_template` may contain other templates as described
above.  The system comes with some pre-existing sub-templates to make
kickstart and preseed generation easier.  The following templates are
available had have some parameters that drive them.  The required
parameters can be applied through profiles or the
:ref:`rs_model_machine` profile.  The templates contain comments with
how to use and parameters to set.

.. index::
  pair: SubTemplate; Update DRP BootEnv

Update Digital Rebar Provisioner BootEnv
++++++++++++++++++++++++++++++++++++++++

This sub-template updates the :ref:`rs_model_machine` object's BootEnv
to the parameter, **next_boot_env**.  If **next_boot_env** is not
defined, the BootEnv will be set to *local*.  This template uses the
**GenerateToken** function to securely update Digital Rebar Provision.
To use, add the following to the post install section of the kickstart
or net-post-install.sh template.

  ::

    {{ template "update-drp-local.tmpl" . }}

An example :ref:`rs_model_profile` that sets the next BootEnv would
be:

  ::

    Name: post-install-bootenv
    Params:
      next_boot_env: cores-live


.. index::
  pair: SubTemplate; Web Proxy

.. _rs_arch_proxy_server:

Web Proxy
+++++++++

This sub-template sets up the environment variables and conditionally
the apt repo to use a web proxy.  The sub-template uses the
**proxy-servers** parameter.  The place the template in the
post-install section of the kickstart or the net-post-install.sh
script.

  ::

    {{ template "web-proxy.tmpl" . }}


An example :ref:`rs_model_profile` that sets proxies would look like
this yaml.

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

**This section is deprecated, it is being replaced by the more general
package-repositories functionality**

It is possible to use the exploded ISOs as repositories for
post-installation work.  This can be helpful when missing internet
connectivity.  To cause the local repos to replace the public repos,
set the *local_repo* parameter to *true*.  This will force them to be
changed.  There is one for ubuntu/debian-based systems,
**ubuntu-drp-only-repos.tmpl** and one for centos/redhat-based
systems, **centos-drp-only-repos.tmpl**.  The place the template in
the post-install section of the kickstart or the net-post-install.sh
script.

  ::

    {{ template "ubuntu-drp-only-repos.tmpl" . }}
    {{ template "centos-drp-only-repos.tmpl" . }}


An example :ref:`rs_model_profile` that sets proxies would look like this yaml.

  ::

    Name: local-repos
    Params:
      local-repo: true

.. index::
  pair: SubTemplate; Package Repositories

Package Repositories
++++++++++++++++++++

As an alternative to rolling your own support for local annd remote
package repositrory management, you can write your templates to use
our package repository support.  This support consists of three parts:

1. Support in the template rendering engine for a parameter named
   "package-repositories", which contains a list of package
   repositories that are available for the various Linux distros we
   support.
2. The .Repos, .MachineRepos, and .InstallRepos functions that are
   available at template expansion time.  These return a list of Repo
   objects, and re described in more detail in the Template section.
3. The .Install and .Lines functions available on each Repo object.

The package-repositories Param
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

The special "package-repositories" parameter must be present in-scope
of the current Machine in order for .Repos, .MachineRepos, and
.InstallRepos to operate correctly -- that is, it must be present
either in the global profile, a profile attached to the machine's
current Stage,a profile attached to a machine, or directly on the
machine as a machine parameter.

A commented example of a "package-repositories" parameter in YAML format:

  ::

    - tag: "centos-7-install" # Every repository needs a unique tag.
      # A repository can be used by multiple operating systems.
      # The usual example of this is the EPEL repository, which
      # can be used by all of the RHEL variants of a given generation.
      os:
        - "centos-7"
      # If installSource is true, then the URL points directly
      # to the location we should use for all OS install purposes
      # save for fetching kernel/initrd pairs from (for now, we will
      # still assume that they will live on the DRP server).
      # When installSounrce is true, the os field must contain a single
      # entry that is an exact match for the bootenv's OS.Name field.
      installSource: true
      # For redhat-ish distros when installSource is true,
      # this URL must contain distro, component, and arch components,
      # and as such they do not need to be further specified.
      url: "http://mirrors.kernel.org/centos/7/os/x86_64"
    - tag: "centos-7-everything"
      # Since installSource is not true here,
      # we can define several package sources at once by
      # providing a distribution and a components section,
      # and having the URL point at the top-level directory
      # where everything is housed.
      # DRP knows how to expand repo definitions for CentOS and
      # ScientificLinux provided that they follow the standard
      # mirror directory layout for each distro.
      os:
        - centos-7
      url: "http://mirrors.kernel.org/centos"
      distribution: "7"
      components:
        - atomic
        - centosplus
        - configmanagement
        - cr
        - dotnet
        - extras
        - fasttrack
        - os
        - rt
        - sclo
        - updates
    - tag: "debian-9-install"
      os:
        - "debian-9"
      installSource: true
      # Debian URLs always follow the same rules, no matter
      # whether the OS install flag is set.  As such,
      # you must always also specify the distribution and
      # at least the main component, although you can also
      # specify other components.
      url: "http://mirrors.kernel.org/debian"
      distribution: stretch
      # If the location of the remote kernel and initrd files cannot be found
      # at the location you would get by appending url and the kernel/initd
      # filenames from the BootEnv, you need to use the bootloc field to
      # override where dr-provision should try to get them from.
      # Kernels and initrds must be located directly at this path.
      bootloc: "http://mirrors.kernel.org/debian/dists/stretch/main/installer-amd64/current/images/netboot/debian-installer/amd64/"
      components:
        - main
        - contrib
        - non-free
    - tag: "debian-9-updates"
      os:
        - "debian-9"
      url: "http://mirrors.kernel.org/debian"
      distribution: stretch-updates
      components:
        - main
        - contrib
        - non-free
    - tag: "debian-9-backports"
      os:
        - "debian-9"
      url: "http://mirrors.kernel.org/debian"
      distribution: stretch-backports
      components:
        - main
        - contrib
        - non-free
    - tag: "debian-9-security"
      os:
        - "debian-9"
      url: "http://security.debian.org/debian-security/"
      securitySource: true
      distribution: stretch/updates
      components:
        - contrib
        - main
        - non-free

The default package-repositories param in drp-community-content
contains working examples for every boot environment supported by
drp-community-content.

Repo Object
^^^^^^^^^^^

As mentioned above, the template-level .Repos, .MachineRepos, and
.InstallRepos return a list of Repo objects that can be used for
further template expansion.  The Repo object contains its own fields
and functions that can be used for template expansion:

===================    ===========
Expansion              Description
===================    ===========
.Tag                   The tag that uniquely identifies one repository definition.  The template-level .Repos function takes a list of tags and returns repos that exactly match them.
.OS                    A list of operating systems (in distro-release format) that this repository supports. The template-level .MachineRepos function matches this field against the current Machine.OS field to determine which templates are applicable to a Machine.
.URL                   The URL to the top of the repository in question.  For yum-style repos, it can either point directly to a specific repository (in which case .Distribution and .Components must not be present), or point to a location that contains an appropriately mirrored repo tree for the OS in question (in which case it cannot be used as an InstallSource or a SecuritySource, and .Distribution and .Components must be set)  For apt-style repos, it must point to the top level of the repository (the level that has "dists" and "pool" as subdirectories), and .Distribution and .Components must always be defined.
.PackageType           An optional field that can be used to determine what kind of packages the repository returns.  It is normally autodetected based on the operating system the repo is being used in.
.RepoType              The type repository this is.  It is optional, and is normally inferred based on the operating system the repo is being used in.
.InstallSource         A boolean value that determines whether this repository should be used as a package source during OS installation. You should have at most one of these per OS install you wish to support.
.SecuritySource        A boolean value that determines whether this repository should be used as a source of security updates that should be applied during an OS install.
.Distribution          A string that corresponds to the OS release version or codename.  This must be present for apt-style repos.
.Components            A list of strings that map to any sub-repositories available as part of this repository.  Examples are "main","contrib", and "non-free" for apt-based repos.
.R                     A helper function that refers back to the top-level template rendering context.
.JoinedComponents      A helper function that joins the .Components list into a space-seperated string.
.UrlFor <component>    A helper function that returns an appropriately formatted URL for the passed Component.
.Install               A helper function that returns the Repo in a format suitable for inclusion in an unattented OS installation file (kickstart, preseed, etc.)  The format returned is currently hardcoded depending on the OS type of the Machine.  That restriction will be lifted in future versions of dr-provision.
.Lines                 A helper function that returns the Repo an a format suitable for direct inclusion into a repo definition file (sources.list, /etc/yum.repos.d/.repo, etc).  The format returned is currently hardcoded based on the OS type of the Machine.  That restriction will be lifted in future versions of dr-provision.
===================    ===========

Expanding Package Repositories
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

To expand the repos suitable for OS installation, use::

    {{range $repo := .InstallRepos}}{{$repo.Install}}{{end}}

To expand the repos suitable for post-install package management, use::

    {{range $repo := .MachineRepos}}{{$repo.Lines}}{{end}}


.. index::
  pair: SubTemplate; Set Hostname

.. _rs_st_set_hostname:

Set Hostname
++++++++++++

To set the hostname on the post-installed system, include this
template.  It will work for ubuntu and centos-based systems.  The
place the template in the post-install section of the kickstart or the
net-post-install.sh script.  The template uses the
:ref:`rs_model_machine` built in parameters.

  ::

    {{ template "set-hostname.tmpl" . }}

.. index::
  pair: SubTemplate; Remote Root Access

.. _rs_st_remote_root_access:

Remote Root Access
++++++++++++++++++

This templates installs an authorized_keys file in the root user's
home directory.  Multiple keys may be provided.  The template also
sets the **/etc/ssh/sshd_config** entry *PermitRootLogin*.  The
default setting is *without-password* (keyed access only), but other
values are available, *no*, *yes*, *forced-commands-only*.

  ::

    {{ template "root-remote-access.tmpl" . }}

An example :ref:`rs_model_profile` that sets the keys and *PermitRootLogin* would look like this yaml.

  ::

    Name: root-access
    Params:
      access-keys:
        key1:  ssh-rsa abasbaksl;gksj;glasgjasyyp
        key2:  ssh-rsa moreblablabalkhjlkasjg
      access_ssh_root_mode: yes

.. index::
  pair: SubTemplate; Digital Rebar Integration

Digital Rebar Integration
+++++++++++++++++++++++++

This template will join the newly installed node into Digital Rebar.
This template requires the use of the :ref:`rs_st_remote_root_access`
and :ref:`rs_st_set_hostname` subtemplates as well.  To use, include
these in the kickstart post install section or the net-post-install.sh
script.  The **join-to-dr.tmpl** requires setting the *join_dr*
parameter to *true* and credentials to access Digital Rebar.  Digital
Rebar's Endpoint is specified with the *CommandURL* parameter,
e.g. https://70.2.3.5.  The username and password used to access
Digital Rebar is specified with *rebar-machine_key*.  This should be
the machine key in the rebar-access role in the system deployment.  It
is necessary to make sure that the rebar root access key is added to
the **access-keys** parameter.  To get these last two values, see the
commands below.


  ::

    {{ template "set-hostname.tmpl" . }}
    {{ template "root-remote-access.tmpl" . }}
    {{ template "join-to-dr.tmpl" . }}


An example :ref:`rs_model_profile`.

  ::

    # Contains parameters for join-to-dr.tmpl and root-remote-access.tmpl
    Name: dr-int
    Params:
      access-keys:
        key1:  ssh-rsa abasbaksl;gksj;glasgjasyyp
      dr_join: true
      CommandURL: https://70.2.3.5
      rebar-machine_key: machine_install:109asdga;hkljhjha3aksljdga

To get the values for the ssh key and the *rebar-machine_key*, check
the *rebar-access* role's attributes or run the following commands.

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
~~~~~~

The Subnet Object defines the configuration of a single subnet for the
DHCP server to process.  Multiple subnets are allowed.  The Subnet can
represent a local subnet attached to a local interface (Broadcast
Subnet) to the Digital Rebar Provision server or a subnet that is
being forwarded or relayed (Relayed Subnet) to the Digital Rebar
Provision server.

The subnet is uniquely identified by its **Name**.  The subnet defines
a CIDR-based range with a specific subrange to hand out for nodes that
do NOT have explicit reservations (**ActiveStart** thru
**ActiveEnd**).  The subnet also defines the **NextServer** in the PXE
chain.  This is usually an IP associated with Digital Rebar Provision,
but if the provisioner is disabled, this can be any next hop server.
The lease times for both reserved and unreserved clients as specified
here (**ReservedLeaseTime** and **ActiveLeaseTime**).  The subnet can
also me marked as only working for explicitly reserved nodes
(**ReservedOnly**).

The subnet also allows for the specification of DHCP options to be
sent to clients.  These can be overridden by
:ref:`rs_model_reservation` specific options.  Some common options
are:

========  ====  =================================
Type      #     Description
========  ====  =================================
IP        3     Default Gateway
IP        6     DNS Server
IP        15    Domain Name
String    67    Next Boot File - e.g. ipxe.pxe
========  ====  =================================

golang template expansion also works in these fields.  This can be
used to make custom request-based reply options.

For example, this value in the Next Boot File option (67) will return
a file based upon what type of machine is booting.  If the machine
supports, iPXE then an iPXE boot image is sent, if the system is
marked for legacy bios, then ipxe.pxe is returned, otherwise return
a 64-bit UEFI iPXE boot loader:

  ::

    {{if (eq (index . 77) "iPXE") }}default.ipxe{{else if (eq (index . 93) "0")}}ipxe.pxe{{else}}ipxe.efi{{end}}


NOTE: Option 67 is optional.  When using DRP as the DHCP server,
it will generate a bootfile like the above template expansion.

The data element for the template expansion as represented by the '.'
above is a map of strings indexed by an integer.  The integer is the
option number from the DHCP request's incoming options.  The IP
addresses and other data fields are converted to a string form (dotted
quads or base 10 numerals).

The final elements of a subnet are the **Strategy** and **Pickers**
options.  These are described in the :ref:`rs_api` JSON description.
They define how a node should be identified (**Strategy**) and the
algorithm for picking addresses (**Pickers**).  The strategy can only
be set to **MAC** currently.  This will use the MAC address of the
node as its DHCP identifier.  Others may show up in time.

.. _rs_model_pickers:

Pickers
~~~~~~~

**Pickers** defines an ordered list of methods to determine the
address to hand out.  Currently, this will default to the list:
*hint*, *nextFree*, and *mostExpired*.  The following options are
available for the list.

* **hint** - which will try to reuse the address that the DHCP packet is
  requesting, if it has one.  If the request does not have a requested
  address, "hint" will fall through to the next strategy. Otherwise,
  it will refuse to try ant reamining strategies whether or not it can
  satisfy the request.  This should force the client to fall back to
  DHCPDISCOVER with no requsted IP address. "hint" will reuse expired
  leases and unexpired leases that match on the requested address,
  strategy, and token.
* **nextFree** - Within the subnet's pool of Active IPs, choose the next
  free making sure to loop over all addresses before reuse.  It will
  fall through to the next strategy if it cannot find a free IP.
  "nextFree" only considers addresses that do not have a lease,
  whether or not the lease is expired.
* **mostExpired** - If no free address is available, use the most expired
  address first.
* **none** - Do NOT hand out an address and refuse to try any remaining
  strategies

All of the address allocation strategies do not consider any addresses
that are reserved, as lease creation will be handled by the
reservation instead.


.. index::
  pair: Model; Reservation

.. _rs_model_reservation:

Reservation
~~~~~~~~~~~

The Reservation Object defines a mapping between a token and an IP
address.  The token is defined by the assigned strategy.  Similar to
:ref:`rs_model_subnet`, the only current strategy is **MAC**.  This
will use the MAC address of the incoming requests as the identity
token.  The reservation allows for the optional specification of
specific options and a next server that override or augment the
options defined in a subnet.  Because the reservation is an explicit
binding of the token to an IP address, the address can be handed out
without the definition of a subnet.  This requires that the
reservation have the Netmask Option (Option 1) specified.  In general,
it is a good idea to define a subnet that will cover the reservation
with default options and parameters, but it is not required.

.. index::
  pair: Model; Lease

.. _rs_model_lease:

Lease
~~~~~

The Lease Object defines the ephemeral mapping of a token, as defined
by the reservation's or subnet's strategy, and an IP address assigned
by the reservation or pulled form the subnet's pool.  The lease
contains the Strategy used for the token and the expiration time.  The
contents of the lease are immutable with the exception of the
expiration time.

.. index::
  pair: Model; Interface

.. _rs_model_interface:

Interface
~~~~~~~~~

The Interface Object is a read-only object that is used to identify
local interfaces and their addresses on the Digital Rebar Provision
server.  This is useful for determining what subnets to create and
with what address ranges.  The :ref:`rs_ui_subnets` part of the
:ref:`rs_ui` uses this to populate possible subnets to create.


.. _rs_additional_models:

Additional Models
-----------------

These models control additional parts and actions of the system.

.. index::
  pair: Model; User

.. _rs_model_user:

User
~~~~

The User Object controls access to the system.  The User object
contains a name and a password hash for validating access.
Additionally, the User :ref:`rs_api` can be used to generate
time-based, function restricted tokens for use in :ref:`rs_api` calls.
The :ref:`rs_model_template` provides a helper function to generate
these for restricted machine access in the discovery and post-install
process.

The User Object is usually created with an unset password.  This
allows for the User have no access but still access the system through
constructed tokens.  The :ref:`rs_cli` has commands to set the
password for a user.

More on access tokens, user creation, and an control in
:ref:`rs_operation`.


.. index::
  pair: Model; Prefs

.. _rs_model_prefs:

Prefs
~~~~~

Most configuration is handle through the :ref:`rs_model_profile`
system, but there are a few modifiable options that can be changed
over time in the server (outside of command line flags).  These are
preferences.  The preferences are key value pairs where both the key
and the value are strings.  The use internally may be an integer, but
the specification through the :ref:`rs_api` is by string.

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

These are not objects in the system but represent files and
directories in the server space.

.. index::
  pair: Model; Files

.. _rs_model_file:

Files
~~~~~

File server has a managed filesystem space.  The :ref:`rs_api` defines
methods to upload, destroy, and get these files outside of the normal
TFTP and HTTP path.  The TFTP and HTTP access paths are read-only.
The only way to modify this space is through the :ref:`rs_api` or
direct filesystem access underneath Digital Rebar Provision.  The
filesystem space defaults to */var/lib/tftpboot*, but can be
overridden by the command line flag *--file-root*,
e.g. *--file-root=`pwd`/drp-data* when using *--isolated* on install.
These directories can be directly manipulated by administrators for
faster loading times.

This space is also used by the :ref:`rs_model_bootenv` import process
when "exploding" an ISO for use by :ref:`rs_model_machine`.

.. note:: Templates are **NOT** rendered to the file system.  They are
          in-memory generated on the fly content.

.. index::
  pair: Model; Isos

.. _rs_model_iso:

Isos
~~~~

The ISO directory in the file server space is managed specially by the
ISO :ref:`rs_api`.  The API handles upload and destroy functionality.
The API also handles notification of the :ref:`rs_model_bootenv`
system to "explode" ISOs that are needed by :ref:`rs_model_bootenv`
and marking the :ref:`rs_model_bootenv` as available.

ISOs can be directly placed into the **isos** directory in the file
root, but the using :ref:`rs_model_bootenv` needs to be modified or
deleted and re-added to force the ISO to be exploded for use.

