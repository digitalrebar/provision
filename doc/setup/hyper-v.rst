.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Setup

.. _rs_setup_hyperv:

Hyper-V Setup Instructions
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Overview
--------

This document will help define one possible method to setup a Hyper-V environment to test Digital Rebar Provision with. This assumes you are running Hyper-V on a Windows 10 Professional desktop, but the process would be similar for a dedicated Hyper-V Hypervisor (with the exception of the NAT section). This is not the only way Hyper-V can be configured for use with Digital Rebar, this is only one example.

Please feel free to provide any feedback or alternative setup methods.


Architecture
------------

The goal is to enable the Hyper-V NAT feature on the default "Internal Switch" and then attach a new VM which will run DRP to that. We refer to this as the "DRP VM", this would run any Linux based OS (for this example I will install Debian 10) and would have two network interfaces. The second network interface would be attached to a "Target Switch" which is a isolated Hyper-V switch. Then we would have a VM which can be used as a target for DRP, which we refer to as the "Target VM". This VM has no OS installed, and a single network connection to the "Target Switch". Here is a simplified diagram: ::

     +--------------+
     |              |
     |   Internal   |
     |    Switch    |
     |              |
     +--------------+
             |
             |
             |
             |
             |
     +--------------+
     |     eth0     | (ex: 192.168.121.11)
     |              |
     |    DRP VM    |
     |              |
     |     eth1     | (ex: 10.20.30.1)
     +--------------+
             |
             |
             | Target Switch
             |
             |
     +--------------+
     |     eth0     | (ex: 10.20.30.4)
     |              |
     |  Target VM   |
     |              |
     +--------------+


This method has a few advantages: 

* You can have many target VMs and they are easy to create/destroy
* You do not need to have the DRP VM attached to a specific network interface, meaning no configuration changes are needed if your laptop uses different Ethernet or WiFi adapters depending on location.
* The DRP endpoint will have Internet access, and can selectively provide Internet access to the target VMs using standard Linux tools.
* We don't need anything on the host beyond Hyper-V and PuTTY to access hosts.

Prerequisites
-------------

You will need:

* Windows 10 Professional, Enterprise, or Education
* PowerShell (Typically included with Windows 10)
* `PuTTY <https://www.chiark.greenend.org.uk/~sgtatham/putty/latest.html>`_
* `Debian 10 Install Media with Firmware <https://cdimage.debian.org/cdimage/unofficial/non-free/cd-including-firmware/10.4.0+nonfree/amd64/iso-dvd/>`_

Hyper-V has `specific hardware and BIOS/Firmware requirements <https://docs.microsoft.com/en-us/virtualization/hyper-v-on-windows/reference/hyper-v-requirements>`_ as defined by Microsoft. If your system meets those requirements, you can then follow Microsoft's `documentation to Install Hyper-V on your system <https://docs.microsoft.com/en-us/virtualization/hyper-v-on-windows/quick-start/enable-hyper-v>`_ for more details.

Create the Internal Network and NAT
-----------------------------------

Next we will need to setup a new ``Internal Network`` switch and add Hyper-V NAT to it. `Microsoft has documentation on this process <https://docs.microsoft.com/en-us/virtualization/hyper-v-on-windows/user-guide/setup-nat-network>`_ but I'll simplify the commands here as much as possible.

The goal is to setup the IPv4 network ``192.168.121.0/24`` as the NAT network and ``192.168.121.1`` as the default gateway on the ``Internal Network`` switch. If you do not want to use those addresses because they will conflict with networks you are normally on, you will need to change them here. However, those defaults should be fine for most cases. Note that Hyper-V does not provide DHCP or DNS services, but those can be provided by Linux so it doesn't really matter in this use case.

Here is the process to create the switch and NAT:

#. Open a PowerShell as Administrator
#. Run the command ``Get-NetAdapter``
#. Run the command ``New-VMSwitch -SwitchName "Internal Network" -SwitchType Internal``
#. Run the command ``Get-NetAdapter`` again and note the new line. The we need the IfIndex of the new line as this is your ner switch.
#. Run the command ``New-NetIPAddress -IPAddress 192.168.121.1 -PrefixLength 24 -InterfaceIndex <ifIndex>`` replacing ``<IfIndex>`` with the number we found in the previous step.
#. The the command ``New-NetNat -Name "InternalNetworkNAT" -InternalIPInterfaceAddressPrefix 192.168.121.0/24``

As long as no errors are encountered, you can now close the Elevated PowerShell prompt.

These commands may also work on dedicated Hyper-V Hypervisor installs, but check with your network and/or sysadmin teams for more information about your specific environment before implementing NAT in a shared environment.

Creating the rest of the environment using the Hyper-V Manager
--------------------------------------------------------------

Now we can build the rest of the environment using the Hyper-V Manager as opposed to the PowerShell CLI. This makes some tasks much easier. First, let's create the target switch:

