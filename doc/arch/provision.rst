.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Provision Models


.. _rs_provisioning_models:

Provisioning Models
<<<<<<<<<<<<<<<<<<<

These models work together to manage any and all lifecycle needs for
managing Machines in *dr-provision*. This includes:

- Keeping track of what machines are being managed by *dr-provision*.
- Controlling what OS environment any given Machine will boot to over
  the network.
- Managing the order in which Tasks will be run on Machines.
- Managing the Context in which machine tasks will be run.
- Making sure that any files that are needed to complete the
  provisioning process are available and valid.

.. _rs_data_template:

Template
--------

Template expansion underlies just about everything that DigitalRebar
Provision does.  All template expansion in *dr-provision* happens
accroding to the rules defined by the `golang text/template
<https://golang.org/pkg/text/template/#hdr-Actions>`_ package, which
(along with this document) is required reading if you want to build
content for *dr-provision*. Template objects define common template
content that different parts of *dr-provision* can reuse as they see
fit.  Template objects contain the following fields:

- **ID**: A unique identifier for the Template.
- **Contents**: The contents of the Template.  It must be parseable as a
  go text/template.

.. _rs_data_templateinfo:

TemplateInfo
------------

Closely related to the :ref:`rs_data_template` is the TemplateInfo
object, which is included as part of the :ref:`rs_data_bootenv`,
:ref:`rs_data_stage`, and :ref:`rs_data_task` objects.  TemplateInfo
objects have the following fields:

- **Name**: The name of this TemplateInfo.

- **Path**: A string that will be expanded as if it were a
  :ref:`rs_data_template` to generate a path for the template.

- **Contents**: If present, a string that will be expanded as if it were a
  :ref:`rs_data_template` to generate the file that will be made
  available at the location indicated by the Path field.  Contents
  must be empty or not present if ID is set.

- **ID**: If present, the ID of the :ref:`rs_data_template` that will be
  used to generated the file that will be made available at the
  location indicated by the Path field.  ID must be empty or not
  present if Contents is set.

.. _rs_data_render:

Rendering Templates
-------------------

Whenever *dr-provision* needs to render something as a template (whether
or not it is a Template object, a TemplateInfo object, or just a
string that might contain a template), it always does so in the
context of a RenderData object, which provides a slew of useful helper
functions along with references to the applicable objects.  RenderData
is what *dr-provision* uses for `dot` or `{{ . }}` when executing a
template.  RenderData has the following fields:

- **Machine**: the Machine that we are rendering templates for.  Except
  for rendering the unknownBootEnv, all template rendering that
  *dr-provision* does happens against a machine and one of a BootEnv, a
  Task, or a Stage.  If Machine is present, the following helpers are
  present on RenderData:

  - **.Machine.Path** returns a machine-specific Path fragment (based on
    the Machine UUID) that can be used to store or refer to machine
    specific information via the static file server or via TFTP. It is
    particularly useful for ensuring that templates are expanded into
    a unique file space for each machine by using it in a TemplateInfo
    Path field.

  - **.Machine.Address** returns the IP address of the Machine as
    recorded in the Lease or Reservation.

  - **.Machine.HexAddress** returns the IP address of the Machine in hex
    format, suitable for use by anything expecting a hex encoded IP
    address.

  - **.Machine.Url** returns a machine specific http URL that can be used to
    access machine specific information via http.

  - **.ParamExists <key>** returns true if the specified key is a valid
    parameter available for this rendering.

  - **.Param <key>** returns the value for the specified key for this
    rendering.  .Param and .ParamExists always look parameters up in the following order:

    1. Params set directly on a Machine.

    2. Params set on the Profiles that have been added to a Machine,
       in the order of that Machine's Profiles list.

    3. Params set on the Profiles added to the Stage that the Machine
       is currently in, in the order of that Stage's Profile list.

    4. The current default Profile.

    5. The default value defined as part of the JSON schema for the Param.

    Param returns values as simple strings! For complex output, look at
    .ParamAsJSON and .ParamAsYAML below.

  - **.ParamCompose <key>** like Param above, except that it returns the
    composed value of a parameter across *all sources* instead of choosing the
    first from the Param precedence.

  - **.ParamExpand <key>** returns the value for the specified key for this
    rendering, but then re-expands the value again through the renderer.
    If the value is an array or object, the value is recursively processed to
    expand internal strings.  Additional similar functions include:

    - **.StringExpand** does templating on the contents of a string value
    
    - **.ObjectExpand** using an arbitrary object and does an expand on it recursively.
      If the object does not understand expansion, it will be returned as is.
      If the object is an array or slice or map, the system will recurse on
      each element.

  - **.ParamComposeExpand <key>** like ParamExpand above, except that it returns the
    composed value of a parameter across *all sources* instead of choosing the
    first from the ParamExpand precedence.

  - **.ParamAsJSON <key>** returns the value for the specified key for this
    rendering preserved in JSON formatting.  This is important for templates
    that rely on ``jq`` or other commands that need consistent formatting

    Note: .ParamAsJSON will use the .Param lookup order above.

  - **.ParamAsYAML <key>** returns the value for the specified key for this
    rendering preserved in YAML formatting.  This is important for configuration
    files and templates that need consistent formatting

    Note: .ParamAsYAML will use the .Param lookup order above.

  - **.Repos <tag>, <tag>,...** returns Repos (as defined by the
    package-repositories param currently in scope) with the matching
    tags.

  - **.MachineRepos** will return a list of OS package repositories that
    can be used to install packages on the Machine.  The repos
    returned will be for .Machine.OS

  - **.InstallRepos** will return at most one OS package repository that
    can be used to install an OS from, and at most one OS package
    repository that contains security updates to apply during OS
    install.

