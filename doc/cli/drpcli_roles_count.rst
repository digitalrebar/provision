drpcli roles count
------------------

Count all roles

Synopsis
~~~~~~~~

This will count all roles by default. You can narrow down the count
returned using index filters. Use the “indexes” command to get the
indexes available for roles.

To filter by indexes, you can use the following stanzas:

-  *index* Eq *value* This will return items Equal to *value* according
   to *index*
-  *index* Ne *value* This will return items Not Equal to *value*
   according to *index*
-  *index* Lt *value* This will return items Less Than *value* according
   to *index*
-  *index* Lte *value* This will return items Less Than Or Equal to
   *value* according to *index*
-  *index* Gt *value* This will return items Greater Than *value*
   according to *index*
-  *index* Gte *value* This will return items Greater Than Or Equal to
   *value* according to *index*
-  *index* Re *re2 compatible regular expression* This will return items
   in *index* that match the passed-in regular expression We use the
   regular expression syntax described at
   https://github.com/google/re2/wiki/Syntax
-  *index* Between *lower* *upper* This will return items Greater Than
   Or Equal to *lower* and Less Than Or Equal to *upper* according to
   *index*
-  *index* Except *lower* *upper* This will return items Less Than
   *lower* or Greater Than *upper* according to *index*
-  *index* In *comma,separated,list,of,values* This will return any
   items In the set passed for the comma-separated list of values.
-  *index* Nin *comma,separated,list,of,values* This will return any
   items Not In the set passed for the comma-separated list of values.

You can chain any number of filters together, and they will pipeline
into each other as appropriate.

::

   drpcli roles count [filters...] [flags]

Options
~~~~~~~

::

     -h, --help   help for count

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

-  `drpcli roles <drpcli_roles.html>`__ - Access CLI commands relating
   to roles
