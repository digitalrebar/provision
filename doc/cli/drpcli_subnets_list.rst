drpcli subnets list
-------------------

List all subnets

Synopsis
~~~~~~~~

This will list all subnets by default. You can narrow down the items
returned using index filters. Use the “indexes” command to get the
indexes available for subnets.

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
into each other as appropriate. After the above filters have been
applied, you can further tweak how the results are returned using the
following meta-filters:

-  ‘reverse’ to return items in reverse order
-  ‘limit’ *number* to only return the first *number* items
-  ‘offset’ *number* to skip *number* items
-  ‘sort’ *index* to sort items according to *index*

::

   drpcli subnets list [filters...] [flags]

Options
~~~~~~~

::

     -h, --help          help for list
         --limit int     Maximum number of items to return (default -1)
         --offset int    Number of items to skip before starting to return data (default -1)
         --slim string   Should elide certain fields.  Can be 'Params', 'Meta', or a comma-separated list of both.

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

     -c, --catalog string      The catalog file to use to get product information (default "https://repo.rackn.io")
     -d, --debug               Whether the CLI should run in debug mode
     -E, --endpoint string     The Digital Rebar Provision API endpoint to talk to (default "https://127.0.0.1:8092")
     -f, --force               When needed, attempt to force the operation - used on some update/patch calls
     -F, --format string       The serialzation we expect for output.  Can be "json" or "yaml" (default "json")
     -x, --noToken             Do not use token auth or token cache
     -P, --password string     password of the Digital Rebar Provision user (default "r0cketsk8ts")
     -r, --ref string          A reference object for update commands that can be a file name, yaml, or json blob
     -T, --token string        token of the Digital Rebar Provision access
     -t, --trace string        The log level API requests should be logged at on the server side
     -Z, --traceToken string   A token that individual traced requests should report in the server logs
     -U, --username string     Name of the Digital Rebar Provision user to talk to (default "rocketskates")

SEE ALSO
~~~~~~~~

-  `drpcli subnets <drpcli_subnets.html>`__ - Access CLI commands
   relating to subnets
