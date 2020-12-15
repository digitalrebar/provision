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
  * Classification (single entry for now)
  * Workflow Validation
  * Callback Complete into Workflow
  * FlexiFlow Chain Workflow

# GREG: Review this - general naming + the classifier pieces those are wrong

The `discover` layout would have the following stages with parameters with their defaults:

  * universal-discover-start-callback
    * Parameter setting in stage: callback/action = universal-discover-start
    * Parameter setting in stage: callback/event = true
  * universal-discover-pre-flexiflow
    * Parameter setting in stage: flexiflow/list-parameter = universal-discover-pre-flexiflow
    * Parameter default: universal-discover-pre-flexiflow = []
  * The default stages from the `discover` workflow
  * universal-discover-post-flexiflow
    * Parameter setting in stage: flexiflow/list-parameter = universal-discover-post-flexiflow
    * Parameter default: universal-discover-post-flexiflow = []
  * universal-discover-classification
    * Parameter setting in stage: classify/stage-list-parameter = universal/discover-classification-list
    * Parameter default: universal/discover-classification-list = []
  * universal-discover-post-validation
    * Parameter setting in stage: validation/list-parameter = universal-discover-post-validation
    * Parameter default: universal-discover-post-validation = []
  * universal-discover-complete-callback
    * Parameter setting in stage: callback/action = universal-discover-complete
    * Parameter setting in stage: callback/event = true
  * flexiflow-chain-workflow

