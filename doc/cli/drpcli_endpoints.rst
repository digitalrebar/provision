drpcli endpoints
----------------

Access CLI commands relating to endpoints

Synopsis
~~~~~~~~

Access CLI commands relating to endpoints

Options
~~~~~~~

::

     -h, --help   help for endpoints

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
-  `drpcli endpoints action <drpcli_endpoints_action.html>`__ - Display
   the action for this endpoint
-  `drpcli endpoints actions <drpcli_endpoints_actions.html>`__ -
   Display actions for this endpoint
-  `drpcli endpoints add <drpcli_endpoints_add.html>`__ - Add the
   endpoints param *key* to *blob*
-  `drpcli endpoints create <drpcli_endpoints_create.html>`__ - Create a
   new endpoint with the passed-in JSON or string key
-  `drpcli endpoints destroy <drpcli_endpoints_destroy.html>`__ -
   Destroy endpoint by id
-  `drpcli endpoints exists <drpcli_endpoints_exists.html>`__ - See if a
   endpoints exists by id
-  `drpcli endpoints get <drpcli_endpoints_get.html>`__ - Get a
   parameter from the endpoint
-  `drpcli endpoints indexes <drpcli_endpoints_indexes.html>`__ - Get
   indexes for endpoints
-  `drpcli endpoints list <drpcli_endpoints_list.html>`__ - List all
   endpoints
-  `drpcli endpoints meta <drpcli_endpoints_meta.html>`__ - Gets
   metadata for the endpoint
-  `drpcli endpoints params <drpcli_endpoints_params.html>`__ -
   Gets/sets all parameters for the endpoint
-  `drpcli endpoints remove <drpcli_endpoints_remove.html>`__ - Remove
   the param *key* from endpoints
-  `drpcli endpoints runaction <drpcli_endpoints_runaction.html>`__ -
   Run action on object from plugin
-  `drpcli endpoints set <drpcli_endpoints_set.html>`__ - Set the
   endpoints param *key* to *blob*
-  `drpcli endpoints show <drpcli_endpoints_show.html>`__ - Show a
   single endpoints by id
-  `drpcli endpoints update <drpcli_endpoints_update.html>`__ - Unsafely
   update endpoint by id with the passed-in JSON
-  `drpcli endpoints wait <drpcli_endpoints_wait.html>`__ - Wait for a
   endpoint’s field to become a value within a number of seconds
