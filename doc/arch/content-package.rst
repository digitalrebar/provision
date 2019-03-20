.. Copyright (c) 2019 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
 pair: Digital Rebar Provision; Content Packages

Content Packages
----------------

Content packages (also known as content bundles, or content packs) are
the primary way of adding new functionality to dr-provision.  Content
packages consist of a collection of objects (tasks, bootenvs,
templates, stages, workflows, etc.), along with important metadata
about the content package.  An example of what a content package looks
like in YAML is:

.. code-block:: yaml

    ---
    meta:
      Author: ""
      CodeSource: ""
      Color: ""
      Copyright: ""
      Description: Default objects that must be present
      DisplayName: ""
      DocUrl: ""
      Documentation: ""
      Icon: ""
      License: ""
      Name: BasicStore
      Order: ""
      Overwritable: true
      Prerequisites: ""
      RequiredFeatures: ""
      Source: ""
      Tags: ""
      Type: basic
      Version: 3.12.0
      Writable: false
    sections:
      bootenvs:
        ignore:
          Available: false
          BootParams: ""
          Bundle: BasicStore
          Description: The boot environment you should use to have unknown machines boot
            off their local hard drive
          Documentation: ""
          Endpoint: ""
          Errors: []
          Initrds: []
          Kernel: ""
          Meta:
            color: green
            feature-flags: change-stage-v2
            icon: circle thin
            title: Digital Rebar Provision
          Name: ignore
          OS:
            Codename: ""
            Family: ""
            IsoFile: ""
            IsoSha256: ""
            IsoUrl: ""
            Name: ignore
            SupportedArchitectures: {}
            Version: ""
          OnlyUnknown: true
          OptionalParams: []
          ReadOnly: false
          RequiredParams: []
          Templates:
          - Contents: |
              DEFAULT local
              PROMPT 0
              TIMEOUT 10
              LABEL local
              {{.Param "pxelinux-local-boot"}}
            ID: ""
            Meta: null
            Name: pxelinux
            Path: pxelinux.cfg/default
          - Contents: |
              #!ipxe
              chain {{.ProvisionerURL}}/${netX/mac}.ipxe && exit || goto chainip
              :chainip
              chain tftp://{{.ProvisionerAddress}}/${netX/ip}.ipxe || exit
            ID: ""
            Meta: null
            Name: ipxe
            Path: default.ipxe
          Validated: false
        local:
          Available: false
          BootParams: ""
          Bundle: BasicStore
          Description: The boot environment you should use to have known machines boot
            off their local hard drive
          Documentation: ""
          Endpoint: ""
          Errors: []
          Initrds: []
          Kernel: ""
          Meta:
            color: green
            feature-flags: change-stage-v2
            icon: radio
            title: Digital Rebar Provision
          Name: local
          OS:
            Codename: ""
            Family: ""
            IsoFile: ""
            IsoSha256: ""
            IsoUrl: ""
            Name: local
            SupportedArchitectures: {}
            Version: ""
          OnlyUnknown: false
          OptionalParams: []
          ReadOnly: false
          RequiredParams: []
          Templates:
          - Contents: |
              DEFAULT local
              PROMPT 0
              TIMEOUT 10
              LABEL local
              {{.Param "pxelinux-local-boot"}}
            ID: ""
            Meta: null
            Name: pxelinux
            Path: pxelinux.cfg/{{.Machine.HexAddress}}
          - Contents: |
              #!ipxe
              exit
            ID: ""
            Meta: null
            Name: ipxe
            Path: '{{.Machine.Address}}.ipxe'
          - Contents: |
              DEFAULT local
              PROMPT 0
              TIMEOUT 10
              LABEL local
              {{.Param "pxelinux-local-boot"}}
            ID: ""
            Meta: null
            Name: pxelinux-mac
            Path: pxelinux.cfg/{{.Machine.MacAddr "pxelinux"}}
          - Contents: |
              #!ipxe
              exit
            ID: ""
            Meta: null
            Name: ipxe-mac
            Path: '{{.Machine.MacAddr "ipxe"}}.ipxe'
          Validated: false
      params:
        pxelinux-local-boot:
          Available: false
          Bundle: BasicStore
          Description: The method pxelinux should use to try to boot to the local disk
          Documentation: |2-

            On most systems, using 'localboot 0' is the proper thing to do to have
            pxelinux try to boot off the first hard drive.  However, some systems
            do not behave properlydoing that, either due to firmware bugs or
            malconfigured hard drives.  This param allows you to override 'localboot 0'
            with another pxelinux command.  A useful reference for alternate boot methods
            is at https://www.syslinux.org/wiki/index.php?title=Comboot/chain.c32
          Endpoint: ""
          Errors: []
          Meta: {}
          Name: pxelinux-local-boot
          ReadOnly: false
          Schema:
            default: localboot 0
            type: string
          Secure: false
          Validated: false
      roles:
        superuser:
          Available: false
          Bundle: BasicStore
          Claims:
          - action: '*'
            scope: '*'
            specific: '*'
          Description: ""
          Documentation: ""
          Endpoint: ""
          Errors: []
          Meta: {}
          Name: superuser
          ReadOnly: false
          Validated: false
      stages:
        local:
          Available: false
          BootEnv: local
          Bundle: BasicStore
          Description: Stage to boot into the local BootEnv.
          Documentation: ""
          Endpoint: ""
          Errors: []
          Meta:
            color: green
            icon: radio
            title: Digital Rebar Provision
          Name: local
          OptionalParams: []
          Profiles: []
          ReadOnly: false
          Reboot: false
          RequiredParams: []
          RunnerWait: false
          Tasks: []
          Templates: []
          Validated: false
        none:
          Available: false
          BootEnv: ""
          Bundle: BasicStore
          Description: Noop / Nothing stage
          Documentation: ""
          Endpoint: ""
          Errors: []
          Meta:
            color: green
            icon: circle thin
            title: Digital Rebar Provision
          Name: none
          OptionalParams: []
          Profiles: []
          ReadOnly: false
          Reboot: false
          RequiredParams: []
          RunnerWait: false
          Tasks: []
          Templates: []
          Validated: false

