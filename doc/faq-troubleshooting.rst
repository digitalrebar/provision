.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. DigitalRebar Provision documentation under Digital Rebar master license
.. index::
  pair: DigitalRebar Provision; FAQ
  pair: DigitalRebar Provision; Troubleshooting

.. _rs_faq:

FAQ / Troubleshooting
~~~~~~~~~~~~~~~~~~~~~

The following section is designed to answer frequently asked questions and help troubleshoot DigitalRebar Provision installs.

.. _rs_bind_error:

Bind Error
----------

DigitalRebar Provision will fail if it cannot attach to one of the required ports.

* Typical Error Message: "listen udp4 :67: bind: address already in use"
* Additional Information: The conflicted port will be included in the error between colons (e.g.: `:67:`)
* Workaround: If you do not need the conflicting service, you can disable that service
* Resolution: Stop the offending service on your system.  Typical corrective actions are:

  * 67 - dhcp.  Correct with `sudo pkill dnsmasq`

See the port mapping list on start-up for a complete list.

.. _rs_gen_cert:

Generate Certificate
--------------------

Sometimes the cert/key pair in the github tree is corrupt or not sufficient for the environment.  You can run the following command to rebuild a local cert/key pair.

  ::

    sudo openssl req -new -x509 -keyout server.key -out server.crt -days 365 -nodes


You may need to install the openssl tools.

Rocket Skates?
--------------

Rocket Skates was the working name for Digital Rebar Provision during initial development.  Since they are fast and powerful boots, it seemed like a natural name for a Cobbler replacement.

.. figure::  doc/images/rocket.jpg
   :align:   right
   :width: 320 px
   :alt: Code name Rocket Skates
   :target: https://www.pexels.com/photo/aerospace-engineering-exploration-launch-34521/

