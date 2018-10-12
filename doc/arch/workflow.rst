.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Workflow

.. _rs_workflow:

Workflow
========

dr-provision implements a basic Workflow system to help automate the
various tasks needed to provision and decommission systems.  The
workflow system is the results of several other components in
dr-provision interacting.  The rest of this section goes over those
parts in more detail.


The Bits and Bobs
^^^^^^^^^^^^^^^^^

.. _rs_workflow_tasks:

Tasks
-----

The basic unit of work that dr-provision sequences to drive Machines
through a workflow is a :ref:`rs_data_task`.  Individual Tasks are
executed against Machines by creating a Job for them.

Tasks contain individual Templates that are expanded for a Machine
whenever a Job is created.  Each of these individual Templates can
expand to either a script to be executed (if the Path parameter is
empty or not present), or a file to be placed on the filesystem at the
location indicated by template-expanding the Path parameter.(if the
Path parameter is not empty).

.. _rs_workflow_jobs:

Jobs
----

A :ref:`rs_data_job` is used to track the execution history of Tasks
against a specific Machine.  A Job is created every time a Task is
executed against a machine -- Jobs keep track of their state of
execution.  The history of what has been executed (including all log
output from scripts) is stored as a chain of Jobs, and the exit status
of the Job determines what a machine agent will do next.

.. _rs_workflow_stages:

Stages
------

A :ref:`rs_data_stage` is used to provide a list of Tasks that should
be run on a Machine along with the BootEnv the tasks should be run in.

.. _rs_workflow_mc_agent:

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

Change Stage Map (DEPRECATED)
-----------------------------

The change-stage/map parameter defines what stage to change to when
you finish all the tasks in the current stage.  The change stage map
is a map whose keys correspond to the stage the machine is currently
in and whose values indicate the next stage to transition to and what
to have the runner do on the stage transition.

The change-stage/map mechanism has been replaced by the Workflow
mechanism, but will be maintained for the forseeable future.  You are
encouraged to migrate to using Workflows.

.. _rs_workflows:

Workflows
---------

A :ref:`rs_data_workflow` is used to provide a list of Stages that a
Machine should run through to get to a desired end state.  When a
Workflow is added to a Machine, it renders the BootEnv and Stage for
the Machine read-only, replaces the task list on the Machine with one
that will step through all the Stages, BootEnvs, and Tasks needed to
drive the machine through the Workflow.

How They Work Together
^^^^^^^^^^^^^^^^^^^^^^

.. _rs_workflow_agent:

Machine Agent (client side)
---------------------------

The Machine Agent runs on the Client and is responsible for executing
tasks and rebooting the Machine as needed. It is structured as a
finite state machine for increased reliability and auditability.  The
Machine Agent always starts in the AGENT_INIT state.


AGENT_INIT
  Initializes the Agent with a fresh copy of the Machine
  data, marks the current Job for the machine as `failed` if it is
  `created` or `running`,and creates an event stream that recieves
  events for that Machine from dr-provision.  If an error was recorded,
  the Agent prints it to stderr and then clears it out.

  If an error occurrs during this, the agent will sleep for a bit and
  transition back to AGENT_INIT, otherwise it will transition to
  AGENT_WAIT_FOR_RUNNABLE.

AGENT_WAIT_FOR_RUNNABLE
  Waits for the Machine to be both Available
  and Runnable. Once it is, the Agent transitions to AGENT_REBOOT if
  the machine changed BootEnv, AGENT_EXIT if the Agent recieved a
  termination signal, AGENT_INIT if there was an error waiting for the
  state change, and AGENT_RUN_TASK otherwise.

AGENT_RUN_TASK
  Tries to create a new Job to run on the machine.

  If there was an error creating the Job, transitions back to
  AGENT_INIT.

  If there was no job created, the Agent transitions to
  AGENT_CHANGE_STATE if the Machine does not have a Workflow, and
  AGENT_WAIT_FOR_CHANGE_STAGE if it does.

  If a Job was created, the Agent attempts to execute all the steps in
  the Task for which the Job was created, and updates the Job
  depending on the exit status of the steps.

  If there was an error executing the Job, the agent will transition
  back to AGENT_INIT.

  If the Job signalled that a reboot is needed, the Agent transitions
  to AGENT_REBOOT.

  If the Job signalled that the system should be powered off, the
  Agent transitions to AGENT_POWEROFF.

  If the Job signalled that the Agent should stop processing Jobs, the
  Agent transitions to AGENT_EXIT.

  Otherwise, the Agent transitions to AGENT_WAIT_FOR_RUNNABLE.

