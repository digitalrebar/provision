drpcli plugin\_providers meta
=============================

Gets metadata for the plugin\_provider

Synopsis
--------

Gets metadata for the plugin\_provider

::

    drpcli plugin_providers meta [id] [flags]

Options
-------

::

      -h, --help   help for meta

Options inherited from parent commands
--------------------------------------

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
--------

-  `drpcli plugin\_providers <drpcli_plugin_providers.html>`__ - Access
   CLI commands relating to plugin\_providers
-  `drpcli plugin\_providers meta
   add <drpcli_plugin_providers_meta_add.html>`__ - Atomically add
   [key]:[val] to the metadata on [plugin\_providers]:[id]
-  `drpcli plugin\_providers meta
   get <drpcli_plugin_providers_meta_get.html>`__ - Get a specific
   metadata item from plugin\_provider
-  `drpcli plugin\_providers meta
   remove <drpcli_plugin_providers_meta_remove.html>`__ - Remove the
   meta [key] from [plugin\_providers]:[id]
-  `drpcli plugin\_providers meta
   set <drpcli_plugin_providers_meta_set.html>`__ - Set metadata
   [key]:[val] on [plugin\_providers]:[id]
