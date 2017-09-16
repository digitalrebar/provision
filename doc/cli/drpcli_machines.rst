drpcli machines
===============

Access CLI commands relating to machines

Synopsis
--------

Access CLI commands relating to machines

Options
-------

::

      -h, --help   help for machines

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
-  `drpcli machines action <drpcli_machines_action.html>`__ - Display
   the action for this machine
-  `drpcli machines actions <drpcli_machines_actions.html>`__ - Display
   actions for this machine
-  `drpcli machines addprofile <drpcli_machines_addprofile.html>`__ -
   Add profile to the machine's profile list
-  `drpcli machines bootenv <drpcli_machines_bootenv.html>`__ - Set the
   machine's bootenv
-  `drpcli machines create <drpcli_machines_create.html>`__ - Create a
   new machine with the passed-in JSON or string key
-  `drpcli machines destroy <drpcli_machines_destroy.html>`__ - Destroy
   machine by id
-  `drpcli machines exists <drpcli_machines_exists.html>`__ - See if a
   machine exists by id
-  `drpcli machines get <drpcli_machines_get.html>`__ - Get a parameter
   from the machine
-  `drpcli machines list <drpcli_machines_list.html>`__ - List all
   machines
-  `drpcli machines params <drpcli_machines_params.html>`__ - Gets/sets
   all parameters for the machine
-  `drpcli machines patch <drpcli_machines_patch.html>`__ - Patch
   machine with the passed-in JSON
-  `drpcli machines processjobs <drpcli_machines_processjobs.html>`__ -
   For the given machine, process pending jobs until done.
-  `drpcli machines
   removeprofile <drpcli_machines_removeprofile.html>`__ - Remove a
   profile from the machine's list
-  `drpcli machines runaction <drpcli_machines_runaction.html>`__ - Set
   preferences
-  `drpcli machines set <drpcli_machines_set.html>`__ - Set the
   machine's param to
-  `drpcli machines show <drpcli_machines_show.html>`__ - Show a single
   machine by id
-  `drpcli machines stage <drpcli_machines_stage.html>`__ - Set the
   machine's stage
-  `drpcli machines update <drpcli_machines_update.html>`__ - Unsafely
   update machine by id with the passed-in JSON
-  `drpcli machines wait <drpcli_machines_wait.html>`__ - Wait for a
   machine's field to become a value within a number of seconds
