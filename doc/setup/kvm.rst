.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Setup

.. _rs_setup_kvm:

KVM Setup Instructions
~~~~~~~~~~~~~~~~~~~~~~

Overview
--------

This document will help define one possible method to setup a KVM hypervisor host
as a Digital Rebar endpoint, running KVM, and provisioning Virtual Machines on
that KVM host.

The configuration of this setup utilizes the RackN ``kvm-test`` plugin which is
currently in *Alpha* or *Experimental* stages.  It has limited functionality, but
does provide useful capabilities for managing the KVM environment, and appears
to be relatively stable.

This environment is tested and running as a primary development platform for one
of the RackN engineers

Please feel free to provide any feedback or alternative setup methods.


Architecture
------------

The setup of the RackN ``kvm-test`` *Plugin Provider* requires that the DRP Endpoint service
(``dr-provision``) run directly on the same machine that the KVM hypervisor services are
also running on.  Please see the Future ToDos at the bottom of this document regarding this.

The setup of this environment utilizes the DRP Endpoint running on the bare hardware machine,
alongside the KVM services.  All Virtual Machines will be connected to a single bridge
(called ``kvm-test`` in this document), which provides the connectivity to the Hypervisor
and outbound *external* network connectivity for the Virtual Machines.

The Hypervisor will provide NAT/Masquerading services and act as the Default Gateway to the
*External* network (either your local network, or outbound to the Internet).

Here is a simplified diagram:

  ::

     +--------------+
     |   External   |
     |              +------------ ( Internet )
     |    Switch    |
     +--------------+
             |
             |
             |
             |
             |
     +--------------+
     |     eth0     | (ex: 10.10.10.10)
     |              |
     | DRP  Machine |   * must be a physical server
     |   KVM Qemu   |   * performs IP gateway NAT/Masquerading
     |              |
     |   kvm-test   | (ex: 192.168.124.1)
     +--------------+
             |
             |
             |                     KVM Bridge (kvm-test - 192.168.124.0/24)
             |
             +-----------------------------------+------------------------------------+------- ...
             |                                   |                                    |
     +-----------+                          +----+------+                          +--+--+
     |   eth0    | (ex: 192.168.124.11)     |   eth0    | (ex: 192.168.124.12)     |     |
     |           |                          |           |                          |     |
     |   VM 01   |                          |   VM 02   |                          | ... |
     |           |                          |           |                          |     |
     +-----------+                          +-----------+                          +-----+


This method has a few advantages:

* You can have many target VMs and they are easy to create/destroy
* The DRP Endpoint controls the KVM network bridge setup/tear down
* The DRP endpoint may have Internet access, and can selectively provide Internet access to the target VMs using standard Linux tools
* Virtual Machine creation/destruction can be performed from Digital Rebar APIs / Control
* Virtual Machine power operations can be performed form Digital Rebar APIs / Control
* KVM can replicate many different machine architectures and configurations, Legacy BIOS, UEFI BIOS, etc.
* The only primarly limitaiton to number of VMs is physical memory, disk space, and your oversubscription comfort levels

Prerequisites
-------------

This test environment is based on CentOS 7.8 and *Libvirt 5.0.0* and *QEMU-KVM 2.12.0*.  It
is entirely likely that this setup will run fine under other Linux Distros.  You may need to
adjust package selections or other distro specific *isms* to the below recipe.

You will need:

* Physical System with a `processor that supports virtualization <https://www.linux-kvm.org/page/Processor_support>`_
* CentOS 7.8 - **minimal** installation (**do not** do additional Package Group installation)
* Shell/SSH access to the physical system

The remainder of this document assumes you have CentOS 7.8 installed on your physical server,
and you are accessing the system either from the console, or an SSH session.


Install KVM/Libvirt
-------------------

