drpcli templates show
=====================

Show a single template by id

Synopsis
--------

This will show a templates.

It is possible to specify the id in the request by the using normal key or by
index.

Functional Indexs:

-  ID = string

When using the index name, use the following form:

-  Index:Value

Example:

-  e.g: ID:fred

::

    drpcli templates show [id]

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

-  `drpcli templates <drpcli_templates.html>`__ - Access CLI commands
   relating to templates
