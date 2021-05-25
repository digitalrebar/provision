drpcli system ha
----------------

Access CLI commands to get the state of high availability

Synopsis
~~~~~~~~

Access CLI commands to get the state of high availability

Options
~~~~~~~

::

     -h, --help   help for ha

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

-  `drpcli system <drpcli_system.html>`__ - Access CLI commands relating
   to system
-  `drpcli system ha active <drpcli_system_ha_active.html>`__ - Get the
   machine ID of the current active node in the consensus system
-  `drpcli system ha dump <drpcli_system_ha_dump.html>`__ - Dump the
   detailed state of all members of the consensus system.
-  `drpcli system ha enroll <drpcli_system_ha_enroll.html>`__ - Have the
   endpoint at [endpointUrl] join the cluster.
-  `drpcli system ha
   failOverSafe <drpcli_system_ha_failOverSafe.html>`__ - Check to see
   if at least one non-observer passive node is caught up
-  `drpcli system ha id <drpcli_system_ha_id.html>`__ - Get the machine
   ID of this endpoint in the consensus system
-  `drpcli system ha
   introduction <drpcli_system_ha_introduction.html>`__ - Get an
   introduction from an existing cluster, save it in [file]
-  `drpcli system ha join <drpcli_system_ha_join.html>`__ - Join a
   cluster using the introduction saved in [file]
-  `drpcli system ha leader <drpcli_system_ha_leader.html>`__ - Get the
   machine ID of the leader in the consensus system
-  `drpcli system ha peers <drpcli_system_ha_peers.html>`__ - Get basic
   info on all members of the consensus system
-  `drpcli system ha remove <drpcli_system_ha_remove.html>`__ - Remove
   the node with provided Consensus Id from the cluster
-  `drpcli system ha state <drpcli_system_ha_state.html>`__ - Get the HA
   state of the system.
