drpcli stages
=============

Access CLI commands relating to stages

Synopsis
--------

Access CLI commands relating to stages

Options
-------

::

      -h, --help   help for stages

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
-  `drpcli stages action <drpcli_stages_action.html>`__ - Display the
   action for this stage
-  `drpcli stages actions <drpcli_stages_actions.html>`__ - Display
   actions for this stage
-  `drpcli stages addprofile <drpcli_stages_addprofile.html>`__ - Add
   profile to the machine's profile list
-  `drpcli stages addtask <drpcli_stages_addtask.html>`__ - Add task to
   the stage's task list
-  `drpcli stages bootenv <drpcli_stages_bootenv.html>`__ - Set the
   stage's bootenv
-  `drpcli stages create <drpcli_stages_create.html>`__ - Create a new
   stage with the passed-in JSON or string key
-  `drpcli stages destroy <drpcli_stages_destroy.html>`__ - Destroy
   stage by id
-  `drpcli stages exists <drpcli_stages_exists.html>`__ - See if a
   stages exists by id
-  `drpcli stages indexes <drpcli_stages_indexes.html>`__ - Get indexes
   for stages
-  `drpcli stages list <drpcli_stages_list.html>`__ - List all stages
-  `drpcli stages meta <drpcli_stages_meta.html>`__ - Gets metadata for
   the stage
-  `drpcli stages removeprofile <drpcli_stages_removeprofile.html>`__ -
   Remove a profile from the machine's list
-  `drpcli stages removetask <drpcli_stages_removetask.html>`__ - Remove
   a task from the stage's list
-  `drpcli stages runaction <drpcli_stages_runaction.html>`__ - Run
   action on object from plugin
-  `drpcli stages show <drpcli_stages_show.html>`__ - Show a single
   stages by id
-  `drpcli stages update <drpcli_stages_update.html>`__ - Unsafely
   update stage by id with the passed-in JSON
-  `drpcli stages wait <drpcli_stages_wait.html>`__ - Wait for a stage's
   field to become a value within a number of seconds
