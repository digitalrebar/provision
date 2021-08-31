drpcli stages
-------------

Access CLI commands relating to stages

Synopsis
~~~~~~~~

Access CLI commands relating to stages

Options
~~~~~~~

::

     -h, --help   help for stages

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
-  `drpcli stages action <drpcli_stages_action.html>`__ - Display the
   action for this stage
-  `drpcli stages actions <drpcli_stages_actions.html>`__ - Display
   actions for this stage
-  `drpcli stages add <drpcli_stages_add.html>`__ - Add the stages param
   *key* to *blob*
-  `drpcli stages addprofile <drpcli_stages_addprofile.html>`__ - Add
   profile to the stage’s profile list
-  `drpcli stages addtask <drpcli_stages_addtask.html>`__ - Add task to
   the stage’s task list
-  `drpcli stages bootenv <drpcli_stages_bootenv.html>`__ - Set the
   stage’s bootenv
-  `drpcli stages count <drpcli_stages_count.html>`__ - Count all stages
-  `drpcli stages create <drpcli_stages_create.html>`__ - Create a new
   stage with the passed-in JSON or string key
-  `drpcli stages destroy <drpcli_stages_destroy.html>`__ - Destroy
   stage by id
-  `drpcli stages etag <drpcli_stages_etag.html>`__ - Get the etag for a
   stages by id
-  `drpcli stages exists <drpcli_stages_exists.html>`__ - See if a
   stages exists by id
-  `drpcli stages get <drpcli_stages_get.html>`__ - Get a parameter from
   the stage
-  `drpcli stages indexes <drpcli_stages_indexes.html>`__ - Get indexes
   for stages
-  `drpcli stages list <drpcli_stages_list.html>`__ - List all stages
-  `drpcli stages meta <drpcli_stages_meta.html>`__ - Gets metadata for
   the stage
-  `drpcli stages params <drpcli_stages_params.html>`__ - Gets/sets all
   parameters for the stage
-  `drpcli stages patch <drpcli_stages_patch.html>`__ - Patch stage by
   ID using the passed-in JSON Patch
-  `drpcli stages remove <drpcli_stages_remove.html>`__ - Remove the
   param *key* from stages
-  `drpcli stages removeprofile <drpcli_stages_removeprofile.html>`__ -
   Remove a profile from the stage’s profile list
-  `drpcli stages removetask <drpcli_stages_removetask.html>`__ - Remove
   a task from the stage’s list
-  `drpcli stages runaction <drpcli_stages_runaction.html>`__ - Run
   action on object from plugin
-  `drpcli stages set <drpcli_stages_set.html>`__ - Set the stages param
   *key* to *blob*
-  `drpcli stages show <drpcli_stages_show.html>`__ - Show a single
   stages by id
-  `drpcli stages update <drpcli_stages_update.html>`__ - Unsafely
   update stage by id with the passed-in JSON
-  `drpcli stages uploadiso <drpcli_stages_uploadiso.html>`__ - This
   will attempt to upload the ISO from the specified ISO URL.
-  `drpcli stages wait <drpcli_stages_wait.html>`__ - Wait for a stage’s
   field to become a value within a number of seconds
