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
    # <UDF name="drp_version" Label="Version to Install" default="stable" example="tip, stable, v3.13, ..." />
    # <UDF name="drp_password" Label="Admin Password" default="r0cketsk8ts" example="password" />
    
    ### Install DRP from Tip
    curl -fsSL get.rebar.digital/tip | bash -s -- install --systemd --version=$DRP_VERSION --drp-password=$DRP_PASSWORD --ipaddr=[DRP public address]
    
    ### Now open the right firewall ports for DRP
    firewall-cmd --permanent --add-port=8092/tcp
    firewall-cmd --permanent --add-port=8091/tcp
    firewall-cmd --reload
    
    ### Install Content and Configure Discovery
    drpcli catalog item install task-library --version=$DRP_VERSION
    drpcli catalog item install drp-community-content --version=$DRP_VERSION
    drpcli workflows create '{"Name": "discover-linode", "Stages":
      ["discover", "network-firewalld", "runner-service", "complete"]
    }'
    drpcli profiles set global param "network/firewalld-ports" to '[
      "22/tcp", "6443/tcp", "8379/tcp",  "8380/tcp", "10250/tcp"
    ]'
    drpcli prefs set defaultWorkflow discover-linode unknownBootEnv discovery

Once the system is online, you can access DRP using https://[DRP public address]:8092.


Join a machine to a DRP Endpoint in Linode
------------------------------------------

Once you have a DRP endpoint installed in Linode

  ::

    #!/bin/bash
    # <UDF name="drp_ip" Label="IP Address of the DRP Endpoint" default="" example="192.168.1.100" />
    # <UDF name="drp_port" Label="Provisioning Port of the DRP Endpoint (not API port)" default="8091" example="8091" />
    # <UDF name="open_ports" Label="Ports to open on the machine" default="22 2379 2380 6443 10250" example="22 6443 10250" />

    for PORT in ${OPEN_PORTS}; do
       firewall-cmd --permanent --add-port=${PORT}/tcp
    done 
    firewall-cmd --reload    

    timeout 300 bash -c 'while [[ "$(curl -fsSL -o /dev/null -w %{http_code} ${DRP_IP}:${DRP_PORT}/machines/join-up.sh)" != "200" ]]; do sleep 5; done' || false
    
    curl -fsSL ${DRP_IP}:${DRP_PORT}/machines/join-up.sh | sudo bash --
    