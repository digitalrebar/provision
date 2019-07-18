Linode.com
==========

.. index::
  pair: Digital Rebar Provision; Linode

.. _rs_setup_linode:

Existing machines can use the `join-up.sh` script to join a DRP endpoint without having to go through a netboot.  The Linode integration uses this feature to manage virtual machines.

These instructions can also be adapted to work on GCE or other cloud infrastructure.

The following video shows the process: https://youtu.be/e4Vp3quzlhM

Stackstripts
------------

You can use public Stackscripts to quickly setup your DRP endpoint and attach nodes.  They perform the same basic tasks as the code below.

  1. Create an endpoint using `zehicle/drp` and supply your version and password.
  1. Verify that the DRP Endpoint is available at the `https://[endpoint_ip]:8092`
  1. Use `zehicle/drp-node` to attach nodes to the endpoint.  You'll need to know the endpoint_ip.
  1. profit.


Install DRP in Linode
---------------------

You can use the Amazon Linux AMI.  While more is recommended, make sure that you have at least 1 GB of RAM.  You should open ports `8091` and `8092` to access the DRP server.


  ::

    #!/bin/bash

    ### Install DRP from Tip
    curl -fsSL get.rebar.digital/tip | bash -s -- install --systemd --version=tip --drp-password=r0cketsk8ts

    ### Now open the right firewall ports for DRP
    firewall-cmd --permanent --add-port=8092/tcp
    firewall-cmd --permanent --add-port=8091/tcp
    firewall-cmd --reload

    ### Install Content and Configure Discovery
    drpcli contents upload catalog:task-library-tip
    drpcli contents upload catalog:drp-community-content-tip
    drpcli workflows create '{"Name": "discover-linode", "Stages":
      ["discover", "runner-service", "complete"]
    }'
    drpcli prefs set defaultWorkflow discover-linode unknownBootEnv discovery

    ### Capture Node Info 
    drpcli profiles create '{"Name":"linode"}'
    drpcli profiles set linode param cloud/provider to "LINODE"
    drpcli machines set linode param cloud/instance-id to "\"${LINODE_ID}\""
    drpcli profiles set linode param cloud/username to "${LINODE_LISHUSERNAME}"
    drpcli profiles set linode param cloud/instance-type to "\"${LINODE_RAM}\""
    drpcli profiles set linode param cloud/placement/availability-zone to "\"${LINODE_DATACENTERID}\""

Once the system is online, you can access DRP using https://[DRP public address]:8092.


Join a machine to a DRP Endpoint in Linode
------------------------------------------

Once you have a DRP endpoint installed in Linode

  ::

    #!/bin/bash
    curl -fsSL ${DRP_IP}:${DRP_PORT}/machines/join-up.sh | sudo bash --