#. Open the Hyper-V Manager
#. Select your system from the list on the left
#. Click on "Virtual Switch Manager" in the right most panel
#. Click on "New Virtual Switch" on the left panel, and you are presented with a list of options on the right panel. We want either an "Internal" or "Private" switch. If you make an "Internal" switch it is possible to assign an IP address to the host computer to this switch and then access target VMs directly using tools like PuTTY. A "Private" switch can only be accessed by other VMs attached to that switch. If unsure, choose "Private" and press "Create Virtual Switch".
#. A new dialog will appear allowing you to name the switch and set options. Give your switch a descriptive name (in this case I called mine "drp stable targets") and press "OK".
#. You can now close the Hyper-V Switch Manager

Now we need to create the Virtual Machine which will run DRP. This is probably the most complicated sequence.

#. From the Hyper-V Manager click on "New..." from the left most panel and choose "Virtual Machine"
#. Press "Next" on the Before You Begin screen
#. Give your VM a name (for example "DRP Endpoint") and press Next
#. Choose Generation 2 and press Next
#. The default of 1024 MB of RAM is more than sufficient for DRP. Uncheck Dynamic Memory and press Next
#. Choose the "Internal Network" switch we created earlier from the drop down for the network connection and press Next
#. You will need to create a virtual hard disk for this VM. The Operating installation will be small, about 2GB, but boot environments can be large. For example, if you plan on using CentOS 8 that requires about 15G of space. Choose what makes sense for your system (60G is what I used in this example) and press Next.
#. Choose "Install an Operating System from a bootable image file" and then select the Debian 10 Firmware ISO you downloaded from the prerequisites section and press Next.
#. You will be presented with a summary, if everything looks as expected press Finish.
#. The Virtual Machine will be created and then you will return to the Hyper-V Manager
#. Select the new Virtual Machine from the middle panel and choose "Settings" from the right panel
#. When the settings dialog appears, select "Security" from the left panel and uncheck "Enable Secure Boot" from the right panel and press OK. (Note: you can use secure boot, but cannot use the Windows option when trying to use DRP since we boot a Linux environment. You can use the "Microsoft UEFI Certificate Authority" option for the VMs if Secure Boot is desired.)
#. From the Hyper-V Manager, press the "Connect" option on the right panel. Then you can press "Start" on the new Virtual Machine Connection window. This will begin the Debian install process.
#. After a few seconds, the Debian Installer Boot Menu should appear. Press Enter to continue.
#. Choose your language, and press "Continue"
#. Choose your location, and press "Continue"
#. Choose your keyboard layout, and press "Continue"
#. The installer will detect the virtual media and load some additional components. It will eventually try to detect networking and timeout. This is expected because there is no DHCP services provided by Hyper-V. Press "Continue"
#. Select "Configure Network Manually" and press "Continue"
#. Enter the IP Address as ``192.168.121.11`` and press Enter
#. The default subnet mask of ``255.255.255.0`` is correct, just press Enter
#. The default gateway of ``192.168.121.1`` is correct, press Enter
#. For nameservers, you need to specify some that will work almost anywhere. I recommend ``1.1.1.1 8.8.8.8 9.9.9.9`` but you can also use your corporate DNS servers if needed. Enter whatever will work for your environment and press Enter.
#. For a hostname, input what you would like and press Enter (do not use spaces, dash is OK)
#. For the domain name, you can leave it blank and press Enter
#. On the next screen you will be prompted for the root password. Simply leave the values blank and click Continue (this will automatically enable sudo for the user account we are about to create)
#. Next enter your name and press Enter
#. A username will be generated, you can accept this as is or replace it and press Enter
#. You will then be prompted for a password. The password must meet minimum complexity requirements, you will be told if it does not. Type your desired password in both fields and press Continue
#. Choose your timezone, and press Continue
#. The disk configuration tool will start up, it is recommended that you choose ``Guided - use entire disk and setup LVM`` and press Continue
#. There should only be the single disk selected, press Continue
#. Choose the default of ``All files in one partition`` and press Continue
#. Select ``Yes`` on the partition screen and press Continue
#. The default amount of disk space to use is the maximum, press Continue
#. You will be asked if you want to force UEFI installation, select Yes and press Continue
#. You will then be presented with a summary of disk configuration, choose Yes and press Continue
#. At this point the disk will be configured and the base system installed, it should only take a few moments
#. Once the base install is complete, you will be asked if there are other media you wish to scan. Select "No" and click Continue
#. Choose your country for mirror selection and press Continue
#. The default mirror is usually acceptable, press Continue
#. You hopefully do not have any HTTP Proxy information, so just press Continue when prompted (if you do require a proxy you will not be able to update packages if your HTTP proxy is unavailable)
#. The package manager will download data from the mirror and prepare to apply updates and additional software. This should only take a few moments
#. When asked if you want to participate in the survey, choose whichever option you like and press Continue
#. You will then be asked for software packages to be installed. Uncheck everything and then check "SSH Server" and "standard system utilities" and press Continue
#. The additional software will download and install, this should only take a few moments
#. You will then be told the installation is complete, press Continue to reboot into the new system
#. Within a few seconds you should be at a default login screen, which looks like this.
#. Now let's add the second Network interface to the VM connected to the "Target Switch". From the Virtual Machine Connection window, go to File and the Settings. The "Add Hardware" panel will open in the settings screen by default. Select "Network Adapter" and press "Add".
#. Select the "DRP Targets" switch we created earlier from the pulldown and then press "OK"
#. You can now close the Virtual Machine Connection windows (the VM will remain running)
#. At this point, you should be able to connect to the instance via PuTTY which will make cut and paste much easier. Open PuTTY and connect to ``192.168.121.11`` and login with the account you created. 
#. Once you login, let's install some additional tools with ``sudo apt update && sudo apt install -y iptables unbound nano git curl bsdtar p7zip-full``
#. Now we can configure the 2nd network interface. Run ``sudo nano /etc/interfaces.d/eth1`` and input the data shown in :ref:`interfaces.d-eth1`
#. You can adjust the IP address and netmask to your taste. The interface is completely isolated if your switch was configured to be Private, so no need to worry about IP address collisions. Only the DRP Endpoint and the Targets will be able to access it. You can then save with Ctrl+O followed by Enter and then quit by pressing Ctrl+X
#. Enable the second network interface by running ``sudo ifup eth1``
#. Now we can configure unbound to provide DNS for your private network. Edit the configuration file by running ``sudo nano /etc/unbound/unbound.conf.d/targets.conf`` and entering the following shown in :ref:`unbound-target`
#. Again, adjust your IP information to match what you put in the network configuration in previous steps if necessary. Save your changes like before with Ctrl+W followed by Enter. Then quit the editor with Ctrl+X. Then restart unbound to read the new configuration with ``sudo systemctl restart unbound``
#. Now we can install DRP with the following command: ``curl -fsSL get.rebar.digital/stable | bash -s -- install --systemd``
#. Next we get the discovery and sledgehammer environments downloaded: ``drpcli bootenvs uploadiso sledgehammer``
#. Then we configure the default workflow for new machines to use the discovery workflow: ``drpcli prefs set defaultWorkflow discover-base unknownBootEnv discovery``
#. Now you should be able to connect to the DRP endpoint at https://192.168.122.11:8092 and login to the UI
#. From here, click on the Subnets option under Networking on the left hand panel, and then press the "Add" button at the top of the right panel. You will be presented with a list of "eth0" and "eth1" with details. Click the "Use Interface" on the card with the eth1 details.
#. A details page will appear, scroll down to the DNS server option and set that to the same IP address as the default gateway. Then scroll to bottom of the page (the default options are fine) and press "Add". This will cause DRP to now serve DHCP leases to the VMs connected to the Target Network switch.