- **[Sprig functions]** are string, math, file and flow functions for golang
  templates from the `Sprig Function Library <http://masterminds.github.io/sprig/>`_.
  They can be added to pipeline evaluation to perform useful template
  rendering operations.

- **Provisioner** items are specific to the Provisioner

  - **.ProvisionerAddress** returns the IP address to access
    the Provisioner based upon the requesting IP address.

  - **.ProvisionerURL**  returns a URL to access the file server part of the server using the
    requesting IP address as a basis.

  - **.ApiURL** returns a URL to access the api server part of the server using the
    requesting IP address as a basis.

  - **.Info.** models.Info structure

- **Env**: The BootEnv that we are rendering templates for, if applicable.
  Unless the BootEnv has the OnlyUnknown flag set, RenderData will
  also include a Machine.  If Env is present, the following helpers will also
  be present on RenderData:

  - **.Env.PathFor <proto> <partial>** is a helper that makes it easier to
    build paths that the client side shuld expect.  proto should be
    either **http** or **tftp**, and partial is a partial path
    relative to the root of a package repository.

  - **.Env.JoinInitrds <proto>** joins together a list of initrds in a way that
    is applicable for the passed in proto.

  - **.BootParams** returns a rendered version of .Env.BootParams.  It will be rendered
    against the current RenderData.

  - **.Env.OS.FamilyName**: The contents of .Env.OS.Family if present,
    otherwise the result of splitting .Env.OS.Name by hyphens and
    taking the first part.

  - **.Env.OS.FamilyVersion**: The contents of .Env.OS.Version if
    present, otherwise the result of splitting .Env.OS.Name by hyphens
    and taking the second part.

  - **.Env.OS.FamilyType**: The type of .Env.OS.FamilyName. rhel for
    distros based on RHEL, debian for distros based on Debian,
    otherwise the same as .Env.OS.FamilyName.  More return types will
    be added upon request.

  - **.Env.OS.VersionEq <testVersion>**: Splits testVersion and
    .Env.OS.FamilyVersion into pieces seperated by a period.  Returns
    true if .Env.OS.FamilyVersion has at least as many pieces as
    testVersion and all the pieces they have in common are numerically
    equal.

- **Task**: the Task we are rendering templates for, if applicable.
  RenderData will include a Machine.

- **Stage**: the Stage we are rendering templates for, if
  applicable. RenderData will include a Machine.

RenderData includes the following helper methods:

- **.ProvisionerAddress** returns an IP address that is on the provisioner
  that is the most direct access to the machine.
- **.ProvisionerURL** returns an HTTP URL to access the base file server
  root
- **.ApiURL** returns an HTTPS URL to access the Digital Rebar Provision
  API
- **.GenerateToken** generates either a **known token** or an **unknown
  token** for use by the template to update objects in Digital Rebar
  Provision.  The tokens are valid for a limited time as defined by
  the **knownTokenTimeout** and **unknownTokenTimeout**
  :ref:`rs_model_prefs` respectively.  The tokens are also restricted
  to the function the can perform.  The *known token* is limited to
  only reading and updating the specific machine the template is being
  rendered for.  If a machine is not present during the render, an
  *unknown token* is generated that has the ability to query and
  create machines.  These are used by the install process to indicate
  that the install is finished and that the *local* BootEnv should be
  used for the next boot and during the discovery process to create
  the newly discovered machine.
- **.GenerateInfiniteToken** works like **.GenerateToken**, but creates
  a token with a 3 year timeout.
- **.ParseURL <segment> <url>** parses the specified URL and return the
  segment requested.  Supported segments can be one of *scheme* (eg "https"),
  *host* (eg "drp.example.com:8092"), or *path* (eg "/api/v3/machines").
  (note: *host* does not separate name and port)
- **template <string> .** includes the template specified by the string.
  String can NOT be a variable and note that template does NOT have a dot
  (.) in front.
- **.CallTemplate <string> .** works like **template** but allows for
  template expansion inside the string to allow for dynamic template
  references.  Note that CallTemplate does have dot (.) in frount.

.. _rs_data_param:

Param
-----

