drpcli bootenvs show
====================

Show a single bootenv by id

Synopsis
--------

This will show a bootenvs.

You may specify the id in the request by the using normal key or by
index.

Functional Indexs:

-  Available = boolean
-  Name = string
-  OnlyUnknown = boolean

When using the index name, use the following form:

-  Index:Value

Example:

-  e.g: OnlyUnknown:fred

::

    drpcli bootenvs show [id]

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

-  `drpcli bootenvs <drpcli_bootenvs.html>`__ - Access CLI commands
   relating to bootenvs
