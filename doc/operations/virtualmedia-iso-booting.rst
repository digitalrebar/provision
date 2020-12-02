.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Platform documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Platform; Virtual Media ISO Booting

.. _rs_virtualmedia_iso_booting:

Virtual Media ISO Booting
~~~~~~~~~~~~~~~~~~~~~~~~~

.. warning:: The Virtual Media ISO mount and booting method is currently a
             *Technical Preview*, as of December 2020.  This feature will be
             available in the v4.6.0 version (as a Tech Preview).  The
             documentation in this guide may be out of sync with usage as this
             feature is productionized.  Contact RackN support for any questions.

.. note:: This capability will be changing as new use case requirements and
          field testing is performed.  If you have any feature enhancements
          or issues with this feature, please contact RackN support.

Introduction
============

The primary method for Digital Rebar Platform (DRP) to obtain control and manage the lifecycle
of systems is via the in-band DHCP/PXE boot path.  In some cases, this path may not be available
due to network and security policy, or by operational principle.

The primary in-band DHCP/PXE path is dramatically faster, scales to 10s of thousands of systems,
and is in general a much more reliable path.  The Baseboard Management Controller (BMC, aka "IPMI",
aka "redfish controller") on many enterprise grade ssytems often lock up and require soft or hard
resets, is extremely slow to process the ISO virtual media boot path, and inconsistent from vendor
to vendor.

Should you still need to mount ISOs through the BMC, this Operations document will help you
to dynamically generate the Boot ISO and control your hardware via the BMC Virtual Media booting
path.

Here is a short-ish (10 min) video demonstrating this feature:

.. youtube:: YkGbi2jKM18
   :width: 100%


Supported Platforms
-------------------

Due to the customized and unique design of each BMC in different manufacturers hardware platforms,
the Virtual Media boot functions only work on limited platform types.  Here is the current list of
support in DRP by vendor.

==========  ============  ===================================================================
Vendor      Working?      Notes
==========  ============  ===================================================================
Dell        Yes           Tested with iDRAC 8 platforms.  Untested on iDRAC 9.
HPE         Maybe         Tooling is in place to support HPE platforms, but not tested yet.
Lenovo      Not Yet       No integration completed yet, future versions should be supported.
Other       N/A           Cisco UCS, Supermicro, and possibly other platforms may be added in
                          the future.  At this time there is no projected delivery dates for
                          these platforms.
==========  ============  ===================================================================


Prerequisites
-------------

For *Technical Preview* use of the VirtualMedia ISO Booting method, you must ensure the
following prerequisites are met:

  * DRP Version: *v4.6.0-alpha00.50* or newer
  * ``drp-community-content``: *v4.6.0-alpha00.154* or newer (*)
  * ``ipmi plugin``: *v4.6.0-alpha00.59* or newer
  * IPMI actions configured and working via the IPMI plugin
  * ``sledgehammer-builder``: *v4.6.0-alpha00.154* or newer (**)

.. note:: ``(*)`` A new version of the Sledgehammer BootEnv iso/tar must also be updated with
          the DRP Community Content update.  There are enhancements in Sledgehammer to
          support this feature (eg ``drpcli bootenvs uploadiso sledgehammer``).

.. note:: ``(**)`` Sledgehammer Builder is only required if you intend to create customized versions
          of Sledgehammer for the Boot ISO image, otherwise this content pack is not required.

The General Availability (GA) release v4.6.0 of all respective components should support
this capability, however it may still be marked as a *Technical Preview* feature and not
a production supported feature.

