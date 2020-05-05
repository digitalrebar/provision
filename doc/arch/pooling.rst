.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Pooling Architecture

.. _rs_pooling_arch:

Pooling Architecture
<<<<<<<<<<<<<<<<<<<<

Digital Rebar Provision provides a mechanism to group machines into pools with defined actions on the transitions.  This
functionality can be used to provide ``cloud-like`` operations on bare metal servers.

This section will focus on the architecture of the pooling system.  The operational methods can be found at :ref:`rs_pooling_ops`.

Pools provide three basic functions.

* Machine grouping into pools
* Machine allocation to track usage or consumption
* Machine manipulation upon transitions in the pool

Pools are dynamic and exist because machines are in pools.  Optionally, pool objects can be added to define additional
control points and information.

API calls can be made blocking or non-blocking.  Blocking calls will wait for a timeout or for all the machines to achieve
workflow complete before returning.

Pool Functions
--------------

Grouping by Pools
=================

Machines can be group into pools.  A machine can only be in one pool at a time.  All machines start in the default pool
unless the pool variables are set upon machine creation.  The machines can be managed by API (CLI or UX) calls to move
the machines between pools.

Machines may be ``added`` or ``removed`` from pools.  Machines may be selected by specific machine values, e.g. UUID
or Name, or filters.  The API methods allow for the requesting of a number of machines and how many are required to succeed.
Machines are only moved if the minimum number of machines can be found.

Pools are loosely hierarchical.  The ``default`` pool is the root pool.   Any pool with out an explicit parent will
consider the ``default`` pool as its parent.  Pool parentage can be used when adding machines to and removing machines from
a pool.  Unless explicitly declared in the API, the pool's parent will be used as the source or destination of the operation.

Machine Allocation
==================

Machines can be allocated and released from inside a pool.  This allows for API calls to ask for machines without having
to know the specific sets of machines. Machines may be selected by specific machine values, e.g. UUID or Name, or filters.
The API methods allow for the requesting of a number of machines and how many are required to succeed.
Machines are only moved if the minimum number of machines can be found.

Allocating a machine reserves that machine and makes it not available for additional allocation.  Release a machine removes
the reservation on the machine and makes it available for allocation within the pool again.

Machine Transitions
===================

As machines transition into or out of the pools or allocated or released from with in a pool, these machines can be altered
as the system operates on them.  The transitions can change workflow, add or remove parameters, and add or remove profiles.

All transitions are supported, transitioning into a pool, allocating a machine, releasing a machine, and transitioning
out of the pool.  A machine leaving a pool actually handles two transitions, the removal from the current pool actions
and the entry into the new pool actions.

These actions can be added on the API call for direct control.  The actions can also be added to the pool object that
allows for simpler calls.


Pool Implementation
-------------------

The following sections define the implementation details of Pooling.

Machine Object
==============

To track the machine's status in the pool, three fields are on the machine object.

* Pool - This string field declares what pool the machine is in.
* PoolAllocated - This boolean field indicates if the machine has been allocated or not.
* PoolStatus - This string field describes the machine's state in the pool.

Additionally, the pool system needs to know when a workflow is complete.  There is also a boolean
field ``WorkflowComplete`` that indicates if the current workflow is complete.

Pool Status can have the following values:

* Joining - Running a workflow defined when entering the pool, but not complete yet.
* HoldJoin - Running a workflow defined when entering the pool that has a failed task and needs remediation before continuing.
* Free - Waiting for allocation within the pool.
* Building - Running a workflow defined when allocated, but not complete yet.
* HoldBuild - Running a workflow defined when allocated that has a failed task and needs remediation before continuing.
* InUse - Waiting for release but has completed all allocation actions
* Destroying - Running a workflow defined when released, but not complete yet.
* HoldDestroy - Running a workflow defined when released that has a failed task and needs remediation before continuing.
* Leaving - Running a workflow defined when removed from the pool, but not complete yet.
* HoldLeave - Running a workflow defined when removed that has a failed task and needs remediation before continuing.

Pool Object
===========

The pool object is an optional construct for the system.  The pool object is needed when the parent field is needed.  The
object is also needed when defining transition pieces.

The pool object defines ``EnterActions``, ``ExitActions``, ``AllocateActions``, and ``ReleaseActions``.  These hold a
structure that allow you to define:

* Workflow - the workflow to use for this transition.
* AddProfiles - A list of profiles to add to the machine - If already present, this is not an error.
* RemoveProfiles - A list of profiles to remove from the machine. If not present on the machine, this is not an error.
* AddParameters - A map of key/value pairs that get set as parameters on the machine.  If already present, this will replace the current value.
* RemoveParameters - A list of parameters to remove from the machine.  If not present on the machine, this is not an error.

There is a pending function that is *NOT* enabled for autofilling pools that are empty.  This will change in the coming
releases.

