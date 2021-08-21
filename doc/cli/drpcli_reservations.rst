drpcli reservations
-------------------

Access CLI commands relating to reservations

Synopsis
~~~~~~~~

Access CLI commands relating to reservations

Options
~~~~~~~

::

     -h, --help   help for reservations

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
-  `drpcli reservations action <drpcli_reservations_action.html>`__ -
   Display the action for this reservation
-  `drpcli reservations actions <drpcli_reservations_actions.html>`__ -
   Display actions for this reservation
-  `drpcli reservations count <drpcli_reservations_count.html>`__ -
   Count all reservations
-  `drpcli reservations create <drpcli_reservations_create.html>`__ -
   Create a new reservation with the passed-in JSON or string key
-  `drpcli reservations destroy <drpcli_reservations_destroy.html>`__ -
   Destroy reservation by id
-  `drpcli reservations etag <drpcli_reservations_etag.html>`__ - Get
   the etag for a reservations by id
-  `drpcli reservations exists <drpcli_reservations_exists.html>`__ -
   See if a reservations exists by id
-  `drpcli reservations indexes <drpcli_reservations_indexes.html>`__ -
   Get indexes for reservations
-  `drpcli reservations list <drpcli_reservations_list.html>`__ - List
   all reservations
-  `drpcli reservations meta <drpcli_reservations_meta.html>`__ - Gets
   metadata for the reservation
-  `drpcli reservations
   runaction <drpcli_reservations_runaction.html>`__ - Run action on
   object from plugin
-  `drpcli reservations show <drpcli_reservations_show.html>`__ - Show a
   single reservations by id
-  `drpcli reservations update <drpcli_reservations_update.html>`__ -
   Unsafely update reservation by id with the passed-in JSON
-  `drpcli reservations wait <drpcli_reservations_wait.html>`__ - Wait
   for a reservation’s field to become a value within a number of
   seconds
