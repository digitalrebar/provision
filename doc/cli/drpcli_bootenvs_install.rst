drpcli bootenvs install
=======================

Install a bootenv along with everything it requires

Synopsis
--------

bootenvs install assumes a directory with two subdirectories: bootenvs/
templates/

bootenvs must contain [bootenvFile] templates must contain any templates
that the requested bootenv refers to.

bootenvs install will try to upload any required ISOs if they are not
already present in DigitalRebar Provision. If [isoPath] is specified, it
will use that directory to to check and download ISOs into, otherwise it
will use isos/ If the ISO is not present, we will try to download it if
the bootenv specifies a location to download the ISO from. If we cannot
find an ISO to upload, then the bootenv will still be uploaded, but it
will not be available until the ISO is uploaded using isos upload.git

::

    drpcli bootenvs install [bootenvFile] [isoPath] [flags]

Options
-------

::

      -h, --help            help for install
          --skip-download   Whether to try to download ISOs from their upstream

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

-  `drpcli bootenvs <drpcli_bootenvs.html>`__ - Access CLI commands
   relating to bootenvs
