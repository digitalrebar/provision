drpcli bootenvs
---------------

Access CLI commands relating to bootenvs

Synopsis
~~~~~~~~

Access CLI commands relating to bootenvs

Options
~~~~~~~

::

     -h, --help   help for bootenvs

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
-  `drpcli bootenvs action <drpcli_bootenvs_action.html>`__ - Display
   the action for this bootenv
-  `drpcli bootenvs actions <drpcli_bootenvs_actions.html>`__ - Display
   actions for this bootenv
-  `drpcli bootenvs create <drpcli_bootenvs_create.html>`__ - Create a
   new bootenv with the passed-in JSON or string key
-  `drpcli bootenvs destroy <drpcli_bootenvs_destroy.html>`__ - Destroy
   bootenv by id
-  `drpcli bootenvs exists <drpcli_bootenvs_exists.html>`__ - See if a
   bootenvs exists by id
-  `drpcli bootenvs fromAppleNBI <drpcli_bootenvs_fromAppleNBI.html>`__
   - This will attempt to translate an Apple .nbi directory into a
   bootenv and an archive.
-  `drpcli bootenvs indexes <drpcli_bootenvs_indexes.html>`__ - Get
   indexes for bootenvs
-  `drpcli bootenvs install <drpcli_bootenvs_install.html>`__ - Install
   a bootenv along with everything it requires
-  `drpcli bootenvs list <drpcli_bootenvs_list.html>`__ - List all
   bootenvs
-  `drpcli bootenvs meta <drpcli_bootenvs_meta.html>`__ - Gets metadata
   for the bootenv
-  `drpcli bootenvs runaction <drpcli_bootenvs_runaction.html>`__ - Run
   action on object from plugin
-  `drpcli bootenvs show <drpcli_bootenvs_show.html>`__ - Show a single
   bootenvs by id
-  `drpcli bootenvs update <drpcli_bootenvs_update.html>`__ - Unsafely
   update bootenv by id with the passed-in JSON
-  `drpcli bootenvs uploadiso <drpcli_bootenvs_uploadiso.html>`__ - This
   will attempt to upload the ISO from the specified ISO URL.
-  `drpcli bootenvs wait <drpcli_bootenvs_wait.html>`__ - Wait for a
   bootenv’s field to become a value within a number of seconds
