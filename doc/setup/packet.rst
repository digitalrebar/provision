Packet.net
==========

.. index::
  pair: Digital Rebar Provision; packet

.. _rs_setup_packet:

.. note:: you can use the code `RACKN100` to get a $25 signup discount!

Install DRP in Packet
---------------------

Create a machine in Packet.net.  For initial testing, a t1 is sufficient.  We recommend using Centos as the DRP host O/S.

From the SSH console, you'll need to run the install script adding the public IP address of the system.

To begin, execute the following commands in a shell or terminal:

  ::

     export DRP_IP=[PUBLIC IP]
     mkdir drp ; cd drp
     curl -fsSL get.rebar.digital/stable | bash -s -- --systemd --ipaddr=$DRP_IP install
     # wait for install
     systemctl daemon-reload &&  systemctl start dr-provision
     systemctl daemon-reload &&  systemctl enable dr-provision
     echo "You can now visit https://$DRP_IP:8092"


Once DRP is running, perform added confirmation steps with the CLI:

  ::

     # install plugin provider from catalog
     drpcli catalog item install packet-ipmi
     # install discovery image and set default
     drpcli bootenvs uploadiso sledgehammer
     drpcli prefs set defaultWorkflow discover-packet unknownBootEnv discovery defaultBootEnv sledgehammer defaultStage discover
     # optional
     drpcli bootenvs uploadiso ubuntu-18.04-install &
     drpcli bootenvs uploadiso centos-7-install &
     echo "You can now use https://$DRP_IP:8092"


Join a machine to a DRP Endpoint in Packet
------------------------------------------

To add machines to the DRP endpoint, we highly recommend using the Plugin process below.  If you'd like to create machines yourself, then the following steps are required to configure the Packet IPXE system.

When creating a machine, you must choose the `custom ipxe` O/S type and set the `IPXE Script URL` to be `http://[DRP IP]:8091/default.ipxe` where you use the IP address of the DRP server.  Packet will confirm this URL is accessible.

In addition, you should set `Persist PXE on Reboot` to true under additional settings.  This is very important if you want to be able to reprovision or reboot systems and maintain DRP control.

If you use the DRP plugin below, it will make these setting automatically.


Using the Plugin to add a machine to a DRP Endpoint in Packet
--------------------------------------------------------------

If you are using Packet.net as a reference platform.  The following steps are specific to that platform.

From the UX
~~~~~~~~~~~

The following steps add the Packet-IPMI plugin:

  1. From the Plugins page, Add the Packet-IPMI plugin
  1. From the Packet-IPMI plugin panel, set your API key and Project UUID
  1. Save

After the plugin is created, revisit the Packet-IPMI plugin page.  It will display a `machine create` line that allows you to quickly create machines that have the correct Param settings for adding machines.

From the CLI
~~~~~~~~~~~~

Make sure you set your information in the exports!

  ::

     # configure plugin
     # UPDATE THESE: Packet Project Information for Plugin
     export PACKET_API_KEY="[packet_api_key]"
     export PACKET_PROJECT_ID="[packet_project_id]"
     drpcli plugins create '{ "Name": "packet-ipmi",
       "Params": {
         "packet/api-key": "$PACKET_API_KEY",
         "packet/project-id": "$PACKET_PROJECT_ID"
       },
       "Provider": "packet-ipmi"
      }'
     # verify it worked - should return true
     drpcli plugins show packet-ipmi | jq .Available

.. note:: The URLs provided for plugin downloads will change overtime for newer versions