The new version of the Sledgehammer BootEnv ISO (technically it's a tarball), contains
a new embedded tempalte ``cdboot.iso`` which is used in the dynamic generated per-machine
custom ISO.  In addition, Sledgehammer has been modified to allow injection of Static IP
assignments to the Sledgehammer environment (in addition to the default DHCP IP address
assignment method).


The Process
===========

Process Overview
----------------

Digital Rebar Platform (DRP) will dynamically generate a custom bootable ISO that should
successfully boot both Legacy and UEFI BIOS systems.  The ISO utilizes GRUB2 as the boot
method, which has extensive and broad support for a large number of systems (specifically
the UEFI boot stack pieces).

To set a specific machine (or machines) to begin generating a dynamic generated boot ISO,
the Machine must have the Param ``boot-virtual-iso`` set to the Boolean value *True*.

This param will then cause the ``dr-provision`` service to automatically generate the ISO
when the Machine(s) change BootEnv **in to** the Sledgehammer BootEnv.  There are several
other Params that are used to inject custom details in to the dynamically generated ISO,
which are described below.

In addition, DRP will automatically mount the dynamically generated ISO path location on
the BMC, and attempt to set the Next Boot directive to boot via the VirtualMedia path.

The current usage path is as follows:

  * Set the ``boot-virtual-iso`` Param on the target machine(s)
  * Use ``drpcli`` to verify the dynamically generated ISO is set to mount to the Machine(s) BMC
  * Set any supporting customization Params (eg ``network-data`` static IP addressing)
  * Set the Machine(s) BootEnv to Sledgehammer (if currently in Sledgehammer, change it to ``local``, then back to ``Sledgehammer``).

The IPMI plugin must be installed, and have been configured for the Machine(s) that you
are controlling.  This is typically done in a *discovery* phase workflow with the use
of the ``ipmi-inventory`` Stage.

More complete details in the below sections.


Enable Dynamic ISO Generation
-----------------------------

For any machine(s) that you will attach Virtual Media ISOs to, you must set the following
Params:

  * ``boot-virtual-iso`` to the Boolean value *true*

  ::

    # example of setting machine 'mach-01' (a UUID can be used) to have dynamically generated ISOs
    drpcli machines set Name:mach-01 param boot-virtual-iso to true

Once this value is set on the machine, and all appropriate Prerequisites fulfilled (listed above),
then the ``dr-provision`` service will dynamically generate a custom ISO when the Machine is
changed **in to** the Sledgehammer BootEnv.

The ISO will be built and cached at on the DRP server under the ``tftpboot`` path in a directory
named ``dynamic_isos``.  In a typical "production" install mode, this is found at the following
fully qualified path in a separate directory for each Machine, with the Machines UUID value:

  * ``/var/lib/dr-provision/tftpboot/dynamic_isos/<MACHINE_UUID>``


Verify the Virtual Media ISO Mount
----------------------------------

For the system to boot from the Virtual Media ISO - the Baseboard Management Controller (BMC)
must be instructed to mount the ISO.  This path is not very well standardized between the
different manufacturers.  Currently, the Redfish protocol is the only supported method for
mounting the Virtual Media.  However, there is no standardized Redfish path for setting
the "bootonce" via VirtualMedia to the BMC.  The IPMI plugin must be correctly configured to
support the vendor specific capabilities to control the BMC (eg iDRAC or iLO) to support
the reboot once to VirtualMedia.

.. note:: BootEnv transitions will automatically attempt to set the VirtualMedia mount path,
          and set the appropriate boot once directive.

To verify the media mount path, the ``drpcli`` command line tool has been extended to support
several Redfish query paths to check/set/verify/mount/unmount media on the BMCs Virtual Media
mount point.

  ::

    # verify the mount path (again, a Machine UUID can be specified instead of Name)
    drpcli machines runaction Name:mach-01 statusVirtualMedia ipmi/mode redfish

An example output showing the automatically generated VirtualMedia mount path:

  ::

    {
      "Image": "http://10.10.10.10:8091/dynamic_isos/aebf8b66-276f-4234-87b4-a0d79075d76f/sledgehammer/boot.iso",
      "Inserted": true
    }

The BMC web portal and other tools should also reflect this status.


Customize the Network Configuration
-----------------------------------

By default, the custom ISO will attempt to utilize DHCP on the first iterated network
interface (eth0).  If this behavior is not desired and needs changed; the use of the
``network-data`` Param structure can control the values.

.. warning::  This ``network-data`` structure MUST be added to the machine prior to the
              machine generating the custom ISO.  Do not transition BootEnvs until the
              correct values have been set in this Param.

Here is an example Param stanza to define static IP assignment to a machine, as an example:

  ::

    {
      "prov": {
        "address": "10.10.10.100",
        "gateway": "10.10.10.1",
        "interface": "eth10",
        "netmask": "255.255.255.0"
      }
    }

It may also be possible (but is as of yet untested), to set a VLAN tag value for environments
using VLAN tagging with the addition of the key/value pair ``"vlan": 1020``.

.. note:: This ``network-data`` structure can be used with the standard Sledgehammer PXE 
          in-band boot path, and should allow you to set static IP assignments for
          Sledgehammer, disabling IP address acquisition via DHCP.


Boot From the VirtualMedia ISO
------------------------------

Once you have enabled the generation of the dynamic ISOs, set any specific ``network-data`` values
required, and verified the VirtualMedia mount, you can now boot the system.

To boot in to the dynamic custom ISO - transition the Machine in to the *Sledgehammer* BootEnv.

If the Machine is already in the *Sledgehammer* BootEnv, you will need to change the machine to
something else (eg ``local`` BootEnv), then back to Sledgehammer.

Here is the example of transitioning a machine that is currently in Sledgehammer, to local, then
setting the Workflow to ``discover-base``; which defines Sledgehammer as the bootenv.

  ::

    # remove workflow for machine named 'mach-01'
    drpcli machines workflow Name:mach-01 ""

    # set the BootEnv to local - expect up to 60 second wait for command to complete
    drpcli machines bootenv Name:mach-01 local

    # set workflow to 'discover-base', which specifies Sledgehammer bootenv
    # again - wait up to 60 seconds for command to complete
    drpcli machines workflow Name:mach-01 discover-base

.. note::  The bootenv transition changes now trigger a dynamic ISO regeneration, and the
           VirtualMedia attach commands to the BMC.  These operations are all slow, and
           take time.  Patience, young Skywalker... 

It is advised that you should watch the physical or virtual console to verify the machine
boot process.  You should see the system boot in to the Sledgehammer dynamically generated
ISO.  The process looks VERY similar to the boot process of the standard in-band DHCP/PXE
boot process.


Notes and Troubleshooting
=========================

Here is a list of notes or debugging processes to help if there are issues with
the VirtualMedia booting process.

Restore Default In-Band Management Path
---------------------------------------

If a machine object has been modified to use the out-of-band dynamically generated
custom ISO, it can be returned to proper in-band management by simply removing the
``boot-virtual-iso`` Param from the machine, for example:

  ::

    # remove the boot-virtual-iso param from machine Named 'mach-01'
    drpcli machines remove Name:mach-01 param boot-virtual-media

In addition, the ``network-data`` param may or may not need to be removed.  If
moving back to DHCP IP address based PXE booting, then typically this param should
be removed.  However, the DHCP/PXE boot path process for in-band management of the
system will still honor the settings in this param when Sledgehammer boots.

If complete clean up is required, you may also want to remove the dynamically generated
ISO images in the ``tftpboot/dynamic_isos/`` directory path.  Note that ISOs are stored
in a sub-directory with the Machines UUID as the directory name.


Performance Impact
------------------

Any command and control functions implemented directly to the Baseboard Management Controller
(BMC) are generally extremely slow.  Many commands described above will block and wait for 30
to 60 seconds before the command completes.

Additionally, with the ``boot-virtual-iso`` set to ``true``, specific BootEnv changes force the
``dr-provision`` service to dynamically generate a new custom ISO.  This process can be CPU and I/O
intensive, especially if many machines are transitioned at once.

There is currently no sizing guidelines to for large scale infrastructure use of this feature.
However, expect additional CPU and disk I/O impact.


Verifying the Boot to VirtualMedia
----------------------------------

This process attempts to automatically set the VirtualMedia boot process and attach the dynamic
generated ISO to the BMC VirtualMedia mount point.  There are several ways to verify this
has happened, including use of the vendor specific tooling, vendor BMC Web service, Redfish
calls, etc.  In addition, the ``drpcli`` client tool has support to manipulate and verify
the boot process.

  ::

    # verify the status - note this can take a long time to complete
    drpcli machines runaction Name:mach-01 statusVirtualMedia ipmi/mode redfish

In addition, observing the Boot POST process of the Machine in question should yield visual
clues.  For example, Dell systems with iDRAC 8 BMCs would show output like:

  * ``IPMI: Boot to Virtual CD Requested``


VirtualMedia Mount Options
--------------------------

The new actions in the IPMI plugin support manipulating the VirtualMedia mount paths, here
are examples of different usage scenarios:

**Mount ISO**

  ::

    # mount the dynamically generated ISO for the machine specified by UUID
    # also set the boot once from virtual media option
    drpcli machines runaction bb1eadf9-4b5e-46a7-a577-d07e2a33138f mountVirtualMedia ipmi/mode redfish ipmi/virtual-media-url http://10.10.10.10:8091/dynamic_iso/bb1eadf9-4b5e-46a7-a577-d07e2a33138f/sledgehammer/boot.iso ipmi/virtual-media-boot true

**Unmount ISO**

  ::

    # by machine Name reference:
    drpcli machines runaction Name:mach-01 unmountVirtualMedia ipmi/mode redfish

**Perform Power Reboot via Redfish**

  ::

    # powercycle machine by name, using Redfish
    drpcli machines runaction Name:mach-01 powercycle ipmi/mode redfish

**Get Current Power Status**

  ::

    # get current power status using the default IPMI mode (redfish, ipmi protocol, or vendor specific)
    drpcli machines runaction Name:mach-01 powerstatus

    # get it specifically via the Redfish API
    drpcli machines runaction Name:mach-01 powerstatus ipmi/mode redfish


Validate Dynamic ISO Generated
------------------------------

Virtual Media ISOs are generated and stored under the ``tftpboot`` directory structure, in
the ``dynamic_isos`` directory.  Each dynamic ISO for a Machine is stored in a sub-directory
with the Machine's UUID.  In a standard production install, this would be:

  * ``/var/lib/dr-provision/tftpboot/dynamic_isos/<MACHINE_UUID>/``

After the Machine has transitioned into Sledgehammer, the ISO will be stored in this
directory path, and the directory tree will look like the following:

  ::

    cd /var/lib/dr-provision/tftpboot/dynamic_isos

    tree bb1eadf9-4b5e-46a7-a577-d07e2a33138f/
    bb1eadf9-4b5e-46a7-a577-d07e2a33138f/
    ├── local
    └── sledgehammer
        └── boot.iso

    2 directories, 1 file

Mounting the ISO and reviewing it's contents should show:

  ::

    mount bb1eadf9-4b5e-46a7-a577-d07e2a33138f/sledgehammer/boot.iso /mnt
    tree /mnt
    /mnt
    ├── boot
    │   └── grub
    │       ├── fonts
    │       │   └── unicode.pf2
    │       ├── grub.cfg
    │       ├── i386-pc
    │       │   ├── acpi.mod
    │       │   ├── <...snip...>
    │       │   └── zfs.mod
    │       └── roms
    ├── boot.catalog
    ├── EFI
    │   └── BOOT
    │       ├── BOOT.conf
    │       ├── BOOTIA32.EFI
    │       ├── BOOTX64.EFI
    │       ├── fonts
    │       │   ├── TRANS.TBL
    │       │   └── unicode.pf2
    │       ├── grub.cfg
    │       ├── grubia32.efi
    │       ├── grubx64.efi
    │       ├── mmia32.efi
    │       ├── mmx64.efi
    │       └── TRANS.TBL
    ├── stage1.img
    └── vmlinuz0

    8 directories, 292 files

The customizations to network configuration are written in to the GRUB boot config
file, which can be verified as follows:

  ::

    $ sudo cat /mnt/boot/grub/grub.cfg
    if [ ${grub_platform} == "efi" ]; then
      set root=(cd0)
      set linuxcmd=linuxefi
      set initrdcmd=initrdefi
    else
      set root=(cd)
      set linuxcmd=linux
      set initrdcmd=initrd
    fi
    timeout=0
    # There are 15 lines of 80 comments after for padding.

    # replace here
    menuentry "Sledgehammer" {
      $linuxcmd /vmlinuz0 BOOTIF=discovery rootflags=loop root=live:/sledgehammer.iso rootfstype=auto ro liveimg rd_NO_LUKS rd_NO_MD rd_NO_DM provisioner.web=http://10.10.10.10:8091 rs.uuid=bb1eadf9-4b5e-46a7-a577-d07e2a33138f      provisioner.ip=10.10.10.199/24   provisioner.gw=10.10.10.254   provisioner.interface="eth10"      -- console=ttyS0,115200 console=tty0
      $initrdcmd /stage1.img
      boot
    }

The relevant customizations from the ``network-data`` structure are converted to
the Sledgehammer *menuentry* stanza values (eg *provisioner.ip*, *provisioner.gw*, etc.).

.. note:: There are also a large number of "padding" pound sign characters, which is
          required for absurd and arcane GRUB reasons.  Do not change them.  You have
          been warned.

