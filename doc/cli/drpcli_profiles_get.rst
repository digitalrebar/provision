drpcli profiles get
===================

Get a parameter from the profile

Synopsis
--------

A helper function to return the value of the parameter on the profile

::

    drpcli profiles get [id] param [key]

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