Params are how *dr-provision* provides validation and a last-ditch
default value for data that we use during template expansion.
Strictly speaking, you do not have to define a Param in order to use
it during template expansion, but *dr-provision* will not be able to
enforce that param data is syntactically valid.  A Param object has
the following fields:

- **Name**: The unique name of the Param.  Any time you update a Profile
  or add, remove, or change a parameter value on another object,
  *dr-provision* will check to see if a Param exists for the
  corresponding parameter key.

- **Schema**: A JSON object that contains a valid
  `JSONSchema <http://json-schema.org/>`_ (draft v4 or higher) that
  describes what a valid value for the Param looks like.  You may also
  provide a default value for the Param using the `default` stanza in
  the JSON schema.

- **Secure**: Data managed in this param must be handled in a secure
  fashion.  It will never be passed in cleartext over the API without
  proper Role based authorization, will be stored in an encrypted
  wrapper, and will only be made available in an unencrypted form for
  schema validation on the server, performing plugin actions, and
  running Tasks on a machine.

.. _rs_data_param_secure:

Secure Params
~~~~~~~~~~~~~

Secure param management is a licensed feature.  You must have a
license with the **secure-params** feature enabled to be able to
create and retrieve secure param values.  SecureData uses a simple
encryption mechanism based on the NACL Box API (as implemented by
libsodium, golang.org/x/crypto/nacl/box, tweetnacl-js, PyNaCl, and
many others), using curve25519 and xsalsa20 for crypto, and poly1305
for message verification.


Secure params are handled by
the API and stored on the backend using a SecureData struct, which has
the following fields:

- **Payload**: The encrypted payload.  When marshalled to JSON, this
  should be converted to a base64 encoded string.

- **Nonce**: 24 cryptographically random bytes.  When marshalled to
  JSON, this should be converted into a base64 encoded string.

- **Key**: a 32 byte curve25519 ephemeral public key.  When marshalled
  to JSON, this should be converted to a base64 encoded string.

When a Param has the Secure flag, the following additional steps must be
taken to set and get values for this param on objects that hold params.

.. _rs_data_setting_param_secure:

Setting Secure Param Values
===========================

1. Get the peer public key for the object you want to set a secure param on
   from its `pubkey` endpoint.  These endpoints are at
   `/api/v3/<objectType>/<objectID>/pubkey` -- as an example, the
   pubkey endpoint for the global profile is
   `/api/v3/profiles/global/pubkey`.  Access to these API endpoints
   requires an appropriate Claim with the **updateSecure** action.
   These API endpoints return a JSON string containing the base64
   encoding of an array containing 32 bytes.

2. Generate local ephemeral curve25519 public and private keys using a
   cryptographically secure random number source.

3. Generate a 24 byte nonce using a cryptographically secure random
   number source.

4. Encrypt the JSON-marshalled param using the nonce, the peer public
   key, and the ephemeral private key.

5. Generate a SecureData struct with **Key** set to the ephemeral
   public key, **Nonce** set to the generated nonce, and **Payload**
   set to the encrypted data.

6. Use the SecureData struct in place of the raw param value when
   making API calls to add, set, or update params.

.. _rs_data_retrieving_param_secure:

Retrieving Decrypted Secure Data Values
=======================================

In order to retrieve decrypted secure data values, you must have an
appropriate Claim with the **getSecure** action.  That will allow you
to make GET requests to the params API endpoints for param-carrying
objects with the `decode=true` query parameter.  That will cause the
frontend to decrypt any encryped parameter values before returning
from the API call.

.. _rs_data_task:

Task
----

Tasks in *dr-provision* represent the smallest discrete unit work that
the machine agent can use to perform work on a specific machine.  The
machine agent creates and executes a Job for each Task on the
machine. Tasks have the following fields:

- **Name**: The unique name of the task.

- **RequiredParams**: A list of parameters that are required to be present
  (directly or indirectly) on a Machine to use this Task.  It is used
  to verify that a Machine has all the parameters it needs to be able
  to execute this Task.

- **OptionalParams**: A list of parameters that the Task may use if
  present (directly or indirectly) on a Machine.

- **Templates**: A list of TemplateInfos that will be rendered into Job
  Actions when the machine agent starts exeuting this Task as a Job.

- **Prerequisites**: A list of Tasks that must be run in the current BootEnv
  before this task can be run.  dr-provision will not allow a cyclical
  prerequisite -- task cannot have themselves as prerequisites, either directly
  or indirectly.

Rendering a Task for a Machine
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The Templates for a Task are rendered for a specific Machine whenever
the Actions for the Job for that particular task/machine combo are
requested.

All referenced templates can refer to each other by their ID (if
referring to a Template object directly), or by the TemplateInfo Name
(if the TemplateInfo object), in addition to all the Template objects
by ID.

Template Prerequisite Expansion
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

When a Task is added to a Task list, its fully expanded list of
prerequisite tasks are expanded, any tasks in that expanded list that
already appear in the machine task list in the same BootEnv are
discarded, and the resultant set of prerequisite tasks are inserted
just before the Task to be inserted.

