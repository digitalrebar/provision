drpcli jobs
-----------

Access CLI commands relating to jobs

Synopsis
~~~~~~~~

Access CLI commands relating to jobs

Options
~~~~~~~

::

     -h, --help   help for jobs

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

     -c, --catalog string          The catalog file to use to get product information (default "https://repo.rackn.io")
     -C, --colors string           The colors for JSON and Table/Text colorization.  8 values in the for 0=val,val;1=val,val2... (default "0=32;1=33;2=36;3=90;4=34,1;5=35;6=95;7=32;8=92")
     -d, --debug                   Whether the CLI should run in debug mode
     -D, --download-proxy string   HTTP Proxy to use for downloading catalog and content
     -E, --endpoint string         The Digital Rebar Provision API endpoint to talk to (default "https://127.0.0.1:8092")
     -f, --force                   When needed, attempt to force the operation - used on some update/patch calls
     -F, --format string           The serialization we expect for output.  Can be "json" or "yaml" or "text" or "table" (default "json")
     -N, --no-color                Whether the CLI should output colorized strings
     -H, --no-header               Should header be shown in "text" or "table" mode
     -x, --noToken                 Do not use token auth or token cache
     -P, --password string         password of the Digital Rebar Provision user (default "r0cketsk8ts")
     -J, --print-fields string     The fields of the object to display in "text" or "table" mode. Comma separated
     -r, --ref string              A reference object for update commands that can be a file name, yaml, or json blob
     -T, --token string            token of the Digital Rebar Provision access
     -t, --trace string            The log level API requests should be logged at on the server side
     -Z, --traceToken string       A token that individual traced requests should report in the server logs
     -j, --truncate-length int     Truncate columns at this length (default 40)
     -u, --url-proxy string        URL Proxy for passing actions through another DRP
     -U, --username string         Name of the Digital Rebar Provision user to talk to (default "rocketskates")

SEE ALSO
~~~~~~~~

-  `drpcli <drpcli.html>`__ - A CLI application for interacting with the
   DigitalRebar Provision API
-  `drpcli jobs actions <drpcli_jobs_actions.html>`__ - Get the actions
   for this job
-  `drpcli jobs count <drpcli_jobs_count.html>`__ - Count all jobs
-  `drpcli jobs create <drpcli_jobs_create.html>`__ - Create a new job
   with the passed-in JSON or string key
-  `drpcli jobs destroy <drpcli_jobs_destroy.html>`__ - Destroy job by
   id
-  `drpcli jobs etag <drpcli_jobs_etag.html>`__ - Get the etag for a
   jobs by id
-  `drpcli jobs exists <drpcli_jobs_exists.html>`__ - See if a jobs
   exists by id
-  `drpcli jobs indexes <drpcli_jobs_indexes.html>`__ - Get indexes for
   jobs
-  `drpcli jobs list <drpcli_jobs_list.html>`__ - List all jobs
-  `drpcli jobs log <drpcli_jobs_log.html>`__ - Gets the log or appends
   to the log if a second argument or stream is given
-  `drpcli jobs meta <drpcli_jobs_meta.html>`__ - Gets metadata for the
   job
-  `drpcli jobs plugin_action <drpcli_jobs_plugin_action.html>`__ -
   Display the action for this job
-  `drpcli jobs plugin_actions <drpcli_jobs_plugin_actions.html>`__ -
   Display actions for this job
-  `drpcli jobs purge <drpcli_jobs_purge.html>`__ - Purge jobs in excess
   of the job retention preferences
-  `drpcli jobs runplugin_action <drpcli_jobs_runplugin_action.html>`__
   - Run action on object from plugin
-  `drpcli jobs show <drpcli_jobs_show.html>`__ - Show a single jobs by
   id
-  `drpcli jobs update <drpcli_jobs_update.html>`__ - Unsafely update
   job by id with the passed-in JSON
-  `drpcli jobs wait <drpcli_jobs_wait.html>`__ - Wait for a job’s field
   to become a value within a number of seconds
