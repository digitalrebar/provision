drpcli machines create
======================

Create a new machine with the passed-in JSON or string key

Synopsis
--------

As a useful shortcut, you can pass '-' to indicate that the JSON should
be read from stdin.

In either case, for the Machine, BootEnv, User, and Profile objects, a
string may be provided to create a new empty object of that type. For
User, BootEnv, Machine, and Profile, it will be the object's name.

::

    drpcli machines create [json]

Options inherited from parent commands
--------------------------------------

::

      -d, --debug             Whether the CLI should run in debug mode
      -E, --endpoint string   The Digital Rebar Provision API endpoint to talk to (default "https://127.0.0.1:8092")
      -F, --format string     The serialzation we expect for output.  Can be "json" or "yaml" (default "json")
      -P, --password string   password of the Digital Rebar Provision user (default "r0cketsk8ts")
      -T, --token string      token of the Digital Rebar Provision access
      -U, --username string   Name of the Digital Rebar Provision user to talk to (default "rocketskates")

SEE ALSO
--------

-  `drpcli machines <drpcli_machines.html>`__ - Access CLI commands
   relating to machines
