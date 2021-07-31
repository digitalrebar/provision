.. Copyright (c) 2021 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Platform documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Platform; LPAR
  pair: Digital Rebar Platform; ppc64le

.. _rs_lpar_ops:

LPAR Operations
===============

This section will describe the support of LPAR VMs.

Support Changes
---------------

In the v4.7.0 and beyond releases, RackN DRP supports the ppc64le architecture for Linux operating systems.

This includes the server, client, and booting virtual machines.  The DRP endpoint and client
has been update to run in the ppc64le linux environment.  The install script also supports this as well.
These system work the same way as the other architectures.

The new feature to support LPAR vms in DRP.  The LPAR VM once created can be network booted, discovered, and installed.
The following section will describe the configuration changes to support LPAR VMs.

NOTE: You must get the latest sledgehammer for ppc64le.


LPAR Configuration
------------------

The LPAR VM can only network boot from BOOTP.  This requires that a reservation be created in DRP to boot the machine.  Once the LPAR VM has been created,
the HMC can be queried to get the MAC address.  The only requirement for the Reservation is to set the bootfile to `core.elf`.

In the following example, there is an LPAR VM with a mac address, *fa8b9f59bd20* and its address is *129.40.108.5*.

  ::

    Addr: 129.40.108.5
    Token: 'fa:8b:9f:59:bd:20'
    Options:
      - Code: 67
        Value: core.elf
    Strategy: MAC


With that in place, issuing a netboot command for the LPAR VM from the HMC will cause the system to boot sledgehammer and discover the system.

  ::

    lpar_netboot -f -t ent -m fa8b9f59bd20 -s auto -d auto b1p052_Target-8c0bd1b0-0000084e default_profile Server-8247-22L-SN212169A


To add IPMI-like operations for the LPAR VM, setting the following parameters will allow the LPAR VM to be controlled by IPMI operations.

* ipmi/address = HMC address
* ipmi/username = HMC access username
* ipmi/password = HMC access password
* ipmi/lpar-id = The UUID of the LPAR VM from the HMC.  e.g: 48EB2A9F-6028-41C6-8D2F-6A539361BB29

With these in place, you can powercycle, powerstatus, poweron, and poweroff the LPAR VM.

Currently, only sledgehammer and centos8 have been updated to support ppc64le.