.. _interfaces.d-eth1:

Contents of ``/etc/interfaces.d/eth1``
--------------------------------------

::

 auto eth1
  iface eth1 inet static
     address 10.20.30.1
     netmask 24

.. _unbound-target:

Contents of ``/etc/unbound/unbound.conf.d/target.conf``
-------------------------------------------------------

::

  server:
    interface: 10.20.30.1
    access-control: 10.20.30.0/24 allow
    access-control: ::1 allow

Creating DRP Target VMs using Hyper-V Manager
---------------------------------------------

In order to effectively use DRP, you will need at least one target VM. The process is similar to before, but has much less steps:

#. Open the Hyper-V Manager
#. Select "New" then "Virtual Machine" from the actions panel
#. Press "Next" on the "Before you begin" screen
#. Give the machine a name (this is only shown in Hyper-V) and press Next
#. Select Generation 2 and press Next
#. Choose an amount of memory for your target machine and press next. Note that CentOS 8 requires at least 2048MB, Ubuntu 20.04 requires at least 3172MB, and Windows requires at least 4096MB.
#. For the network connection choose your "DRP Targets" switch and press Next
#. For the new virtual hard disk, choose storage size and press next. 40GB is reasonable for most test instances.
#. For the installation option screen, choose "Install and operating system from a network-based installation server" and press Finish
#. Once the VM is created, click on it in Hyper-V manager and select "Settings" from the action pane.
#. Click on the "Security" section. If you have a DRP license, you can use Secure Boot but must change the template to "Microsoft UEFI Certificate Authority" (even if you are planning on running a Windows VM because the discovery process boots Linux). If you do not have a DRP license, uncheck "Enable Secure Boot" and press OK.
#. At this point you can start the VM. It should boot, obtain a DHCP lease from DRP, boot into the sledgehammer discovery image, and appear in the machines list on the UI. You can repeat this process for as many VMs as your system resources can support.

Advanced Networking
-------------------

At this point, everything will function as expected. The VMs themselves will not be able to communicate with your host system or the Internet by default. If you wish to make this possible, you can run the following commands from the DRP Endpoint via PuTTY::

  sudo sysctl net.ipv4.ip_forward=1
  sudo iptables -t nat -A POSTROUTING -s 10.20.30.0/24 -j MASQUERADE

These commands will allow for network address translation of the traffic from the target VMs network to reach both the host and the Internet. These commands are not persistent and will need to be re-run each time you reboot the DRP endpoint VM. DNS will function and resolve names correctly even if these commands are not run due to the configuration of unbound we did during setup.