Create a CentOS ``yum`` repository with the latest Qemu/KVM and Libvirt packages:

  ::

    cat << EOF-REPO > /etc/yum.repos.d/CentOS-Virt.repo

    # CentOS-Virt

    [virt-libvirt-latest]
    name=CentOS-$releasever - Virt Libvirt Latest
    #mirrorlist=http://mirrorlist.centos.org/?release=$releasever&arch=$basearch&repo=virt&infra=$infra
    #baseurl=http://mirror.centos.org/centos/$releasever/virt/$basearch/
    baseurl=http://mirror.centos.org/centos/7/virt/x86_64/libvirt-latest/
    gpgcheck=1
    gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-CentOS-7

    [virt-kvm-common]
    name=CentOS-$releasever - Virt KVM Common
    #mirrorlist=http://mirrorlist.centos.org/?release=$releasever&arch=$basearch&repo=virt&infra=$infra
    #baseurl=http://mirror.centos.org/centos/$releasever/virt/$basearch/
    baseurl=http://mirror.centos.org/centos/7/virt/x86_64/kvm-common/
    gpgcheck=1
    gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-CentOS-7

    EOF-REPO

Install the packages:

  ::

      yum makecache
      yum install qemu-kvm libvirt libvirt-python libguestfs-tools bridge-utils iptables iptables-services util-linux unbound curl wget jq


Install Digital Rebar
---------------------

Install Digital Rebar based on the :ref:`rs_quickstart` or :ref:`rs_install` documentation,
***WITH THE FOLLOWING DIFFERENCE***

* use "production" mode (do not use ``--isolated`` install flag)
* add ``--systemd`` to enable the SystemD startup unit files

Once you have a basic DRP Endpoint up and running, add the ``kvm-test`` Plugin Provider
via the *Catalog* in the UX, or via the following command line usage:

* ``drpcli catalog item install kvm-test``


Install EFI Firmware Code for Virtual Machines
----------------------------------------------

The default configuration and machine creating by the Digital Rebar plugin will
build VMs with UEFI BIOS and firmware.  This requires installation of the CODE
and VARS to support the EFI mode in the VMs.  There are many ways to integrate
EFI in to KVM/Libvirt guests.  This is only one, and is required to be setup
this way for this environment.

  ::

    wget https://s3-us-west-2.amazonaws.com/get.rebar.digital/artifacts/fw.tar.gz -O /tmp/fw.tar.gz
    cd /
    tar -xzvf /tmp/fw.tar.gz

If you are building guest VMs direclty with Libvirt/Qemu tools, you are free to use
any other UEFI firmware packages or solutions.  Additionally, you can create
virtual machines with Standard / Legacy BIOS boot mode and not utilize UEFI at all.


Create the ``kvm-test`` Virtual Machine Subnet
----------------------------------------------

Now create a DHCP Subnet (pool) for the Virtual Machines to utilize.  If you have
previously created a DHCP pool in DRP, you may need to either destroy that pool,
or create this new pool for your VMs.

  ::

    cat << EOF_SUBNET > $HOME/kvm-test-net.yaml
    {
      "Name": "kvm-test",
      "ActiveEnd": "192.168.124.200",
      "ActiveLeaseTime": 3600,
      "ActiveStart": "192.168.124.11",
      "Enabled": true,
      "Meta": {},
      "NextServer": "",
      "OnlyReservations": false,
      "Options": [
        { "Code": 3,  "Description": "Default GW",  "Value": "192.168.124.1" },
        { "Code": 6,  "Description": "DNS Servers", "Value": "8.8.8.8" },
        { "Code": 15, "Description": "Domain",      "Value": "kvm-test.local" },
        { "Code": 1,  "Description": "Netmask",     "Value": "255.255.255.0" },
        { "Code": 28, "Description": "Broadcast",   "Value": "192.168.124.255" }
      ],
      "Pickers": [ "hint", "nextFree", "mostExpired" ],
      "Proxy": false,
      "ReservedLeaseTime": 7200,
      "Strategy": "MAC",
      "Subnet": "192.168.124.1/24",
      "Unmanaged": false
    }
    EOF_SUBNET

    drpcli subnets create $HOME/kvm-test-net.yaml

.. note:: See the DNS notes section below for additional options related to name server
          lookup configurations.

Once our DHCP Subnet is created, we can now define the KVM Plugin configuration,
which works cooperatively with the DHCP Subnet.


Instantiate the ``kvm-test`` Plugin and Configuration
-----------------------------------------------------

The ``kvm-test`` Plugin Provider adds the ability to manage KVM networks.  We have
to instantiate a Plugin with the network details to create/enable the bridges
for this setup.  Create a YAML file with the following configuration to specify
the storage pool, and the network (subnet) configurations:

