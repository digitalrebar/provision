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

      -d, --debug             Whether the CLI should run in debug mode
      -E, --endpoint string   The Digital Rebar Provision API endpoint to talk to (default "https://127.0.0.1:8092")
      -f, --force             When needed, attempt to force the operation - used on some update/patch calls
      -F, --format string     The serialzation we expect for output.  Can be "json" or "yaml" (default "json")
      -P, --password string   password of the Digital Rebar Provision user (default "r0cketsk8ts")
      -T, --token string      token of the Digital Rebar Provision access
      -U, --username string   Name of the Digital Rebar Provision user to talk to (default "rocketskates")

SEE ALSO
--------

-  `drpcli <drpcli.html>`__ - A CLI application for interacting with the
   DigitalRebar Provision API
-  `drpcli tasks create <drpcli_tasks_create.html>`__ - Create a new
   task with the passed-in JSON or string key
-  `drpcli tasks destroy <drpcli_tasks_destroy.html>`__ - Destroy task
   by id
-  `drpcli tasks exists <drpcli_tasks_exists.html>`__ - See if a task
   exists by id
-  `drpcli tasks list <drpcli_tasks_list.html>`__ - List all tasks
-  `drpcli tasks patch <drpcli_tasks_patch.html>`__ - Patch task with
   the passed-in JSON
-  `drpcli tasks show <drpcli_tasks_show.html>`__ - Show a single task
   by id
-  `drpcli tasks update <drpcli_tasks_update.html>`__ - Unsafely update
   task by id with the passed-in JSON
