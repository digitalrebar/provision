.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Operational Runbooks


Runbooks
++++++++

Backup
------
For the following commands it is assumed you already have ``drbup`` installed and in your PATH. For full documentation on ``drbup`` please see :ref:`the documentation <rs_drbup>`.

* Create a backup of a running DRP endpoint. Backs up the system to a single file that can be shipped off to a separate server. In this process the dr-provision service will be stopped while the system is backed up. Once the backup is complete the service will be restarted.

  ::

    drbup backup --source /var/lib/dr-provision --dest /root/dr-provision_backup.tar.bz2



* Restore the backup file, and start the dr-provision service.

  ::

    systemctl stop dr-provision.service
    drbup restore --source /root/dr-provision_backup.tar.bz2 --dest /var/lib/dr-provision --start-service

* Sync local DRP to remote DRP2

  ::

    systemctl stop dr-provision.service
    drbup sync --remote --source /var/lib/dr-provision --dest /var/lib/ --remote-host admin@remote-drp2.internal
    systemctl start dr-provision.service

.. note:: Local DRP and remote DRP2 must have shared ssh keys for the user, otherwise manual intervention will be required by an operator to supply credentials.