.. Copyright (c) 2021 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Platform documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Platform; Image Deploy Storage Examples

.. _rs_imagedeploy_storage:

Image Deploy Storage Examples
=============================

The Digital Rebar Platform (DRP) plugin named `image-deploy` can handle very complex
disk storage partitioning setups.  If no configuration is provided in the `curtin/partitions`
Param, then a default storage partitioning scheme will be used on the first bootable
disk found.

To change the behavior of the storage partitioning scheme, the operator must specify the
setup of the storage partitioning in the `curtin/partitions` param.  The setup is often
complex, and will generally require a fair amount of trial and error before it is
working successfully.

This document outlines several **example** configurations as potential starting points
for building your own customized paritioning scheme.  Please note that each of the
below examples may have **very specific** requiremens on the number, size, and type
of disks for the example to work correctly.  If these pre-conditions are not
carefully observed, you are almost guaranteed to not build a successful filesystem
structure.

.. note:: When you build a successful new partitioning layout, please consider forwarding
          to RackN Support for inclusing in this document.  Or better yet, open a Pull
          Request to submit the new setup!!

Additional documentation on the Image Deploy capabilities within DRP can be found at the
following locations:

  * Full Image Deploy Plugin documentation: :ref:`rs_cp_image_deploy`
  * :ref:`rs_imagedeploy`
  * Also check the various :ref:`rs_knowledge_base`


.. _id_stor_1_windows_single:

Windows Single Disk
-------------------

This example is for a *wim-like* root fs for windows, but other custom partitions could be used.

Notes:

  * Uses a single disk
  * Legacy BIOS with MSDOS partition format
  * Specifies 139 GB ``C:\`` drive
  * Adjust ``disk0-part1`` size directive as appropriate for your disks

.. code-block:: yaml

    storage:
      version: 1
      config:
        - id: disk0
          type: disk
          ptable: msdos
          path: /dev/sda
          name: main_disk
          wipe: superblock
        - id: disk0-part1
          type: partition
          number: 1
          device: disk0
          size: 139G
          flag: boot
        - id: disk0-part1-format-root
          type: format
          fstype: ntfs
          quiet: True
          volume: disk0-part1
        - id: disk0-part1-mount-root
          type: mount
          path: /
          device: disk0-part1-format-root

.. _id_stor_1_linux_single_paritioned:

Linux Single Disk, Separate Partitions
--------------------------------------

Linux LVM configuration with separate ``root``, ``boot``, ``var``, and ``swap`` partitions.

Notes:

  * Single disk configuration
  * Separate filesystem partitions for ``root``, ``boot``, ``var``, and ``swap``
  * Formats disks as ``XFS`` type filesystems
  * Adjust sizes in each partition stanza

.. code-block:: yaml

   storage:
     version: 1
     config:
       - id: sda
         type: disk
         ptable: gpt
         path: /dev/sda
         name: main_disk
         grub_device: true
       - id: bios_boot_partition
         type: partition
         size: 1MB
         device: sda
         flag: bios_grub
       - id: boot_part
         name: boot_part
         type: partition
         size: 8GB
         device: sda
         flag: boot
       - id: lvm_part
         type: partition
         size: 40G
         device: sda
       - id: volgroup1
         name: volgroup1
         type: lvm_volgroup
         devices:
           - lvm_part
       - id: root_part
         name: root_part
         size: 10G
         type: lvm_partition
         volgroup: volgroup1
       - id: swap_part
         name: swap_part
         type: lvm_partition
         volgroup: volgroup1
         size: 4G
       - id: var_part
         name: var_part
         type: lvm_partition
         volgroup: volgroup1
       - id: boot_fs
         type: format
         fstype: xfs
         volume: boot_part
       - id: root_fs
         name: storage
         type: format
         fstype: xfs
         volume: root_part
       - id: var_fs
         name: storage
         type: format
         fstype: xfs
         volume: var_part
       - id: swap_fs
         name: storage
         type: format
         fstype: swap
         volume: swap_part
       - id: root_mount
         type: mount
         path: /
         device: root_fs
       - id: boot_mount
         type: mount
         path: /boot
         device: boot_fs
       - id: var_mount
         type: mount
         path: /var
         device: var_fs
       - id: swap_mount
         type: mount
         path: swap
         device: swap_fs


