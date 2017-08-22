drpcli reservations update
==========================

Unsafely update reservation by id with the passed-in JSON

Synopsis
--------

As a useful shortcut, pass '-' to indicate that the JSON should
be read from stdin

::

    drpcli reservations update [id] [json] [flags]

Options
-------

::

      -h, --help   help for update

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

-  `drpcli reservations <drpcli_reservations.html>`__ - Access CLI
   commands relating to reservations
