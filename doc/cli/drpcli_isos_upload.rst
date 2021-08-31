drpcli isos upload
------------------

Upload the isos [src] as [dest]

Synopsis
~~~~~~~~

The DRP files API allows exploding a compressed file, using bsdtar,
after it has been uploaded. This can be very helpful when multiple files
or a full directory tree are being uploaded.

This is a two stage process enabled by the –explode flag. The first
stage simply uploads the compressed file to the target path and
location. The second stage explodes the file in that path.

For example: *drpcli files upload my.zip as mypath/my.zip –explode*

The above command will upload the *my.zip* file to the files *mypath*
location. It will also expand all the files in *my.zip* into */mypath*
after upload. All paths in *my.zip* will be preserved and created
relative to */mypath*.

::

   drpcli isos upload [src] as [dest] [flags]

Options
~~~~~~~

::

         --explode   After upload, file will be untarred/unzipped in file's local path
     -h, --help      help for upload

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

     -c, --catalog string          The catalog file to use to get product information (default "https://repo.rackn.io")
     -C, --colors string           The colors for JSON and Table/Text colorization.  8 values in the for 0=val,val;1=val,val2... (default "0=32;1=33;2=36;3=90;4=34,1;5=35;6=95;7=32;8=92")
     -d, --debug                   Whether the CLI should run in debug mode
     -D, --download-proxy string   HTTP Proxy to use for downloading catalog and content
     -E, --endpoint string         The Digital Rebar Provision API endpoint to talk to (default "https://127.0.0.1:8092")
     -X, --exit-early              Cause drpcli to exit if a command results in an object that has errors
     -f, --force                   When needed, attempt to force the operation - used on some update/patch calls
     -F, --format string           The serialization we expect for output.  Can be "json" or "yaml" or "text" or "table" (default "json")
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

-  `drpcli isos <drpcli_isos.html>`__ - Access CLI commands relating to
   isos
