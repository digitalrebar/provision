.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; FAQ
  pair: Digital Rebar Provision; Troubleshooting

.. _rs_faq:

FAQ / Troubleshooting
~~~~~~~~~~~~~~~~~~~~~

The following section is designed to answer frequently asked questions and help troubleshoot Digital Rebar Provision installs.

.. _rs_bind_error:

Bind Error
----------

Digital Rebar Provision will fail if it cannot attach to one of the required ports.

* Typical Error Message: "listen udp4 :67: bind: address already in use"
* Additional Information: The conflicted port will be included in the error between colons (e.g.: `:67:`)
* Workaround: If the conflicting service is not required, simply disable that service
* Resolution: Stop the offending service on the system.  Typical corrective actions are:

  * 67 - dhcp.  Correct with `sudo pkill dnsmasq`

See the port mapping list on start-up for a complete list.

.. _rs_gen_cert:

Generate Certificate
--------------------

Sometimes the cert/key pair in the github tree is corrupt or not sufficient for the environment.  The following command can be used to rebuild a local cert/key pair.

  ::

    sudo openssl req -new -x509 -keyout server.key -out server.crt -days 365 -nodes

It may be necessary to install the openssl tools.

Add SSH Keys to Authorized Keys
-------------------------------

To have provisioned operating systems (including discovery/sledgehammer) add keys, you should set the ``access_keys`` parameter with a hash of the desired keys.  This can be accomplished by editing the root access profile to add your key(s) and then update the profile via the CLI.

  ::

    vi assets/profiles/root-access.yaml
    ./drpcli profiles update root-access - < assets/profiles/root-access.yaml
    
NOTE: By default, these changes are targeted at the ``root-access`` profile and you will need to add that profile to selected machines for your keys to be injected.

If you want this parameter applied to all machines by default, then you should change ``root-access`` to ``global`` in the yaml file and command line.  

  ::

    cp assets/profiles/root-access.yaml assets/profiles/global.yaml
    # remember to change root-access to global!
    vi assets/profiles/global.yaml
    ./drpcli profiles update global - < assets/profiles/global.yaml

Turn on autocomplete for the CLI
--------------------------------

The DRP CLI will automatically create the autocomplete file if a path is provided.  You must make sure to use the correct path!  The example below is for Ubuntu.

  ::
  
    ./drpcli autocomplete /etc/bash_completion.d/drpcli
    
Log out and log back in to take effect or run:

  * ``. /etc/bash_completion`` # On Ubuntu
  * ``. /etc/profile.d/bash_completion.sh`` # On Centos
  * ``. /usr/local/etc/bash_completion`` # On OS X with bash 4 installed.
    
Turn Up the Debug
-----------------

To get additional debug from dr-provision, set debug preferences to increase the logging.  See :ref:`rs_model_prefs`.

Rocket Skates?
--------------

Rocket Skates was the working name for Digital Rebar Provision during initial development.  Since they are fast and powerful boots, it seemed like a natural name for a Cobbler replacement.

.. figure::  images/rocket.jpg
   :align:   right
   :width: 320 px
   :alt: Code name Rocket Skates
   :target: https://www.pexels.com/photo/aerospace-engineering-exploration-launch-34521/

