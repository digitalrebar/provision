drpcli params
=============

Access CLI commands relating to params

Synopsis
--------

Access CLI commands relating to params

Options
-------

::

      -h, --help   help for params

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
-  `drpcli params create <drpcli_params_create.html>`__ - Create a new
   param with the passed-in JSON or string key
-  `drpcli params destroy <drpcli_params_destroy.html>`__ - Destroy
   param by id
-  `drpcli params exists <drpcli_params_exists.html>`__ - See if a
   params exists by id
-  `drpcli params indexes <drpcli_params_indexes.html>`__ - Get indexes
   for params
-  `drpcli params list <drpcli_params_list.html>`__ - List all params
-  `drpcli params meta <drpcli_params_meta.html>`__ - Gets metadata for
   the param
-  `drpcli params show <drpcli_params_show.html>`__ - Show a single
   params by id
-  `drpcli params update <drpcli_params_update.html>`__ - Unsafely
   update param by id with the passed-in JSON
-  `drpcli params wait <drpcli_params_wait.html>`__ - Wait for a param's
   field to become a value within a number of seconds
