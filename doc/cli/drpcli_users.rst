drpcli users
------------

Access CLI commands relating to users

Synopsis
~~~~~~~~

Access CLI commands relating to users

Options
~~~~~~~

::

     -h, --help   help for users

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

     -c, --catalog string      The catalog file to use to get product information (default "https://repo.rackn.io")
     -d, --debug               Whether the CLI should run in debug mode
     -E, --endpoint string     The Digital Rebar Provision API endpoint to talk to (default "https://127.0.0.1:8092")
     -f, --force               When needed, attempt to force the operation - used on some update/patch calls
     -F, --format string       The serialzation we expect for output.  Can be "json" or "yaml" (default "json")
     -x, --noToken             Do not use token auth or token cache
     -P, --password string     password of the Digital Rebar Provision user (default "r0cketsk8ts")
     -r, --ref string          A reference object for update commands that can be a file name, yaml, or json blob
     -T, --token string        token of the Digital Rebar Provision access
     -t, --trace string        The log level API requests should be logged at on the server side
     -Z, --traceToken string   A token that individual traced requests should report in the server logs
     -U, --username string     Name of the Digital Rebar Provision user to talk to (default "rocketskates")

SEE ALSO
~~~~~~~~

-  `drpcli <drpcli.html>`__ - A CLI application for interacting with the
   DigitalRebar Provision API
-  `drpcli users action <drpcli_users_action.html>`__ - Display the
   action for this user
-  `drpcli users actions <drpcli_users_actions.html>`__ - Display
   actions for this user
-  `drpcli users create <drpcli_users_create.html>`__ - Create a new
   user with the passed-in JSON or string key
-  `drpcli users destroy <drpcli_users_destroy.html>`__ - Destroy user
   by id
-  `drpcli users exists <drpcli_users_exists.html>`__ - See if a users
   exists by id
-  `drpcli users indexes <drpcli_users_indexes.html>`__ - Get indexes
   for users
-  `drpcli users list <drpcli_users_list.html>`__ - List all users
-  `drpcli users meta <drpcli_users_meta.html>`__ - Gets metadata for
   the user
-  `drpcli users password <drpcli_users_password.html>`__ - Set the
   password for this id
-  `drpcli users passwordhash <drpcli_users_passwordhash.html>`__ - Get
   a password hash for a password
-  `drpcli users runaction <drpcli_users_runaction.html>`__ - Run action
   on object from plugin
-  `drpcli users show <drpcli_users_show.html>`__ - Show a single users
   by id
-  `drpcli users token <drpcli_users_token.html>`__ - Get a login token
   for this user with optional parameters
-  `drpcli users update <drpcli_users_update.html>`__ - Unsafely update
   user by id with the passed-in JSON
-  `drpcli users wait <drpcli_users_wait.html>`__ - Wait for a user’s
   field to become a value within a number of seconds
