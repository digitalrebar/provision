drpcli contents
===============

Access CLI commands relating to contents

Synopsis
--------

Access CLI commands relating to contents

Options
-------

::

      -h, --help   help for contents

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
-  `drpcli contents bundle <drpcli_contents_bundle.html>`__ - Bundle a
   directory into a single file, specifed by [file]. [meta fields]
   allows for the specification of the meta data.
-  `drpcli contents create <drpcli_contents_create.html>`__ - Create a
   new content with the passed-in JSON or string key
-  `drpcli contents destroy <drpcli_contents_destroy.html>`__ - Destroy
   content by id
-  `drpcli contents exists <drpcli_contents_exists.html>`__ - See if a
   content exists by id
-  `drpcli contents list <drpcli_contents_list.html>`__ - List all
   contents
-  `drpcli contents show <drpcli_contents_show.html>`__ - Show a single
   content by id
-  `drpcli contents unbundle <drpcli_contents_unbundle.html>`__ -
   Unbundle a [file] into the local directory.
-  `drpcli contents update <drpcli_contents_update.html>`__ - Unsafely
   update content by id with the passed-in JSON