Task Control Operations
-----------------------

To help with Task and Workflow debugging or advanced flow control, there are a
couple of parameters that can help with debug messages and task control.

Single Step
~~~~~~~~~~~

Setting the Machine parameter ``task-single-step`` to ``true`` will cause
the system to make the machine not runnable after each task.  Removing the
parameter or setting it to ``false`` will cause the system to resume normal
operations.

In UX v4.6, this capability can be set from the the machine debug view.


Stop At
~~~~~~~

Setting the Machine parameter ``task-stop-at`` to a list of tasks or
action entries will cause the system to stop prior to those tasks.  Setting
the machine to runnable will cause the tasks to continue.

In UX v4.6, this capability is used to enable break points in the machine debug view.

Retry
~~~~~

Setting the Machine parameter ``task-retry`` to a map of key/value pairs
where the keys are task names and the values are integer number of times to
retry the failing task.  The system will retry that number of times before
failing.  The system will cause an exponential backoff starting at 2 seconds
between retries.  There is no maximum wait time.

The Machine field ``RetryTaskAttempt`` tracks the retry for the current task.
This is reset to 0 for every new task.

For example, if the flaky task is ``always-fails`` then setting the ``task-retry`` parameter to ``{"always-fails": 3}`` will cause DRP to retry the task three times.


Error Handlers
~~~~~~~~~~~~~~

Setting the Machine parameter ``task-error-handlers`` to a map of key/value pairs
where the keys are task names and the values are lists of task names.  When a task
fails and all the retries have been attempted, the map will be consulted and if
present, the task list will be executed to and if successful, the failing task
will be restarted and retried again with a reset ``RetryTaskAttempt``.  Error
handling tasks can not be error handled.

The Machine field ``TaskErrorStacks`` is used to track the previous Task list and
CurrentTask to return to after running the error lists.


.. _rs_data_profile:

Profile
-------

Profiles are named collections of parameters that can be used to
provide common sets of parameters across multiple Machines.  Profile
objects have the following fields:

- **Name**: The unique name of the Profile.
- **Params**: a map of param name -> param value pairs for this Profile.

.. _rs_data_stage:

Stage
-----

Stages are used to define a set of Tasks that must be run in a
specific order, potentially in a specific BootEnv.  Stages contain the
following fields:

- **Name**: The unique name of the Stage.

- **Templates**: A list of TemplateInfos that will be template-expanded
  for a Machine whenever it transitions to a new Stage.

- **RequiredParams**: A list of parameters that are required to be present
  (directly or indirectly) on a Machine to use this Stage.  It is used
  to verify that a Machine has all the parameters it needs to be able
  to boot using this Stage.

- **OptionalParams**: A list of parameters that the Stage may use if
  present (directly or indirectly) on a Machine.

- **BootEnv**: The boot environment that the Stage must run in.  If this
  field is empty or blank, the assumption is that the Stage will
  function no matter what environment the machine was booted in.
  Changing the Stage of a Machine will always change the boot
  environment of the machine to the one that the stage needs, if any.

- **Profiles**: This is a list of Profile names that will be used for param
  resolution at template expansion time.  These profiles have a higher
  priority than the default profile,and a lower priority than profiles
  attached to a Machine directly.

- **Tasks**: This is a list of Task names that will replace the Tasks list
  on a Machine whenever the Machine switches to using this Stage.

- **Reboot**: DEPRECATED. This flag indicates whether or not the
  Machine must be rebooted if a Machine switches to this Stage.
  Generally, if this flag is set the Stage will also have a specific
  BootEnv defined as well.  While this flag is still honored, the
  runner will automatically reboot the machine as needed to satisfy
  the BootEnv of the Stage.

- **RunnerWait**: DEPRECATED. This flag used to indicate that the
  machine agent should wait for more Tasks to be added to the Machine
  once it finishes runnning the Tasks for this Stage.  The runner will
  currently always wait unless it is explicitly told to exit by an
  entry in the change-stage/map (also deprecated), or by the exit
  status of a Task.

Rendering a Stage for a Machine
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The Stage for a Machine is rendered *dr-provision* starts up, whenever a
Machine changes to a different Stage, or whenever a Stage referred to
by a machine changes.

All of the templates referred to by the Templates section of the
Stage will be rendered as static files available over the http and
tftp services of the provisioner at the paths indicated by each entry
in the Templates section.  All referenced templates can refer to each
other by their ID (if referring to a Template object directly), or by
the TemplateInfo Name (if the TemplateInfo object), in addition to all
the Template objects by ID.

.. _rs_data_context:

Context
-------

Contexts allow dr-provision to encapsulate the idea that we can run
Tasks for a Machine in several different places, of which the default
is on the Machine itself.  Contexts have the following fields:

- **Name**: The unique name of the Context.

- **Engine**: The name of the Plugin that provides the functionality
  needed to manage the execution environment that Tasks run in on
  behalf of a given Machine in the Context.  An Engine could be a
  Plugin that interfaces with Docker or Podman locally, Kubernetes,
  Rancher, or AWS, or any number of other things.

