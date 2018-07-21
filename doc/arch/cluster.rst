.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Cluster

.. _rs_cluster_pattern:

Multi-Machine Cluster Pattern
=============================

A clear pattern has emerged for Digital Rebar Provision to build application clusters.  This pattern allows machines to coordinate activities using atomic storage operations against the Digital Rebar Provision API.  Having a safe, atomic shared storage area allows scripts on machines working in parallel to share data and synchronize activities.

The critcal elements of this pattern are using 1) a shared cluster profile on all the machines in the cluster and 2) using JSON PATCH with a reference model to ensure atomic updates.

:ref:`component_task_library` content has a set of :ref:`component_task_library_cluster` that implement this behavior in a generic way.

Cluster Profile
---------------

The pattern using a shared Profile that has been assigned to all Machines in the cluster.  The profile is self-referential: it must contain the name of the profile in a parameter so that machine action will be aware of the shared profile.

For example, if we are using the Profile ``example`` to create a cluster, then we need to include the Param ``cluster-profile: example`` in the Profile.  While this may appear redundant, it is essential for the machines to find the profile when they are operating against it.  Typically, all cluster scripts start with a "does my cluster profile exist" stanza:

  ::

    {{if .ParamExists "cluster-profile" -}}
    CLUSTER_PROFILE={{.Param "cluster-profile"}}
    PROFILE_TOKEN={{.GenerateProfileToken (.Param "cluster-profile") 7200}}
    {{else -}}
    echo "Missing cluster-profile on the machine!"
    exit 1
    {{end -}}

The Digital Rebar API has special behaviors allow machines to modify these templates including an extention for Golang template rendering (see :ref:`rs_data_render`) to include ``.GenerateProfileToken``.  This special token must be used when updating the shared template.

Atomic Updates with JSON PATCH
------------------------------

The Digital Rebar CLI and UX use JSON PATCH (https://tools.ietf.org/html/rfc6902) instead of PUT extensively.  PATCH allows atomic field-level updates by including tests in the update.  This means that simulataneous upates do not create "last in" race conditions.  Instead, the update will fail in a predictable way that can be used in scripts.

The DRPCLI facilitates use of PATCH for atomic operations by allowing scripts to pass in a reference (aka pre-modified) object.  If the ``-r`` reference object does not match then the update will be rejected.

This allows machines take actions that require synchronization among the cluster such as electing leaders or waiting on operations to finish on other machines.  For example, it is typical for clustered machines to poll a param waiting for it to be set by a cluster leader.

The following example shows code that runs on all maachines but only succeeds for the cluster leader.  It assumes the Param ``my/leader`` is set to default to "none".

  ::

    {{template "setup.tmpl" .}}
    cl=$(get_param "my/leader")
    while [[ $cl = "none" ]]; do
      drpcli -r "$cl" -T "$PROFILE_TOKEN" profiles set "cluster-name" param "my/leader" to "$RS_UUID" 2>/dev/null >/dev/null && break
      # sleep is is a hack but it allows for backoffs
      sleep 1
      # get the cluster info
      cl=$(get_param "my/leader")
    done

Collecting Data in Profile
--------------------------

As automation is run on the cluster profile, it is common to update Params in the cluster profile.  This is a good place to make information available for operators or other machines in the script instead of storing on the specific machine elected leader.

Showing Progress and Leaders
----------------------------

To show cluster install progress, common practice is to set the Machine.Meta icon and color.  This provides a fast reference back to operators about the state of the cluster without having to open the profile.

  :: 
    drpcli machines update $RS_UUID "{\"Meta\":{\"color\":\"purple\", \"icon\": \"anchor\"}}" | jq .Meta

Scripts are often updated so that the elected leader(s) have a distinct icon or color from the rest of the cluster.