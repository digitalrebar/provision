drpcli
======

A CLI application for interacting with the DigitalRebar Provision API

Synopsis
--------

A CLI application for interacting with the DigitalRebar Provision API

Options
-------

::

      -d, --debug               Whether the CLI should run in debug mode
      -E, --endpoint string     The Digital Rebar Provision API endpoint to talk to (default "https://127.0.0.1:8092")
      -f, --force               When needed, attempt to force the operation - used on some update/patch calls
      -F, --format string       The serialzation we expect for output.  Can be "json" or "yaml" (default "json")
      -h, --help                help for drpcli
      -P, --password string     password of the Digital Rebar Provision user (default "r0cketsk8ts")
      -r, --ref string          A reference object for update commands that can be a file name, yaml, or json blob
      -T, --token string        token of the Digital Rebar Provision access
      -t, --trace string        The log level API requests should be logged at on the server side
      -Z, --traceToken string   A token that individual traced requests should report in the server logs
      -U, --username string     Name of the Digital Rebar Provision user to talk to (default "rocketskates")

.. _rs_cli_filters:

Filters
-------

The CLI supports :ref:`rs_api_filters` on the command line by including `Field=Value` in the command list.  It is possible to have several filters applied in a single line as an AND operation.

Since simple Params are automatically mapped as fields, you can select objects by their .Properties.  If the object has .Params then the keys are also indexed as search keys.  For example `drpcli machines list rack/name=DC01AAA` will return all the machines with the Param rack/name set to DC01AAA.  You can also search into the .Meta field using `Meta.icon=lock`.

.. note:: Machines have a special shortcut `Name:[value]` that allows operators to select machines by name instead of UUID.  For example: `drpcli machines show Name:my.fqdn.com`.  While this uses filters in the background, it is not a general purpose filter.

SEE ALSO
--------

-  `drpcli autocomplete <drpcli_autocomplete.html>`__ - Generate CLI
   Command Bash AutoCompletion File (may require 'bash-completion' pkg
   be installed)
-  `drpcli bootenvs <drpcli_bootenvs.html>`__ - Access CLI commands
   relating to bootenvs
-  `drpcli certs <drpcli_certs.html>`__ - Access CLI commands relating
   to certs
-  `drpcli contents <drpcli_contents.html>`__ - Access CLI commands
   relating to content
-  `drpcli events <drpcli_events.html>`__ - DigitalRebar Provision Event
   Commands
-  `drpcli extended <drpcli_extended.html>`__ - Access CLI commands
   relating to extended
-  `drpcli files <drpcli_files.html>`__ - Access CLI commands relating
   to files
-  `drpcli gohai <drpcli_gohai.html>`__ - Get basic system information
   as a JSON blob
-  `drpcli info <drpcli_info.html>`__ - Access CLI commands relating to
   info
-  `drpcli interfaces <drpcli_interfaces.html>`__ - Access CLI commands
   relating to interfaces
-  `drpcli isos <drpcli_isos.html>`__ - Access CLI commands relating to
   isos
-  `drpcli jobs <drpcli_jobs.html>`__ - Access CLI commands relating to
   jobs
-  `drpcli leases <drpcli_leases.html>`__ - Access CLI commands relating
   to leases
-  `drpcli logs <drpcli_logs.html>`__ - Access commands relating to logs
-  `drpcli machines <drpcli_machines.html>`__ - Access CLI commands
   relating to machines
-  `drpcli params <drpcli_params.html>`__ - Access CLI commands relating
   to params
-  `drpcli plugin\_providers <drpcli_plugin_providers.html>`__ - Access
   CLI commands relating to plugin\_providers
-  `drpcli plugins <drpcli_plugins.html>`__ - Access CLI commands
   relating to plugins
-  `drpcli prefs <drpcli_prefs.html>`__ - List and set DigitalRebar
   Provision operational preferences
-  `drpcli profiles <drpcli_profiles.html>`__ - Access CLI commands
   relating to profiles
-  `drpcli reservations <drpcli_reservations.html>`__ - Access CLI
   commands relating to reservations
-  `drpcli roles <drpcli_roles.html>`__ - Access CLI commands relating
   to roles
-  `drpcli stages <drpcli_stages.html>`__ - Access CLI commands relating
   to stages
-  `drpcli subnets <drpcli_subnets.html>`__ - Access CLI commands
   relating to subnets
-  `drpcli system <drpcli_system.html>`__ - Access CLI commands relating
   to system
-  `drpcli tasks <drpcli_tasks.html>`__ - Access CLI commands relating
   to tasks
-  `drpcli templates <drpcli_templates.html>`__ - Access CLI commands
   relating to templates
-  `drpcli tenants <drpcli_tenants.html>`__ - Access CLI commands
   relating to tenants
-  `drpcli users <drpcli_users.html>`__ - Access CLI commands relating
   to users
-  `drpcli version <drpcli_version.html>`__ - Digital Rebar Provision
   CLI Command Version
-  `drpcli workflows <drpcli_workflows.html>`__ - Access CLI commands
   relating to workflows


.. _rs_bundlize_note:

Note on Contents Bundlize
-------------------------

The CLI offers the "bundlize" special command to extract data from an endpoint.  It can be used to back data or recover objects that were developed on a live endpoint and should now be moved into a bundle.

The special syntax for bundlize allows operators to name the objects that they want to extract using the following convention `[object type]:[name]` where the type is the plural object type (e.g. `machines`) and the name is the index of the object.

For example: `drpcli contents bundlize example.yaml workflows:discover` will create a file with a single workflow object.

