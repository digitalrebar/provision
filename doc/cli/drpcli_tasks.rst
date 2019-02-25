drpcli tasks
============

Access CLI commands relating to tasks

Synopsis
--------

Access CLI commands relating to tasks

Options
-------

::

      -h, --help   help for tasks

Options inherited from parent commands
--------------------------------------

::

      -c, --catalog string      The catalog file to use to get product information (default "https://repo.rackn.io")
      -d, --debug               Whether the CLI should run in debug mode
      -E, --endpoint string     The Digital Rebar Provision API endpoint to talk to (default "https://127.0.0.1:8092")
      -f, --force               When needed, attempt to force the operation - used on some update/patch calls
      -F, --format string       The serialzation we expect for output.  Can be "json" or "yaml" (default "json")
      -x, --noToken             Do not use token auth or token cache
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
-  `drpcli tasks action <drpcli_tasks_action.html>`__ - Display the
   action for this task
-  `drpcli tasks actions <drpcli_tasks_actions.html>`__ - Display
   actions for this task
-  `drpcli tasks create <drpcli_tasks_create.html>`__ - Create a new
   task with the passed-in JSON or string key
-  `drpcli tasks destroy <drpcli_tasks_destroy.html>`__ - Destroy task
   by id
-  `drpcli tasks exists <drpcli_tasks_exists.html>`__ - See if a tasks
   exists by id
-  `drpcli tasks indexes <drpcli_tasks_indexes.html>`__ - Get indexes
   for tasks
-  `drpcli tasks list <drpcli_tasks_list.html>`__ - List all tasks
-  `drpcli tasks meta <drpcli_tasks_meta.html>`__ - Gets metadata for
   the task
-  `drpcli tasks runaction <drpcli_tasks_runaction.html>`__ - Run action
   on object from plugin
-  `drpcli tasks show <drpcli_tasks_show.html>`__ - Show a single tasks
   by id
-  `drpcli tasks update <drpcli_tasks_update.html>`__ - Unsafely update
   task by id with the passed-in JSON
-  `drpcli tasks wait <drpcli_tasks_wait.html>`__ - Wait for a task's
   field to become a value within a number of seconds
