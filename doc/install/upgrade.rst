.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Upgrade

.. _rs_upgrade:

Upgrade and Downgrade DRP
~~~~~~~~~~~~~~~~~~~~~~~~~

Upgrading Digital Rebar Platform (DRP) is generally pretty straight forward.  The
``dr-provision`` binary is simply replaced with a newer version.  Upgrades can be
performed on a running system with DRP provided tools, via the Portal, or by stopping
the ``dr-provision`` service, replacing the binary, and staring it back up again.

However, note that there are a few general rules to consider for upgrades:

  * always upgrade the DRP service first
  * upgrade content packs and plugins after the DRP service upgrade

.. warning:: We HIGHLY SUGGEST that you backup your existing install as a safey measure.


.. _rs_backup_instructions:

Backup
======

It's always a good policy to backup any important data, configuration, and
content information that may be related to an application before an upgrade.
We strongly encourage you to backup your content prior to doing any upgrade activity.


Additional References
---------------------

RackN has some basic backup management process and scripts that it maintains.  Please
see the following references for these:

  * 4.4.5 and newer See: :ref:`rs_backup_restore`
  * 4.4.4 and older: Stop ``dr-provision`` and make a tarball containing the complete
    contents of the `--base-root` (usually /var/lib/dr-provision)

Isolated Install
----------------

For "isolated" modes (eg. originally installed with something like
``install.sh | bash -s -- install --isolated``) , perform the following tasks:

#. log in to your Provision server as the user you performed the original install as
#. copy the drp-data directory to a backup location:
   ::

     D=`date +%Y%m%d-%H%M%S`
     cp -r drp-data drp-data.backup.$D


Production Install
------------------

For "production" install modes (no ``--isolated`` flag provided to ``install.sh``), perform the following tasks

#. log in to  your Digital Rebar Provision server
#. for DRP version 3.0.5 or older:
   ::

     D=`date +%Y%m%d-%H%M%S`
     mkdir backups
     sudo cp -r /var/lib/dr-provision backups/dr-provision.backup.$D
     sudo cp -r /var/lib/tftpboot backups/dr-provision.backup.$D

#. for DRP version 3.1.0 or newer:
   ::

     D=`date +%Y%m%d-%H%M%S`
     mkdir backups
     sudo cp -r /var/lib/dr-provision backups/dr-provision.backup.$D


Upgrade Steps
=============

Isolated or Production Installed Modes
--------------------------------------

The basic steps are the same for both Isolated and Production install modes, the only
primary difference is adding the ``--isolated`` flag for the ``install.sh`` script if
you are using the installer upgrade path.  If you are performing in-product upgrades
(CLI or Web Portal), the steps are the same.

There are a few different paths for upgrading a DRP Endpoint in.  Note that how you
originally installed the service may dictate how you need to upgrade.  For some
environments where several security meausres have been taken (non-root install, for
example), you may only be able to upgrade via the in-product path.  Upgrades via the
installer script will require root permissions due to either low port binding
requirements, or ``setcap`` setup requirements.

1. In-product upgrade from the command line (CLI)

  * See current versions in the catalog:

    ::

       drpcli catalog item show drp

  * Install with Catalog reference (**requires Internet connetion**), or in multiple
    steps, if you need to acquire the Zip file, and then perform an air-gap (**non
    Internet connected system**) upgrade.

    * **Internet Connected** - Single Step Upgrade via CLI:

      ::

        # note the item name (drp) has been combined with the selected version
        # without a version, the latest 'stable' release will be used

        drpcli system upgrade catalog:drp-v4.3.0-alpha2.18

    * **Air-gap** (non-internet connected DRP Endpoint) - Download zip file with an
      internet connected host then move to air-gap (non Internet connected) DRP
      endpoint, then perform the upgrade

      *  Download the DRP zip file for the version you want to upgrade to:

        ::

          drpcli catalog item download drp --version v4.3.0-alpha2.18

      * Upgrade from the downloaded zip file:

        ::

          drpcli system upgrade drp.zip

      * Verify newly installed version

        ::

          drpcli info get | jq '.version'

2. Upgrade via the Web Portal

  * navigate to the ``Catalog`` menu item
  * find the ``Digital Rebar Provision`` entry
  * select the version you want to Upgrade (or downgrade) to in the *v.Available*
    (short for *Available Versions*) column
  * click on the green *Install* icon to the right, or the blue button top center
    that says *Install*
  * note that the DRP version zip file has to be downloaded from the RackN hosted
    catalog location, then uploaded to the DRP Endpoint, it may take a few minutes

