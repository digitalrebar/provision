Moved: KRIB (Kubernetes Rebar Integrated Bootstrapping)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

NOTICE: Moving to :ref:`component_krib`.  Please update references.

KRIB (Kubernetes Rebar Integrated Bootstrapping)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

License: KRIB is APLv2

This document provides information on how to use the Digital Rebar *KRIB* content add-on.  Use of this content will enable the operator to install Kubernetes in either a Live Boot (immutable infrastructure pattern) mode, or via installed to local hard disk OS mode.

KRIB uses the `kubeadm <https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/>`_ cluster deployment methodology coupled with Digital Rebar enhancements to help proctor the cluster Master election and secrets management.  With this content pack, you can install Kubernetes in a zero-touch manner.

KRIB does also support production, highly available (HA) deployments, with multiple masters.  To enable this configuration, we've chosen to manage the TLS certificates and etcd installation in the Workflow instead of using the kubeadm process.

This document assumes you have your Digital Rebar Provisioning endpoint fully configured, tested, and working.  We assume that you are able to properly provision Machines in your environment as a base level requirement for use of the KRIB content add-on use.

KRIB Video References
---------------------

The following videos have been produced or presented by RackN related to the Digital Rebar KRIB solution.

* `KRIB Zero Config Kubernetes Cluster channel <https://www.youtube.com/watch?v=SYOHI8DfRMo&list=PLXPBeIrpXjfhKqmTvxI5-0CmgUh82dztr&index=1>`_ on YouTube.
* `KubeCon: Zero Configuration Pattern on Bare Metal <https://youtu.be/Psm9aOWzfWk>`_ on YouTube - RackN presentation at 2017 KubeCon/Cloud NativeCon in Austin TX

Online Requirements
-------------------

KRIB uses community `kubeadm` for installation.  That process relies on internet connectivity to download containers and other components.

Immutable -vs- Local Install Mode
---------------------------------

The two primary deployment patterns that the Digital Rebar KRIB content pack supports are:

#. Live Boot (immutable infrastructure pattern - references [#]_ [#]_)
#. Local Install (standard install-to-disk pattern)

The *Live Boot* mode uses an in-memory Linux image based on the Digital Rebar Sledgehammer (CentOS based) image.  After each reboot of the Machine, the node is reloaded with the in-memory live boot image.  This enforces the concept of *immutable infrastructure* - every time a node is booted, deployed, or needs updating, simply reload the latest Live Boot image with appropriate fixes, patches, enhancements, etc.

The Local Install mode mimics the traditional "install-to-my-disk" method that most people are familiar with.

KRIB Basics
-----------

KRIB is a Content Pack addition to Digital Rebar Provision.  It uses the :ref:`rs_cluster_pattern` which provides atomic guarantees.  This allows for Kubernetes master(s) to be dynamically elected, forcing all other nodes to wait until the kubeadm on the elected master to generate an installation token for the rest of the nodes.  Once the Kubernetes master is bootstrapped, the Digital Rebar system facilitates the security token hand-off to rest of the cluster so they can join without any operator intervention.

Elected -vs- Specified Master
-----------------------------

By default, the KRIB process will dynamically elect a Master for the Kubernetes cluster.  This masters simply win the *race-to-master* election process and the rest of the cluster will coalesce around the elected master.

If you wish to specify a specific machines to be the designated masters, you can do so by setting a *Param* in the cluster *Profile* to the specific *Machine* that will be come the master.  To do so, set the ``krib/cluster-masters``  *Param* to a JSON structure with the Name, UUID and IP of the machines to become masters.  You may add this *Param* to the *Profile* in the below specifications, as follows:

  ::

    # JSON reference to add to the Profile Params section
    "krib/cluster-masters": [{"Name":"<NAME>", "Uuid":"<UUID>", "Address": "<ADDRESS>"}]

    # or drpcli command line option
    drpcli profiles set my-k8s-cluster param krib/cluster-master to <JSON>

The Kubernetes Master will be built on this Machine specified by the *<UUID>* value.

.. note:: This *MUST* be in the cluster profile because all machines in the cluster must be able to see this parameter.

Install KRIB
------------

KRIB is a Content Pack and is installed in the standard method as any other Contents.   We need the ``krib.json`` content pack to fully support KRIB and install the helper utility contents for stage changes.


CLI Install
===========


