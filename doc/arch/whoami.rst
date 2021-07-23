.. Copyright (c) 2021 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Who Am I


.. _rs_provisioning_whoami:

Machine Matching (the Who Am I API)
<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

In order to match discovered machines with ones already in the database, 
Digital Rebar uses a matching process called "whoami".  This process uses
DRPCLI to create a fingerprint based on several criteria about the machine
stored in the `Machine.Fingerprint` array.  The API then scores the fingerprint
against known machines' fingerprints. If a match is found, then the API returns
the machines UUID.

.. _rs_cli_whoami:

DRPCLI whoami
-------------

The command ``DRPCLI machines whoami`` performs the matching process and returns the Machine'
UUID if a match is made.


Fingerprint Critieria
---------------------

.. _rs_fingerprint:

Fingerprint Critieria
---------------------

The follow items compose the machine's fingerprint and other properties.
Each has a different contribution to the matching score.  A score of 100 must
be reached for a match to be made.

Score calculates how closely the passed in Whoami matches a candidate Machine.
In the current implementation, Score awards points based on the following
criteria:

* 25 points if the Machine has an SSNHash that matches the one in the Whoami
* 25 points if the Machine has a CSNHash that matches the one in the Whoami
* 50 points if the Machine has a SystemUUID that matches the one in the Whoami
* 0 to 100 points varying depending on how many memory DIMMs from the machine fingerprint are present in Whoami.
* 0 to 100 points varying depending on how many HardwareAddrs from the Machine are present in Whoami.
* 1000 points if the machine UUID matches OnDiskUUID
* 500 points if Cloud Type & Cloud ID matches from cloud-init JSON spec.

If the score is less than 100 at the end of the scoring process, it is rounded down
to zero.  The intent is to be resilient in the face of hardware changes:
  * SSNHash, CSNHash, and SystemUUID come from the motherboard.
  * MemoryIds are generated deterministically from the DIMMs installed in the system
  * MacAddrs comes from the physical Ethernet devices in the system.