3. Upgrade with the ``install.sh`` script

  * Stop dr-provision, using the system method of choice

    ::

      sudo systemctl stop dr-provision

    or

    ::

      sudo service dr-provision stop

    or

    ::

      sudo pkill dr-provision

  * Install new code - Use the same install technique as the first install, but
    change ``install`` to ``upgrade`` option.  (Reference: :ref:`rs_install`)

    ::

      # you will want to use additional options if you specified them
      # in your original 'install' steps (eg "--systemd --startup")
      #
      # your original install should have saved a copy of the install.sh
      # script as '/usr/local/bin/drp-install.sh' for this purpose.
      #
      # if an Isolated install was performed originally, add '--isolated'

      drp-install.sh upgrade <Other_Options>

      # or, re-get the installer code if it's not available

      curl -s get.rebar.digital/stable | bash -s -- upgrade <Other_Options>

  * Start up dr-provision

    ::

      systemctl start dr-provision

    or

    ::

      service dr-provision start

    or

    Manually restart as per your standard *Isolated* mode install directions.


.. _rs_upgrade_container:

Container Upgrade Process
-------------------------

As of DRP version v4.3.0, container based installs do not support in-product
upgrade path, the original container must be upgraded via the container
management system.  RackN releases it's container with a separate data
volume for storing the backing write layers of the *dr-provision* service.

By default the DRP service container will be named ``drp``, and the backing
volume will be named ``drp-data``.  Note that you can change these with the
install time flags if desired.

The upgrade process entails:

  * stop dr-provision service to flush all writable data to disk
  * kill the container on the container host
  * start a new container, re-attaching the backing data volume

The installer scripts (``install.sh``) supports these operations.  Review the
script options with the ``--help`` flag for the most up to date information on
usage.

