drpcli profiles create
======================

Create a new profile with the passed-in JSON

Synopsis
--------

As a useful shortcut, you can pass '-' to indicate that the JSON should
be read from stdin

::

    drpcli profiles create [json]

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

-  `drpcli profiles <drpcli_profiles.html>`__ - Access CLI commands
   relating to profiles
