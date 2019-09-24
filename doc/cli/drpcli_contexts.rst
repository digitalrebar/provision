drpcli contexts
---------------

Access CLI commands relating to contexts

Synopsis
~~~~~~~~

Access CLI commands relating to contexts

Options
~~~~~~~

::

     -h, --help   help for contexts

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

     -c, --catalog string          The catalog file to use to get product information (default "https://repo.rackn.io")
     -d, --debug                   Whether the CLI should run in debug mode
     -D, --download-proxy string   HTTP Proxy to use for downloading catalog and content
     -E, --endpoint string         The Digital Rebar Provision API endpoint to talk to (default "https://127.0.0.1:8092")
     -f, --force                   When needed, attempt to force the operation - used on some update/patch calls
     -F, --format string           The serialzation we expect for output.  Can be "json" or "yaml" (default "json")
     -x, --noToken                 Do not use token auth or token cache
     -P, --password string         password of the Digital Rebar Provision user (default "r0cketsk8ts")
     -r, --ref string              A reference object for update commands that can be a file name, yaml, or json blob
     -T, --token string            token of the Digital Rebar Provision access
     -t, --trace string            The log level API requests should be logged at on the server side
     -Z, --traceToken string       A token that individual traced requests should report in the server logs
     -U, --username string         Name of the Digital Rebar Provision user to talk to (default "rocketskates")

SEE ALSO
~~~~~~~~

-  `drpcli <drpcli.html>`__ - A CLI application for interacting with the
   DigitalRebar Provision API
-  `drpcli contexts create <drpcli_contexts_create.html>`__ - Create a
   new context with the passed-in JSON or string key
-  `drpcli contexts destroy <drpcli_contexts_destroy.html>`__ - Destroy
   context by id
-  `drpcli contexts exists <drpcli_contexts_exists.html>`__ - See if a
   contexts exists by id
-  `drpcli contexts indexes <drpcli_contexts_indexes.html>`__ - Get
   indexes for contexts
-  `drpcli contexts list <drpcli_contexts_list.html>`__ - List all
   contexts
-  `drpcli contexts meta <drpcli_contexts_meta.html>`__ - Gets metadata
   for the context
-  `drpcli contexts show <drpcli_contexts_show.html>`__ - Show a single
   contexts by id
-  `drpcli contexts update <drpcli_contexts_update.html>`__ - Unsafely
   update context by id with the passed-in JSON
-  `drpcli contexts wait <drpcli_contexts_wait.html>`__ - Wait for a
   context’s field to become a value within a number of seconds
