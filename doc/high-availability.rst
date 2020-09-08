.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; High Availability

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

Aside from settings listed later in this section, configuration flags and startup options for the dr-provision
services participating in an HA cluster should be identical.  It is not required that the servers participating
in the HA cluster have identical versions, but they must be running on the same OS and system architecture types.
If you try to add a server version to a cluster that is incompatible, it will fail with an error telling
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

To start bootstrapping an HA cluster, start by installing what you want to be the default active dr-provision node.
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
    # You can generate a usable token with:
    #
    #    $ drpcli users token rocketskates ttl 3y |jq -r '.Token'
    Environment=RS_HA_TOKEN=your-token

    # RS_HA_PASSIVE is an intial flag (not used after synchronization) to identify the active endpoint.
    Environment=RS_HA_PASSIVE=false

Once that file is created, reload the config and restart dr-provision::

    $ systemctl daemon-reload
    $ systemctl restart dr-provision

When dr-provision comes back up, it will be running on the IP address you set aside as the HA IP address.

The Initially Passive Nodes
===========================

WARNING: Do not start a passive endpoint(s) in "normal mode."  When installing a passive endpoint, the active
endpoint _must_ be available when the endpoing is started.

Perform the same installation steps you used for the initially active node, but change the `RS_HA_PASSIVE` line
to false in the `/etc/systemd/system/dr-provision.service.d/20-ha.conf` file

  ::

    Environment=RS_HA_PASSIVE=true

which will cause the node to come up as a passive node when you start it up.  The first time you start up the node,
it will replicate all of the runtime data from the active mode, which (depending on your network bandwidth and
how busy the active node is) may take awhile.  You can monitor the progress of the replication by
watching the output of ```journalctl -fu dr-provision``` --- when it says "Stream switch to realtime streaming" the
passive node is fully caught up to the active node.

Switching from Active to Passive
--------------------------------

To switch a dr-provision instance between states, an API call will need to be done.  **drpcli** can be used to
send that API call.  Issuing a **POST** request with empty JSON object to **/api/v3/system/active** and
**/api/v3/system/passive** will cause the system to transition to active or passive, respectively.

As of right now, there are no other mechanisms (automated or manual) for changing HA state on a node.

.. note:: When doing a practice failover, the active endpoint should be stopped first.

To stop the active endpoint (becomes passive):

  ::

    // deactivate endpoint (goes into passive mode)
    drpcli system passive

To promote a passive endpoint to active

  ::

    // activate endpoint (goes into active mode)    
    drpcli system active

.. note:: Prior to v4.5.0, Signals were used to shift state.  SIGUSR2 was used to go from active to passive and
  SIGUSR1 was used to go from passive to active.

Install Video
-------------

This video was created at the time of v4.3 beta: https://youtu.be/xM0Zr3iL5jQ.  Please check for more recent updates.


Troubleshooting
---------------

Log Verification
================

It is normal to see ``Error during replication: read tcp [passive IP]:45786->[cluster IP]:8092: i/o timeout`` on the passive endpoints logs when the active endpoint is killed or switches to passive mode.  This is an indication that the active endpoint has stopped sending updates.


Transfer Start-up Time
======================

It may take up to a minute for a passive endpoint to come online after it has recieved ``-USR1`` signals.

Network Interface Locked
========================

It is possible for the HA interface to become locked if you have to stop and restart the service during configuration testing.  To clear the interface, use ```ip addr del [ha ip] dev [ha interface]```

This happens because Digital Rebar is attaching to (and detaching from) the cluster IP.  If this process is interrupted, then the association may not be correctly removed.

WAL File Checksums
==================

When operating correctly, all the WAL files should match on all endpoints.  You can check the signature of the wal files using `hexdump -C`

For example:

  :: 

    cd /var/lib/dr-provision/wal
    hexdump -C base.0

Active Endpoint File ha-state is Passive:true
=============================================

Digital Rebar uses the ``ha-state.json`` file in it's root directory (typically ``/var/lib/dr-provision``) to track transitions from active to passive state.

.. note:: removing this file incorrectlycan cause very serious problems!  This is a last resort solution.

The ``ha-state.json`` file has a single item JSON schema that changes from true to false depending on the endpoint HA state.  This file can be updated or change to force a reset.  The dr-provision server must be restarted afterwards.

  ::

    {"Passive":false}


When making this changes, stop ALL dr-provision servers in the HA cluster.  Fix the state files for all servers.  Start the selected Active endpoint first.  After it is running, start the passive endpoints.
