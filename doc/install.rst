.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Rocket Skates documentation under Digital Rebar master license
.. index::
  pair: Rocket Skates; Install

.. _rs_install:

Install
~~~~~~~


Running Server
--------------

To run a local copy that will use the local filesystem as a storage area, do the following:

* cd test-data
* sudo ../rocket-skates

Please review `--help` for options like disabling services, logging or paths.

If you get a bind error, you may need to kill dnsmasq using `sudo pkill dnsmasq`

.. note:: I need the sudo to bind the tftp port.  This is configurable, i.e.  *--tftp-port=30000*  

The following pieces endpoints are available:

* https://127.0.0.1:8092/swagger-ui - swagger-ui to explore the API
* https://127.0.0.1:8092/swagger.json - API Swagger JSON file
* https://127.0.0.1:8092/api/v3 - Raw api endpoint
* https://127.0.0.1:8091 - Static files served from the test-data/tftpboot directory
* udp 69 or *--tftp-port* - Static files served from the test-data/tftpboot directory through the tftp protocol
* udp 67 - DHCP Server listening socket - will only server addresses when once configured.  By default, silent.

.. note:: If your SSL certificate is not valid, then follow the :ref:`rs_gen_cert` steps.

.. note:: On OSX, you may have to add a route for broadcast addresses to work.  This can be done with the following
comand.  The 192.168.100.1 is the IP address of the interface that you want to send messages through.

  ::

    sudo route add 255.255.255.255 192.168.100.1


.. _rs_gen_cert:

Generate Certificate
--------------------

Sometimes the cert/key pair in the github tree is corrupt or not sufficient for the environment.  You can run
the following command to rebuild a local cert/key pair.

  ::

    sudo openssl req -new -x509 -keyout server.key -out server.crt -days 365 -nodes


You may need to install the openssl tools.
