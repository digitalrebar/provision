drpcli profiles
===============

Access CLI commands relating to profiles

Synopsis
--------

Access CLI commands relating to profiles

Options
-------

::

      -h, --help   help for profiles

Options inherited from parent commands
--------------------------------------

::

      -c, --catalog string      The catalog file to use to get product information (default "https://repo.rackn.io")
      -d, --debug               Whether the CLI should run in debug mode
      -E, --endpoint string     The Digital Rebar Provision API endpoint to talk to (default "https://127.0.0.1:8092")
      -f, --force               When needed, attempt to force the operation - used on some update/patch calls
      -F, --format string       The serialzation we expect for output.  Can be "json" or "yaml" (default "json")
      -P, --password string     password of the Digital Rebar Provision user (default "r0cketsk8ts")
      -r, --ref string          A reference object for update commands that can be a file name, yaml, or json blob
      -T, --token string        token of the Digital Rebar Provision access
      -t, --trace string        The log level API requests should be logged at on the server side
      -Z, --traceToken string   A token that individual traced requests should report in the server logs
      -U, --username string     Name of the Digital Rebar Provision user to talk to (default "rocketskates")

SEE ALSO
--------

-  `drpcli <drpcli.html>`__ - A CLI application for interacting with the
   DigitalRebar Provision API
-  `drpcli profiles action <drpcli_profiles_action.html>`__ - Display
   the action for this profile
-  `drpcli profiles actions <drpcli_profiles_actions.html>`__ - Display
   actions for this profile
-  `drpcli profiles add <drpcli_profiles_add.html>`__ - Add the profiles
   param *key* to *blob*
-  `drpcli profiles create <drpcli_profiles_create.html>`__ - Create a
   new profile with the passed-in JSON or string key
-  `drpcli profiles destroy <drpcli_profiles_destroy.html>`__ - Destroy
   profile by id
-  `drpcli profiles exists <drpcli_profiles_exists.html>`__ - See if a
   profiles exists by id
-  `drpcli profiles get <drpcli_profiles_get.html>`__ - Get a parameter
   from the profile
-  `drpcli profiles indexes <drpcli_profiles_indexes.html>`__ - Get
   indexes for profiles
-  `drpcli profiles list <drpcli_profiles_list.html>`__ - List all
   profiles
-  `drpcli profiles meta <drpcli_profiles_meta.html>`__ - Gets metadata
   for the profile
-  `drpcli profiles params <drpcli_profiles_params.html>`__ - Gets/sets
   all parameters for the profile
-  `drpcli profiles remove <drpcli_profiles_remove.html>`__ - Remove the
   param *key* from profiles
-  `drpcli profiles runaction <drpcli_profiles_runaction.html>`__ - Run
   action on object from plugin
-  `drpcli profiles set <drpcli_profiles_set.html>`__ - Set the profiles
   param *key* to *blob*
-  `drpcli profiles show <drpcli_profiles_show.html>`__ - Show a single
   profiles by id
-  `drpcli profiles update <drpcli_profiles_update.html>`__ - Unsafely
   update profile by id with the passed-in JSON
-  `drpcli profiles wait <drpcli_profiles_wait.html>`__ - Wait for a
   profile's field to become a value within a number of seconds
