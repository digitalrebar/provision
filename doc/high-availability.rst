.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; High Availiability

.. _rs_high_availability:


High Availability
~~~~~~~~~~~~~~~~~

High Availability is a new feature that will be present in dr-provision starting with version 4.3.0.
To begin with, the working model will be Active/Passive with unidirectional data replication and
manual failover.  Automated failover and failback may be added in future versions of dr-provision.
Currently, high-availability should be considered a technical preview -- the details of how it is implemented
and how to manage it are likely to change based on customer feedback.


Prerequisites
-------------

There are a few conditions that need to be met in order to set up an HA cluster of dr-provision machines:

#. A fast network between the nodes you want to run dr-provision on.  Data replication between HA nodes
   uses synchronous log replay and file transfer, so a fast network (at least gigabit Ethernet) is required to
   not induce too much lag into the system as a whole.

#. A virtual IP address that can be moved from node to node.  dr-provision will handle adding and removing
   the virtual IP from a specified interface and sending out gratuitous ARP packets on failover events.  If the
   nodes forming the HA cluster are not in the same layer 2 broadcast domain or subnet, then you will need to
   arrange for the appropriate route changes required.

#. Enough storage space on each node to store a complete copy of all replicated data.  dr-provision will wind up
   replicating all file, ISO, job log, content bundle, plugin, and writable data to all nodes participating in the
   HA cluster.  The underlying storage should be fast enough that write speeds do not become a bottleneck -- we
   recommend backing the data with a fast SSD or NVME device.

#. A high-availability entitlement in your license.  High-availability is a licensed enterprise feature.  If you
   are not sure if your license includes high-availability support, contact support@rackn.com.

Configuration
-------------

Aside from settings listed later in this section, any and configuration flags and startup options for the dr-provision
services participating in an HA cluster should be identical.  It is not required that the servers participating in the
HA cluster have identical versions, but they must be running on the same OS and system architecture types.
If you try to add a server version to a cluster that is incompatible, it will fail to be added with an error telling
you what to do to resolve the issue.

High Availability Startup Options
=================================

--static-ip (or the environment variable RS_STATIC_IP)
  Not specifically a high-availability startup option, if it is configured it must be different
  on each server.

--drp-id (or the environment variable RS_DRP_ID)
  Also not specifically a high-availability startup option, this must be different on each server.

--ha-id (or the environment variable RS_HA_ID)
  Must be the same on all nodes participating in an HA cluster.

--ha-enabled (or the environment variable RS_HA_ENABLED)
  Must be true on all nodes participating in an HA cluster.

--ha-address (or the environment variable RS_HA_ADDRESS)
  This is the IP address and netmask of the virtual IP that the active cluster member will use
  for communication.  It must be in CIDR format (aa.bb.cc.dd/xx)

--ha-passive (or the environment variable RS_HA_PASSIVE)
  This must be true on the nodes that should start as passive nodes by default.  In practice, this means
  every node after the initial node.

--ha-token (or the environment variable RS_HA_TOKEN)
  This is the authentication token that HA nodes use to authenticate and communicate with each other.
  It should be identical across the nodes, and it should be a superuser auth token with a long lifetime.
  With the default usernames, you can generate such a token with::

      drpcli users token rocketskates ttl 3y

  and then extracting the Token field from the resulting JSON.

Bootstrapping
-------------

This bootstrapping documentation will assume that you are working with dr-provision running as a native service
managed by systemd on a Linux server.

The Initially Active Node
=========================

To start bootstrapping an HA cluster, start by installing what you want ot be the default active dr-provision node.
Once it is up and running, create a file named /etc/systemd/system/dr-provision.service.d/20-ha.conf with
the following contents::

    [Service]

    # RS_HA_ENABLED tells dr-provision to operate in high-availability mode.
    Environment=RS_HA_ENABLED=true

    # RS_HA_INTERFACE is the network interface that dr-provision will add/remove the
    # virtual IP address to.  This interface should be one that machines being managed by
    # dr-provision can access.
    Environment=RS_HA_INTERFACE=kvm-test

    # RS_HA_ADDRESS is the IP address and netmask in CIDR format that all communication to
    # and from dr-provision will use.
    Environment=RS_HA_ADDRESS=192.168.124.200/24

    # RS_HA_ID is the cluster ID.  This must be the same for all members participating in the cluster.
    Environment=RS_HA_ID=8c:ec:4b:ea:d9:fe

    # RS_HA_TOKEN is a long-lived access token that the cluster nodes will use to authenticate with each other.
    # You can generate a useable token with:
    #
    #    $ drpcli users token rocketskates ttl 3y |jq -r '.Token'
    Environment=RS_HA_TOKEN=your-token

Once that file is created, reload the config snre restart dr-provision::

    $ systemctl daemon-reload
    $ systemctl restart dr-provision

When dr-provision comes back up, it will be running on the IP address you set aside as the HA IP address.

The Initially Passive Nodes
===========================

Perform the same installation steps you used for the initally active node, but add one extra line to
the /etc/systemd/system/dr-provision.service.d/20-ha.conf file::

    Environment=RS_HA_PASSIVE=true

which will cause the node to come up as a passive node when you start it up.  The first time you start up the node,
it will replicate all of the runtime data from the active mode, which (depending on your network bandwidth and
how busy the active node is) may take awhile.  You can monitor the progress of the replication by
watching the output of ```journalctl -fu dr-provision``` --- when it says "Stream switch to realtime streaming" the
passive node is fully caught up to the active node.

Switching from Active to Passive
--------------------------------

To switch a dr-provision instance from active to passive, send it the USR2 signal.  To switch it to active, send it the
USR1 signal.  As of right now, there are no other mechanisms (automated or manual) for changing HA state on a node.