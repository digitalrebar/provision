drpcli pools
------------

Access CLI commands relating to pools

Synopsis
~~~~~~~~

Access CLI commands relating to pools

Options
~~~~~~~

::

     -h, --help   help for pools

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

     -c, --catalog string          The catalog file to use to get product information (default "https://repo.rackn.io")
     -d, --debug                   Whether the CLI should run in debug mode
     -D, --download-proxy string   HTTP Proxy to use for downloading catalog and content
     -E, --endpoint string         The Digital Rebar Provision API endpoint to talk to (default "https://127.0.0.1:8092")
     -f, --force                   When needed, attempt to force the operation - used on some update/patch calls
     -F, --format string           The serialization we expect for output.  Can be "json" or "yaml" or "text" or "table" (default "json")
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
-  `drpcli pools action <drpcli_pools_action.html>`__ - Display the
   action for this pool
-  `drpcli pools actions <drpcli_pools_actions.html>`__ - Display
   actions for this pool
-  `drpcli pools active <drpcli_pools_active.html>`__ - List active
   pools
-  `drpcli pools count <drpcli_pools_count.html>`__ - Count all pools
-  `drpcli pools create <drpcli_pools_create.html>`__ - Create a new
   pool with the passed-in JSON or string key
-  `drpcli pools destroy <drpcli_pools_destroy.html>`__ - Destroy pool
   by id
-  `drpcli pools exists <drpcli_pools_exists.html>`__ - See if a pools
   exists by id
-  `drpcli pools indexes <drpcli_pools_indexes.html>`__ - Get indexes
   for pools
-  `drpcli pools list <drpcli_pools_list.html>`__ - List all pools
-  `drpcli pools manage <drpcli_pools_manage.html>`__ - Manage machines
   in pools
-  `drpcli pools runaction <drpcli_pools_runaction.html>`__ - Run action
   on object from plugin
-  `drpcli pools show <drpcli_pools_show.html>`__ - Show a single pools
   by id
-  `drpcli pools status <drpcli_pools_status.html>`__ - Get Pool status
-  `drpcli pools update <drpcli_pools_update.html>`__ - Unsafely update
   pool by id with the passed-in JSON
-  `drpcli pools wait <drpcli_pools_wait.html>`__ - Wait for a pool’s
   field to become a value within a number of seconds
