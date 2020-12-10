drpcli catalog_items
--------------------

Access CLI commands relating to catalog_items

Synopsis
~~~~~~~~

Access CLI commands relating to catalog_items

Options
~~~~~~~

::

     -h, --help   help for catalog_items

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
-  `drpcli catalog_items count <drpcli_catalog_items_count.html>`__ -
   Count all catalog_items
-  `drpcli catalog_items create <drpcli_catalog_items_create.html>`__ -
   Create a new catalog_item with the passed-in JSON or string key
-  `drpcli catalog_items destroy <drpcli_catalog_items_destroy.html>`__
   - Destroy catalog_item by id
-  `drpcli catalog_items etag <drpcli_catalog_items_etag.html>`__ - Get
   the etag for a catalog_items by id
-  `drpcli catalog_items exists <drpcli_catalog_items_exists.html>`__ -
   See if a catalog_items exists by id
-  `drpcli catalog_items indexes <drpcli_catalog_items_indexes.html>`__
   - Get indexes for catalog_items
-  `drpcli catalog_items list <drpcli_catalog_items_list.html>`__ - List
   all catalog_items
-  `drpcli catalog_items show <drpcli_catalog_items_show.html>`__ - Show
   a single catalog_items by id
-  `drpcli catalog_items update <drpcli_catalog_items_update.html>`__ -
   Unsafely update catalog_item by id with the passed-in JSON
-  `drpcli catalog_items wait <drpcli_catalog_items_wait.html>`__ - Wait
   for a catalog_item’s field to become a value within a number of
   seconds