.. _id_stor_1_linux_with_lvm:

Linux, LVM Multi-Partition Layout
---------------------------------

Linux LVM multi partition layout.

Notes:

  * Single disk
  * Logical Volume Manager (LVM) with multi-partition layout
  * Partition table format MSDOS

.. code-block:: yaml

   storage:
     version: 1
     config:
       - id: sda
         type: disk
         ptable: msdos
         path: /dev/sda
         name: main_disk
       - id: sda1
         type: partition
         size: 3GB
         device: sda
         flag: boot
       - id: sda_extended
         type: partition
         size: 5G
         flag: extended
         device: sda
       - id: sda2
         type: partition
         size: 2G
         flag: logical
         device: sda
       - id: sda3
         type: partition
         size: 3G
         flag: logical
         device: sda
       - id: volgroup1
         name: vg1
         type: lvm_volgroup
         devices:
           - sda2
           - sda3
       - id: lvmpart1
         name: lv1
         size: 1G
         type: lvm_partition
         volgroup: volgroup1
       - id: lvmpart2
         name: lv2
         type: lvm_partition
         volgroup: volgroup1
       - id: sda1_root
         type: format
         fstype: ext4
         volume: sda1
       - id: lv1_fs
         name: storage
         type: format
         fstype: ext4
         volume: lvmpart1
       - id: lv2_fs
         name: storage
         type: format
         fstype: ext3
         volume: lvmpart2
       - id: sda1_mount
         type: mount
         path: /
         device: sda1_root
       - id: lv1_mount
         type: mount
         path: /srv/data
         device: lv1_fs
       - id: lv2_mount
         type: mount
         path: /srv/backup
         device: lv2_fs


.. _id_stor_1_ubuntu_two_disk_md:

Linux - Ubuntu with Software RAID
---------------------------------

Linux Ubuntu system utilizing Software RAID (``md``) across two disks.

Notes:

  * Requires two disks
  * Software Raid (``md``) volumes across ``sda`` and ``sdb`` partitions
  * Supports Legacy BIOS Boot mode
  * Supports UEFI Boot mode
  * Specifies backup EFI boot partitions
  * Separate ``root``, ``boot``, EFI, ``var``, and ``swap`` filesystems
  * Formats as ``XFS`` type filesystems

.. code-block:: yaml

   storage:
     version: 1
     config:
     - grub_device: 1
       id: sda
       name: main_disk
       path: /dev/sda
       ptable: gpt
       type: disk wipe: superblock - device: sda
       id: boot_efi_part
       size: 200MB
       type: partition
     - device: sda
       id: boot_part
       name: boot_part
       size: 8GB
       type: partition
     - device: sda
       id: sda4
       size: 30GB
       type: partition
     - device: sda
       id: bios_boot_partition
       size: 1MB
       type: partition
     - id: sdb
       name: second_disk
       path: /dev/sdb
       ptable: gpt
       type: disk
       wipe: superblock
     - device: sdb
       id: backup_boot_efi_part
       size: 200MB
       type: partition
     - device: sdb
       id: backup_boot_part
       name: backup_boot_part
       size: 8GB
       type: partition
     - device: sdb
       id: sdb4
       size: 30GB
       type: partition
     - device: sdb
       id: backup_bios_boot_partition
       size: 1MB
       type: partition
     - devices:
       - sda4
       - sdb4
       id: mddevice
       name: md1
       raidlevel: 1
       type: raid
     - devices:
       - boot_part
       - backup_boot_part
       id: mdboot
       name: md0
       raidlevel: 1
       type: raid
     - fstype: ext4
       id: md_root
       type: format
       volume: mddevice
     - devices:
       - mddevice
       id: volgroup1
       name: volgroup1
       type: lvm_volgroup
     - id: root_part
       name: root_part
       size: 10G
       type: lvm_partition
       volgroup: volgroup1
     - id: swap_part
       name: swap_part
       size: 4G
       type: lvm_partition
       volgroup: volgroup1
     - id: var_part
       name: var_part
       type: lvm_partition
       volgroup: volgroup1
     - flag: boot
       fstype: xfs
       id: boot_fs
       type: format
       volume: mdboot
     - fstype: vfat
       id: boot_efi_fs
       type: format
       volume: boot_efi_part
     - fstype: vfat
       id: backup_boot_efi_fs
       type: format
       volume: backup_boot_efi_part
     - fstype: xfs
       id: root_fs
       name: storage
       type: format
       volume: root_part
     - fstype: xfs
       id: var_fs
       name: storage
       type: format
       volume: var_part
     - fstype: swap
       id: swap_fs
       name: storage
       type: format
       volume: swap_part
     - device: root_fs
       id: root_mount
       path: /
       type: mount
     - device: boot_fs
       id: boot_mount
       path: /boot
       type: mount
     - device: boot_efi_fs
       id: boot_efi_mount
       path: /boot/efi
       type: mount
     - device: var_fs
       id: var_mount
       path: /var
       type: mount
     - device: swap_fs
       id: swap_mount
       path: swap
       type: mount


