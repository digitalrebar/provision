.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; High Availability

.. _rs_high_availability:

High Availability
#################

There are two strategies available for implementing high availability in dr-provision: automated failover using Raft for
consensus and liveness checking, and manual failover via synchronous replication.  The former is new in 4.6.0.
The latter has been available since v4.3.0, and will continue to remain available for the foreseeable future.

Prerequisites
~~~~~~~~~~~~~

There are a few conditions that need to be met in order to set up an HA cluster of dr-provision nodes:

#. A fast network between the nodes you want to run dr-provision on.  Data replication between HA nodes
   uses synchronous log replay and file transfer, so a fast network (at least gigabit Ethernet) is required to
   not induce too much lag into the system as a whole.

#. Enough storage space on each node to store a complete copy of all replicated data.  dr-provision will wind up
   replicating all file, ISO, job log, content bundle, plugin, and writable data to all nodes participating in the
   HA cluster.  The underlying storage should be fast enough that write speeds do not become a bottleneck -- we
   recommend backing the data with a fast SSD or NVME device.

#. A high-availability entitlement in your license.  High-availability is a licensed enterprise feature.  If you
   are not sure if your license includes high-availability support, contact support@rackn.com.

#. A virtual IP address that client and external traffic can be directed to.  If using dr-provision's internal
   IL address management (i.e not using the --load-balanced command line option), dr-provision will handle adding and
   removing the virtual IP from a specified interface and sending out gratuitous ARP packets on failover events, and
   the nodes forming the HA cluster must be on the same layer 2 network.  If using an external load balancer,
   then the virtual IP must point to the load balancer, and that address will be used by everything outside of the
   cluster to communicate with whichever cluster node is the active one.

Synchronous Replication
~~~~~~~~~~~~~~~~~~~~~~~

Synchronous Replication is the feature used to implement high availability with manual failover along with
streaming realtime backups from dr-provision v4.3.0 thru v4.6.0.  It will continue to be present going forward
as the realtime backup and protocol, but the shift towards consensus with automated failover will be the path
going forward for high availability.

When operating in synchronous replication mode, there must be at least 2 servers -- one active, and at least 1
passive node.

Consensus via Raft
~~~~~~~~~~~~~~~~~~

Consensus via Raft is the feature used to implement high availability with automated failover.  This mode requires that
you have at least 3 servers in a cluster.  More servers in a cluster are also permitted, but there must be an odd number
to prevent the cluster from deadlocking in the case of communication failures that isolate half of the cluster from the
other half.  Consensus via raft also requires a stable IP address:port that can be used for the replication protocol.

Contraindications
~~~~~~~~~~~~~~~~~

The high-availability support code in dr-provision assumes a model where either:

* There is a single IP available for the HA cluster.  This requires one of the following two items:

  * The machines are in the same layer 2 broadcast domain to allow for moving the HA IP address via gratuitous AR

  * An external load balancer is responsible for holding the virtual IP address and directing all traffic to the
    current active node.

* The writable storage that dr-provision uses is not shared (via NFS, iSCSI, drbd, or whatever) between servers running
  dr-provision.

* If running in automated failover, there must be at least 3 servers in the cluster, and there must be an odd number
  of servers in the cluster.

If none of the above are true, then you cannot use dr-provision in high-availability mode.

* If you are running on separate broadcast domains, you will need to either ensure that there is an alternate mechanism for
  ensuring that packets destined for the HA IP address get routed correctly, or accept that provisioning operations
  will fail from the old network until the clients are informed of the new IP address.

* If you are using a shared storage mechanism (NFS, DRBD, iSCSI, or some form of distributed filesystem), then you should
  not use our integrated HA support, as it will lead to data corruption.  You should also make sure you never run more than
  one instance of dr-provision on the same backing data at the same time, or the data will be corrupted.

It is possible to use shared storage that replicates extended attributes for the tftproot space.  This will reduce transfer
times for replication, but only some distributed filesystems or shared devices support extended attribute sharing.

Configuration
~~~~~~~~~~~~~

Aside from settings listed later in this section, configuration flags and startup options for the dr-provision
services participating in an HA cluster should be identical.  It is not required that the servers participating
in the HA cluster have identical versions, but they must be running on the same OS and system architecture types.
If you try to add a server version to a cluster that is incompatible, it will fail with an error telling
you what to do to resolve the issue.