AGENT_WAIT_FOR_STAGE_CHANGE
  Waits for the Machine to be Available,
  and for any of the following fields on the Machine to change:

  - CurrentTask
  - Tasks
  - Runnable
  - BootEnv
  - Stage

  Once those conditions are met, follows the same rules as
  AGENT_WAIT_FOR_RUNNABLE.

AGENT_CHANGE_STAGE
  Checks the change-stage/map to determine what
  (and how) to transition to the next Stage when AGENT_RUN_TASK does
  not get a Job to run from dr-provision.

  The Agent first tries to retrieve the change-stage/map Param for the
  Machine from dr-provision.  If it fails due to connection issues,
  the Agent will transition to AGENT_INIT.  If there is no
  change-stage map, the Agent uses an empty one.

  If there is a key in the change-stage/map for the current Stage, the
  Agent saves the corresponding value as val for further processing.

  If there is no next entry for the current Stage in the
  change-stage/map and the Machine is in a BootEnv that ends in
  -install, the Agent assumes that val is "local", otherwise the Agent
  transitions to AGENT_WAIT_FOR_STAGE_CHANGE.

  The Agent splits val into nextStage and targetState on the first ':'
  character in val.

  If targetState is empty, it is set according to the following rules:

  - If the BootEnv for nextStage is not empty or different from the
    current BootEnv, targetState is set to "Reboot"

  - Otherwise targetState is set to "Success"

  The Agent changes the machine Stage to the value indicated by
  nextStage.  If an error occurs during that process, the Agent
  transitions to AGENT_INIT.

  If targetState is "Reboot", the agent transitions to AGENT_REBOOT.
  if targetState is "Stop", the agent transitions to AGENT_EXIT.
  If targetState is "Shutdown", the agent transitions to AGENT_POWEROFF.
  If targetState is anything else, the agent transitions to AGENT_WAIT_FOR_RUNNABLE.

AGENT_EXIT
  Exits the Agent.

AGENT_REBOOT
  Reboots the system.

AGENT_POWEROFF
  Cleanly shuts the system down.

.. _rs_workflow_reboot:

Reboot! Using Agent State Changes in Scripts
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

These functions are implemented in the community content shared template accessed by adding `{{ template "setup.tmpl" }}` in your content. 

By adding this library, you can call the functions ```exit, exit_incomplete, exit_reboot, exit_shutdown, exit_stop, exit_incomplete_reboot, exit_incomplete_shutdown``` to access the agent states.

Script authors can also force behaviors using specialized ``exit`` code in their routines. While ``exit 0`` provides a regular clean exit, non-0 values provide enhanced functionality:

  * ``exit 16`` stops the script
  * ``exit 32`` triggers a shutdown
  * ``exit 64`` triggers a reboot
  * ``exit 128`` means task is incomplete
  * ``exit 192`` means task is incomplete AND system should reboot
  * ``exit 160`` means task is incomplete AND system should shutdown.

The codes are based on intpretation of bit position left as a trivial exercise to the reader until someone updates the documentation.

.. _rs_workflow_server:

dr-provision (server side)
--------------------------

In dr-provision, the machine Agent relies on these API endpoints to perform its work:

- GET from `/api/v3/machines/<machine-uuid>` to get a fresh copy of
  the Machine during AGENT_INIT.

- PATCH to `/api/v3/machines/<machine-uuid>` to update the machine
  Stage and BootEnv during the AGENT_CHANGE_STAGE.

- GET from `/api/v3/machines/<machine-uuid>/params/change-stage/map`
  to fetch the change-stage/map for the system during
  AGENT_CHANGE_STAGE.

- POST to `/api/v3/jobs` to retrieve the next Job to run during
  AGENT_RUN_TASK.

- PATCH to `/api/v3/jobs/<job-uuid>` to update Job status during
  AGENT_RUN_TASK and during AGENT_INIT.

- PUT to `/api/v3/jobs/<job-uuid>/log` to update the job log during
  AGENT_RUN_TASK.

- UPGRADE to `/api/v3/ws` to create the EventStream websocket that
  recieves Events for the Machine from dr-provision.  Each Event
  contains a copy of the Machine state at the point in time that the
  event was created.

.. _rs_workflow_next job:

Retrieving the next Job
~~~~~~~~~~~~~~~~~~~~~~~

