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

    ./drpcli users token [username]
  	export TOKEN=[token]
  	curl -H "Authorization: Bearer $TOKEN" --insecure https//[endpoint url]/api/v3/info

You can also use user security

  ::

  	curl --user rocketskates:r0cketsk8ts --insecure https//[endpoint url]/api/v3/info


.. note: the ``--insecure`` flag is needed if you are using self-signed certificates.

Uploading ISO, File or Plugin Providers
---------------------------------------

For binary items, Digital Rebar Provision expects either and "application/octet-stream" or "multipart/form-data" for the POST.  Octet is helpful for direct sending from programs like the CLI.  Multipart is helpful when sending files from a webapp where the browser is directly responsible for handling file uploads.

  ::

  	./drpcli users token [username]
  	export TOKEN=[token]
    curl -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: multipart/form-data" -F "file=@[filepath]/[filename]" --insecure https://[endpoint url]/api/v3/isos/[filename]