.. _id_stor_1_centos_standard:

Linux - CentOS Standard Install
-------------------------------

Mimics the standard *Anaconda* CentOS partitioning scheme.

Notes:

  * Single disk
  * Supports BIOS and UEFI Boot modes
  * Uses LVM volumes
  * Separate ``root``, ``boot``, ``var``, and ``swap`` partitions
  * ``boot`` is 8GB ``XFS`` (not contained in the LVM partition)
  * ``root`` (10GB), ``var`` (26GB), and ``swap`` (4GB) are contained in a 40 GB LVM partition
  * Formats to ``XFS`` filesystem type

.. code-block:: yaml

   storage:
     version: 1
     config:
     - grub_device: true
       id: sda
       name: main_disk
       path: /dev/sda
       ptable: gpt
       type: disk
     - device: sda
       flag: bios_grub
       id: bios_boot_partition
       size: 1024MB
       type: partition
     - device: sda
       id: boot_efi_part
       size: 1024MB
       type: partition
     - device: sda
       id: boot_part
       name: boot_part
       size: 8GB
       type: partition
     - device: sda
       id: lvm_part
       size: 40G
       type: partition
     - devices:
       - lvm_part
       id: volgroup1
       name: volgroup1
       type: lvm_volgroup
     - id: root_part
       name: root_part
       size: 10G
       type: lvm_partition
       volgroup: volgroup1
     - id: swap_part
       name: swap_part
       size: 4G
       type: lvm_partition
       volgroup: volgroup1
     - id: var_part
       name: var_part
       type: lvm_partition
       volgroup: volgroup1
     - fstype: xfs
       id: boot_fs
       type: format
       volume: boot_part
     - fstype: vfat
       id: boot_efi_fs
       type: format
       volume: boot_efi_part
     - fstype: xfs
       id: root_fs
       name: storage
       type: format
       volume: root_part
     - fstype: xfs
       id: var_fs
       name: storage
       type: format
       volume: var_part
     - fstype: swap
       id: swap_fs
       name: storage
       type: format
       volume: swap_part
     - device: root_fs
       id: root_mount
       path: /
       type: mount
     - device: boot_fs
       id: boot_mount
       path: /boot
       type: mount
     - device: boot_efi_fs
       id: boot_efi_mount
       path: /boot/efi
       type: mount
     - device: var_fs
       id: var_mount
       path: /var
       type: mount
     - device: swap_fs
       id: swap_mount
       path: swap
       type: mount


.. _id_stor_1_centos_raid:

Linux - CentOS with RAID
------------------------

CentOS Linux with Software RAID configuration.

Notes:

  * Two disks
  * Uses Software RAID (LVM)
  * Supports BIOS and UEFI boot modes
  * Formats to XFS filesystem types
  * Creates a backup EFI partition
  * Formats to ``EXT4`` filesystem type
  * Creates two software raid mirror volumes

