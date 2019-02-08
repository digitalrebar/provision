drpcli tenants
==============

Access CLI commands relating to tenants

Synopsis
--------

Access CLI commands relating to tenants

Options
-------

::

      -h, --help   help for tenants

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
-  `drpcli tenants create <drpcli_tenants_create.html>`__ - Create a new
   tenant with the passed-in JSON or string key
-  `drpcli tenants destroy <drpcli_tenants_destroy.html>`__ - Destroy
   tenant by id
-  `drpcli tenants exists <drpcli_tenants_exists.html>`__ - See if a
   tenants exists by id
-  `drpcli tenants indexes <drpcli_tenants_indexes.html>`__ - Get
   indexes for tenants
-  `drpcli tenants list <drpcli_tenants_list.html>`__ - List all tenants
-  `drpcli tenants meta <drpcli_tenants_meta.html>`__ - Gets metadata
   for the tenant
-  `drpcli tenants show <drpcli_tenants_show.html>`__ - Show a single
   tenants by id
-  `drpcli tenants update <drpcli_tenants_update.html>`__ - Unsafely
   update tenant by id with the passed-in JSON
-  `drpcli tenants wait <drpcli_tenants_wait.html>`__ - Wait for a
   tenant's field to become a value within a number of seconds
