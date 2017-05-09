drpcli reservations exists
==========================

See if a reservation exists by id

Synopsis
--------

This will detect if a reservations exists.

You may specify the id in the request by the using normal key or by
index.

Functional Indexs:

-  Addr = IP Address
-  NextServer = IP Address
-  Strategy = string
-  Token = string

When using the index name, use the following form:

-  Index:Value

Example:

-  e.g: Token:fred

::

    drpcli reservations exists [id]

Options inherited from parent commands
--------------------------------------

::

      -d, --debug             Whether the CLI should run in debug mode
      -E, --endpoint string   The Rocket-Skates API endpoint to talk to (default "https://127.0.0.1:8092")
      -F, --format string     The serialzation we expect for output.  Can be "json" or "yaml" (default "json")
      -P, --password string   password of the Rocket-Skates user (default "r0cketsk8ts")
      -T, --token string      token of the Rocket-Skates access
      -U, --username string   Name of the Rocket-Skates user to talk to (default "rocketskates")

SEE ALSO
--------

-  `drpcli reservations <drpcli_reservations.html>`__ - Access CLI
   commands relating to reservations
