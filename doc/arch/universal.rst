.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Universal Workflow Architecture

.. _rs_universal_arch:

Universal Workflow Architecture
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

This section will address architecture of the Universal Workflow system.  The operational methods for Universal Workflow
is described at :ref:`rs_universal_ops`.

Universal Workflow implements the following goals:

  * Provide accumulated best practices for managing systems through full life-cycle
  * Provide all features of the RackN system without requiring cloned workflows and stages
  * Provide profile-based configuration for altering a solid baseline

Universal Workflow is designed to allow for sets of workflows to be chained together to generate a dynamic system for
building and maintaining systems.

Each Universal workflow follows a common format and layout.  This layout allows for extension within the workflow
and chaining among workflows.  This is not required, but helps define a standard common form.

Requirements
============

The universal workflow system requires some additional content packages to function.  Minimally, universal requires:

  * drp-community-content
  * callback
  * classify
  * flexiflow
  * validation
  * task-library

Some workflows require additional elements, e.g. image-deploy.  At
the moment, these are defined as requirements, but may move into their owning content packages or plugins.

Universal Application
=====================

Universal Workflow revolves around the concept of an application.  The application can be anything that a set of machines
can be driven to.  An application can be defined to install ESXi or install Centos8 or install CentOS8 after configuring
the hardware and afterwards add a webserver.  The application is a helper to define a set of things to do.  In almost
all cases, these tasks are reused across different applications.

For example, the default content provides the `centos-8` application. This is provided as a profile that sets the
`universal/application` parameter to `centos-8` and other values to do a kickstart-based centos-8 linux install.
This includes doing a full inventory and discovery, hardware, bios, and raid configuration, and linux installation.
All these paths will default to minimal actions unless other profiles are provided or "classified" onto the system.
To use it, apply the `universal-application-centos-8` profile to a machine and start the `universal-discover` workflow
on the system.  Once the machine's WorkflowComplete field becomes true, the system will be an installed centos-8 server.

Users can define their own applications and extensions as needed to deploy what they need.  More on this in the
:ref:`rs_universal_ops` section.

Layout
======

Each chainable workflow should have to following general layout.  The goal of this layout is to enable tracking of the
machine through process, validation of the workflow, and allow for extension to the defaults.

To allow for notifications, the workflow should start and end (second to last) with a callback stage that defines the
action and that an event should also be generated.  The stage should have the following parameters set:

  * callback/action - The value should be <workflow>-start or <workflow>-complete.
  * callback/event - Should the callback generate an event or not.  Boolean parameter defaulting to false.

To allow for validation, the workflow should have a validation stage that is similar to `validation-post-discover`.
The validation content pack has some available examples.  The universal workflow will provide default tasks for validation
but additional tasks can be added as needed.

To allow for extension and customization for things not covered by the default universal usage, each workflow should
have a pre- and post- FlexiFlow stage to allow for additional tasks as users need.

To allow for dynamic path changes, each universal workflow should have a classifier and allow for additions of additional
classifiers to alter the machine's state for changing the path through the workflows.

The basic layout is:

  * Callback Entry into Workflow
  * FlexiFlow Pre Tasks Extensions
  * Workflow Specific Stages
  * FlexiFlow Post Tasks Extensions
  * Classification
  * Workflow Validation
  * Callback Complete into Workflow
  * FlexiFlow Chain Workflow

The `discover` layout would have the following stages with parameters with their defaults:

  * universal-discover-start-callback
    * Parameter setting in stage: callback/action = universal-discover-start
    * Parameter setting in stage: callback/event = true
  * universal-discover-pre-flexiflow
    * Parameter setting in stage: flexiflow/list-parameter = universal-discover-pre-flexiflow
    * Parameter default: universal-discover-pre-flexiflow = []
  * The default stages for `discover` are the inventory and discovery components.
  * universal-discover-post-flexiflow
    * Parameter setting in stage: flexiflow/list-parameter = universal-discover-post-flexiflow
    * Parameter default: universal-discover-post-flexiflow = []
  * universal-discover-classification
    * Parameter setting in stage: classify/stage-list-parameter = universal/discover-classification-list
    * Parameter default: universal/discover-classification-list = [ "universal-discover-classification-base" ]
  * universal-discover-post-validation
    * Parameter setting in stage: validation/list-parameter = universal-discover-post-validation
    * Parameter default: universal-discover-post-validation = []
  * universal-discover-complete-callback
    * Parameter setting in stage: callback/action = universal-discover-complete
    * Parameter setting in stage: callback/event = true
  * universal-chain-workflow

Defaults and Overrides
======================

The parameters can be override by the profiles to update the various parts.  In general, all the lists and parameters
default to empty.  This is true for all except classification.  Each classifier starts with a default classifier for each
workflow.  This workflow references these parameters to define universal/discover-classification-base-data and
universal/discover-classification-base-functions.  For discover, this defaults to a set of classification actions
that do the following:

  * set up hardware param - this sets the universal/hardware parameter to a derived string
  * set universal/application - this converts the rack/build parameter to the universal/application parameter
  * apply universal application profile - this converts the universal/application parameter into a profile name and applies that profile to the machine.
  * a set of tests to find hardware profiles - this converts a set of parameters into hardware specific profile names and applies them.

These patterns are tested and if found applied.  This way a hardware specific profile is applied.

  * universal-bom-<rack/bom>-<universal/hardware>-<universal/application>
  * universal-hw-<rack/bom>-<universal/hardware>-<universal/application>
  * universal-bom-<universal/hardware>-<universal/application>
  * universal-hw-<universal/hardware>-<universal/application>
  * universal-bom-<universal/hardware>
  * universal-hw-<universal/hardware>
  * universal-bom-<universal/application>
  * universal-hw-<universal/application>

All the other classifiers default to no actions.

See the operations documentation for examples and usage.

Workflow Chaining
=================

The other main goal of the universal workflow system is to allow for workflow chaining.  This allows building up of
a consistent set of workflows that do pieces of infrastructure management and chain them based upon the universal application.

The `universal-chain-workflow` stage uses the `universal/application` parameter to lookup in the `universal/workflow-chain-map`
parameter to figure out what the next workflow is after the current one.  Additionally, there is a parameter
`universal/workflow-chain-index-override` that allows the lookup to occur overriding `universal/application`.
The whole map can be overridden by the `universal/workflow-chain-override`.  All this gets a single map
that is used to look up the current workflow and see if there is a next workflow to set.  If nothing is found, the
workflow continues to completion.  Otherwise, the new workflow is applied.

There is a special case for the `universal-hardware` workflow.  Using the `universal/maintenance-mode` parameter, the
workflow chain will be forced to the `universal-local` workflow that drives the system back the currently installed disk.
This allows for hardware maintenance without being destructive to install application.  The process unsets maintenance mode
as part of the processing.

