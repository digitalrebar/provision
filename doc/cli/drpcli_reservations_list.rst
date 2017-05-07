drpcli reservations list
========================

List all reservations

Synopsis
--------

This will list all reservations by default.

You may specify:

-  Offset = integer, 0-based inclusive starting point in filter data.
-  Limit = integer, number of items to return

Functional Indexs:

-  Addr = IP Address
-  Token = string
-  Strategy = string
-  NextServer = IP Address

Functions:

-  Eq(value) = Return items that are equal to value
-  Lt(value) = Return items that are less than value
-  Lte(value) = Return items that less than or equal to value
-  Gt(value) = Return items that are greater than value
-  Gte(value) = Return items that greater than or equal to value
-  Between(lower,upper) = Return items that are inclusively between
   lower and upper
-  Except(lower,upper) = Return items that are not inclusively between
   lower and upper

Example:

-  NextServer=fred - returns items named fred
-  NextServer=Lt(fred) - returns items that alphabetically less than
   fred.
-  NextServer=Lt(fred)&Available=true - returns items with Name less
   than fred and Available is true

::

    drpcli reservations list [key=value] ...

Options
-------

::

          --limit int    Maximum number of items to return (default -1)
          --offset int   Number of items to skip before starting to return data (default -1)

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
