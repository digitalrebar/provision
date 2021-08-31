.. Copyright (c) 2021 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; CLI

.. _rs_drpcli:

Digital Rebar CLI (drpcli)
--------------------------

A CLI application for interacting with the DigitalRebar Provision API

Synopsis
~~~~~~~~

drpcli is a general-purpose command for interacting with a dr-provision
endpoint. It has several subcommands which have their own help.

It also has several environment variables that control aspects of its
operation:

-  RS_OBJECT_ERRORS_ARE_FATAL: Have drpcli exit with a non-zero exit
   status if a returned object has an Errors field that is not empty.
   Normally it will only exit with a non-zero exit status when the API
   returns with an error or fatal status code.

-  RS_ENDPOINTS: A space-seperated list of URLS that drpcli should try
   to communicate with. The first one that authenticates will be used.

-  RS_ENDPOINT: The URL that drpcli should try to communicate. Ignored
   if RS_ENDPOINTS exists in the environment. Default to
   https://127.0.0.1:8092

-  RS_URL_PROXY: The HTTP proxy drpcli should use when communicating
   with the dr-provision endpoint. It functions like the standard
   http_proxy environment variable.

-  RS_TOKEN: The token to use for authentication with the dr-provision
   endpoint. Overrides RS_KEY.

-  RS_CATALOG: The URL to use to fetch the artifact catalog. All
   commands in the ‘drpcli catalog’ group of commands use this. Defaults
   to https://repo.rackn.io

-  RS_FORMAT: The output format drpcli will use. Defaults to json

-  RS_PRINT_FIELDS: The fields of an object to display in text or table
   format. Defaults to all of them.

-  RS_DOWNLOAD_PROXY: The http proxy to use when downloading bootenv ISO
   files.

-  RS_NO_HEADER: Controls whether to print column headers in text or
   table output mode.

-  RS_NO_COLOR: Controls whether output to a terminal should be
   stripped.

-  RS_COLORS: Controls the 8 ANSI colors that should be used in
   colorized output.

-  RS_TRUNCATE_LENGTH: The max length of an individual column in text or
   table mode.

-  RS_KEY: The default username:password to use when missing a token.

Options
~~~~~~~

::

     -c, --catalog string          The catalog file to use to get product information (default "https://repo.rackn.io")
     -C, --colors string           The colors for JSON and Table/Text colorization.  8 values in the for 0=val,val;1=val,val2... (default "0=32;1=33;2=36;3=90;4=34,1;5=35;6=95;7=32;8=92")
     -d, --debug                   Whether the CLI should run in debug mode
     -D, --download-proxy string   HTTP Proxy to use for downloading catalog and content
     -E, --endpoint string         The Digital Rebar Provision API endpoint to talk to (default "https://127.0.0.1:8092")
     -X, --exit-early              Cause drpcli to exit if a command results in an object that has errors
     -f, --force                   When needed, attempt to force the operation - used on some update/patch calls
     -F, --format string           The serialization we expect for output.  Can be "json" or "yaml" or "text" or "table" (default "json")
     -h, --help                    help for drpcli
     -N, --no-color                Whether the CLI should output colorized strings
     -H, --no-header               Should header be shown in "text" or "table" mode
     -x, --no-token                Do not use token auth or token cache
     -P, --password string         password of the Digital Rebar Provision user (default "r0cketsk8ts")
     -J, --print-fields string     The fields of the object to display in "text" or "table" mode. Comma separated
     -r, --ref string              A reference object for update commands that can be a file name, yaml, or json blob
     -T, --token string            token of the Digital Rebar Provision access
     -t, --trace string            The log level API requests should be logged at on the server side
     -Z, --trace-token string      A token that individual traced requests should report in the server logs
     -j, --truncate-length int     Truncate columns at this length (default 40)
     -u, --url-proxy string        URL Proxy for passing actions through another DRP
     -U, --username string         Name of the Digital Rebar Provision user to talk to (default "rocketskates")

SEE ALSO
~~~~~~~~

-  `drpcli agent <drpcli_agent.html>`__ - Manage drpcli running as an
   agent
-  `drpcli autocomplete <drpcli_autocomplete.html>`__ - Generate CLI
   Command Bash AutoCompletion File (may require ‘bash-completion’ pkg
   be installed)
