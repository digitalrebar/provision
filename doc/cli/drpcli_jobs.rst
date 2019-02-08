drpcli jobs
===========

Access CLI commands relating to jobs

Synopsis
--------

Access CLI commands relating to jobs

Options
-------

::

      -h, --help   help for jobs

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
-  `drpcli jobs actions <drpcli_jobs_actions.html>`__ - Get the actions
   for this job
-  `drpcli jobs create <drpcli_jobs_create.html>`__ - Create a new job
   with the passed-in JSON or string key
-  `drpcli jobs destroy <drpcli_jobs_destroy.html>`__ - Destroy job by
   id
-  `drpcli jobs exists <drpcli_jobs_exists.html>`__ - See if a jobs
   exists by id
-  `drpcli jobs indexes <drpcli_jobs_indexes.html>`__ - Get indexes for
   jobs
-  `drpcli jobs list <drpcli_jobs_list.html>`__ - List all jobs
-  `drpcli jobs log <drpcli_jobs_log.html>`__ - Gets the log or appends
   to the log if a second argument or stream is given
-  `drpcli jobs meta <drpcli_jobs_meta.html>`__ - Gets metadata for the
   job
-  `drpcli jobs plugin\_action <drpcli_jobs_plugin_action.html>`__ -
   Display the action for this job
-  `drpcli jobs plugin\_actions <drpcli_jobs_plugin_actions.html>`__ -
   Display actions for this job
-  `drpcli jobs runplugin\_action <drpcli_jobs_runplugin_action.html>`__
   - Run action on object from plugin
-  `drpcli jobs show <drpcli_jobs_show.html>`__ - Show a single jobs by
   id
-  `drpcli jobs update <drpcli_jobs_update.html>`__ - Unsafely update
   job by id with the passed-in JSON
-  `drpcli jobs wait <drpcli_jobs_wait.html>`__ - Wait for a job's field
   to become a value within a number of seconds
