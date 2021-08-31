drpcli plugins
--------------

Access CLI commands relating to plugins

Synopsis
~~~~~~~~

Access CLI commands relating to plugins

Options
~~~~~~~

::

     -h, --help   help for plugins

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

-  `drpcli <drpcli.html>`__ - A CLI application for interacting with the
   DigitalRebar Provision API
-  `drpcli plugins action <drpcli_plugins_action.html>`__ - Display the
   action for this plugin
-  `drpcli plugins actions <drpcli_plugins_actions.html>`__ - Display
   actions for this plugin
-  `drpcli plugins add <drpcli_plugins_add.html>`__ - Add the plugins
   param *key* to *blob*
-  `drpcli plugins count <drpcli_plugins_count.html>`__ - Count all
   plugins
-  `drpcli plugins create <drpcli_plugins_create.html>`__ - Create a new
   plugin with the passed-in JSON or string key
-  `drpcli plugins destroy <drpcli_plugins_destroy.html>`__ - Destroy
   plugin by id
-  `drpcli plugins etag <drpcli_plugins_etag.html>`__ - Get the etag for
   a plugins by id
-  `drpcli plugins exists <drpcli_plugins_exists.html>`__ - See if a
   plugins exists by id
-  `drpcli plugins get <drpcli_plugins_get.html>`__ - Get a parameter
   from the plugin
-  `drpcli plugins indexes <drpcli_plugins_indexes.html>`__ - Get
   indexes for plugins
-  `drpcli plugins list <drpcli_plugins_list.html>`__ - List all plugins
-  `drpcli plugins meta <drpcli_plugins_meta.html>`__ - Gets metadata
   for the plugin
-  `drpcli plugins params <drpcli_plugins_params.html>`__ - Gets/sets
   all parameters for the plugin
-  `drpcli plugins patch <drpcli_plugins_patch.html>`__ - Patch plugin
   by ID using the passed-in JSON Patch
-  `drpcli plugins remove <drpcli_plugins_remove.html>`__ - Remove the
   param *key* from plugins
-  `drpcli plugins runaction <drpcli_plugins_runaction.html>`__ - Run
   action on object from plugin
-  `drpcli plugins set <drpcli_plugins_set.html>`__ - Set the plugins
   param *key* to *blob*
-  `drpcli plugins show <drpcli_plugins_show.html>`__ - Show a single
   plugins by id
-  `drpcli plugins update <drpcli_plugins_update.html>`__ - Unsafely
   update plugin by id with the passed-in JSON
-  `drpcli plugins uploadiso <drpcli_plugins_uploadiso.html>`__ - This
   will attempt to upload the ISO from the specified ISO URL.
-  `drpcli plugins wait <drpcli_plugins_wait.html>`__ - Wait for a
   plugin’s field to become a value within a number of seconds
