drpcli async_actions
--------------------

Access CLI commands relating to async_actions

Synopsis
~~~~~~~~

Access CLI commands relating to async_actions

Options
~~~~~~~

::

     -h, --help   help for async_actions

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
-  `drpcli async_actions action <drpcli_async_actions_action.html>`__ -
   Display the action for this async_action
-  `drpcli async_actions actions <drpcli_async_actions_actions.html>`__
   - Display actions for this async_action
-  `drpcli async_actions count <drpcli_async_actions_count.html>`__ -
   Count all async_actions
-  `drpcli async_actions create <drpcli_async_actions_create.html>`__ -
   Create a new async_action with the passed-in JSON or string key
-  `drpcli async_actions destroy <drpcli_async_actions_destroy.html>`__
   - Destroy async_action by id
-  `drpcli async_actions etag <drpcli_async_actions_etag.html>`__ - Get
   the etag for a async_actions by id
-  `drpcli async_actions exists <drpcli_async_actions_exists.html>`__ -
   See if a async_actions exists by id
-  `drpcli async_actions indexes <drpcli_async_actions_indexes.html>`__
   - Get indexes for async_actions
-  `drpcli async_actions list <drpcli_async_actions_list.html>`__ - List
   all async_actions
-  `drpcli async_actions meta <drpcli_async_actions_meta.html>`__ - Gets
   metadata for the async_action
-  `drpcli async_actions purge <drpcli_async_actions_purge.html>`__ -
   Purge action_actions in excess of the action_action retention
   preferences
-  `drpcli async_actions
   runaction <drpcli_async_actions_runaction.html>`__ - Run action on
   object from plugin
-  `drpcli async_actions show <drpcli_async_actions_show.html>`__ - Show
   a single async_actions by id
-  `drpcli async_actions update <drpcli_async_actions_update.html>`__ -
   Unsafely update async_action by id with the passed-in JSON
-  `drpcli async_actions wait <drpcli_async_actions_wait.html>`__ - Wait
   for a async_action’s field to become a value within a number of
   seconds
