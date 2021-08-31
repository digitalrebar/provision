drpcli pools manage release
---------------------------

release machines to pool

Synopsis
~~~~~~~~

Release frees machines in the selected pool. The machines must be
allocated and in InUse status.

::

   drpcli pools manage release [id ][filter options a=f(v) style] [flags]

Options
~~~~~~~

::

     -h, --help   help for release

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

         --add-parameters string      A JSON string of parameters to add to the machine
         --add-profiles string        Comma separated list of profiles to add to the machine
         --all-machines               Selects all available machines
     -c, --catalog string             The catalog file to use to get product information (default "https://repo.rackn.io")
     -C, --colors string              The colors for JSON and Table/Text colorization.  8 values in the for 0=val,val;1=val,val2... (default "0=32;1=33;2=36;3=90;4=34,1;5=35;6=95;7=32;8=92")
         --count int                  Count of machines to allocate
     -d, --debug                      Whether the CLI should run in debug mode
     -D, --download-proxy string      HTTP Proxy to use for downloading catalog and content
     -E, --endpoint string            The Digital Rebar Provision API endpoint to talk to (default "https://127.0.0.1:8092")
     -X, --exit-early                 Cause drpcli to exit if a command results in an object that has errors
     -f, --force                      When needed, attempt to force the operation - used on some update/patch calls
     -F, --format string              The serialization we expect for output.  Can be "json" or "yaml" or "text" or "table" (default "json")
         --machine-list string        Comma separated list of machines UUID or Field:Value
         --minimum int                Minimum number of machines to return - defaults to count
         --new-workflow string        A workflow to set on the machines
     -N, --no-color                   Whether the CLI should output colorized strings
     -H, --no-header                  Should header be shown in "text" or "table" mode
     -x, --no-token                   Do not use token auth or token cache
     -P, --password string            password of the Digital Rebar Provision user (default "r0cketsk8ts")
     -J, --print-fields string        The fields of the object to display in "text" or "table" mode. Comma separated
     -r, --ref string                 A reference object for update commands that can be a file name, yaml, or json blob
         --remove-parameters string   Comma separated list of parameters to remove from the machine
         --remove-profiles string     Comma separated list of profiles to remove from the machine
     -T, --token string               token of the Digital Rebar Provision access
     -t, --trace string               The log level API requests should be logged at on the server side
     -Z, --trace-token string         A token that individual traced requests should report in the server logs
     -j, --truncate-length int        Truncate columns at this length (default 40)
     -u, --url-proxy string           URL Proxy for passing actions through another DRP
     -U, --username string            Name of the Digital Rebar Provision user to talk to (default "rocketskates")
         --wait-timeout string        An amount of time to wait for completion in seconds or time string (e.g. 30m)

SEE ALSO
~~~~~~~~

-  `drpcli pools manage <drpcli_pools_manage.html>`__ - Manage machines
   in pools