Out of all those endpoints, the one that does the most work is the
`POST /api/v3/jobs` endpoint, which is responsible for figuring out
what (if any) is the next Job that should be provided to the Machine
Agent.  It encapsulates the following logic:

#. dr-provision recieves an incoming POST on `/api/v3/jobs` that
   contains a Job with just the Machine filled out.

   If the Machine does not exist, the endpoint returns an
   Unprocessable Entity HTTP status code.

   If the Machine is not Runnable and Available, the endpoint returns
   a Conflict status code.

   If the Machine has no more runnable Tasks (as indicated by
   CurrentTask being greater than or equal to the length of the
   Machine Tasks list), the endpoint returns a No Content status code,
   indicating to the Machine Agent that there are no more tasks to
   run.

#. dr-provision retrieves the CurrentJob for the Machine.  If the
   Machine does not have a CurrentJob, we create a fake one in the
   Failed state and use that as CurrentJob for the rest of this
   process.

#. dr-provision tentatively sets `nextTask` to CurrentTask + 1.

#. If the CurrentTask is set to -1 or points to a `stage:` or
   `bootenv:` entry in the machine Task list, we mark the CurrentTask
   as `failed` if it is not already `failed` or `created`.

#. If CurrentTask is set to -1, we update it to 0 and set `nextTask` to 0.

#. If CurrentTask points to a `stage:` or a `bootenv:` entry in the
   Tasks list, and the Machine is not already in the appropriate Stage
   or BootEnv, we skip the next step. Otherwise we skip past these
   entries in the Tasks list until we get to an entry that refers to a
   Task and update CurrentTask and `nextTask` to point to that entry.

#. Depending on the State of the CurrentJob, we take one of the following actions:

   - "incomplete": This indicates that CurrentJob did not fail, but it
     also did not finish.  dr-provision returns CurrentJob unchanged,
     along with the Accepted status code.

   - "finished": This indicates that the CurrentJob finished without
     error, and dr-provision should create a new Job for the next Task in the
     Tasks list.  dr-provision sets CurrentTask to `nextTask`.

   - "failed": This indicates that the CurrentJob failed.  Since
     updating a Job to the `failed` state automatically makes the
     Machine not Runnable, something else has intervened to make the
     machine Runnable again. dr-provision will create a new Job for
     the current Task in the Tasks list.

#. dr-provision creates a new Job for the Task in the Tasks list
   pointed to by CurrentTask.  If CurrentTask points to a `stage:` or
   a `bootenv:` task entry, the new Job is created in the `finished`
   state, otherwise it is created in the `created` state. The Machine
   CurrentJob is updated with the UUID of the new Job.  The new Job
   and the Machine are saved.

#. If the new Job is in the `created` state, it is returned along with
   Created HTTP status code, otherwise nothing is returned along with
   the NoContent status code.

.. _rs_workflow_changing:

Changing the Workflow on a Machine
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Changing a Workflow on the Machine has the following effects:

- The Stages in the Workflow are expanded to create a new Tasks list.
  Each Stage gets expanded into a List as follows:

  - `stage:<stageName>`
  - `bootenv:<bootEnvName>` if the Stage specifies a non-empty BootEnv.
  - The Tasks list in the Stage

  The Tasks list on the Machine are replaced with the results of the
  above expansion.

- The CurrentTask index is set directly to -1.

- The Stage and BootEnv fields become read-only from the API.
  Instead, they will change in accordance with any `stage:` and
  `bootenv:` elements in the Task list resulting from expanding the
  Stages in the Workflow.  Any Stage changes that happen during
  processing a Workflow do not affect the Tasks list or the
  CurrentTask index.

.. _rs_workflow_removing:

Removing a Workflow from a Machine
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

To remove a workflow from a Machine, set the Workflow field to the
empty string.  The Stage field on the Machine is set to `none`, the
Tasks list is emptied, and the CurrentTask index is set back to -1.

Changing the Stage on a Machine
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Changing a Stage on a Machine has the folowing effects when done via
the API and the Machine does not have a Workflow:

- The Tasks list on the Machine is replaced by the Tasks list on the
  Stage.

- If the BootEnv field on the Stage is not empty, it replaces the
  BootEnv on the Machine.

- The CurrentTask index is set to -1

- If the Machine has a different BootEnv now, it is marked as not Runnable.

Resetting the CurrentTask index to -1
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

If the Machine does not have a Workflow, the CurrentTask index is
simply set to -1.  Otherwise. it is set to the most recent entry that
would not occur in a different BootEnv from the machine's current
BootEnv.
