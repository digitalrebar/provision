.. Copyright (c) 2019 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Runner State

Runner State, Reboots and BootEnvs
==================================

This section discusses how and when Digital Rebar performs reboots based on Workflows and BootEnv changes.

There are two parts to this equation: the runner, aka drpcli in "process jobs" mode, and the DRP endpoint.  The DRP endpoint is maintaining the state with regard to the machine object.  This includes the task list and the fields (the important ones for now are stage, bootenv, and workflow).

When a runner starts (and the machine is runnable), the runner will start processing tasks.  These tasks are the normal ones that you think about in the stages, but also special ones injected by DRP when the workflow, stage, or bootenv are set.

When you look at a task list, you might see something like this example task list from a VirtualBox discovery workflow.
 	::
    0: stage:discover
    1: bootenv:sledgehammer
    2: gohai
    3: ssh-access
    4: stage:ipmi-configure
    5: ipmi-configure
    6: stage:virtualbox-discover
    7: virtualbox-discover-uuid
    8: stage:sledgehammer-wait

That workflow is: `discover`->`ipmi-configure`->`virtualbox-discover`->`sledgehammer-wait` for testing some stuff in ipmi.

When DRP is told to set a workflow on a machine, it decomposes the workflow into tasks. There are some special tasks ``stage:xxxx`` and ``bootenv:yyyy`` that represent stage changes in the workflow or bootenv changes in the workflow as indicated by a stage that has a different workflow from the previous one.

In this scenario, the machine boots into sledgehammer then the runner starts and asks DRP for its tasks.  We’ve already set the machine’s workflow to discover.  This built that task list and set the stage and bootenv to the values specified by the first stage in the workflow. (discover and sledgehammer respectively).

Once the runner started in sledgehammer, the runner asks DRP for tasks:

  #. The first task is set stage to discover.  DRP does this and sees it is already done, it then moves to the next task.  that task is set bootenv to sledgehammer.
  #. It sees that is already done and moves to the next one, gohai.  It can’t do that one and returns to the runner the info to do the gohai task.
  # Once done, the runner asks for more tasks after updating the status.
  #. Once DRP gets to the stage:ipmi-configure task, it will change the stage, and return that to the runner.  The runner goes okay, and asks for more tasks.  This allows the runner to do something on the stage change if necessary.
  # Then we keep going through the list until both the runner and DRP (as directed by the runner) think there are no more tasks to run.

In this case, the bootenv doesn’t change - so the runner doesn’t reboot anything.

Say we have this workflow: ``prep-install`` -> ``centos-7-install`` -> ``runner-service`` -> ``finish-install`` -> ``complete``

The task list looks like this:
  ::
    0: stage:centos-7-install
    1: bootenv:centos-7-install
    2: set-hostname
    3: centos-drp-only-repos
    4: ssh-access
    5: stage:runner-service
    6: drpcli-install
    7: stage:finish-install
    8: bootenv:local
    9: stage:complete

This workflow will wipe the disks of a system, then install centos 7,  install a runner into the image, finish the install, reboot, and the get marked complete by the runner in newly booted os.  In this case, 

  #. the first stage sets the bootenv to sledgehammer (if we are there, it is fine nothing happens).  If drpcli sees this as a change, it will attempt to kexec or reboot the node into that bootenv.
  # In our case of a discovered node, the system is sitting in sledgehammer so nothing happens.  The runner and DRP move through task list cleaning the disks until the bootenv change to centos-7-install.
  # At this point, the runner sees the bootenv change and reboots/kexecs the system into that new bootenv.
  # The centos-7-install bootenv installs the machine from the kickstart templates and during the post-install phase starts a runner in the system chroot.
  # The runner pulls tasks and continues updating the system.
  # This continues until the ``bootenv:local``
  # drpcli notices the bootenv change and prepares to reboot/kexec the system, but in this case does NOT.

The runner has a historical anomaly for this case.  If the bootenv’s name ends in ``-install``, the runner exits instead of rebooting/kexecing.  This is to allow the kickstart / preseed based OSes “finish” their processing and reboot or kexec themselves.

Once the system boots from the local disk, the runner starts and processes the last task which has DRP set the stage to `complete`.  The runner goes idle at that point. (edited) 

Empty BootEnv
-------------

Another thing to note is that a stage with an empty bootenv (“”) means use the currently set bootenv without change.  So, you could also boot into sledgehammer, then use stages that never change bootenv or set bootenv.  This machine doesn’t reboot through the process.

Many of the stages don’t specify bootenv because they can be run in many different bootenvs.  This way they continue working where ever they are run.  In some cases, stages have a specific requirement about a bootenv (like the install ones or some of the machine prep/update ones assume sledgehammer because of tooling or machine state (like disks not mounted)).