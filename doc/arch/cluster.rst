.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Cluster

.. _rs_cluster_pattern:

Multi-Machine Cluster Pattern v4.6+
===================================

In version 4.6, Digital Rebar multi-machine clustering patterns were updated to
leverage contexts for cluster management.  If you are using pre v4.6 then please see :ref:`rs_cluster_pattern45`.


The following components are central to the design:

* a shared cluster data in a shared profile
* coordination by a logical cluster manager machine (generally in a context)
* use of WorkflowComplete to synchronize tasks across the cluster
* params to identify cluster rules for manager, leader and worker.


Design Pattern
--------------

Clusters are multi-machine systems where the install and update operations coordinated by Digital Rebar workflows.  The design allows for central coordination via a cluster manager workflow(s) and shared data via a shared cluster profile.

Since coordinating cluster activity can be complex, the design must be as simple as possible.  This means limiting side-effects and ensuring idempotency when possible.

To keep the pattern simple, workflows for machine specific actions operate using the normal workflow pattern on individual machines.  Data collection and order of operation actions are managed by the cluster manager with simple "wait for" operations as a backup or failsafe.

The shared data in the cluster profile uses a special Profile Token enabled by Digital Rebar that allows machines to require write access to profile in addition to normal read-only access.

The cluster manager is a Digital Rebar construct and not expected to be an active part of the cluster workflow.  For this reason, it is generally implemented as a logical machine using a Context.  The cluster manager is responsible for coordinating activity in the cluster by 1) ensuring shared data is correct, 2) starting workflows on the cluster members and 3) running administrative tasks as needed.

Ideally, this results in member workflows that can be operated individually outside of the cluster and a declarative control loop that is basically "run these workflow on machines matching these filter."

Cluster Roles
-------------

The following roles are used throughout this pattern:

* manager (identified by cluster/manager: true) - the logical machine used to coordinate operators for the cluster
* leader (identified by cluster/leader: true) - an optional cluster defined role generally needed by platforms for cluster operations.  Different roles can be defined as needed.
* worker - an optional cluster defined role generally needed by platforms for cluster operations.  Different roles can be defined as needed.

Note: cluster/manager and cluster/leader both default to false.

Leader and worker are provided for convenience.  For more complex operators, automation designers may choose to define their own roles; however, they are common neutral terms and the assumed roles for this documentation.

.. _rs_cluster_pattern_shared:

Shared Data via Cluster Profile
-------------------------------

The Cluster Profile is shared Profile that has been assigned to all Machines in the cluster including the Manager.  The profile is self-referential: it must contain the name of the profile in a parameter so that machine action will be aware of the shared profile.

The Digital Rebar API has special behaviors allow machines to modify these templates including an extention for Golang template rendering (see :ref:`rs_data_render`) to include ``.GenerateProfileToken``.  This special token must be used when updating the shared template.

For example, if we are using the Profile ``example`` to create a cluster, then we need to include the Param ``cluster/profile: example`` in the Profile.  While this may appear redundant, it is essential for the machines to find the profile when they are operating against it.

Typically, all cluster scripts start with a "does my cluster profile exist" stanza from the ``cluster-utilities`` template.  

Which is typically included in any cluster related task with the following:

  ::

    {{ template "cluster-utilities.tmpl" .}}

The cluster-utilities template has the following code to initialize variables for cluster tasks.

  ::

    {{ if .ParamExists "cluster/profile" }}
    CLUSTER_PROFILE={{.Param "cluster/profile"}}
    PROFILE_TOKEN={{.GenerateProfileToken (.Param "cluster/profile") 7200}}
    echo "  Cluster Profile is $CLUSTER_PROFILE"
    {{ else }}
    echo "  WARNING: no cluster profile defined!  Run cluster-initialize task."
    {{ end }}


Adding Data to the Cluster Profile
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

As data collects on the cluster profile from the manager or other members, it is common to update shared data as Params in the cluster profile.  This makes the data available to all members in the cluster.  

  ::

    drpcli -T $PROFILE_TOKEN profiles add $CLUSTER_PROFILE param "myval" to "$OUTPUT"

Developers should pay attention to timing with Param data.  Params that are injected during template rendering (e.g.: ``{{ .Param "myval" }}``) are only evaluated when the job is created and will not change during a task run (aka a job).

If you are looking for data the could be added or changed inside a job then you should use the DRPCLI to retrieve the information from the shared profile with the ``-T $PROFILE_TOKEN`` pattern.

.. _rs_cluster_pattern_patch:

Resolving Potential Race Conditions via Atomic Updates with JSON PATCH
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

When cases where multiple machines write data into the Cluster Profile, there is a potential for race conditions.  The following strategy is used to address these cases in a scalable way.  No additional tooling is required.