High Availability Startup Options
---------------------------------

--static-ip (or the environment variable RS_STATIC_IP)
  Not specifically a high-availability startup option, if it is configured it must be different
  on each server.

--drp-id (or the environment variable RS_DRP_ID)
  Also not specifically a high-availability startup option, this must be different on each server.

--ha-id (or the environment variable RS_HA_ID)
  Must be the same on all nodes participating in an HA cluster.

--ha-enabled (or the environment variable RS_HA_ENABLED)
  Must be included on all nodes participating in an HA cluster.

--ha-address (or the environment variable RS_HA_ADDRESS)
  This is the IP address and netmask of the virtual IP that the active cluster member will use
  for communication.  It must be in CIDR format (aa.bb.cc.dd/xx) when not using an external load
  balancer, and a raw IP address when using an external load balancer.

--ha-interface (or the environment variable RS_HA_INTERFACE)
  This is the Ethernet interface that the ha address should be added to and removed from when
  dr-provision transitions between active and passive.  Only applicable when not using an external
  load balancer.

--ha-passive (or the environment variable RS_HA_PASSIVE)
  This must be true on the nodes that should start as passive nodes by default.  In practice, this means
  every node after the initial node.

--ha-join (or the environment variable RS_HA_JOIN)
  The URL of the active node that should be contacted when starting replication as a passive node in
  a synchronous replication cluster.  If not present, this defaults to https://$RS_HA_ADDRESS:$RS_API_PORT/

--ha-token (or the environment variable RS_HA_TOKEN)
  This is the authentication token that HA nodes use to authenticate and communicate with each other.
  It should be identical across the nodes, and it should be a superuser auth token with a long lifetime.
  With the default usernames, you can generate such a token with::

      drpcli users token rocketskates ttl 3y

  and then extracting the Token field from the resulting JSON.

--ha-interface-script (or the environment variable RS_HA_INTERFACE_SCRIPT)
  This is the full path to the script that should be run whenever dr-provision needs to add or remove the
  ha address to the ha interface.  If not set, dr-provision defaults to using ``ip addr add`` and ``ip addr del``
  internally on Linux, and ``ifconfig`` on Darwin.  You can use the following example as a starting point::

    #/usr/bin/env bash
    # $1 is the action to perform.  "add" and "remove" are the only ones supported for now.
    # $2 is the network interface to operate on.  It will be set to the value of --ha-interface.
    # $3 is the address to add or remove.  It will be set to the value of --ha-address.
    case $1 in
       add)    sudo ip addr add "$3" dev "$2";;
       remove) sudo ip addr del "$3" dev "$2";;
       *) echo "Unknown action $1"; exit 1;;
    esac

  Customize to taste to suit your preferred method of getting authority to add and remove addresses
  to interfaces.

--ha-consensus-addr (or the environment variable RS_HA_CONSENSUS_ADDR)
  This is the address:port that this node will use for all consensus traffic.  It must be accessible
  by all the nodes that will participate in the cluster, and it will both originate TCP connections and listen
  for incoming traffic on this address:port combination.

ha-state.json
~~~~~~~~~~~~~

As of version 4.6.0, the ha-state.json file will be the proxy Source of Truth for all high availability
settings.  Settings in ha-state.json take precedence over any from the commandline or environment, and they
will be automatically updated as conditions change as a result of HA-related API requests and general cluster
status changes.  A sample ha-state.json looks like this::

    {
      "ActiveUri": "",
      "ApiUrl": "",
      "ConsensusAddr": "",
      "ConsensusEnabled": false,
      "ConsensusID": "ab0f7bec-5c48-45c3-8970-b3543ec2e9d4",
      "ConsensusJoin": "",
      "Enabled": false,
      "HaID": "",
      "LoadBalanced": false,
      "Observer": false,
      "Passive": false,
      "Roots": [],
      "Token": "",
      "Valid": true,
      "VirtAddr": "",
      "VirtInterface": "",
      "VirtInterfaceScript": ""
    }

ActiveUrl
---------

ActiveUrl is the URL that external services and clients should use to talk to the dr-provision cluster.
It is automatically populated when a cluster is created wither via API or by booting with the appropriate
command-line options and a missing or invalid ha-state.json.  This setting must be the same across all
members participating in a cluster, and in a consensus cluster that is enforced by the consensus protocol.