-  `drpcli bootenvs <drpcli_bootenvs.html>`__ - Access CLI commands
   relating to bootenvs
-  `drpcli catalog <drpcli_catalog.html>`__ - Access commands related to
   catalog manipulation
-  `drpcli catalog_items <drpcli_catalog_items.html>`__ - Access CLI
   commands relating to catalog_items
-  `drpcli certs <drpcli_certs.html>`__ - Access CLI commands relating
   to certs
-  `drpcli contents <drpcli_contents.html>`__ - Access CLI commands
   relating to content
-  `drpcli contexts <drpcli_contexts.html>`__ - Access CLI commands
   relating to contexts
-  `drpcli debug <drpcli_debug.html>`__ - Gather [type] of debug
   information and save it to [target]
-  `drpcli endpoints <drpcli_endpoints.html>`__ - Access CLI commands
   relating to endpoints
-  `drpcli events <drpcli_events.html>`__ - DigitalRebar Provision Event
   Commands
-  `drpcli extended <drpcli_extended.html>`__ - Access CLI commands
   relating to extended
-  `drpcli files <drpcli_files.html>`__ - Access CLI commands relating
   to files
-  `drpcli fingerprint <drpcli_fingerprint.html>`__ - Get the machine
   fingerprint used to determine what machine we are running on
-  `drpcli gohai <drpcli_gohai.html>`__ - Get basic system information
   as a JSON blob
-  `drpcli info <drpcli_info.html>`__ - Access CLI commands relating to
   info
-  `drpcli interfaces <drpcli_interfaces.html>`__ - Access CLI commands
   relating to interfaces
-  `drpcli isos <drpcli_isos.html>`__ - Access CLI commands relating to
   isos
-  `drpcli jobs <drpcli_jobs.html>`__ - Access CLI commands relating to
   jobs
-  `drpcli leases <drpcli_leases.html>`__ - Access CLI commands relating
   to leases
-  `drpcli logs <drpcli_logs.html>`__ - Access commands relating to logs
-  `drpcli machines <drpcli_machines.html>`__ - Access CLI commands
   relating to machines
-  `drpcli objects <drpcli_objects.html>`__ - Access CLI commands
   relating to objects
-  `drpcli params <drpcli_params.html>`__ - Access CLI commands relating
   to params
-  `drpcli plugin_providers <drpcli_plugin_providers.html>`__ - Access
   CLI commands relating to plugin_providers
-  `drpcli plugins <drpcli_plugins.html>`__ - Access CLI commands
   relating to plugins
-  `drpcli pools <drpcli_pools.html>`__ - Access CLI commands relating
   to pools
-  `drpcli prefs <drpcli_prefs.html>`__ - List and set DigitalRebar
   Provision operational preferences
-  `drpcli profiles <drpcli_profiles.html>`__ - Access CLI commands
   relating to profiles
-  `drpcli proxy <drpcli_proxy.html>`__ - Run a local UNIX socket proxy
   for further drpcli commands. Requires RS_LOCAL_PROXY to be set in the
   env.
-  `drpcli reservations <drpcli_reservations.html>`__ - Access CLI
   commands relating to reservations
-  `drpcli roles <drpcli_roles.html>`__ - Access CLI commands relating
   to roles
-  `drpcli stages <drpcli_stages.html>`__ - Access CLI commands relating
   to stages
-  `drpcli subnets <drpcli_subnets.html>`__ - Access CLI commands
   relating to subnets
-  `drpcli support <drpcli_support.html>`__ - Access commands related to
   RackN Tech Support
-  `drpcli system <drpcli_system.html>`__ - Access CLI commands relating
   to system
-  `drpcli tasks <drpcli_tasks.html>`__ - Access CLI commands relating
   to tasks
-  `drpcli templates <drpcli_templates.html>`__ - Access CLI commands
   relating to templates
-  `drpcli tenants <drpcli_tenants.html>`__ - Access CLI commands
   relating to tenants
-  `drpcli users <drpcli_users.html>`__ - Access CLI commands relating
   to users
-  `drpcli version <drpcli_version.html>`__ - Digital Rebar Provision
   CLI Command Version
-  `drpcli version_sets <drpcli_version_sets.html>`__ - Access CLI
   commands relating to version_sets
-  `drpcli workflows <drpcli_workflows.html>`__ - Access CLI commands
   relating to workflows
