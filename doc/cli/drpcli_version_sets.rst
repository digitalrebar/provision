drpcli version_sets
-------------------

Access CLI commands relating to version_sets

Synopsis
~~~~~~~~

Access CLI commands relating to version_sets

Options
~~~~~~~

::

     -h, --help   help for version_sets

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
-  `drpcli version_sets action <drpcli_version_sets_action.html>`__ -
   Display the action for this version_set
-  `drpcli version_sets actions <drpcli_version_sets_actions.html>`__ -
   Display actions for this version_set
-  `drpcli version_sets create <drpcli_version_sets_create.html>`__ -
   Create a new version_set with the passed-in JSON or string key
-  `drpcli version_sets destroy <drpcli_version_sets_destroy.html>`__ -
   Destroy version_set by id
-  `drpcli version_sets exists <drpcli_version_sets_exists.html>`__ -
   See if a version_sets exists by id
-  `drpcli version_sets indexes <drpcli_version_sets_indexes.html>`__ -
   Get indexes for version_sets
-  `drpcli version_sets list <drpcli_version_sets_list.html>`__ - List
   all version_sets
-  `drpcli version_sets
   runaction <drpcli_version_sets_runaction.html>`__ - Run action on
   object from plugin
-  `drpcli version_sets show <drpcli_version_sets_show.html>`__ - Show a
   single version_sets by id
-  `drpcli version_sets update <drpcli_version_sets_update.html>`__ -
   Unsafely update version_set by id with the passed-in JSON
-  `drpcli version_sets wait <drpcli_version_sets_wait.html>`__ - Wait
   for a version_set’s field to become a value within a number of
   seconds
