drpcli contents
===============

Access CLI commands relating to content

Synopsis
--------

Access CLI commands relating to content

Options
-------

::

      -h, --help   help for contents

Options inherited from parent commands
--------------------------------------

::

      -d, --debug               Whether the CLI should run in debug mode
      -E, --endpoint string     The Digital Rebar Provision API endpoint to talk to (default "https://127.0.0.1:8092")
      -f, --force               When needed, attempt to force the operation - used on some update/patch calls
      -F, --format string       The serialzation we expect for output.  Can be "json" or "yaml" (default "json")
      -P, --password string     password of the Digital Rebar Provision user (default "r0cketsk8ts")
      -r, --ref string          A reference object for update commands that can be a file name, yaml, or json blob
      -T, --token string        token of the Digital Rebar Provision access
      -t, --trace string        The log level API requests should be logged at on the server side
      -Z, --traceToken string   A token that individual traced requests should report in the server logs
      -U, --username string     Name of the Digital Rebar Provision user to talk to (default "rocketskates")

SEE ALSO
--------

-  `drpcli <drpcli.html>`__ - A CLI application for interacting with the
   DigitalRebar Provision API
-  `drpcli contents bundle <drpcli_contents_bundle.html>`__ - Bundle the
   current directory into [file]. [meta fields] allows for the
   specification of the meta data.
-  `drpcli contents bundlize <drpcli_contents_bundlize.html>`__ - Bundle
   the specified object into [file]. [meta fields] allows for the
   specification of the meta data. [objects] define which objects to
   record.
-  `drpcli contents convert <drpcli_contents_convert.html>`__ - Expand
   the content bundle [file] into DRP as read-write objects
-  `drpcli contents create <drpcli_contents_create.html>`__ - Add a new
   content layer to the system
-  `drpcli contents destroy <drpcli_contents_destroy.html>`__ - Remove
   the content layer [id] from the system.
-  `drpcli contents document <drpcli_contents_document.html>`__ - Expand
   the content bundle [file] into documentation
-  `drpcli contents exists <drpcli_contents_exists.html>`__ - See if
   content layer referenced by [id] exists
-  `drpcli contents list <drpcli_contents_list.html>`__ - List the
   installed content bundles
-  `drpcli contents show <drpcli_contents_show.html>`__ - Show a single
   content layer referenced by [id]
-  `drpcli contents unbundle <drpcli_contents_unbundle.html>`__ - Expand
   the content bundle [file] into the current directory
-  `drpcli contents update <drpcli_contents_update.html>`__ - Replace a
   content layer in the system.
-  `drpcli contents upload <drpcli_contents_upload.html>`__ - Upload a
   content layer into the system, replacing the earlier one if needed.