ApiUrl
------

ApiUrl is the URL used to contact the current node.  It is automatically populated on every start of the current node.
It is specific to an individual node.

ConsensusAddr
-------------

ConsensusAddr is the address:port that all consensus traffic will go over on this node.  It is initially populated
by the --ha-consensus-addr commandline flag.  It is specific to an individual node.

ConsensusEnabled
----------------

ConsensusEnabled indicates whether this node can participate in a consensus cluster.  It is automatically set
to true when ConsensusAddr is not empty.  It must be true on all nodes of a consensus cluster, but can be
different when using synchronous replication.

ConsensusID
-----------

ConsensusID is set when loading an invalid ha-state.json for the first time, and must not be changed afterwards.
It is what the node uses to uniquely identify itself to other cluster nodes, and it must be unique.

ConsensusJoin
-------------

ConsensusJoin is the URL for the current consensus cluster leader, if any.  It is automatically updated by
the consensus replication protocol, and should not be manually edited.

Enabled
-------

Enabled is set when either form of high availability is enabled on this node.  It corresponds to the --ha-enabled
command line option.

HaID
----

HaID is the shared high-availability ID of the cluster.  This setting must be the same across all
members participating in a cluster, and in a consensus cluster that is enforced by the consensus protocol.
It corresponds to the --ha-id commandline option.

LoadBalanced
------------

LoadBalanced indicates that the HA address is managed by an external load balancer instead of by dr-provision.
This setting must be the same across all members participating in a cluster, and in a consensus cluster that is
enforced by the consensus protocol.  It coresponds to the --ha-load-balanced command line option.

Observer
--------

Observer indicates that this node can participate in a consensus cluster, but cannot become the active dr-provision
node.  It is intended to be set when you are setting up a server to act as a consensus tiebreaker, realtime backup,
repoting endpoint, or similar use.

Passive
-------

Passive indicates that this node is not the active node in the cluster.  All nodes but the current active
node must be Passive, and in a consensus cluster that is enforced by the consensus replication protocol.
It corresponds to the --ha-passive commandline option.

Roots
-----

Roots is the list of current trust roots for the consensus protocol.  All consensus traffic is secured via TLS
1.3 mutual authentication, and the self-signed certificates in this list are uses as the trust roots for that
mutual auth process.  Individual trust roots are valid for 3 months, and are rotated every month.

Token
-----

Token is the authentication token that can be used for nodes participating in the same cluster to talk to
each other's APIs. In both cluster types, Token will be rotated on a regular basis.

Valid
-----

Valid indicates that the state stored in ha-state.json is valid.  If state is not valid, it is populated with
matching parameters from the command line options, otherwise it takes precedence over command line options.

VirtAddr
--------

VirtAddr is the address that all external traffix to the cluster should sue to communicate to the cluster.
If LoadBalanced is true, it should be a raw IP address, otherwise it should be a CIDR address in address/prefix
form.  It must be the same on all nodes in a cluster, and corresponds to the --ha-address command line option.

VirtInterface
-------------

If LoadBalanced is false, VirtInterface is the name of the network interface that VirtAddr will be added or
removed from.  It is specific to each node, and corresponds to the --ha-interface commandline option.

VirtInterfaceScript
-------------------

If present, this is the name of the script that will be run whenever we need to add or remove VirtAddr
to VirtInterface.It is specific to each node, and corresponds to the --ha-interface-script commandline option.

Bootstrapping Consensus via Raft (v4.6.0 and later)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

In 4.6 and later, you can bootstrap, add nodes to, and remove nodes from a consensus cluster using `drpcli` without
needing to stop nodes for manual reconfiguration or mess with systemd config files.  This is the preferred method of
high availability.

Self-enroll the initial active node
-----------------------------------

To start the initial active node, you can use the `drpcli system ha enroll` command to have it
enroll itself.  The form of the command to run is as follows::

    drpcli system ha enroll $RS_ENDPOINT username password \
        ConsensusAddr address:port \
        Observer true/false \
        VirtInterface interface \
        VirtInterfaceScript /path/to/script \
        HaID ha-identifier \
        LoadBalanced true/false \
        VirtAddr virtualaddr

The last 3 of those settings can only be specified during self-enroll, and even then they can only be specified
if the system you are self-enrolling is not already in a synchronous replication cluster.

