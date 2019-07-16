Amazon Web Services (AWS)
=========================

.. index::
  pair: Digital Rebar Provision; AWS

.. _rs_setup_aws:

Existing machines can use the `join-up.sh` script to join a DRP endpoint without having to go through a netboot.  The AWS integration uses this feature to manage virtual machines.

These instructions can also be adapted to work on GCE or other cloud infrastructure.


Install DRP in AWS
------------------

You can use the Amazon Linux AMI.  While more is recommended, make sure that you have at least 1 GB of RAM.  You should open ports `8091` and `8092` to access the DRP server.


  ::

    #!/bin/bash
    sudo yum install bsdtar -y
    curl -fsSL get.rebar.digital/tip | bash -s -- install --systemd --version=tip --drp-password=r0cketsk8ts

    ### Install Content and Configure Discovery
    drpcli contents upload catalog:task-library-tip
    drpcli contents upload catalog:drp-community-content-tip
    drpcli workflows create '{"Name": "discover-aws", "Stages":
      ["discover","aws-discover", "runner-service", "complete"]
    }'
    drpcli prefs set defaultWorkflow discover-aws unknownBootEnv discovery


    ### Optional: add some Kubernetes magic
    drpcli plugin_providers upload certs from catalog:certs-tip
    drpcli contents upload catalog:krib-tip
    drpcli profiles create '{"Name":"krib", "Meta": {
        "render": "krib", "reset-keeps": "krib/cluster-profile,etcd/cluster-profile",
      }
    }'
    drpcli profiles set krib param "etcd/cluster-profile" to "krib"
    drpcli profiles set krib param "krib/cluster-profile" to "krib"
    drpcli workflows create '{"Name":"krib-aws", "Stages": [
        "ssh-access", "docker-install", "kubernetes-install","etcd-config","krib-config","krib-helm","krib-live-wait"
      ]
    }'


Once the system is online, you can access DRP using https://[DRP public address].


Join a machine to a DRP Endpoint in AWS
---------------------------------------

Once you have a DRP endpoint installed in AWS

  ::

    #!/bin/bash
    curl -fsSL [DRP Address]:8091/machines/join-up.sh | sudo bash --


The machines started using this process will register with their internal IP address.  By including the `aws-discover` stage, the machines will log their external IP address to the `cloud/public-ipv4` parameter.