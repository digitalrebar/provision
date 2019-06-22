drpcli subnets
--------------

Access CLI commands relating to subnets

Synopsis
~~~~~~~~

Access CLI commands relating to subnets

Options
~~~~~~~

::

     -h, --help   help for subnets

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

-  `drpcli <drpcli.html>`__ - A CLI application for interacting with the
   DigitalRebar Provision API
-  `drpcli subnets action <drpcli_subnets_action.html>`__ - Display the
   action for this subnet
-  `drpcli subnets actions <drpcli_subnets_actions.html>`__ - Display
   actions for this subnet
-  `drpcli subnets create <drpcli_subnets_create.html>`__ - Create a new
   subnet with the passed-in JSON or string key
-  `drpcli subnets destroy <drpcli_subnets_destroy.html>`__ - Destroy
   subnet by id
-  `drpcli subnets exists <drpcli_subnets_exists.html>`__ - See if a
   subnets exists by id
-  `drpcli subnets get <drpcli_subnets_get.html>`__ - Get dhcpOption
   [number]
-  `drpcli subnets indexes <drpcli_subnets_indexes.html>`__ - Get
   indexes for subnets
-  `drpcli subnets leasetimes <drpcli_subnets_leasetimes.html>`__ - Set
   the leasetimes of a subnet
-  `drpcli subnets list <drpcli_subnets_list.html>`__ - List all subnets
-  `drpcli subnets meta <drpcli_subnets_meta.html>`__ - Gets metadata
   for the subnet
-  `drpcli subnets nextserver <drpcli_subnets_nextserver.html>`__ - Set
   next non-reserved IP
-  `drpcli subnets pickers <drpcli_subnets_pickers.html>`__ - assigns IP
   allocation methods to a subnet
-  `drpcli subnets range <drpcli_subnets_range.html>`__ - set the range
   of a subnet
-  `drpcli subnets runaction <drpcli_subnets_runaction.html>`__ - Run
   action on object from plugin
-  `drpcli subnets set <drpcli_subnets_set.html>`__ - Set the given
   subnet’s dhcpOption to a value
-  `drpcli subnets show <drpcli_subnets_show.html>`__ - Show a single
   subnets by id
-  `drpcli subnets update <drpcli_subnets_update.html>`__ - Unsafely
   update subnet by id with the passed-in JSON
-  `drpcli subnets wait <drpcli_subnets_wait.html>`__ - Wait for a
   subnet’s field to become a value within a number of seconds
