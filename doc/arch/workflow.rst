.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Workflow


Workflow
========

dr-provision implements a basic Workflow system to help automate the
various tasks needed to provision and decommission systems.  The
workflow system is the results of several other components in
dr-provision interacting.  The rest of this section goes over those
parts in more detail.


The Bits and Bobs
^^^^^^^^^^^^^^^^^

Tasks
-----

The basic unit of work that dr-provision sequences to drive Machines
through a workflow is a :ref:`rs_data_task`.  Individual Tasks are
executed against Machines by creating a Job for them.

Tasks contain individual Templates that are expanded for a Machine
whenever a Job is created.

Jobs
----

A :ref:`rs_data_job` is used to track the execution history of Tasks
against a specific Machine.  A Job is created every time a Task is
executed against a machine -- Jobs keep track of their state of
execution.  The history of what has been executed (including all log
output from scripts) is stored as a chain of Jobs, and the exit status
of the Job determines what a machine agent will do next.


Stages
------

A :ref:`rs_data_stage` is used to provide a list of Tasks that should
be run on a Machine along with the BootEnv the tasks should be run in.

Machine Agent
-------------

The Machine Agent is responsible for creating Jobs, writing out or
executing any JobActions, streaming job Logs back to dr-provision for
archival purposes, and updating the Job state based on the exit state
of any Actions.  The Machine Agent is also responsible for rebooting,
powering off, and changing to a different Stage as indicated by the
Job exit status or the change stage map.  Unless directed to exit, the
Machine Agent watches the event stream for the Machine it is running
on and will execute new tasks as they come to be available.

Change Stage Map
----------------

The change-stage/map parameter defines what stage to change to when
you finish all the tasks in the current stage.  The change stage map
is a map whose keys correspond to the stage the machine is currently
in and whose values indicate the next stage to transition to and what
to have the runner do on the stage transition.

How They Work Together
^^^^^^^^^^^^^^^^^^^^^^

