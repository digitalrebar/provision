drpcli pools manage
-------------------

Manage machines in pools

Synopsis
~~~~~~~~

Manage machines in pools

Options
~~~~~~~

::

         --add-parameters string      A JSON string of parameters to add to the machine
         --add-profiles string        Comma separated list of profiles to add to the machine
         --all-machines               Selects all available machines
         --count int                  Count of machines to allocate
     -h, --help                       help for manage
         --machine-list string        Comma separated list of machines UUID or Field:Value
         --minimum int                Minimum number of machines to return - defaults to count
         --new-workflow string        A workflow to set on the machines
         --remove-parameters string   Comma separated list of parameters to remove from the machine
         --remove-profiles string     Comma separated list of profiles to remove from the machine
         --wait-timeout string        An amount of time to wait for completion in seconds or time string (e.g. 30m)

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

-  `drpcli pools <drpcli_pools.html>`__ - Access CLI commands relating
   to pools
-  `drpcli pools manage add <drpcli_pools_manage_add.html>`__ - add
   machines to pool
-  `drpcli pools manage allocate <drpcli_pools_manage_allocate.html>`__
   - allocate machines to pool
-  `drpcli pools manage release <drpcli_pools_manage_release.html>`__ -
   release machines to pool
-  `drpcli pools manage remove <drpcli_pools_manage_remove.html>`__ -
   remove machines to pool