- **Image**: The name of the template that the Engine should use to create
  execution environments for this Context when Tasks should run on behalf
  of a Machine.  Images must contain all the tools needed to run the Tasks
  that are designed to run in them, as well as a version of drpcli
  with a context-aware `machines processjobs` command.

Plugins that provide Engines
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

By convention, a Plugin that provides support for using Contexts must support the
following Actions:

- **imageUpload**: This Action runs against the Plugin, and takes two arguments:

  - **context/image-name** which is a string containing the name of
    the Image being uploaded

  - **context/image-path** which is a string containing the location
    the image is being uploaded from.

  This action must result in the artifact at context/image-path being
  made available to create execution environments from using
  context/image-name.

  In general, context/image-path should either be a URL that points to
  the location of the artifact, or a relative path indicating where
  the artifact is located at on dr-provision's static file server.

  By convention, images stored on the static file server should be
  stored at `files/contexts/<plugin-name>/<image-name>`, with no extra
  file extensions.  Doing so will allow the **agentStart** Action to
  automatically load the Image on demand, otherwise the image must
  have already been uploaded before executing the **agentStart**
  action.

  Examples:

  `drpcli plugins runaction docker-context imageUpload
  context/image-name foo:latest context/image-path
  https://your.image.repo/path/to/foo:latest` will have the
  `docker-context` plugin load the image `foo:latest` from the
  upstreeam repo.

  `drpcli plugins runaction docker-context imageUpload
  context/image-name foo:latest context/image-path
  files/contexts/docker-context/foo:latest` will have the
  `docker-context` plugin load the image `foo:latest` from
  dr-provision's static file server at
  `files/contexts/docker-context/foo:latest`

- **imageExists**: This Action runs against the Plugin, and takes one
  argument:

  - **context/image-name**, which is the name of the image we testing
    to see if it already exists.

  This action returns `true` if the image exists, `false` if it does
  not.

  Example:

  `drpcli plugins runaction docker-context imageExists
  context/image-name foo:latest` will return `true` if `foo:latest` is
  present, and `false` if it is not.

- **imageRemove**: This Action runs against the Plugin, and takes one
  argument:

  - **context/image-name**, which is the name of the image to remove.

  This action returns true if either the image was removed or the
  image did not exist, false and an error otherwise.

  Example:

  `drpcli plugins runaction docker-context imageRemove
  context/image-name foo:latest` will have `docker-context` remove the
  image `foo:latest` if it exists.

- **agentStart**: This Action runs against the Machine, and takes two
  arguments:

  - **context/image-name**, which is the name of the image to use

  - **context/name**, which is the name of the Context that should be
    used to run Tasks for the Machine the Action was invoked with.

  This action causes the Engine to start a new execution context based
  on the Image, and arrange for `drpcli processjobs` to be called with
  the following environment variables set:

  - **RS_UUID**: the UUID of the Machine the Action was invoked with.

  - **RS_TOKEN**: an auth token suitable for Machines running Tasks
    against themselves.  Unless otherwise required by the jobs to be
    run in the Context, this token should be the results of calling
    `/api/v3/machines/<uuid>/token?ttl=3y` when the agentStart Action
    is called.

  - **RS_CONTEXT**: The name of the Context that the Agent should use
    when listening for machine state changes and job creation
    requests.

  - **RS_ENDPOINTS**: A space-seperated list of all the endpoints
    that the Agent should try to use when connecting to the
    dr-provision endpoint.  The Agent will stop at the first one that it
    successfully connects to.

  In general, you should only need to invoke this action manually when
  testing new images or testing a new plugin.  In the normal course of
  operation, plugins handle starting and stopping agents automatically
  based on changes on machine Context fields.

  Example:

  `drpcli machines runaction Name:test agentStart context/image-name
  foo:latest context/name foo` will start an agent for the machine
  running in the `foo` context on the `foo:latest` image.

- **agentStop**: This Action runs against the Machine, and takes one
  argument:

  - **context/name**, which is the name of the Context used in the
    corresponding agentStart action.

  This action tears down the execution context created by the
  corresponding agentStart Action.  In general, you should only need
  to invoke this action manually in the course of testing a new
  plugin, plugins handle stopping agents automatically based on
  Machine Context field state changes.

  Example:

  `drpcli machines runaction Name:test agentStop context/name foo`
  will stop the agent running in the `foo` context for the machine.

Additionally, all plugins that provide support for Contexts must
subscribe to the event stream that dr-provision emits.  They must
watch for machine delete, create, update, and save events to set up
and tear down execution environments as appropriate.

.. _rs_data_bootenv:

BootEnv
-------

Boot Environments (or BootEnv for short) are what DigitalRebar
Provision uses to model a network boot environment.  Each BootEnv
contains the following fields:

- **Name**: The name of the boot environment.  Each bootenv must have a
  unique name, and bootenvs that are responsible for booting into an
  environment that will install an OS on a machine must end in
  `-install`.