Using the Command Line (`drpcli`) utility configured to your endpoint, use this process:

  ::

  	# Get code
  	git clone https://github.com/digitalrebar/provision-content
  	cd krib

    # KRIB content install
    drpcli contents bundle krib.yaml
    drpcli contents upload krib.yaml

UX Install
==========

In the UX, follow this process:

#. Open your DRP Endpoint: (eg. https://127.0.0.1:8092/ )
#. Authenticate to your Endpoint
#. Login with your ```RackN Portal Login``` account (upper right)
#. Go to the left panel "Content Packages" menu
#. Select ``Kubernetes (KRIB: Kubernetes Rebar Immutable Bootstrapping)`` from the right side panel (you may need to select *Browser for more Content* or use the *Catalog* button)
#. Select the *Transfer* button for both content packs to add the content to your local Digital Rebar endpoint


Configuring KRIB
----------------

The basic outline for configuring KRIB follows the below steps:

#. create a *Profile* to hold the *Params* for the KRIB configuration (you can also clone the ``krib-example`` profile)
#. add a *Param* of name ``krib/cluster-profile`` to the *Profile* you created
#. add a *Param* of name ``etcd/cluster-profile`` to the *Profile* you created
#. apply the Profile to the Machines you are going to add to the KRIB cluster
#. change the Workflow on the Machines to ``krib-live-cluster`` for memory booting or ``krib-install-cluster`` to install to Centos.  You may clone these reference workflows to build custom actions.
#. installation will start as soon as the Workflow has been set.

