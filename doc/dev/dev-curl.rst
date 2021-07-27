.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Using CURL with the API

.. _rs_dev_curl:

Using CURL with the API
~~~~~~~~~~~~~~~~~~~~~~~

In many cases, it is helpful to test the API using CURL commands.

.. note: This topic is covered in several areas, please try and consolidate them into this page1


Here is the basic structure of a curl command using an auth token

  ::

    export RS_TOKEN=$(./drpcli users token [username] | jq -r .Token)
    curl -H "Authorization: Bearer $RS_TOKEN" --insecure https//[endpoint url]/api/v3/info

You can also use user security

  ::

    curl --user rocketskates:r0cketsk8ts --insecure https//[endpoint url]/api/v3/info


.. note: the ``--insecure`` flag is needed if you are using self-signed certificates.

.. _rs_dev_curl_iso:

Uploading ISO, File or Plugin Providers
---------------------------------------

For binary items, Digital Rebar Provision expects either and "application/octet-stream" or "multipart/form-data" for the POST.  Octet is helpful for direct sending from programs like the CLI.  Multipart is helpful when sending files from a webapp where the browser is directly responsible for handling file uploads.

  ::

    export RS_TOKEN=$(./drpcli users token [username] | jq -r .Token)
    curl -X POST  --insecure \
      -H "Authorization: Bearer $RS_TOKEN" \
      -H "Content-Type: multipart/form-data" \
      -F "file=@[filepath]/[filename]"\
      https://[endpoint url]/api/v3/isos/[filename]

.. _rs_dev_patch:

Using PATCH with CURL
---------------------

For updates with the Digital Rebar API, PATCH is strongly prefered over PUT.

Here's an example of updating an object using PATCH

  ::

    curl -X PATCH --insecure \
        -H "Authorization: Bearer $RS_TOKEN" \
        -H "Content-Type: application/json" \
        -d '[{"op":"replace", "path":"/Fingerprint/CloudInstanceID", "value":"cloud:314159265"}]' \
        {{.ApiURL}}/api/v3/machines/$RS_UUID