As the above example implies, YAML is the preferred format for
shipping around content packages, as it is generally easier to read
and edit than JSON is, especially when longer multi-line templates are
present.  You can also use the `drpcli contents` commands to
`unbundle` a cntent back for easier editing, and then `bundle` it back
up for processing and uploading to dr-provision.

Metadata
========

All content packages must have a meta section, which contains a
variety of different string values.

Operational Fields
~~~~~~~~~~~~~~~~~~

These metadata fields have meaning to dr-provision directly, and
control how dr-provision will process the content package whenever it
is loaded.

Name
<<<<

The name of the content package.  All content packages must have a
name, and names are not allowed to collide in a running instance of
dr-provision.  The name of the content package should either be a
single word or a short hypenated series of words for ease of command
line usage.

Version
<<<<<<<

The version of the content package.  Versions are roughly `Semver
compliant <https://semver.org/>`, except that we allow a leading
lower-case v and disgregard everything including and after the first
hyphen.  Version is optional, and if it is missing it is considered to
be 0.0.0.

RequiredFeatures
<<<<<<<<<<<<<<<<

A space separated list of features that dr-provision must provide for
the content package to function properly.  If you try to load a
content package onto a version of dr-provision that does not include a
required feature, the load will fail with an error indicating what
features are missing.  This field should be left blank if the content
package does not rely on any particular features of dr-provision.

Prerequisites
<<<<<<<<<<<<<

A comma separated list of other content packages that must be present
on the system for this content package to load.  Each entry in the
prerequisites list must either be the name of a content package, or
the name of a content package followed by a colon (:) and a space
separated list of version constrints.  If the field is left blank,
then this content pack is not considered to rely on any other content
packs.

Here are a couple of examples:

.. code-block:: yaml
   ---
   meta:
     Name: one
     Version: v1.2.3
   ---
   meta:
     Name: two
     Version: v1.2.3
     Prerequisites: one
   ---
   meta:
     Name: three
     Version: v1.2.4
     Prerequisites: 'one: >=1.0, two: <2.0.0'

Version Constraints
>>>>>>>>>>>>>>>>>>>

Prerequisite version constraints are processed according to the
following rules:

* `<1.0.0` Less than `1.0.0`
* `<=1.0.0` Less than or equal to `1.0.0`
* `>1.0.0` Greater than `1.0.0`
* `>=1.0.0` Greater than or equal to `1.0.0`
* `1.0.0`, `=1.0.0`, `==1.0.0` Equal to `1.0.0`
* `!1.0.0`, `!=1.0.0` Not equal to `1.0.0`. Excludes version `1.0.0`.
* `>1.0.0 <2.0.0` Greater than `1.0.0` AND less than `2.0.0`, so `1.1.1` and
  `1.8.7` but not `1.0.0` or `2.0.0`
* `<2.0.0 || >=3.0.0` Less than `2.0.0` OR greater than or equal to
  `3.0.0`, so would match `1.x.x` and `3.x.x` but not `2.x.x`

You can combine AND and OR constraints, AND has higher precedence.  It
is not possible to override precedence with parentheses.

Informational Fields
~~~~~~~~~~~~~~~~~~~~

These metadata fields contain information that may be of interest to
users of the content package, but that is not required for dr-provision
to prooperly load or use the content package.

Description
<<<<<<<<<<<

A short (one line) description of what the content bundle provides.

Source
<<<<<<

Where the content package is from.  This is generally either the
author or the organization that produced and maintains the content
package.  Deprecated in favor of Author and CodeSource.

Documentation
<<<<<<<<<<<<<

Longer information about what the content bundle is and what it does.
The documentation field may be either plain text or Restructured Text.

DisplayName
<<<<<<<<<<<

The name of the content package as it will de displayed in the Web UI.

Icon
<<<<

The icon that will be used for the content package in the Web UI.

Color
<<<<<

The color that the icon will be displayed in the Web UI

Author
<<<<<<

The original author of the content package.

License
<<<<<<<

The name of the license that everything in the content package is
distributed as.

Copyright
<<<<<<<<<

The copyright holder of the content package.

CodeSource
<<<<<<<<<<

The location that the content pack was loaded from.

Order
<<<<<

The order in which the content package will be displayed in the Web
UI.

Tags
<<<<

A comma-seperated list of tags appropriate for this content package.
Mainly used by the Web UI for sorting and filtering purposes.

DocUrl
<<<<<<

A URL to external documentation about this content package.


Data
====

Strictly speaking, a content package does not have to define any
objects, though it is little more than a placeholder if it
doesn't. Objects are defined in a content pack in the `sections` part
as follows:

.. code-block:: yaml
  ---
  meta:
    Name: foo
  sections:
    tasks:
      task-name:
        # rest of the tasks fields here
        Name: foo
    bootenvs:
      bootenv-name:
        # rest of the bootenvs fields here
        Name: bar

