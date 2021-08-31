drpcli extended
---------------

Access CLI commands relating to extended

Synopsis
~~~~~~~~

Access CLI commands relating to extended

Options
~~~~~~~

::

     -h, --help           help for extended
     -l, --ldata string   object type for extended data commands

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

     -c, --catalog string          The catalog file to use to get product information (default "https://repo.rackn.io")
     -C, --colors string           The colors for JSON and Table/Text colorization.  8 values in the for 0=val,val;1=val,val2... (default "0=32;1=33;2=36;3=90;4=34,1;5=35;6=95;7=32;8=92")
     -d, --debug                   Whether the CLI should run in debug mode
     -D, --download-proxy string   HTTP Proxy to use for downloading catalog and content
     -E, --endpoint string         The Digital Rebar Provision API endpoint to talk to (default "https://127.0.0.1:8092")
     -X, --exit-early              Cause drpcli to exit if a command results in an object that has errors
     -f, --force                   When needed, attempt to force the operation - used on some update/patch calls
     -F, --format string           The serialization we expect for output.  Can be "json" or "yaml" or "text" or "table" (default "json")
     -N, --no-color                Whether the CLI should output colorized strings
     -H, --no-header               Should header be shown in "text" or "table" mode
     -x, --no-token                Do not use token auth or token cache
     -P, --password string         password of the Digital Rebar Provision user (default "r0cketsk8ts")
     -J, --print-fields string     The fields of the object to display in "text" or "table" mode. Comma separated
     -r, --ref string              A reference object for update commands that can be a file name, yaml, or json blob
     -T, --token string            token of the Digital Rebar Provision access
     -t, --trace string            The log level API requests should be logged at on the server side
     -Z, --trace-token string      A token that individual traced requests should report in the server logs
     -j, --truncate-length int     Truncate columns at this length (default 40)
     -u, --url-proxy string        URL Proxy for passing actions through another DRP
     -U, --username string         Name of the Digital Rebar Provision user to talk to (default "rocketskates")

SEE ALSO
~~~~~~~~

-  `drpcli <drpcli.html>`__ - A CLI application for interacting with the
   DigitalRebar Provision API
-  `drpcli extended action <drpcli_extended_action.html>`__ - Display
   the action for this extended
-  `drpcli extended actions <drpcli_extended_actions.html>`__ - Display
   actions for this extended
-  `drpcli extended add <drpcli_extended_add.html>`__ - Add the param
   *key* to *blob*
-  `drpcli extended count <drpcli_extended_count.html>`__ - Count all
-  `drpcli extended create <drpcli_extended_create.html>`__ - Create a
   new extended with the passed-in JSON or string key
-  `drpcli extended destroy <drpcli_extended_destroy.html>`__ - Destroy
   extended by id
-  `drpcli extended etag <drpcli_extended_etag.html>`__ - Get the etag
   for a by id
-  `drpcli extended exists <drpcli_extended_exists.html>`__ - See if a
   exists by id
-  `drpcli extended get <drpcli_extended_get.html>`__ - Get a parameter
   from the extended
-  `drpcli extended indexes <drpcli_extended_indexes.html>`__ - Get
   indexes for
-  `drpcli extended list <drpcli_extended_list.html>`__ - List all
-  `drpcli extended meta <drpcli_extended_meta.html>`__ - Gets metadata
   for the extended
-  `drpcli extended params <drpcli_extended_params.html>`__ - Gets/sets
   all parameters for the extended
-  `drpcli extended patch <drpcli_extended_patch.html>`__ - Patch
   extended by ID using the passed-in JSON Patch
-  `drpcli extended remove <drpcli_extended_remove.html>`__ - Remove the
   param *key* from
-  `drpcli extended runaction <drpcli_extended_runaction.html>`__ - Run
   action on object from plugin
-  `drpcli extended set <drpcli_extended_set.html>`__ - Set the param
   *key* to *blob*
-  `drpcli extended show <drpcli_extended_show.html>`__ - Show a single
   by id
-  `drpcli extended update <drpcli_extended_update.html>`__ - Unsafely
   update extended by id with the passed-in JSON
-  `drpcli extended uploadiso <drpcli_extended_uploadiso.html>`__ - This
   will attempt to upload the ISO from the specified ISO URL.
-  `drpcli extended wait <drpcli_extended_wait.html>`__ - Wait for a
   extended’s field to become a value within a number of seconds
