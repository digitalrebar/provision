drpcli roles
------------

Access CLI commands relating to roles

Synopsis
~~~~~~~~

Access CLI commands relating to roles

Options
~~~~~~~

::

     -h, --help   help for roles

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

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
~~~~~~~~

-  `drpcli <drpcli.html>`__ - A CLI application for interacting with the
   DigitalRebar Provision API
-  `drpcli roles create <drpcli_roles_create.html>`__ - Create a new
   role with the passed-in JSON or string key
-  `drpcli roles destroy <drpcli_roles_destroy.html>`__ - Destroy role
   by id
-  `drpcli roles exists <drpcli_roles_exists.html>`__ - See if a roles
   exists by id
-  `drpcli roles indexes <drpcli_roles_indexes.html>`__ - Get indexes
   for roles
-  `drpcli roles list <drpcli_roles_list.html>`__ - List all roles
-  `drpcli roles meta <drpcli_roles_meta.html>`__ - Gets metadata for
   the role
-  `drpcli roles show <drpcli_roles_show.html>`__ - Show a single roles
   by id
-  `drpcli roles update <drpcli_roles_update.html>`__ - Unsafely update
   role by id with the passed-in JSON
-  `drpcli roles wait <drpcli_roles_wait.html>`__ - Wait for a role’s
   field to become a value within a number of seconds