- **OnlyUnknown**: a boolean value indicating that this boot environment
  is tailored for use by unknown machines.  Most boot environments
  will not have this flag.

- **Loaders**: A map of indicated PXE boot architecture to the bootloader
  that should be used for that architecture.  Unlike the per-arch Loader
  field, this one is relative to the path the BootEnv is expanded to, and
  distinguishes between amd64 legacy boot and amd64 UEFI boot.

  Allowable keys:

  - **amd64-uefi**, used when the PXE options indicate that the system
    is attempting to PXE boot into a UEFI environment on an amd64 compatible system.

  - **arm64-uefi**, used when the PXE options indicate that the system is attempting
    to aPXE boot in an ARM64 based UEFI environment.

  - **386-pcbios**, used when the PXE options indicate that hte system is attempting to
    PXE boot into into a i386 compatible legacy BIOS environment.

   The values for these keys are interpreted as if they were appended to the appropriate env
   arch-specific path.

- **OS**: an embedded structure that contains some basic information on
  the OS that this BootEnv will boot into, if applicable.  OS contains
  the following fields:

  - **Name**: the name of the OS this BootEnv will boot into or install.
    It must be in the format of `distro-version`.  centos-7, debian-8,
    windows-2012r2, ubuntu-16.04 are all examples of what an OS name
    should look like.

  - **Family**: The family of the OS, if any.

  - **Codename**: The codename of the OS, if any.  Generally only really
    used by Debian, Ubuntu, and realted Linux distributions.

  - **Version**: The version of the OS, if any.

  - **IsoFile**: As an install convienence, DigitalRebar Provision
    contains built-in ISO expansion functionality that can be used to
    provide a local mirror for installing operating systems.  This
    field indicates the name of an install archive (usually a .iso
    file) that should be expanded to provide a local install repo for
    an operating system.

  - **IsoSha256**: If present, the SHA256sum that IsoFile should have.
  - IsoUrl: The URL that IsoFile can be downloaded from.

- **Kernel**: If present, a partial path to the kernel that should be used
  to boot a machine over the network.  The kernel must be specified as
  a relative path -- no leading / or .. characters are allowed.  As an
  example, the Kernel parameter for the community provided
  ubuntu-16.04-install boot environment is
  `install/netboot/ubuntu-installer/amd64/linux`, the path to the
  kernel relative to the root of the Ubuntu install ISO.

- **Initrds**: If present, a list of partial paths to initrds that should
  be loaded along with the Kernel when booting a machine over the
  network. Initrd paths follow the same rules as kernel paths.

- **BootParams**: If present, a string that will undergo template
  expansion as if it were a :ref:`rs_data_template`, and passed as
  arguments to the kernel when it boots.

- **RequiredParams**: A list of parameters that are required to be present
  (directly or indirectly) on a Machine to use this BootEnv.  Only
  applicable to bootenvs that do not have the OnlyUnknown flag set.
  It is used to verify that a Machine has all the parameters it needs
  to be able to boot using this BootEnv.

- **OptionalParams**: A list of parameters that the BootEnv may use if
  present (directly or indirectly) on a Machine.

- **Templates**: A list of templates that will be expanded and made
  available via static HTTP and TFTP for this BootEnv.  Each entry in
  this list must have the following fields:

  All bootenvs should include entries in their Templates list for the
  `pxelinux`, `grub`, and `ipxe` bootloaders.  If the OnlyUnknown
  flag is set, their Paths should expand to an appropriate location to
  be loaded as the fallback config file for each bootloader type,
  otherwise their Paths should expand to an approriate location to be
  used as a boot file for the loader based on the IP address of the
  machine.  Good examples for each are the `discovery
  <https://github.com/digitalrebar/provision-content/blob/v4/content/bootenvs/discovery.yml>`_
  and the `sledgehammer
  <https://github.com/digitalrebar/provision-content/blob/v4/content/bootenvs/sledgehammer.yml>`_
  bootenvs.

Rendering the unknownBootEnv
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The BootEnv for the unknownBootEnv preference is rendered whenever
*dr-provision* starts up or the BootEnv for the preference is changed.
It is the only time that templates are rendered without a Machine
being referenced, which is why BootEnvs that can be rendered this way
must have the OnlyUnknown flag set.

All of the templates referred to by the Templates section of the
BootEnv will be rendered as static files available over the http and
tftp services of the provisioner at the paths indicated by each entry
in the Templates section.  All referenced templates can refer to each
other by their ID (if referring to a Template object directly), or by
the TemplateInfo Name (if the TemplateInfo object), in addition to all
the Template objects by ID.

Rendering a BootEnv for a Machine
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The BootEnv for a Machine is rendered whenever *dr-provision* starts up,
whenever a Machine changes to a different boot environment, or
whenever a boot environment referred to by a machine changes.