.. note:: WARNING: It is important that you retain a copy of the settings used
          from your original container install.  The upgrade process does not
          have any awareness of previous container start settings.  It may be
          possible to parse this from the container environment (eg 'docker
          inspect drp'), but this has not been determined yet.

Example upgrade of a container based service, based on the following install
command line options:

  ::

    ./install.sh install --container --container-restart=always --container-netns=host --container-env="RS_METRICS_PORT=8888 RS_BINL_PORT=1104"

Based on these install options, the upgrade process is as follows:

  ::

    ./install.sh upgrade --container --container-restart=always --container-netns=host --container-env="RS_METRICS_PORT=8888 RS_BINL_PORT=1104"

.. note:: The only material differnece is the use of the 'upgrade' argument to the
          install script for upgrades, instead of 'install' for installation.


.. _rs_downgrade_drp:

Downgrade Steps
===============

Downgrading DRP from one minor relaase to another *REQUIRES* addtional steps - as the
underlaying database that backs the service may very well change between minor point
releases (eg v4.5.x to v4.6.x).  Database changes do not occur between Patch releases
(eg v4.5.5 to v4.5.6).

.. warning:: You should **ALWAYS** perform these downgrade steps if you are moving from
             one higher point release to a lower point release (eg v4.6.x to v4.5.x).

All downgrade steps and examples below are run at the shell of the server where the ``dr-provision``
service is installed and running, as the ``root`` user (generally, unless installed as a
non-privileged user).

.. note:: Downgrade is only supported for major release version v4.x.x to another v4.x.x version.
          No downgrade is supported or possible in the v3.x.x version line.


Backup DRP First
----------------

Please see :ref:`rs_backup_instructions` documentation.


Stop dr-provision Service
-------------------------

The ``dr-provision`` service needs to be stopped for a downgrade procedure, as we must convert
the database records to flat JSON text files.  We call this process "*humanize*", as it turns
the database records in to human readable components.

  ::

    # for systemd "production" install modes:
    systemctl stop dr-provision
    systemctl status dr-provision               # verify it's not running

    # for other modes, you may need to kill it:
    pkill dr-provision
    ps -ef | grep -v grep | grep dr-provision   # should return no process entries


"*Humanize*" the Database
-------------------------

The ``dr-provision`` binary has a special flag ``--humanize`` which converts the current database
format components in to human readable JSON text files.  You must run the same ``dr-provision``
version binary as the database format is using.  In addition, if you have installed DRP in a
location other than the default production install path (``/var/lib/dr-provision``), you must
also specify where the DRP base directory is with the ``--base-root`` flag.

Once the ``dr-provision`` service is stopped, now perform the "*humanize*" step:

  ::

    # depending on install mode, 'dr-provision' may not be in your direct path,
    # locate the proper binary and call it with correct PATH/dr-provision as appropriate

    DRP_ROOT="/var/lib/dr-provision"                    # adjust this accordingly
    dr-provision --humanize --base-root=$DRP_ROOT

To verify that the "*humanize*" step completed propertly, look at the base directory
for (potentially) a new directory named ``digitalrebar``.

The base directory location will vary depending on how your service is installed.
By default this will be in the ``/var/lib/dr-provision`` directory for "default
production" installs.  It will be a directory named ``drp-data`` for "isolated"
mode installs in the Current Working Directory that the install was performed.

An example of "*humanize*" of a DRP v4.6.0 system:

  ::

    # DRP v4.6.0 currently running example:

    root@mach-04:~# cd /var/lib/dr-provision

    root@mach-04:/var/lib/dr-provision# ls
      ha-state.json  job-logs  plugins  replace  runner  saas-content  server.crt  server.key  tftpboot  ux  wal

    # humanize step

    root@mach-04:/var/lib/dr-provision# /usr/local/bin/dr-provision --humanize --base-root=/var/lib/dr-provision
      dr-provision2021/03/27 15:26:18.250522 Processing arguments
      dr-provision2021/03/27 15:26:18.250812 Version: v4.6.0
      dr-provision2021/03/27 15:26:18.251282 Extracting Default Assets
      dr-provision2021/03/27 15:26:19.614140 [2:1]:backend [ warn]: github.com/hashicorp/raft@v1.2.0/raft.go:214
      [2:1]heartbeat timeout reached, starting election: last-leader=
      dr-provision2021/03/27 15:26:19.711402 [2:2]:backend [audit]: github.com/rackn/provision-server/v4/datastack/stack.go:1958
      [2:2]Seeded CommitID: 3

Now verify that the "*humanize*" completed successfully, and that our database records
have been turned in to human readable JSON files on disk:

  ::

    # verify the humanize completed

    root@mach-04:/var/lib/dr-provision# ls
      digitalrebar  ha-state.json  job-logs  plugins	replace  runner  saas-content  secrets	server.crt  server.key	tftpboot  ux  wal

    root@mach-04:/var/lib/dr-provision# ls digitalrebar
      preferences  profiles  users

In the above output, note the presense of the ``digitalrebar`` directory, and subsequently, the
directory structure underneath it.  This is newwly "*humanized*" objects that were stored in
the v4.6.0 database (in this example).


Install Older dr-provision Service
----------------------------------

To install the older ``dr-provision`` service, we will need to manually extract the binary
out of the distributed TAR.GZ file (even though the file ends in ``.zip``).  If you do not
currently have the older binary (eg from a previous backup or another DRP instance in your
environment) you will have to download it.

  ::

    # example of getting the v4.5.6 with drpcli

    root@mach-04:/tmp# drpcli catalog item download drp --version=v4.5.6

    root@mach-04:/tmp# ls -l *zip
      -rw-r--r-- 1 root root 232361496 Mar  27 15:34 drp.zip

The above command automatically parses the RackN distributed JSON Catalog to find the
download location and get the version.  Some older versions will be removed from the
catalog from time to time to keep it to managable size.  In that case, you may need to
acquire it from alternative locations.

One possibility is to directly download from the RackN staging location in an Amazon S3
bucket.

.. warning:: RackN may change the staging locations in the future, please verify with
             the RackN team if you are having download issues via this mechanism.

  ::

    # using a constructed URL to find an older version archive file and sha256 sum file
    VER="v4.2.0"
    wget -O drp-${VER}.zip https://rebar-catalog.s3-us-west-2.amazonaws.com/drp/${VER}.zip
    wget -O drp-${VER}.zip.sha256 https://rebar-catalog.s3-us-west-2.amazonaws.com/drp/${VER}.sha256

Now unroll the archive file ... yes, the format really is a TAR.GZ despite the filename
ending in ``.zip``:

  ::

     tar -xzvf drp.zip

Verify the binary (using our very old v4.2.0 example from above):

  ::

    root@mach-04:/tmp# ls -l bin/linux/amd64/dr-provision
      -rwxrwxr-x 1 2000 2000 75890688 Dec 29  2019 bin/linux/amd64/dr-provision

    root@mach-04:/tmp# bin/linux/amd64/dr-provision --version
      dr-provision2021/03/27 15:43:28.810743 Processing arguments
      dr-provision2021/03/27 15:43:28.810766 Version: v4.2.0

Put it in place:

  ::

    # move old binary aside - adjust path appropriately for your system
    mv /usr/local/bin/dr-provision /usr/local/bin/dr-provision.old

    # copy new binary in place - adjust path appropriately for your system
    cp bin/linux/amd64/dr-provision /usr/local/bin/

.. note:: For installations as non-root user, you may need to adjust ``setcap`` bits
          appropriately on the binary.  Please see :ref:`rs_install_special_permissions`
          for more details.


Start the dr-provision Service
------------------------------

Now start up your DRP service as you would normally:

  ::

    # systemd "production" install:

    systemctl start dr-provision
    systemctl status dr-provision

    # possible startup command for an isolated mode install (this command assumes
    # the setup symbolic links is still in place and pointing at the binary path
    # correctly):

    sudo ./dr-provision --base-root=`pwd`/drp-data --local-content="" --default-content="" > drp.log 2>&1 &

.. note:: The startup process may take some time (up to 15 minutes), if you have a very
          large number of Machines and Jobs Logs, as the JSON data structures are converted
          in to database records.


Verify the service is running the new (old) version via the command line tool:

  ::

    drpcli info get | jq -r '.version'

The returned string should be the version, eg ``v4.5.6``.


Version to Version Notes
========================

In this section, notes about migrating from one release to another will be added.

Release Notes for each version can be found at:  https://github.com/digitalrebar/provision/v4/releases


Install Script Changed
----------------------

There are minor changes to the install script for isolated mode.  Production mode installs are still done and updated the same way.  For isolated, there are some new flags and options.  Please see the commands output for more details or check the updated :ref:`rs_quickstart`.

For current ``install.sh`` script usage information, please run:

  ::

    install.sh --help


For complete details.