.. warning:: Your network configuration here must match the Subnet created above.

  ::

    cat << EOF-PLUGIN > $HOME/kvm-test-plugin.yaml
    ---
    Name: kvm-test
    Description: ""
    Meta: {}
    Params:
      kvm-test/storage-pool: dirpool
      kvm-test/subnet:
        address: 192.168.124.1/24
        domain: kvm-test.local
        gateway: 192.168.124.1
        name: kvm-test
        nameserver: 8.8.8.8
    Provider: kvm-test

    EOF-PLUGIN

    # now create the plugin from the above config file
    drpcli plugins create $HOME/kvm-test-plugin.yaml

.. note:: This is only one possible configuration.  You can specify different addressing,
          DNS servers, etc. to match your requirements.  See the DNS topic below.

Once you create the Plugin, you should  now be able to see the network bridge in the
OS of your DRP Endpoint.  If you do not, restart DRP (``systemctl restart dr-provision``),
and check again:

  ::

    # ip a sh kvm-test
    7: kvm-test: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default qlen 1000
        link/ether fe:54:00:28:20:b7 brd ff:ff:ff:ff:ff:ff
        inet 192.168.124.1/24 brd 192.168.124.255 scope global kvm-test
          valid_lft forever preferred_lft forever

    # brctl show kvm-test
    bridge name	bridge id		STP enabled	interfaces
    kvm-test		8000.fe54002820b7	no		vnet0


Setup IP Forwarding and NAT Masquerading
----------------------------------------

This portion of the configuration sets up the DRP Endpoint machine as a *router* and
IP NAT Masquerading host for the Virtual Machines that DRP will manage and control.
This has the drawback that all VM network traffic will route *through* the DRP Endpoint
Operating System, however, it is a simplified setup that works very well and is
repeatable and independent from the machines external network topology.

Configuring the Virtual Machine networks differently (eg with multiple NICs), and creating
associated Bridges to other network devices is entirely possible, but outside the
scope of this document.

The below enables extremely simplified ``iptables`` rules to perform these tasks.  Other
firewall services can be substituted, as long as you perform these similar capabilities.
We encourage you to also increase the external network rulesets to better protect your
system services.

.. warning:: If you utilize other firewall services (eg ``firewalld``), ensure you perform
             the equivalent setup as below.

This configuration will check the DRP Endpoint Plugin instantiation (as above) to
set the network devices and rules automatically.

  ::

    #!/usr/bin/env bash

    # get DRP kvm-test base network for $1 (kvm-test by default if not defined)
    # saves configurations so they survive reboots
    PATH=$PATH:/usr/local/bin

    NW=$(drpcli plugins show kvm-test | jq -r '.Params."kvm-test/subnet".name')
    NET=${1:-$NW}
    JSON=/tmp/$NET-network.json
    drpcli subnets show $NET > $JSON
    NETWORK=$(cat $JSON | jq -r '.Subnet' | sed 's/\(.*\)\.\(.*\)\.\(.*\)\..*$/\1.\2.\3.0/')
    HOST=$(cat $JSON | jq -r ".Options | .[] | select(.Code==3) | .Value")
    MASK=$(cat $JSON | jq -r ".Options | .[] | select(.Code==1) | .Value")

    echo "
    NETWORK  :: $NETWORK
    HOST     :: $HOST
    MASK     :: $MASK
    "

    systemctl start iptables
    systemctl enable iptables

    iptables -t nat -A POSTROUTING -s "$NETWORK/$MASK" ! -d "$NETWORK/$MASK" -j MASQUERADE
    iptables -I FORWARD 1 -i $NET -j ACCEPT
    iptables -I FORWARD 1 -o $NET -m state --state RELATED,ESTABLISHED -j ACCEPT

    service iptables save

    sysctl net.ipv4.ip_forward=1
    echo "net.ipv4.ip_forward=1" > /etc/sysctl.d/50-ipv4_ip_forward.conf
    exit 0


Virtual Machine Creation
------------------------

The DRP based KVM plugin supports creating and provides limited control over virtual
machine actions (``poweron``, ``poweroff``, ``reboot``, etc.).  To support these
actions, the Machine must have the KVM specific Parameters:

* ``kvm-test/machine: { ... }``
* ``machine-plugin: kvm-test``

These values help the plugin to control the Virtual Machine state correctly.

.. note:: The *create* VM variations below all initially create the Virtual Machine as an
          object in DRP to manage, and the Virtual Machine within KVM; you will need
          to power the Machine(s) on after creation (details below).


