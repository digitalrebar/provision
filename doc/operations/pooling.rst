.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Platform documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Platform; Pooling Operations

.. _rs_pooling_ops:

Pooling Operations
==================

This section will address usage of the pooling system.  The architecture and implementation of the pooling system is
described at :ref:`rs_pooling_arch`.

Pool Actions
------------

The Pool objects follow the Digital Rebar Provision standard API model.  One can list, show, update, and destroy pools.
Additionally, pools can have actions from plugins.  Pools do not have parameters like a Machine or Plugin.

::

  drpcli pools list  # API call - GET /api/v3/pools
  drpcli pools show mypool # API call - GET /api/v3/pools/mypool
  drpcli pools update mypool '{ "Description": "mypool desc" }' # API call - PUT or PATCH /api/v3/pools/mypool with Data
  drpcli pools destroy mypool # API call - DELETE /api/v3/pools/mypool


NOTE: These commands only operate on Pool objects.  Pools can exist without objects just by machines being assigned
to them.

There is an additional command to get the status of the machines inside a pool.  This works pool objects and non-pool
object pools.  The active pools are pools with machines in them without or with a pool object.

::

  drpcli pools active # API call - GET /api/v3/pools-active - NOTE The change.
  drpcli pools status mypool   # API call - GET /api/v3/pools/mypool/status

These allow you to manage machine membership for pool objects and non-pool objects.

::

  drpcli pools manage add mypool # API call - POST /api/v3/pools/mypool/addMachines with Data
  drpcli pools manage remove mypool # API call - POST /api/v3/pools/mypool/removeMachines with Data
  drpcli pools manage allocate mypool # API call - POST /api/v3/pools/mypool/allocateMachines with Data
  drpcli pools manage release mypool # API call - POST /api/v3/pools/mypool/releaseMachines with Data

The CLI takes additional parameters that let you add additional controls.

* --add-parameters - A JSON string parameters to add to the system.  This is a map like the Parameters section of a Machine.
* --remove-parameters - A comma separated list of parameters to remove from the machine.
* --add-profiles - A comma separated list of profiles to add to the machine.
* --remove-profiles - A comma separated list of profiles to remove from the machine.
* --new-workflow - A workflow to assign to the machine when the action is done.
* --count - the number of machines to operate on
* --all-machines - Override count and use all machines in the pool or source pool.
* --machine-list - a comma separated list of filters.  E.g Name:mymachine or just a UUID
* --minimum - The minimum number of machines that must be found for success.  Using count and minimum can allow for partial success.
* --wait-timeout - the number of seconds (or time string, e.g. 10h) to wait for all the machines to achieve the next status.
* --source-pool - the pool to add machines from.  Defaults to the parent of the pool or default if unspecified. (only for add)

Filters can be added to the drpcli command line to reduce the scope.  For example, ipmi/enabled=true would restrict operations
to machines that have ipmi/enabled set to `true`.  The normal and full set of machine filters are available.  Multiple
filters are ANDed together.

For the API, these are POSTed as the following json blob.  Source pool is a query parameter.

::

  {
    "pool/add-parameters": {
      "param1": "string",
      "param2": true
    },
    "pool/remove-parameters": [ "param3", "param4" ],
    "pool/add-profiles": [ "profile1", "profile2" ],
    "pool/remove-profiles": [ "profile3", "profile4" ],
    "pool/workflow": "nextworkflow",
    "pool/count": 30,
    "pool/all-machines": false,
    "pool/machine-list": ["UUID1", "Name:mymachine"],
    "pool/minimum": 1,
    "pool/wait-timeout": "30m",
    "pool/filter": [ "ipmi/enable=true" ]
  }




Pool Objects
------------

Here is an example of a pool object with actions defined.  This assumes that you have already created the required
profiles and workflows.


::

  ---
  Id: my-machines
  Description: Pool that install linux and cleans the machine
  EnterActions:
    Workflow: clean-machine
    AddProfiles:
    - pool-test-1
    AddParameters:
      auto-param-1: true
    RemoveProfiles:
    - pool-test-2
    RemoveParameters:
    - auto-param-2
  AllocateActions:
    Workflow: install-linux
    AddProfiles:
    - pool-test-2
    AddParameters:
      auto-param-2: false
    RemoveProfiles:
      - pool-test-1
    RemoveParameters:
      - auto-param-1
  ReleaseActions:
    Workflow: discover
    AddProfiles:
    - pool-test-1
    AddParameters:
      auto-param-1: false
    RemoveProfiles:
    - pool-test-2
    RemoveParameters:
    - auto-param-2
  ExitActions:
    Workflow: clean-machine
    AddProfiles:
    - pool-test-2
    AddParameters:
      auto-param-2: true
    RemoveProfiles:
    - pool-test-1
    RemoveParameters:
    - auto-param-1