There are many configuration options available, review the ``krib/*`` and ``etcd/*`` parameters to learn more.  

Configure with Terraform
========================

Please review the ``intergrations/krib`` for example Terraform plans.

Configure with the CLI
======================

The configuration of the Cluster includes several reference *Workflow* that can be used for installation.  Depending on which Workflow you use, will determine if the cluster is built via install-to-local-disk or via an immutable pattern (live boot in-memory boot process).   Outside of the Workflow differences, all remaining configuration elements are the same.

You must writeable create a *Profile* from YAML (or JSON if you prefer) with the Params stagemap and param required information. Modify the *Name* or other fields as appropriate - be sure you rename all subsequent fields appropriately.

  ::

    echo '
    ---
    Name: "my-k8s-cluster"
    Description: "Kubernetes install-to-local-disk"
    Params:
      krib/cluster-profile: "my-k8s-cluster"
      etcd/cluster-profile: "my-k8s-cluster"
    Meta:
      color: "purple"
      icon: "ship"
      title: "My Installed Kubernetes Cluster"
    ' > /tmp/krib-config.yaml

    drpcli profiles create - < /tmp/krib-config.yaml

.. note:: The following commands should be applied to all of the Machines you wish to enroll in your KRIB cluster.  Each Machine needs to be referenced by the Digital Rebar Machine UUID.  This example shows how to collect the UUIDs, then you will need to assign them to the ``UUIDS`` variable.  We re-use this variable throughout the below documentation within the shell function named *my_machines*.  We also show the correct ``drpcli`` command that should be run for you by the helper function, for your reference.

Create our helper shell function *my_machines*
  ::

    function my_machines() { for U in $UUIDS; do set -x; drpcli machines $1 $U $2; set +x; done; }

List your Machines to determine which to apply the Profile to
  ::

    drpcli machines list | jq -r '.[] | "\(.Name) : \(.Uuid)"'

IF YOU WANT to make ALL Machines in your endpoint use KRIB, do:
  ::

    export UUIDS=`drpcli machines list | jq -r '.[].Uuid'`

Otherwise - individually add them to the *UUIDS* variable, like:
  ::

    export UUIDS="UUID_1 UUID_2 ... UUID_n"

Add the Profile to your machines that will be enrolled in the cluster

  ::

    my_machines addprofile my-k8s-cluster

    # runs example command:
    # drpcli machines addprofile <UUID> my-k8s-cluster

Change stage on the Machines to initiate the Workflow transition.  YOU MUST select the correct stage, dependent on your install type (Immutable/Live Boot mode or install-to-local-disk mode).  For Live Boot mode, select the stage ``ssh-access`` and for the install-to-local-disk mode select the stage ``centos-7-install``.

  ::

    # for Live Boot/Immutable Kubernetes mode
    my_machines workflow krib-live-cluster

    # for intall-to-local-disk mode:
    my_machines workflow krib-install-cluster

    # runs example command:
    # drpcli machines workflow <UUID> krib-live-cluster
    # or
    # drpcli machines workflow <UUID> krib-install-cluster

Configure with the UX
=====================

The below example outlines the process for the UX.

RackN assumes the use of CentOS 7 BootEnv during this process.  However, it should theoretically work on most of the BootEnvs.  We have not tested it, and your mileage will absolutely vary...

1. create a *Profile* for the Kubernetes Cluster (e.g. ``my-k8s-cluster``) or clone the ``krib-example`` profile.
2. add a *Param* to that *Profile*: ``krib/cluster-profile`` = ``my-k8s-cluster``
2. add a *Param* to that *Profile*: ``etcd/cluster-profile`` = ``my-k8s-cluster``
3. Add the *Profile* (eg ``my-k8s-cluster``) to all the machines you want in the cluster.
4. Change workflow on all the machines to ``krib-install-cluster`` for install-to-local-disk, or to ``krib-live-cluster`` for the Live Boot/Immutable Kubernetes mode

Then wait for them to complete.  You can watch the Stage transitions via the Bulk Actions panel (which requires RackN Portal authentication to view).

.. note:: The reason the *Immutable Kubernetes/Live Boot* mode does not need a reboot is because they are already running *Sledgehammer* and will start installing upon the stage change.

Operating KRIB
--------------

Who is my Master?
=================

If you have not specified who the Kubernetes Master should be; and the master was chosen by election - you will need to determine which Machine is the cluster Master.
  ::

    # returns the Kubernetes cluster Machine UUID
    drpcli profiles show my-k8s-cluster | jq -r '.Params."krib/cluster-masters"'

Use ``kubectl`` - on Master
===========================

You can log in to the Master node as identified above, and execute ``kubectl`` commands as follows:
  ::

      export KUBECONFIG=/etc/kubernetes/admin.conf
      kubectl get nodes


Use ``kubectl`` - from anywhere
===============================

Once the Kubernetes cluster build has been completed, you may use the ``kubectl`` command to both verify and manage the cluster.  You will need to download the *conf* file with the appropriate tokens and information to connect to and authenticate your ``kubectl`` connections. Below is an example of doing this:
  ::

    # get the Admin configuration and tokens
    drpcli profiles get my-k8s-cluster param krib/cluster-admin-conf > admin.conf

    # set our KUBECONFIG variable and get nodes information
    export KUBECONFIG=`pwd`/admin.conf
    kubectl get nodes

Ingress/Egress Traffic and Dashboard Access
===========================================

The Kubernetes dashboard is enabled within a default KRIB built cluster.  However no Ingress traffic rules are set up.  As such, you must access services from external connections by making changes to Kubernetes, or via the :ref:`rs_k8s_proxy`.

These are all issues relating to managing, operating, and running a Kubernetes cluster, and not restrictions that are imposed by Digital Rebar Provision.  Please see the appropriate Kubernetes documentation on questions regarding operating, running, and administering Kubernetes (https://kubernetes.io/docs/home/).

.. _rs_k8s_proxy:

Kubernetes Dashboard via Proxy
==============================

Once you have obtained the ``admin.conf`` configuration file and security tokens, you may use ``kubectl`` in Proxy mode to the Master.  Simply open a separate terminal/console session to dedicate to the Proxy connection, and do:
  ::

    kubectl proxy

Now, in a local web browser (on the same machine you executed the Proxy command) open the following URL:

    https://127.0.0.1:8001/ui


Multiple Clusters
-----------------

It is absolutely possible to build multiple Kubernetes KRIB clusters with this process.  The only difference is each cluster should have a unique name and profile assigned to it.  A given Machine may only participate in a single Kubernetes cluster type at any one time.  You can install and operate both Live Boot/Immutable with install-to-disk cluster types in the same DRP Endpoint.


Footnotes
---------

.. [#] Immutable Infrastructure Reference: `Making Server Deployment 10x Faster – the ROI on Immutable Infrastructure <https://www.rackn.com/2017/10/11/making-server-deployment-10x-faster-roi-immutable-infrastructure/>`_

.. [#] Immutable Infrastructure Reference: `Go CI/CD and Immutable Infrastructure for Edge Computing Management <https://www.rackn.com/2017/09/15/go-cicd-immutable-infrastructure-edge-computing-management/>`_