All of the templates referred to by the Templates section of the
BootEnv will be rendered as static files available over the http and
tftp services of the provisioner at the paths indicated by each entry
in the Templates section.  All referenced templates can refer to each
other by their ID (if referring to a Template object directly), or by
the TemplateInfo Name (if the TemplateInfo object), in addition to all
the Template objects by ID.

.. _rs_provision_bootenv_customize:

Customization using bootenv-customize
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

As a tech preview feature, the 4.4.x series of dr-provision will allow
for programmatic overrides of various aspects of different bootenvs
on a per-machine (via the usual parameter inheritance mechanisms) or
global (via ovverides in the global profile) basis.

This param allows you to overlay dynamic customizations on top of BootEnvs.
Its intended use is to reduce the number of bootenvs you have to create
for what are ultimately trivial reasons.  This Param is structured as a
map of bootenv name -> override values.  You can override the following bootenv fields:

- **Loaders**

- **OS**

- **Templates**

The override values have the same names and semantics of the equivalent fields in the BootEnv,
with the following exceptions:

1. The Templates section is merged in with the Templates section from the BootEnv,
   with identically-named Templates from this Param overriding the ones from
   the BootEnv.

2. You must use the per-arch SupportedArchitectures section to override
   kernel, initrd, iso, boot parameter, and boot loader values from the
   BootEnv.

This feature is still in alpha, and this param definition and the semantics it enables
may change.  In particular, it is not compatible with bootenv purging right now, and
you may experience UI pauses when uploading a new ISO file.  You will also get unexpected
results if you create multiple overrides for the same BootEnv that have the same OS name
but reference different ISOs, or have different kernel/initrd settings.  You should also
not try to override Sledgehammer on a machine by machine basis, and any overrides for
sledgehammer must be matched by identical changes in the OS section of the discovery bootenv.

.. _rs_data_workflow:

Workflow
--------

A Workflow defines a series of Stages that a Machine should go
through.  It replaces the old change-stage/map mechanism of
orchestrating stage changes, which had the following drawbacks:

- change-stage/map is implemented as a map of currentStage ->
  nextStage:Action pairs.  This make it impossible for a machine to go
  through the same stage twice when going through a workflow.

- It was very easy to get the Action that the runner should perform
  wrong, leading to unexpected reboots or apparent hangs while walking
  through the Stages.  This has been replaced by making the Runner be
  smart enough to know that it must reboot on BootEnv changes to a
  machine, and by having the runner always wait for more tasks unless
  it is in an OS install BootEnv, or the Runner is directed to exit by
  a Task exit state.

- The Machine Tasks field only contained tasks for the current Stage,
  making it hard to see at a glance what Tasks will be executed for
  the entire workflow.

Workflows have the following fields:

- **Name**: The unique Name of the workflow.

- **Stages**: A list of Stages that any machine with this Workflow
  must go through.

When the Workflow field on a machine is set, the current task list on
the machine is replaced with the results of expanding each Stage in
the Workflow using the following items:

- stage:stageName
- bootenv:bootEnvName (if the stage has a non-empty BootEnv field)
- task0...taskN (the content of the Stage Tasks field)

Additionally, the Stage and BootEnv fields of the Machine become
read-only, as Stage and BootEnv transitions will occurr as dictated by
the machine Task list, and when the Stage changes it does not affect
the Task list.

.. _rs_data_machine:

Machine
-------

Machines are what DigitalRebar Provison uses to model a system as it
goes through the various stages of the provisioning process. As such,
Machine objects have many fields used for different tasks:

- **Name**: A user-chosen name for the machine.  It must be unique,
  although it can be updated at any point via the API.  It is a good
  idea for the Name field to be the same as the FQDN of the Machine in
  DNS, although nothing enforces that convention.

- **Uuid**: A randomly-chosen v4 UUID that uniquely identifies the
  machine.  It cannot be changed, and it what everything else in
  dr-provison will use to refer to a machine.

- **Address**: The IPv4 address that third-party systems should expect to
  be able to use to access the Machine.  *dr-provision* does not manage
  this field -- specifically, this does not have to be the same as an
  existing Lease or Reservation.

- **BootEnv**: The boot environment the Machine should PXE boot to the
  next time it reboots.  When you change the BootEnv field on a
  machine or change the BootEnv that a Machine wants to use, all
  relavent templates on the provisioner side are rerendered to reflect
  the updates.  The BootEnv field is read-only if the Workflow field
  is set.

- **Params**: A map containing parameter names and their associated
  values.  Params set directly on a machine override params from any
  other source when templates using those params are rendered.

- **Profiles**: An ordered list of profile names that the template render
  process will use to look up values for Params.  At render time,
  Profiles on a machine are looked at in the order that they appear in
  this list, and the first one that is found wins (assuming the Param
  in question is not provided directly on the Machine).

- **OS**: The operating system that the Machine is running.  It is only
  set by *dr-provision* when the Machine is moved into a BootEnv that
  has -install in the name.

- **Secret**: A random string used when generating auth tokens for this
  machine.  Changing this field will invalidate any existing auth
  tokens for this machine.