You also can only specify VirtInterface and VirtInterfaceScript if LoadBalanced is false.

If any errors are returned during that call, they should be addressed and the command retried.
Once the command finished without error, the chosen system will be in a single node Raft cluster
that is ready to have other nodes added to the cluster.

Adding additional nodes
-----------------------

To add additional nodes to an existing cluster, you also use
`drpcli system ha enroll` against the current active node in that cluster::

    drpcli system ha enroll https://ApiURL_of_target target_username target_password \
        ConsensusAddr address:port \
        Observer true/false \
        VirtInterface interface \
        VirtInterfaceScript /path/to/script

This will get the global HA settings from the active node in the cluster, merge those settings with the
per-node settings from the target node and the rest of the settings passed in on the command line, and direct
the target node to join the cluster using the merged configuration.

**NOTE** The current data on the target node will be backed up, and once the target node has joined the
cluster it will mirror all data from the existing cluster.  All backed up data will be inaccessible from that point.

Other consensus commands
------------------------

`drpcli system ha` has several other commands that you can use to examine the state of consensus on a node.

* `drpcli system ha active` will get the Consensus ID of the node that is currently responsible for
  all client communication in a consensus cluster.  It is possible for this value to be unset if the
  active node has failed and the cluster is deciding on a new active node.

* `drpcli system ha dump` will dump the user-visible parts of the backing finite state machine that
  is responsible for keeping track of the state of the cluster.

* `drpcli system ha failOverSafe` will return true if there is at least one node in the cluster that
  is completly up-to-date with the active node, and it will return false otherwise.  You can pass
  a time to wait (up to 5 seconds) for the cluster to be fail over safe as an optional argument.

* `drpcli system ha id` returns the Consensus ID of the node you are takling to.

* `drpcli system ha leader` returns the Consensus ID of the current leader of the Raft cluster.  This can
  be different than the active ID if the cluster is in the middle of determining which cluster member is
  best suited to handling external cluster traffic.

* `drpcli system ha peers` returns a list of all known cluster members.

* `drpcli system ha state` returns the current HA state of an individual node.

Bootstrapping Synchronous Replication (pre-v4.6.0 style)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

This bootstrapping documentation will assume that you are working with dr-provision running as a native service
managed by systemd on a Linux server.

The Initially Active Node
-------------------------

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
---------------------------

WARNING: Do not start a passive endpoint(s) in "normal mode."  When installing a passive endpoint, the active
endpoint _must_ be available when the endpoint is started.

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

Troubleshooting
~~~~~~~~~~~~~~~

Log Verification
----------------

It is normal to see ``Error during replication: read tcp [passive IP]:45786->[cluster IP]:8092: i/o timeout`` on the
passive endpoints logs when the active endpoint is killed or switches to passive mode.  This is an indication that the
active endpoint has stopped sending updates.


Transfer Start-up Time
----------------------

It may take up to a minute for a passive endpoint to come online after it has received ``-USR1`` signals.

Network Interface Locked
------------------------

It is possible for the HA interface to become locked if you have to stop and restart the service during configuration
testing.  To clear the interface, use ```ip addr del [ha ip] dev [ha interface]```

This happens because Digital Rebar is attaching to (and detaching from) the cluster IP.  If this process is interrupted,
then the association may not be correctly removed.

WAL File Checksums
------------------

When operating correctly, all the WAL files should match on all endpoints.  You can check the signature of the wal files
using `hexdump -C`

For example:

  :: 

    cd /var/lib/dr-provision/wal
    hexdump -C base.0 |less

Active Endpoint File ha-state is Passive:true
---------------------------------------------

This only applies for Synchronous Replication, and not Consensus.

Digital Rebar uses the ``ha-state.json`` file in it's root directory (typically ``/var/lib/dr-provision``) to track
transitions from active to passive state.

.. note:: removing this file incorrectly can cause very serious problems!  This is a last resort solution.

The ``ha-state.json`` file has a single item JSON schema that changes from true to false depending on the endpoint HA state.  This file can be updated or change to force a reset.  The dr-provision server must be restarted afterwards.

  ::

    {"Passive":false}


When making this changes, stop ALL dr-provision servers in the HA cluster.  Fix the state files for all servers.
Start the selected Active endpoint first.  After it is running, start the passive endpoints.
