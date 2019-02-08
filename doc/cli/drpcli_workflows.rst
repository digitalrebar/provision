drpcli workflows
================

Access CLI commands relating to workflows

Synopsis
--------

Access CLI commands relating to workflows

Options
-------

::

      -h, --help   help for workflows

Options inherited from parent commands
--------------------------------------

::

      -c, --catalog string      The catalog file to use to get product information (default "https://repo.rackn.io")
      -d, --debug               Whether the CLI should run in debug mode
      -E, --endpoint string     The Digital Rebar Provision API endpoint to talk to (default "https://127.0.0.1:8092")
      -f, --force               When needed, attempt to force the operation - used on some update/patch calls
      -F, --format string       The serialzation we expect for output.  Can be "json" or "yaml" (default "json")
      -P, --password string     password of the Digital Rebar Provision user (default "r0cketsk8ts")
      -r, --ref string          A reference object for update commands that can be a file name, yaml, or json blob
      -T, --token string        token of the Digital Rebar Provision access
      -t, --trace string        The log level API requests should be logged at on the server side
      -Z, --traceToken string   A token that individual traced requests should report in the server logs
      -U, --username string     Name of the Digital Rebar Provision user to talk to (default "rocketskates")

SEE ALSO
--------

-  `drpcli <drpcli.html>`__ - A CLI application for interacting with the
   DigitalRebar Provision API
-  `drpcli workflows action <drpcli_workflows_action.html>`__ - Display
   the action for this workflow
-  `drpcli workflows actions <drpcli_workflows_actions.html>`__ -
   Display actions for this workflow
-  `drpcli workflows create <drpcli_workflows_create.html>`__ - Create a
   new workflow with the passed-in JSON or string key
-  `drpcli workflows destroy <drpcli_workflows_destroy.html>`__ -
   Destroy workflow by id
-  `drpcli workflows exists <drpcli_workflows_exists.html>`__ - See if a
   workflows exists by id
-  `drpcli workflows indexes <drpcli_workflows_indexes.html>`__ - Get
   indexes for workflows
-  `drpcli workflows list <drpcli_workflows_list.html>`__ - List all
   workflows
-  `drpcli workflows meta <drpcli_workflows_meta.html>`__ - Gets
   metadata for the workflow
-  `drpcli workflows runaction <drpcli_workflows_runaction.html>`__ -
   Run action on object from plugin
-  `drpcli workflows show <drpcli_workflows_show.html>`__ - Show a
   single workflows by id
-  `drpcli workflows update <drpcli_workflows_update.html>`__ - Unsafely
   update workflow by id with the passed-in JSON
-  `drpcli workflows wait <drpcli_workflows_wait.html>`__ - Wait for a
   workflow's field to become a value within a number of seconds