- **Runnable**: A flag that indicates whether the machine agent is allowed
  to create and execute Jobs against this Machine.

- **Locked**: A flag that indicates whether user-initiated changes to
  a Machine will be accepted.  When true, any changes that do not
  include change that sets this flag to false will be rejected.
  Changes from non-user sources will still be accepted -- this includes
  changes made while running tasks on a machine.

- **Workflow**: The name of the Workflow that the Machine is going
  through.  If the Workflow field is not empty, the Stage and BootEnv
  fields are read-only.

- **Tasks**: The list of tasks that the Machine should run or that
  have run.  You can add and remove Tasks from this list as long as
  they have not already run, they are not the current running Task, or
  they are beyond the next Stage transition present in the Tasks
  list.

- **CurrentTask**: The index in Tasks of the current running task.  A
  CurrentTask of -1 indicates that none of the Tasks in the current
  Tasks list have run, and a CurrentTask that is equal to the length
  of the Tasks list indicates that all of the Tasks have run.  The
  machine agent always creates Jobs based on the CurrentTask.  If the
  Workflow field is non-empty, setting this field to -1 will instead
  set this field to the most recent Stage in the Tasks list that did
  not initiate a BootEnv change.

- **Stage**: The current Stage the Machine is in.  Changing the Stage
  of a Machine has the following effects:

  - If the new Stage has a new BootEnv, the Machine Runnable flags
    will be set to False and the BootEnv on the Machine will change.

  - If the Machine Workflow field is empty, the Machine Tasks list
    will be replaced by the task list from the new Stage, and
    CurrentTask will be set back to -1.

  Note that the Stage field is read-only when the Workflow field is
  non-empty.

- **Context**: The name of the Context that tasks for the machine
  should execute in.  When this field is the empty string (""), then
  the agent running on the machine itself is the one that should be
  running tasks.  The Context field on a Machine can be set either via
  a machine PATCH or PUT, or as a side effect of processing an entry
  in the Tasks list in the format of 'context:new-context'.

  When the Tasks list on a machine is replaced as a side effect of a
  Workflow or Stage change, the Context field of the Machine will be
  reset to the value of the BaseContext Meta field on the Machine, or
  to the empty ("") string if that Meta field is not present on the
  Machine.

.. _rs_data_job:

Job
---

Jobs are what *dr-provision* uses to track the state of running
individual Tasks on a Machine.  There can be at most one current Job
for a Machine at any given time.  Job objects have the following
fields:

- **Uuid**: The randomly generated UUID of the Job.

- **Previous**: The UUID of the Job that ran prior to this one.  The Job
  history of a Machine can be traced by following the Previous UUIDs
  until you get to the all-zeros UUID.

- **Machine**: The UUID of the Machine that the job was created for.

- **Task**: The name of the Task that the job was created for.

- **Workflow**: The name of the Workflow that the job was created in

- **Stage**: The name of the Stage that the job was created in.

- **BootEnv**: The name of the BootEnv that the job was created in.

- **Context**: The name of the Context the Job was created for.

- **State**: The state of the Job.  State must be one of the following:

  - **created**: this is the state that all freshly-created jobs start at.

  - **running**: Jobs are automatically transitioned to this state by the
    machine agent when it starts executing this job's Actions.

  - **failed**: Jobs are transitioned to this state when they fail for any
    reason.

  - **finished**: Jobs are transitioned to this state when all their
    Actions have completed successfully.

  - **incomplete**: Jobs are transitioned to this state when an Action
    signals that the job must stop and be restarted later as part of
    its action.

- **ExitState**: The final disposition of the Job. Can be one of the
  following:

  - **reboot**: Indicates that the job stopped executing due to the machine
    needing to be rebooted.

  - **poweroff**: Indicates that the job stopped executing because the
    machine needs to be powered off.

  - **stop**: Indicates that the job stopped because an action indicated
    that it should stop executing.

  - **complete**: Indicates that the job finished.

- **StartTime**: The time the job entered the `running` state.

- **EndTime**: The time the Job entered the `finished` or `failed` state.

- **Archived**: Whether it is possible to retrieve the log the Job
  generated while running.

- **Current**: Whether this job is the most recent for a machine or not.

- **CurrentIndex**: The value of the Machine CurrentTask field when this Job was created.

- **NextIndex**: CurrentIndex++

.. _rs_data_job_action:

Job Actions
-----------

Once a Job has been created and transitioned to the running state, the
machine agent will request that the Templates in the Task for the job
be rendered for the Machine and placed into JobActions.  JobActions
have the following fields:

- **Name**: The name of the JobAction.  It is present for informational
  and troubleshooting purposes, and the name does not effect how the
  JobAction is handled.

- **Content**: The result of rendering a specific Template from a Task
  against a Machine.

- **Path**: If present, the Content will be written to the location
  indicated by this field, replacing any previous file at that
  location.  If Path is not present or empty, then the Contents will
  be treated as a shell script and be executed.