.. code-block:: yaml

   storage:
     version: 1
     config:
     - grub_device: 1
       id: sda
       name: main_disk
       path: /dev/sda
       ptable: gpt
       type: disk
       wipe: superblock
     - device: sda
       flag: boot
       id: boot_efi_part
       size: 200MB
       type: partition
     - device: sda
       id: boot_part
       name: boot_part
       size: 8GB
       type: partition
     - device: sda
       id: sda4
       size: 30GB
       type: partition
     - device: sda
       flag: bios_grub
       id: bios_boot_partition
       size: 1MB
       type: partition
     - id: sdb
       name: second_disk
       path: /dev/sdb
       ptable: gpt
       type: disk
       wipe: superblock
     - device: sdb
       flag: boot
       id: backup_boot_efi_part
       size: 200MB
       type: partition
     - device: sdb
       id: backup_boot_part
       name: backup_boot_part
       size: 8GB
       type: partition
     - device: sdb
       id: sdb4
       size: 30GB
       type: partition
     - device: sdb
       flag: bios_grub
       id: backup_bios_boot_partition
       size: 1MB
       type: partition
     - devices:
       - sda4
       - sdb4
       id: mddevice
       name: md1
       raidlevel: 1
       type: raid
     - devices:
       - boot_part
       - backup_boot_part
       id: mdboot
       name: md0
       raidlevel: 1
       type: raid
     - fstype: ext4
       id: md_root
       type: format
       volume: mddevice
     - devices:
       - mddevice
       id: volgroup1
       name: volgroup1
       type: lvm_volgroup
     - id: root_part
       name: root_part
       size: 10G
       type: lvm_partition
       volgroup: volgroup1
     - id: swap_part
       name: swap_part
       size: 4G
       type: lvm_partition
       volgroup: volgroup1
     - id: var_part
       name: var_part
       type: lvm_partition
       volgroup: volgroup1
     - flag: boot
       fstype: xfs
       id: boot_fs
       type: format
       volume: mdboot
     - fstype: vfat
       id: boot_efi_fs
       type: format
       volume: boot_efi_part
     - fstype: vfat
       id: backup_boot_efi_fs
       type: format
       volume: backup_boot_efi_part
     - fstype: xfs
       id: root_fs
       name: storage
       type: format
       volume: root_part
     - fstype: xfs
       id: var_fs
       name: storage
       type: format
       volume: var_part
     - fstype: swap
       id: swap_fs
       name: storage
       type: format
       volume: swap_part
     - device: root_fs
       id: root_mount
       path: /
       type: mount
     - device: boot_fs
       id: boot_mount
       path: /boot
       type: mount
     - device: boot_efi_fs
       id: boot_efi_mount
       path: /boot/efi
       type: mount
     - device: var_fs
       id: var_mount
       path: /var
       type: mount
     - device: swap_fs
       id: swap_mount
       path: swap
       type: mount


.. _id_stor_1_centos_raid_and_nvme:

Linux - CentOS with Software RAID on NVMe Drives
------------------------------------------------

CentOS Linux with Software RAID (LVM) across 5 NVMe drives.

Notes:

  * Requires 2 standard disks for boot, root, swap partitions
  * Requires 5 NVMe disks for RAID 5 ``var`` partition
  * Supports both Legacy and UEFI Boot modes
  * Builds Software RAID (LVM) volumes across 5 NVMe disks
  * Specifies ``XFS`` filesystem type
  * ``root``, ``boot``, and ``swap`` are located on ``sda`` and ``sda`` mirror
  * ``var`` is located on RAID 5 volume on the NVMe disks

