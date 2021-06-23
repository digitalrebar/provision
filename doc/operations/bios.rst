.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Using The BIOS Plugin


.. _rs_operation_bios:

Using the ``bios`` plugin
+++++++++++++++++++++++++

The ``bios`` plugin is used to configure a server to a specific set of BIOS values.  To
use the BIOS plugin, you must be running a supported hardware platform by Digital Rebar.
Supported platforms will have a Catalog content pack named (generally) for the hardware
manufacturer.  For example ``dell-support``, ``hpe-support``, or ``lenovo-support``.


Installation
------------

You must install (from the *Catalog*), the following items:

  * ``bios`` plugin
  * ``hardware-tooling`` content
  * one of the Vendor support content packs (eg ``dell-support``)

The above items can be installed from the *Catalog* in the Portal, or via the following
command line installation:

  ::

    # change 'vendors' variable to supported vendor(s) name (all lowercase, space separated)
    vendors="dell hpe lenovo"
    drpcli catalog item install bios
    drpcli catalog item install hardware-tooling
    for vendor in $vendors; do drpcli catalog item install ${vendor}-support; done


General Process
---------------

The BIOS configuration subsystem of Digital Rebar Platform is designed (and only
tested) to run in the Sledgehammer environment.  RackN does not support use of the
tooling outside of this environment.

.. warning:: Ensure you boot your machines in to Sledgehammer to perform the following process.

The general process for using the BIOS plugin utilizes three separate *Stages*, in
the following usage flow:

  1. use ``bios-inventory`` - collects detailed information on supported values on the system that the Stage runs on; populates the Param ``bios-current-configuration``
  2. use ``bios-baseline`` - populates the ``bios-target-configuration`` compatible Param of the current system's BIOS configuration
  3. Modify the ``bios-target-configuration`` to support any desired BIOS configuration changes
  4. use ``bios-configure`` - applies the values found in the Param ``bios-target-configuration`` to the system

The ``bios-target-configuration`` Param is the settings that will be applied to a system.
Use the output of the ``bios-baseline`` Stage to create a "copy" of the current systems
configuration values.  This can be used after an operator has manually set a number of
values through the vendors BMC configuration settings, or by using it as a starting point,
and modifying the values of the results of the ``bios-baseline`` *Stage*.


The ``bios-current-configuration``
----------------------------------

The ``bios-current-configuration`` records both the current settings of the various
bios values that may be set on a given platform, and it also provides "type definition"
information on the potential data that a given setting may potentially require to for
use.  If you have questions on what values should be placed in the fields, refer to
the ``bios-current-configuration`` for clues, or the Vendors documentation.

Ultimately, however, it is the ``bios-target-configuration`` which is used to apply
updated settings to the systems BIOS.


The ``bios-target-configuration``
---------------------------------

Note that the ``bios-target-configuration`` structure is slightly different than what the
``bios-inventory`` stage produces in the ``bios-current-configuration``.  The ``bios-baseline``
stage produces a validly formatted ``bios-target-configuration`` for actually making changes
to the system BIOS.

Once you have produced a version of the ``bios-target-configuration`` settings values, it may
safely be used *ONLY* with those Vendor specific Model/Version, and Firmware level systems.
Vendors *regularly* change the BIOS supported configuration values and details between model/platform
versions, and potentially from older to newer versions of Firmware that is flashed on the system.

Legacy -vs- UEFI
----------------

If you have the same vendor model and version of hardware, there are often substantial differences
in the BIOS configuration values, dependent on the system being in either *Legacy* or *UEFI BIOS*
boot modes.  You are advised to ensure you re-run ``bios-baseline`` for each of *Legacy* and *UEFI BIOS*
boot modes to verify values for the boot mode are correct.

In particular - the Boot Order semantics are often very different between *Legacy* and *UEFI BIOS* boot
modes.


The ``bios-skip-config``
------------------------

The ``bios-skip-config`` Param allows for all BIOS related tasks to "skip" running, if this Params
value is set to ``true``.  This is often used in Discovery type default workflows to turn off the
BIOS sub-system, without using customized Workflows for different use cases.  Verify that this value
is not set to ``true`` on the machine if you see the BIOS tasks just skipping, and exiting with zero
(success) value.


Video Example
-------------

The following is an example of using the BIOS subsystem in video form.  This video example
shows the above outlined process on a fleet of Dell R730xd server platforms.  Only a very
minimal BIOS setting value has been changed via use of editing the Param value.

Generally speaking, you should extract the ``bios-target-configuration`` produced in this example,
and incorporate it in to a Content Pack.  See the Color Demo training content pack and
video references for more details:

  * https://github.com/digitalrebar/colordemo

.. youtube:: ABCDEFGHI
   :width: 100%
