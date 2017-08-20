drpcli bootenvs
===============

Access CLI commands relating to bootenvs

Synopsis
--------

Access CLI commands relating to bootenvs

Options
-------

::

      -h, --help   help for bootenvs

Options inherited from parent commands
--------------------------------------

::

      -d, --debug             Whether the CLI should run in debug mode
      -E, --endpoint string   The Digital Rebar Provision API endpoint to talk to (default "https://127.0.0.1:8092")
      -f, --force             When needed, attempt to force the operation - used on some update/patch calls
      -F, --format string     The serialzation we expect for output.  Can be "json" or "yaml" (default "json")
      -P, --password string   password of the Digital Rebar Provision user (default "r0cketsk8ts")
      -T, --token string      token of the Digital Rebar Provision access
      -U, --username string   Name of the Digital Rebar Provision user to talk to (default "rocketskates")

SEE ALSO
--------

-  `drpcli <drpcli.html>`__ - A CLI application for interacting with the
   DigitalRebar Provision API
-  `drpcli bootenvs create <drpcli_bootenvs_create.html>`__ - Create a
   new bootenv with the passed-in JSON or string key
-  `drpcli bootenvs destroy <drpcli_bootenvs_destroy.html>`__ - Destroy
   bootenv by id
-  `drpcli bootenvs exists <drpcli_bootenvs_exists.html>`__ - See if a
   bootenv exists by id
-  `drpcli bootenvs install <drpcli_bootenvs_install.html>`__ - Install
   a bootenv along with everything it requires
-  `drpcli bootenvs list <drpcli_bootenvs_list.html>`__ - List all
   bootenvs
-  `drpcli bootenvs patch <drpcli_bootenvs_patch.html>`__ - Patch
   bootenv with the passed-in JSON
-  `drpcli bootenvs show <drpcli_bootenvs_show.html>`__ - Show a single
   bootenv by id
-  `drpcli bootenvs update <drpcli_bootenvs_update.html>`__ - Unsafely
   update bootenv by id with the passed-in JSON
-  `drpcli bootenvs uploadiso <drpcli_bootenvs_uploadiso.html>`__ - This
   will attempt to upload the ISO from the specified ISO URL.