.. code-block:: yaml

   storage:
     version: 1
     config:
     - grub_device: 1
       id: sda
       name: main_disk
       path: /dev/sda
       ptable: gpt
       type: disk
       wipe: superblock
     - device: sda
       flag: boot
       id: boot_efi_part
       size: 200MB
       type: partition
     - device: sda
       id: boot_part
       name: boot_part
       size: 8GB
       type: partition
     - device: sda
       id: sda4
       size: 210GB
       type: partition
     - device: sda
       flag: bios_grub
       id: bios_boot_partition
       size: 1MB
       type: partition
     - id: sdb
       name: second_disk
       path: /dev/sdb
       ptable: gpt
       type: disk
       wipe: superblock
     - device: sdb
       flag: boot
       id: backup_boot_efi_part
       size: 200MB
       type: partition
     - device: sdb
       id: backup_boot_part
       name: backup_boot_part
       size: 8GB
       type: partition
     - device: sdb
       id: sdb4
       size: 210GB
       type: partition
     - device: sdb
       flag: bios_grub
       id: backup_bios_boot_partition
       size: 1MB
       type: partition
     - id: nvme0n1
       name: nvme0n1
       path: /dev/nvme0n1
       ptable: gpt
       type: disk
       wipe: superblock
     - device: nvme0n1
       id: nvme1
       size: 1450GB
       type: partition
     - id: nvme1n1
       name: nvme1n1
       path: /dev/nvme1n1
       ptable: gpt
       type: disk
       wipe: superblock
     - device: nvme1n1
       id: nvme2
       size: 1450GB
       type: partition
     - id: nvme2n1
       name: nvme2n1
       path: /dev/nvme2n1
       ptable: gpt
       type: disk
       wipe: superblock
     - device: nvme2n1
       id: nvme3
       size: 1450GB
       type: partition
     - id: nvme3n1
       name: nvme3n1
       path: /dev/nvme3n1
       ptable: gpt
       type: disk
       wipe: superblock
     - device: nvme3n1
       id: nvme4
       size: 1450GB
       type: partition
     - id: nvme4n1
       name: nvme4n1
       path: /dev/nvme4n1
       ptable: gpt
       type: disk
       wipe: superblock
     - device: nvme4n1
       id: nvme5
       size: 1450GB
       type: partition
     - devices:
       - sda4
       - sdb4
       id: mddevice
       name: md1
       raidlevel: 1
       type: raid
     - devices:
       - boot_part
       - backup_boot_part
       id: mdboot
       name: md0
       raidlevel: 1
       type: raid
     - devices:
       - nvme1
       - nvme2
       - nvme3
       - nvme4
       - nvme5
       id: mddata
       name: md2
       raidlevel: 5
       type: raid
     - fstype: ext4
       id: md_root
       type: format
       volume: mddevice
     - fstype: xfs
       id: md_data
       type: format
       volume: mddata
     - devices:
       - mddevice
       id: volgroup1
       name: sysvg
       type: lvm_volgroup
     - devices:
       - mddata
       id: volgroup2
       name: varvg
       type: lvm_volgroup
     - id: root_part
       name: rootlv
       size: 20G
       type: lvm_partition
       volgroup: volgroup1
     - id: swap_part
       name: swaplv
       size: 8G
       type: lvm_partition
       volgroup: volgroup1
     - id: tmp_part
       name: tmplv
       size: 2G
       type: lvm_partition
       volgroup: volgroup1
     - id: home_part
       name: homelv
       size: 1G
       type: lvm_partition
       volgroup: volgroup1
     - id: var_part
       name: var_part
       type: lvm_partition
       volgroup: volgroup2
     - flag: boot
       fstype: xfs
       id: boot_fs
       type: format
       volume: mdboot
     - fstype: vfat
       id: boot_efi_fs
       type: format
       volume: boot_efi_part
     - fstype: vfat
       id: backup_boot_efi_fs
       type: format
       volume: backup_boot_efi_part
     - fstype: xfs
       id: root_fs
       name: storage
       type: format
       volume: root_part
     - fstype: xfs
       id: tmp_fs
       name: storage
       type: format
       volume: tmp_part
     - fstype: xfs
       id: home_fs
       name: storage
       type: format
       volume: home_part
     - fstype: xfs
       id: var_fs
       name: storage
       type: format
       volume: var_part
     - fstype: swap
       id: swap_fs
       name: storage
       type: format
       volume: swap_part
     - device: root_fs
       id: root_mount
       path: /
       type: mount
     - device: tmp_fs
       id: tmp_mount
       path: /tmp
       type: mount
     - device: home_fs
       id: home_mount
       path: /home
       type: mount
     - device: boot_fs
       id: boot_mount
       path: /boot
       type: mount
     - device: boot_efi_fs
       id: boot_efi_mount
       path: /boot/efi
       type: mount
     - device: var_fs
       id: var_mount
       path: /var
       type: mount
     - device: swap_fs
       id: swap_mount
       path: swap
       type: mount