Create VMs via Portal (UX)
==========================

You can create new virtual machines from the RackN hosted Portal, by visiting the
*Plugins* menu page, and then selecting the ``kvm-test`` plugin that was created
earlier.  On this page, there is a simple dialog input form to specify the Machine
base name, and number of machines to create.  These machines will be created based
on the internally compiled in machine specs, which can not be changed easily at
this time.


Create VMs via Command Line (drpcli)
====================================

The ``drpcli`` command can be used to create a Machine object, and by setting specific
Param values on the Machine object, the system will create VMs.  The below shell script
builds up the appropriate Machine Ojbect JSON information, initiates the create, and
powers on the VM.  You can just run the script to receive a randomized machine name,
or pass a VM name as ARGv1 to the script.

  ::

    cat << EOF_SCRIPT > $HOME/create-vms.sh
    #!/usr/bin/env bash
    # create a DRP managed KVM Virtual Machine

    # optionally specify machine name as ARGv1

    UUID=$(uuidgen)
    MID=$(mktemp -u XXXXXX | tr '[:upper:]' '[:lower:]')
    NAME=${1:-mach-$MID}

    JSON="
    { \"Name\": \"$NAME\",
      \"Params\": {
        \"machine-plugin\": \"kvm-test\",
        \"kvm-test/machine\": {
        \"arch\": \"x86_64\",
        \"bios\": \"bios\",
        \"bridge\": \"kvm-test\",
        \"cores\": 2,
        \"disk-size\": 20,
        \"memory\": 2048,
        \"Name\": \"$NAME\",
        \"pool\": \"default\",
        \"uuid\": \"$UUID\"
        }
      }
    }
    "

    echo ">>> Creating Machine:  $NAME"
    drpcli machines create "$JSON"

    echo ">>> Powering on machine: $NAME"
    drpcli machines action Name:$NAME poweron

    exit 0
    EOF_SCRIPT

    chmod $HOME/create-vm.sh

Example usage of the script, creating VM named *vm-test-01*:

* ``./create-vm.sh vm-test-01``


Create VMs via ``virsh`` or other Methods
=========================================

You are welcome to create Virtual Machines through any other traditional VM creation
mechanism that is supported by your setup.  Just set the Virtual Machine to boot
PXE first, and generally speaking, it should be discovered, and added in to inventory
on the DRP Endpoint to be managed as any other normal machine.

This method allows you greater control over the virtual machine specifications if you
need to test different hardware architectures, components, BIOS setups, etc.

.. note:: Machines created this way can not be power controlled via DRP, unless you add
          Machine Params as specified above.  Reference a Digital Rebar created VM for
          the correct structure of the Params.


Power Control of your VMs
=========================

Power on, off, reboot, etc. controls are enabled through DRP *Plugin* **actions**.  The actions
work as long as the Machine object has the correct Params set on it to allow the Plugin system
to reference properly.

