drpcli files
============

Commands to manage files on the provisioner

Synopsis
--------

Commands to manage files on the provisioner

Options
-------

::

      -h, --help   help for files

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
-  `drpcli files destroy <drpcli_files_destroy.html>`__ - Destroy file
   by id
-  `drpcli files exists <drpcli_files_exists.html>`__ - See if a file
   exists by id
-  `drpcli files list <drpcli_files_list.html>`__ - List all files
-  `drpcli files show <drpcli_files_show.html>`__ - Show a single file
   by id
-  `drpcli files upload <drpcli_files_upload.html>`__ - Upload a local
   file to Digital Rebar Provision
