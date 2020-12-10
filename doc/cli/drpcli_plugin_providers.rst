drpcli plugin_providers
-----------------------

Access CLI commands relating to plugin_providers

Synopsis
~~~~~~~~

Access CLI commands relating to plugin_providers

Options
~~~~~~~

::

     -h, --help   help for plugin_providers

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
-  `drpcli plugin_providers
   count <drpcli_plugin_providers_count.html>`__ - Count all
   plugin_providers
-  `drpcli plugin_providers
   destroy <drpcli_plugin_providers_destroy.html>`__ - Destroy
   plugin_provider by id
-  `drpcli plugin_providers
   download <drpcli_plugin_providers_download.html>`__ - Download the
   plugin_provider named [item] to [dest]
-  `drpcli plugin_providers etag <drpcli_plugin_providers_etag.html>`__
   - Get the etag for a plugin_providers by id
-  `drpcli plugin_providers
   exists <drpcli_plugin_providers_exists.html>`__ - See if a
   plugin_providers exists by id
-  `drpcli plugin_providers
   indexes <drpcli_plugin_providers_indexes.html>`__ - Get indexes for
   plugin_providers
-  `drpcli plugin_providers list <drpcli_plugin_providers_list.html>`__
   - List all plugin_providers
-  `drpcli plugin_providers meta <drpcli_plugin_providers_meta.html>`__
   - Gets metadata for the plugin_provider
-  `drpcli plugin_providers show <drpcli_plugin_providers_show.html>`__
   - Show a single plugin_providers by id
-  `drpcli plugin_providers
   upload <drpcli_plugin_providers_upload.html>`__ - Upload a program to
   act as a plugin_provider
-  `drpcli plugin_providers wait <drpcli_plugin_providers_wait.html>`__
   - Wait for a plugin_provider’s field to become a value within a
   number of seconds
