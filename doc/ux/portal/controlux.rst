.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; UX

.. _rs_controlux:

Control
=======
This section of the UX contains critical information for managing the provisioning of nodes within Digital Rebar. Workflows are the high level instruction path necessary for DRP to perform a provision. Workflows are composed of individual stages which accomplish a specific task.  

Workflow
--------
This section lists all available Workflows currently available to execute. Each Workflow has the following details:

* Lock/Unlock
* Name
* Description
* Stages 

To the right is a list of ALL available stages that can be used to build a new Workflow. 

At the top of the screen are a series of blue boxes offering additional functionality:

* Switch to Change/Stage Map - This screen shows the Workflow processes available and how the process order from one state to the next. Steps can be added to Workflows and a Workflow Wizard is available to create a new Workflow.   
* Refresh - Update the Workflow list with the latest available Workflows 
* Filter - Select which Workflow to list by Available, Key, Name, ReadOnly, and Valid
* Add - Add a new Workflow
* Clone - Clone a Workflow
* Delete - Delete a Workflow

For more details :ref:`rs_workflows`

Stages
------
This section lists all available Stages within the DRP system. Each Stage has the following details:

* Lock/Unlock
* Name 
* Boot Environment
* Description 

At the top of the screen are a series of blue boxes offering additional functionality: 

* Refresh - Update the Stages list with the latest available Stages
* Filter - Select which Stage to list by Available, BootEnv,  Key, Name, ReadOnly, Reboot and Valid
* Add - Add a new Stage
* Clone - Clone a Stage
* Delete - Delete a Stage

For more details :ref:`rs_data_stage`

Tasks
-----
This section lists all available Tasks within the DRP system. Each Task has the following details:

* Lock/Unlock
* Name
* Description 

At the top of the screen are a series of blue boxes offering additional functionality: 

* Refresh - Update the Task list with the latest available Stages
* Filter - Select which Task to list by Available, Key, Name, ReadOnly, and Valid
* Add - Add a new Task
* Clone - Clone a Task
* Delete - Delete a Task

Task Details
------------
During the boot process, tasks provide additional configuration to machines in the form of templates. BootEnvs will use these sets of templates to construct specific jobs for a machine.

Within a task, templates are processed in the order they are assigned, so it’s important to check that templates are attached correctly to a task.

For more details :ref:`rs_data_task`


Jobs
----
This section lists all available Jobs within the DRP system. Each Job has the following details:

* State
* Start Time
* Run Time
* Job UUID
* Machine
* Stage
* Task 

At the top of the screen are a series of blue boxes offering additional functionality:

* Refresh - Update the Jobs list with the latest available Jobs
* Filter - Select which Task to list by Archived, Available, BootEnv, Current, EndTime,  Key, Machine, Previous, ReadOnly, Stage, StartTime, State, Task, Uuid, Workflow, and Valid
* Delete - Delete a Job 


Job Details
-----------
A job defines a machine’s current step in its boot process. After completing a job, the machine creates a new job from the next instruction in the machine’s task list.

Machines will only process one job at a time, and jobs are not created until the instant they are required.


For more details :ref:`rs_data_job`


