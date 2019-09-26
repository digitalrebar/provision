drpcli catalog item install
---------------------------

Installs [item] from the catalog on the current dr-provision endpoint

Synopsis
~~~~~~~~

Installs [item] from the catalog on the current dr-provision endpoint

::

   drpcli catalog item install [item] [flags]

Options
~~~~~~~

::

     -h, --help   help for install

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

         --arch string             Architecture of the item to work with when downloading a plugin provider (default "amd64")
     -c, --catalog string          The catalog file to use to get product information (default "https://repo.rackn.io")
     -d, --debug                   Whether the CLI should run in debug mode
     -D, --download-proxy string   HTTP Proxy to use for downloading catalog and content
     -E, --endpoint string         The Digital Rebar Provision API endpoint to talk to (default "https://127.0.0.1:8092")
     -f, --force                   When needed, attempt to force the operation - used on some update/patch calls
     -F, --format string           The serialzation we expect for output.  Can be "json" or "yaml" (default "json")
     -x, --noToken                 Do not use token auth or token cache
         --os string               OS of the item to work with when downloading a plugin provider (default "darwin")
     -P, --password string         password of the Digital Rebar Provision user (default "r0cketsk8ts")
     -r, --ref string              A reference object for update commands that can be a file name, yaml, or json blob
     -T, --token string            token of the Digital Rebar Provision access
     -t, --trace string            The log level API requests should be logged at on the server side
     -Z, --traceToken string       A token that individual traced requests should report in the server logs
     -U, --username string         Name of the Digital Rebar Provision user to talk to (default "rocketskates")
         --version string          Version of the item to work with (default "stable")

SEE ALSO
~~~~~~~~

-  `drpcli catalog item <drpcli_catalog_item.html>`__ - Commands to act
   on individual catalog items
