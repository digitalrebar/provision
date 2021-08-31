drpcli machines
---------------

Access CLI commands relating to machines

Synopsis
~~~~~~~~

Access CLI commands relating to machines

Options
~~~~~~~

::

     -h, --help   help for machines

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

     -c, --catalog string          The catalog file to use to get product information (default "https://repo.rackn.io")
     -C, --colors string           The colors for JSON and Table/Text colorization.  8 values in the for 0=val,val;1=val,val2... (default "0=32;1=33;2=36;3=90;4=34,1;5=35;6=95;7=32;8=92")
     -d, --debug                   Whether the CLI should run in debug mode
     -D, --download-proxy string   HTTP Proxy to use for downloading catalog and content
     -E, --endpoint string         The Digital Rebar Provision API endpoint to talk to (default "https://127.0.0.1:8092")
     -X, --exit-early              Cause drpcli to exit if a command results in an object that has errors
     -f, --force                   When needed, attempt to force the operation - used on some update/patch calls
     -F, --format string           The serialization we expect for output.  Can be "json" or "yaml" or "text" or "table" (default "json")
     -N, --no-color                Whether the CLI should output colorized strings
     -H, --no-header               Should header be shown in "text" or "table" mode
     -x, --no-token                Do not use token auth or token cache
     -P, --password string         password of the Digital Rebar Provision user (default "r0cketsk8ts")
     -J, --print-fields string     The fields of the object to display in "text" or "table" mode. Comma separated
     -r, --ref string              A reference object for update commands that can be a file name, yaml, or json blob
     -T, --token string            token of the Digital Rebar Provision access
     -t, --trace string            The log level API requests should be logged at on the server side
     -Z, --trace-token string      A token that individual traced requests should report in the server logs
     -j, --truncate-length int     Truncate columns at this length (default 40)
     -u, --url-proxy string        URL Proxy for passing actions through another DRP
     -U, --username string         Name of the Digital Rebar Provision user to talk to (default "rocketskates")

SEE ALSO
~~~~~~~~

-  `drpcli <drpcli.html>`__ - A CLI application for interacting with the
   DigitalRebar Provision API
-  `drpcli machines action <drpcli_machines_action.html>`__ - Display
   the action for this machine
-  `drpcli machines actions <drpcli_machines_actions.html>`__ - Display
   actions for this machine
-  `drpcli machines add <drpcli_machines_add.html>`__ - Add the machines
   param *key* to *blob*
-  `drpcli machines addprofile <drpcli_machines_addprofile.html>`__ -
   Add profile to the machine’s profile list
-  `drpcli machines addtask <drpcli_machines_addtask.html>`__ - Add task
   to the machine’s task list
-  `drpcli machines bootenv <drpcli_machines_bootenv.html>`__ - Set the
   machine’s bootenv
-  `drpcli machines count <drpcli_machines_count.html>`__ - Count all
   machines
-  `drpcli machines create <drpcli_machines_create.html>`__ - Create a
   new machine with the passed-in JSON or string key
-  `drpcli machines currentlog <drpcli_machines_currentlog.html>`__ -
   Get the log for the most recent job run on the machine
-  `drpcli machines deletejobs <drpcli_machines_deletejobs.html>`__ -
   Delete all jobs associated with machine
-  `drpcli machines destroy <drpcli_machines_destroy.html>`__ - Destroy
   machine by id
-  `drpcli machines etag <drpcli_machines_etag.html>`__ - Get the etag
   for a machines by id
-  `drpcli machines exists <drpcli_machines_exists.html>`__ - See if a
   machines exists by id
-  `drpcli machines get <drpcli_machines_get.html>`__ - Get a parameter
   from the machine
-  `drpcli machines indexes <drpcli_machines_indexes.html>`__ - Get
   indexes for machines
-  `drpcli machines inserttask <drpcli_machines_inserttask.html>`__ -
   Insert a task at [offset] from machine’s running task
-  `drpcli machines inspect <drpcli_machines_inspect.html>`__ - Commands
   to inspect tasks and jobs on machines
-  `drpcli machines jobs <drpcli_machines_jobs.html>`__ - Access
   commands for manipulating the current job
-  `drpcli machines list <drpcli_machines_list.html>`__ - List all
   machines
-  `drpcli machines meta <drpcli_machines_meta.html>`__ - Gets metadata
   for the machine
-  `drpcli machines params <drpcli_machines_params.html>`__ - Gets/sets
   all parameters for the machine
-  `drpcli machines patch <drpcli_machines_patch.html>`__ - Patch
   machine by ID using the passed-in JSON Patch
-  `drpcli machines pause <drpcli_machines_pause.html>`__ - Mark the
   machine as NOT runnable
-  `drpcli machines processjobs <drpcli_machines_processjobs.html>`__ -
   For the given machine, process pending jobs until done.
-  `drpcli machines
   releaseToPool <drpcli_machines_releaseToPool.html>`__ - Release this
   machine back to the pool
-  `drpcli machines remove <drpcli_machines_remove.html>`__ - Remove the
   param *key* from machines
-  `drpcli machines
   removeprofile <drpcli_machines_removeprofile.html>`__ - Remove a
   profile from the machine’s profile list
-  `drpcli machines removetask <drpcli_machines_removetask.html>`__ -
   Remove a task from the machine’s list
-  `drpcli machines run <drpcli_machines_run.html>`__ - Mark the machine
   as runnable
-  `drpcli machines runaction <drpcli_machines_runaction.html>`__ - Run
   action on object from plugin
-  `drpcli machines set <drpcli_machines_set.html>`__ - Set the machines
   param *key* to *blob*
-  `drpcli machines show <drpcli_machines_show.html>`__ - Show a single
   machines by id
-  `drpcli machines stage <drpcli_machines_stage.html>`__ - Set the
   machine’s stage
-  `drpcli machines tasks <drpcli_machines_tasks.html>`__ - Access task
   manipulation for machines
-  `drpcli machines update <drpcli_machines_update.html>`__ - Unsafely
   update machine by id with the passed-in JSON
-  `drpcli machines uploadiso <drpcli_machines_uploadiso.html>`__ - This
   will attempt to upload the ISO from the specified ISO URL.
-  `drpcli machines wait <drpcli_machines_wait.html>`__ - Wait for a
   machine’s field to become a value within a number of seconds
-  `drpcli machines whoami <drpcli_machines_whoami.html>`__ - Figure out
   what machine UUID most closely matches the current system
-  `drpcli machines workflow <drpcli_machines_workflow.html>`__ - Set
   the machine’s workflow
