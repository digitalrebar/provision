drpcli debug
------------

Gather [type] of debug information and save it to [target]

Synopsis
~~~~~~~~

This command gathers various different types of runtime profile data
from a running dr-provision server, provided it has the /api/v3/debug or
/api/v3/drp_debug. The types of data that can be gathered are:

::

   profile: CPU utilization profile information.  Tracks how much CPU time is being used in which
            functions, based on sampling which functions are running every 10 ms.  If the
            --seconds flag is unspecified, profile will gather 30 seconds worth of data.

   trace: Execution trace information, including information on where execution is blocked on
          various types of IO and synchronization primitives.  If the --seconds flag is unspecified,
          trace will gather 1 second of data.

   heap: Memory tracing information for all live data in memory. heap is always point-in-time data.

   allocs: Memory tracing of all memory that has been allocated since the start of the program
           This includes memory that has been garbage-collected.  alloc is always point-in-time data.

   block: Stack traces of all goroutines that have blocked on synchronization primitives.
          block is always point-in-time data.

   mutex: Stack traces of all holders of contended mutexes.  mutex is always point-in-time data.

   threadcreate: Stack traces of all goroutines that led to the creation of a new OS thread.
                 threadcreate is always point-in-time data.

   goroutine: Stack traces of all current goroutines. goroutine is always point-in-time data.

   index: Returns the indexes of the stacks with the flags of the object.

::

   drpcli debug [type] [target] [flags]

Options
~~~~~~~

::

     -h, --help            help for debug
         --prefix string   Limits the index call to just this prefix type.
         --seconds int     How much debug data to gather, for types that gather data over time.

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

     -c, --catalog string          The catalog file to use to get product information (default "https://repo.rackn.io")
     -d, --debug                   Whether the CLI should run in debug mode
     -D, --download-proxy string   HTTP Proxy to use for downloading catalog and content
     -E, --endpoint string         The Digital Rebar Provision API endpoint to talk to (default "https://127.0.0.1:8092")
     -f, --force                   When needed, attempt to force the operation - used on some update/patch calls
     -F, --format string           The serialization we expect for output.  Can be "json" or "yaml" or "text" or "table" (default "json")
     -H, --no-header               Should header be shown in "text" or "table" mode
     -x, --noToken                 Do not use token auth or token cache
     -P, --password string         password of the Digital Rebar Provision user (default "r0cketsk8ts")
     -J, --print-fields string     The fields of the object to display in "text" or "table" mode. Comma separated
     -r, --ref string              A reference object for update commands that can be a file name, yaml, or json blob
     -T, --token string            token of the Digital Rebar Provision access
     -t, --trace string            The log level API requests should be logged at on the server side
     -Z, --traceToken string       A token that individual traced requests should report in the server logs
     -j, --truncate-length int     Truncate columns at this length (default 40)
     -u, --url-proxy string        URL Proxy for passing actions through another DRP
     -U, --username string         Name of the Digital Rebar Provision user to talk to (default "rocketskates")

SEE ALSO
~~~~~~~~

-  `drpcli <drpcli.html>`__ - A CLI application for interacting with the
   DigitalRebar Provision API
