drpcli users list
=================

List all users

Synopsis
--------

This will list all users by default. You can narrow down the items
returned using index filters. Use the "indexes" command to get the
indexes available for users.

To filter by indexes, you can use the following stanzas:

-   Eq This will return items Equal to according to
-   Ne This will return items Not Equal to according to
-   Lt This will return items Less Than according to
-   Lte This will return items Less Than Or Equal to according to
-   Gt This will return items Greater Than according to
-   Gte This will return items Greater Than Or Equal to according to
-   Between This will return items Greater Than Or Equal to and Less
   Than Or Equal to according to
-   Except This will return items Less Than or Greater Than according to

You can chain any number of filters together, and they will pipeline
into each other as appropriate. After the above filters have been
applied, you can further tweak how the results are returned using the
following meta-filters:

-  'reverse' to return items in reverse order
-  'limit' to only return the first items
-  'offset' to skip items
-  'sort' to sort items according to

::

    drpcli users list [filters...] [flags]

Options
-------

::

      -h, --help         help for list
          --limit int    Maximum number of items to return (default -1)
          --offset int   Number of items to skip before starting to return data (default -1)

Options inherited from parent commands
--------------------------------------

::

      -d, --debug             Whether the CLI should run in debug mode
      -E, --endpoint string   The Digital Rebar Provision API endpoint to talk to (default "https://127.0.0.1:8092")
      -f, --force             When needed, attempt to force the operation - used on some update/patch calls
      -F, --format string     The serialzation we expect for output.  Can be "json" or "yaml" (default "json")
      -P, --password string   password of the Digital Rebar Provision user (default "r0cketsk8ts")
      -r, --ref string        A reference object for update commands that can be a file name, yaml, or json blob
      -T, --token string      token of the Digital Rebar Provision access
      -U, --username string   Name of the Digital Rebar Provision user to talk to (default "rocketskates")

SEE ALSO
--------

-  `drpcli users <drpcli_users.html>`__ - Access CLI commands relating
   to users
