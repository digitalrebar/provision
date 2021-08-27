.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Backup and Restore

.. _rs_backup_restore:

Backup and Restore
==================


About
-----

In the v4.4.5 release and beyond, the *dr-waltool* program should be used for creating
backups of a running dr-provision system.

This backup can be used for archival purposes, restoration on failure, or movement of DRP from
one machine to another.

Backup
------

To generate a back-up, you will need the architecture appropriate *dr-waltool* on the system that
will receive the backup.  *dr-waltool* runs on linux/amd64 or darwin/amd64.  It may work on
windows/amd64, but it has not been tested.

.. note::

    Backups created with *dr-waltool* can only be used by *dr-provision* binaries of the same
    architecture and with the same or higher version. For example, a backup taken with *dr-waltool* version 4.6.x
    cannot be used by *dr-provision* 4.5.x.

With *dr-waltool* in your path, you can do the following to backup a running DRP endpoint.  *dr-waltool*
respects the same environment variables as *drpcli*.  It does not follow the .drpclirc or same command
line flags.  See *dr-waltool --help* for more information.

  ::

    export RS_ENDPOINT=https://DRP_IP:DRP_PORT # FILL THIS IN
    export RS_KEY=username:password # FILL THIS IN
    dr-waltool backup --destDir=/tmp/backups/dr-provision --artifacts


This will connect to the DRP Endpoint at address *DRP_IP* on port *DRP_PORT* using the username and
password of *username* and *password*, respectively.

It will attempt to use the specified directory, */tmp/backups/dr-provision*, and create it if not present. When creating
an initial backup with dr-waltool, --destDir must not exist.

Since *--artifacts* was specified, all content packs, plugin providers, and the files and isos
directories will be recorded as well.  If *--artifacts* is not specified, only the data objects
will be backed up.

You can also specify the endpoint and credentials on the command line with the *--endpoint* and
either the *--token* or *--key* flags. These must be set and do not have defaults like drpcli does.

.. note::

    This is an incremental backup.  The first one can take a long time as all the artifacts are
    synchronized.  Running the command again on the same directory will on retrieve differences.

Once complete, this directory can be tarred up for storage elsewhere.

  ::

    cd /tmp/backups/dr-provision
    tar -zcvf ../backup-2020-08-25.tgz *

A complete backup should include the *install.sh* and the *dr-provision.zip* used at install time, and this tarball.

.. note:: If dr-waltool is not installed on your drp endpoint you may be running an older version of
          dr-provision server, or you have installed, or upgraded using the web installer. See below
          for solutions.


Restore
-------

To restore a backup on an running system, you will need to stop DRP, extract the backup into the
active directories for the system, and restart DRP.

.. note::

    Backups created with *dr-waltool* can only be used by *dr-provision* binaries of the same
    architecture and with the same or higher version. For example, a backup taken with *dr-waltool* version 4.6.x
    cannot be used by *dr-provision* 4.5.x.

An example of this for a production installed system would like like this.

  ::

    # Assumes this is running on a systemd-based linux server
    # Assumes that backup tarball is in /root and named backup-2020-08-25.tgz
    systemctl stop dr-provision
    cd /var/lib/dr-provision
    rm -rf /var/lib/dr-provision/*  # To make sure nothing remains.
    tar -zxvf /root/backup-2020-08-25.tgz
    cd -
    systemctl start dr-provision

This will restart dr-provision at the previous state.

.. note::

  The *rm* is for completeness.  If you are just restoring a database, then you only need to
  remove the /var/lib/dr-provision/wal, /var/lib/dr-provision/secrets, and
  /var/lib/dr-provision/digitalrebar directories.


Rebuild / Move
--------------

To move a system or rebuild a system, you will need to install the version of DRP you want.  This
must be the same or newer version than your backup.  Do your installation as before.
Once it is installed, you should follow the *Restore* procedure.

.. note::

  When moving an existing system, you should turn off the Subnets on the previous system before
  starting the new DRP endpoint. Duplicate Subnets can cause networking problems.


dr-waltool Missing
------------------

If dr-waltool is missing from your drp endpoint first make sure you are using at least dr-provision version 4.5 or newer.

  ::

    drpcli info get|jq .version

Once you verify you are running at least 4.5 or newer you can grab the dr-waltool doing the following

  ::

    cd /tmp
    drpcli catalog item download drp --version=stable
    bsdtar -xzvf drp.zip
    cp bin/linux/amd64/dr-waltool /usr/local/bin/
    chmod +x /usr/local/bin/dr-waltool

.. note::

  Note the file is named .zip but for historical reasons it is actually a tar file and using "unzip" instead of bsdtar will result in issues.
