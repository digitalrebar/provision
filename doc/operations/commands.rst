.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Using DRBUP

.. _rs_drbup:

Backups Using drbup
===================
.. note::

    Please note this tool is deprecated and is no longer supported. You should be using dr-waltool. Information about that tool can be found here: :ref:`rs_backup_restore`

.. code-block:: bash

    drbup --help
    drbup backup --help
    drbup restore --help
    drbup sync --help

Backup
------

The `backup` command allows an operator to preform a backup of the Digital Rebar Platform.
At a high level what happens is depending on your level of risk dr-provision is stopped.
A tar file is made of the dr-provision root.

`backup` takes several options.

    +----------------------------------------------------------------------------------+
    |Options                                                                           |
    +==========================+=======================================================+
    | -r, -\\-risky            | Does a backup without stopping the DRP Service.       |
    +--------------------------+-------------------------------------------------------+
    | -s, -\\-source PATH      | The full path to the dr-provision home.               |
    |                          | Default: /var/lib/dr-provision                        |
    +--------------------------+-------------------------------------------------------+
    | -d, -\\-dest TEXT        | The file name to use for the drp backup. This         |
    |                          | should include the path. If no path is included       |
    |                          | cwd is assumed.                                       |
    |                          | Default: drp-backup.tar.bz2                           |
    +--------------------------+-------------------------------------------------------+
    |  -\\-take-lvm-snapshot   |  Flag to enable LVM based snapshots                   |
    +--------------------------+-------------------------------------------------------+
    |  -\\-volume-name TEXT    |  Name of the Logical Volume to snapshot. Note: This   |
    |                          |  argument is mutually inclusive with                  |
    |                          |  take_lvm_snapshot                                    |
    +--------------------------+-------------------------------------------------------+
    |  -\\-vg-name TEXT        |  Name of the Volume Group for your Volume. Note:      |
    |                          |  This argument is mutually inclusive with             |
    |                          |  take_lvm_snapshot                                    |
    +--------------------------+-------------------------------------------------------+
    | -\\-snap-mount-point TEXT|  Where to mount the LVM snapshot. Default is:         |
    |                          |  /mnt/drp-backup                                      |
    |                          |  If the path does not exist it will                   |
    |                          |  be created. Note: This argument is mutually          |
    |                          |  inclusive with take_lvm_snapshot                     |
    +--------------------------+-------------------------------------------------------+
    | -\\-snapshot-size INTEGER|  The size in M for the snapshot.                      |
    |                          |  Default: 1000                                        |
    |                          |  Note: This argument is mutually inclusive with       |
    |                          |  take_lvm_snapshot                                    |
    +--------------------------+-------------------------------------------------------+
    |  -x, -\\-exclude    TEXT |   A file name or a path (with the leading slash left  |
    |                          |   off).                                               |
    |                          |   Example: -x foo/foo1.txt -x bar -x myfile.iso       |
    |                          |   At this time only file names or paths are supported.|
    |                          |                                                       |
    +--------------------------+-------------------------------------------------------+

Restore
-------

The `restore` command allows an operator to preform a restore of the Digital Rebar Platform.
At a high level what happens is a tar file is read from a previous backup. Next its laid on
to the file system. This operation does not start the `dr-provision` service.

`restore` only accepts a couple of options

    +----------------------------------------------------------------------------------+
    |Options                                                                           |
    +==========================+=======================================================+
    |                          |                                                       |
    | -s, -\\-source TEXT      | The full path to the dr-provision backup file.        |
    |                          |                                                       |
    +--------------------------+-------------------------------------------------------+
    | -d, -\\-dest TEXT        | The full parent path to the dr-provision home.        |
    |                          | Default: /var/lib"                                    |
    +--------------------------+-------------------------------------------------------+


Sync
----

The `sync` command allows an operator to synchronize a local source (path) with either a local or remote
destination. This functionality is currently provided by rsync. If you should find the functionality
lacking in someway but its supported using rsync args you can pass in args that will be passed directly
to rsync.

    +----------------------------------------------------------------------------------+
    | Usage: drbup sync [OPTIONS]                                                      |
    | Options                                                                          |
    +=====================+============================================================+
    |  -s, -\\-source TEXT|  Source directory. Example: /var/lib/dr-provision          |
    |                     |  [required]                                                |
    +---------------------+------------------------------------------------------------+
    |  -d, -\\-dest TEXT  |  Full path for the destination. The source and dest        |
    |                     |  could be the same when doing a remote sync.               |
    |                     |  [required]                                                |
    +---------------------+------------------------------------------------------------+
    |  -l, -\\-local      |  Flag to enable a local sync. NOTE: This argument is       |
    |                     |  mutually exclusive with remote                            |
    +---------------------+------------------------------------------------------------+
    |  -r, -\\-remote     |  Flag to enable a remote sync. NOTE: This argument is      |
    |                     |  mutually exclusive with local                             |
    +---------------------+------------------------------------------------------------+
    | -\\-remote-host TEXT|  FQDN, host name, IP, or user@host for the remote host     |
    |                     |  to sync to.                                               |
    |                     |  Example: root@remote-host  OR user@192.168.1.10           |
    |                     |  Note: This argument is mutually inclusive with remote     |
    +---------------------+------------------------------------------------------------+
    |  -v, -\\-verbose    |  Prints the rsync output to stdout.                        |
    +---------------------+------------------------------------------------------------+
    |-\\-rsync-option TEXT|  Options to pass to rsync. If no options are provided      |
    |                     |  -avp is used.                                             |
    |                     |  Example:                                                  |
    |                     |  * -\\-rsync-option '-avp'                                 |
    |                     |  * -\\-rsync-option '--dry-run'                            |
    +---------------------+------------------------------------------------------------+
    | -x, -\\-exclude TEXT|  Exclude files matching PATTERN                            |
    |                     |  Example: -x "*.iso" -x "*.img"                            |
    +---------------------+------------------------------------------------------------+
    |  -\\-help           |  Show this message and exit.                               |
    +---------------------+------------------------------------------------------------+



Examples:
---------

* Create backup without stopping dr-provision

.. code-block:: bash

    drbup backup --risky --source /var/lib/dr-provision --dest /root/drp-backup.tar.bz2


.. note:: -\\-source must be a full path to the dr-provision root.


* Create backup by stopping dr-provision, take an lvm snapshot, mount it, and do our backup from the mounted snapshot

.. code-block:: bash

     drbup backup --take-lvm-snapshot --volume-name drp --vg-name mycompany-vg --size 256 --snap-mount-point /mnt/drp-backup


* Create backup and exclude the tftpboot directory from the archive

.. code-block:: bash

    drbup backup -s /var/lib/dr-provision -d /srv/drp-backup.tar.bz2 -x tftpboot


* Restore a backup stored in /srv/backups/drp-backup.tar.bz2 to /var/lib/dr-provision

.. code-block:: bash

    drbup restore --source /srv/backups/drp-backup.tar.bz2 --dest /var/lib/dr-provision


* Sync drp home to remote site

.. code-block:: bash

    drbup sync --remote -s /var/lib/dr-provision -d /var/lib/ --remote-host admin@remote-host.internal


* Pass args directly to rsync when doing a sync.

.. code-block:: bash

    drbup sync --remote -s /var/lib/dr-provision -d /var/lib --remote-host root@remote-host --rsync-option '-avp' --rsync-option '--dry-run'

..

