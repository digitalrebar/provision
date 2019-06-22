drpcli templates
----------------

Access CLI commands relating to templates

Synopsis
~~~~~~~~

Access CLI commands relating to templates

Options
~~~~~~~

::

     -h, --help   help for templates

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
-  `drpcli templates action <drpcli_templates_action.html>`__ - Display
   the action for this template
-  `drpcli templates actions <drpcli_templates_actions.html>`__ -
   Display actions for this template
-  `drpcli templates create <drpcli_templates_create.html>`__ - Create a
   new template with the passed-in JSON or string key
-  `drpcli templates destroy <drpcli_templates_destroy.html>`__ -
   Destroy template by id
-  `drpcli templates exists <drpcli_templates_exists.html>`__ - See if a
   templates exists by id
-  `drpcli templates indexes <drpcli_templates_indexes.html>`__ - Get
   indexes for templates
-  `drpcli templates list <drpcli_templates_list.html>`__ - List all
   templates
-  `drpcli templates meta <drpcli_templates_meta.html>`__ - Gets
   metadata for the template
-  `drpcli templates runaction <drpcli_templates_runaction.html>`__ -
   Run action on object from plugin
-  `drpcli templates show <drpcli_templates_show.html>`__ - Show a
   single templates by id
-  `drpcli templates update <drpcli_templates_update.html>`__ - Unsafely
   update template by id with the passed-in JSON
-  `drpcli templates upload <drpcli_templates_upload.html>`__ - Upload
   the template file [file] as template [id]
-  `drpcli templates wait <drpcli_templates_wait.html>`__ - Wait for a
   template’s field to become a value within a number of seconds
