.. Copyright (c) 2021 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Platform documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Platform; Certificate Operations

.. _rs_cert_ops:

Certificate Operations
======================

This section will describe how to manage the HTTPS certificates that
dr-provision uses for the API port and the HTTPS static file port.


API Certificate Management
--------------------------

In dr-provision versions earlier than 4.6.0, API certificates can only be managed
via the `--tls-cert` and `--tls-key` startup options.

In dr-provision versions 4.6.0 and later, the `--tls-cert` and `--tls-key` startup
options are only used for loading certificates at initial startup time.  After that,
you must use the `drpcli system certs` commands to manage the TLS certificate
that the API uses.

Get API Certificates
~~~~~~~~~~~~~~~~~~~~

You can fetch the current TLS certificate and private key that the API uses with::

    drpcli system certs get server.crt server.key

This will retrieve the TLS certificate and private key that the API is using,
and saves them in x.509 DER encoded form.

Set API Certificates
~~~~~~~~~~~~~~~~~~~~

You can update the TLS certificate and private key that the API will use for
new connections with::

    drpcli system certs set server.crt server.key

This will upload the X.509 DER encoded certificate and private key to dr-provision,
which (assuming that they are valid) will be used for any new connections.  Additionally,
the new certs will be applied cluster-wide if running in an HA cluster.

Static Certificate Management (4.7.0 and higher)
------------------------------------------------

dr-provision 4.7.0 and higher have a static HTTPS server that will be used as an
alternative to serving static files over HTTP whenever feasible.  It uses the
`--static-tls-cert` and `--static-tls-key` startup options to load the static
HTTPS certs initially.  Afterwards, the `drpcli static certs` commands will
manage the static HTTPS certificates in the same way the `drpcli system certs`
commands manage the API certs.