The Digital Rebar CLI and UX use JSON PATCH (https://tools.ietf.org/html/rfc6902) instead of PUT extensively.  PATCH allows atomic field-level updates by including tests in the update.  This means that simulataneous upates do not create "last in" race conditions.  Instead, the update will fail in a predictable way that can be used in scripts.

The DRPCLI facilitates use of PATCH for atomic operations by allowing scripts to pass in a reference (aka pre-modified) object.  If the ``-r`` reference object does not match then the update will be rejected.

This allows machines take actions that require synchronization among the cluster when waiting on operations to finish on other machines.  This requirement is mitigated by the manager pattern

The following example shows code that runs on all maachines but only succeeds for the cluster leader.  It assumes the Param ``my/data`` is set to default to "none".

  ::

    {{template "setup.tmpl" .}}
    cl=$(get_param "my/data")
    while [[ $cl = "none" ]]; do
      drpcli -r "$cl" -T "$PROFILE_TOKEN" profiles set $CLUSTER_PROFILE param "my/data" to "foo" 2>/dev/null >/dev/null && break
      # sleep is is a hack but it allows for backoffs
      sleep 1
      # get the cluster info
      cl=$(get_param "my/data")
    done


.. _rs_cluster_pattern_workflow:

Cluster Manager Workflow
------------------------

The cluster manager workflow is a mostly a typical DRP workflow that runs automation and collects data.  The primarly difference is the addition of multi-machine task(s) that also initiate and coordinate other machines.  A Cluster manager workflow may also have preparatory tasks and subsequent tasks around the cluster construction.

For example, a cluster workflow may setup/verify the shared profile for the cluster.  If security or credentials are required then tasks can be used to collect that information in advance.  After the cluster install reaches a critical, the cluster workflow can perform cluster level configuration such as running Kubernetes or VMware API calls.  In this regard, the cluster manager's outside of the cluster frame of reference is a helpful for cluster operations, checks and synchronization.


.. _rs_cluster_pattern_collect:

Cluster Filter to Collect Members
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The ``cluster/filter`` Param plays a critical role in allowing the cluster manager to collect members of the cluster.  The filter is a DRPCLI string that is applied to a ``DRPCLI machines list`` or ``DRPCLI machines count`` call to indentify the cluster membership.

This process is baked into the helper routines used for the cluster pattern and should be defined in the cluster profile if the default is not sufficient.  By default, the ``cluster/filter`` is set to ``Profiles Eq $CLUSTER_PROFILE`` and will select all the machines attached to the cluster profile including the manager.  Deveopers may choose to define clusters by other criteria such as Pool membership, machine attributes or Endpoint.


This shows how ``cluster/filter`` can be used in a task to collect the cluster members including the manager.  ``--slim`` is used to reduce the return overhead.

  ::
      CLUSTER_MEMBERS="$(drpcli machines list {{.Param "cluster/filter"}} --slim Params,Meta)"


In practice, additional filters are applied to further select machines based on cluster role or capability (see below).

.. _rs_cluster_pattern_startloop:

Starting Workflow on Cluster Members
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

During the multi-machine task(s), a simple loop can be used to start Workflows on the targeted members.

This example shows a loop that selects all members who are cluster leaders (``cluster/leader Eq true``) and omits the cluster manager as a safe guard (``cluster/manager Eq false``).  Then it apples the target workflow and sets an icon on each leader.

  ::

      CLUSTER_LEADERS="$(drpcli machines list cluster/manager Eq false cluster/leader Eq true {{.Param "cluster/filter"}} --slim Params,Meta)"
      UUIDS=$(jq -rc ".[].Uuid" <<< "$CLUSTER_LEADERS")
      for uuid in $UUIDS; do
        echo "  starting k3s leader install workflow on $uuid"
        drpcli machines meta set $uuid key icon to anchor > /dev/null
        drpcli machines workflow $uuid k3s-machine-install > /dev/null
      done

Since these operations are made against another machine, multi-machine task(s) need to be called with an ``ExtraClaims`` definition that allows * actions for the ``scope: machines``.

Working with Cluster Roles
~~~~~~~~~~~~~~~~~~~~~~~~~~

As discussed above, the cluster pattern includes three built in roles: manager, leader and worker (assumed as not-leader and not-manager).  The cluster/leaders are selected randomly during the ``cluster-initialize`` when run on the cluster manager.  The default number of leaders is 1.

Developers can define additional roles by defining and assigning Params to members during the process.  The three built in roles are used for reference.

.. _rs_cluster_pattern_milestones:

Coordinating Activity with Workflow Complete Milestones
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

An important aspect of cluster building is synchronizing activity between the members.  Generally, this is as simple as waiting for all the non-manager cluster members to complete their current workflow.  When all members are ``WorkflowComplete: true`` then it should be safe to start the next sequential or parallel activity.

To make this process easier, the ``cluster-utilities.tmpl`` template has a ``cluster_wait_for_workflow_complete`` function that uses the ``cluster/filter`` to safely wait for all machines in the cluster to complete.  The process assumes that the cluster developer always wants to wait for all workflows to complete before starting another machine activity.


  ::

    {{ template "cluster-utilities.tmpl" .}}

    echo "starting k3s leader install"
    cluster_wait_for_workflow_complete


By design, the cluster wait loop will exit if the cluster manager is set to ``Runnable: false``.  This provides a natural and easy break out method for operators.