If you are creating machines directly via ``virsh``, ``qemu-kvm``, ``qqemu-system-x86_64``, or
other options; you must manually add the Params to the DRP Machine object.  In this case, it is
recommend you create a VM via the Portal/UX, observe the Machines ``kvm-test/*`` param values,
and replicate those.

For VMs created via the Portal/UX or CLI tool, you can use DRP to perform the machine actions.
In the Portal, use the standard *Actions* dialogs on the *Machines* menu page, or on the Machines
detail panel to effect power changes.  For CLI actions, perform the following:

  ::

    VM="MACHINE_NAME"  # change this!

    # get a list of the Actions availble on the Machine
    drpcli machines actions Name:$VM | jq -r '.[].Command'

    # perform an action
    drpcli machines action Name:$VM poweron

.. note:: It appears that the CLI method of passing machine actions for VM power control
          are not very reliable.  Utilize the Portal/UX method if you have issues with
          the CLI method.


Virtual Consoles for VMs
------------------------

There are a couple of options for managing the view/interaction with the console of the VMs.  These
all basically come down to how the VM was created, and the supported console setup inside the
Machine specifications.  The default setup for VMs created by the Portal/UX Plugin operation, or the
CLI tooling utilize a *Spice* display method.  You can find more at the following resources:

* `Spice on your local KVM workstation <https://www.linux-kvm.org/page/SPICE>`_
* `Spice KVM / CentOS Howto <https://wiki.centos.org/HowTos/Spice-libvirt>`_
* `Download Spice Clients and Tools <https://www.spice-space.org/download.html>`_

The Virtual Machine display specifications can be changed to utilize VNC and VNC clients as well.

In all cases, if you are remote from the KVM hypervisor, you will need to either forward Ports to
the KVM Hypervisor, setup a VPN connection to connect to the Hypervisor and your Guest VMs, or
otherwise arrange to have the Spice/VNC ports opened up on the External network interface of your
hypervisor for access directly over the network.

At the shell/command line of the Hypervisor, you can see the emulated serial console of the VM
by use of the ``virsh`` tool, as follows:

  ::

    virsh console <VM_NAME>

If there is no console output, it may appear that you are not connected - hit ``<Enter>`` to
regenerate the Shell login dialog on the TTY in this case.


Extra Disks and NICs for your VMs
---------------------------------

If you require additional storage devices (eg "disks") inside your VM, more than the single
defined NIC, or other hardware changes, you can create VMs from the Hypervisor tooling with these
devices.  Any configuration that allows DRP to have access on one Boot/Control NIC will work.

You will need to add any additional Hypervisor network (bridges, NAT/Masqurading, etc.) necessary
to support additional NICs.


A Note about VM DNS Services
----------------------------

The above configuration has the DRP Subnet set external DNS servers for resolving
name server lookups.  This works well, but you can also set up services like Unbound
in the DRP Endpoint / KVM hypervisor and limit the DNS queries based on Unbound
rules and capabilities.  A simple example that provides pass-through recursive DNS
query and caching support for Unbound is as follows:

  ::

    # assumes package 'unbound' was installed as per above examples

    cat << EOF_UNOUND > /etc/unbound/conf.d/kvm-test.conf
    # kvm-test DRP managed virtual machine network
    server:
      interface: 192.168.124.1
      access-control: 192.168.124.0/24 allow
      access-control: ::1 allow
    EOF_UNBOUND

    systemctl restart unbound

Now you can adjust the Subnet specifiction for the DNS Option (Code 6) to set the
DNS server that your DHCP clients use to:

* ``192.168.124.1``

You can find more information on configuration and management of the *Unbound* services at:

* https://nlnetlabs.nl/projects/unbound/about/


Additional ToDos and Enhancement Ideas
--------------------------------------

The ``kvm-test`` plugin doesn't receive much love from RackN.  It works pretty well in it's limited
use case today.  The primary drawback is the Machine specifications are hard-wired in to the *Plugin
Provider* golang compiled binary.  This makes changing hardware types left to the direct Virtual
Machine creation path within KVM/Qemu/Libvirt, and not via the DRP Plugin mechanism.

If you find the KVM Test plugin useful, and are so inclined, we welcome Pull Requests enhancements
to add value and additional capabilities to the tooling.  You can find the Github repo and code
for ``kvm-test`` at:

* https://github.com/digitalrebar/provision-plugins/tree/v4/cmds/kvm-test

Some places that could use love and enhancement:

* More flexible Machine creation specifications
* A workflow designed to run in the special Self Runner to configured the DRP Endpoint with KVM configuration
* Remote API invocation of KVM hypervisors across the network, so DRP Endpoint isn't required to be on the same KVM host

There are many more small areas that can be enhanced.  `Please contact RackN <http://rackn.com/contact>`_ if
you have any thoughts or questions to make this better!


Example Machine Creation - ``virt-install`
-------------------------------------------

A simple example of using ``virt-install`` to create a Virtual Machine:

  ::

    virt-install \
      --name=winders --ram=4096 --cpu host --hvm --vcpus=2 \
      --os-type=windows --os-variant=win10 \
      --disk /var/lib/libvirt/images/winders.qcow2,size=80,bus=virtio \
      --pxe --network bridge=kvm-test,model=virtio \
      --graphics vnc,listen=127.0.0.1,password=foobar --check all=off &

This creates the Virtual machine and specifies PXE network boot.  The machine should create
and load the Sledgehammer Discovery image.  Note that when the Machine object is created
in Digital Rebar, the name will not match - it will recieve the *dname* based on the PXE
boot NIC MAC address like *d52-54-00-3b-c5-76*.

The console is availble via the VNC protocol on localhost with the password *foobar* in this
example.

