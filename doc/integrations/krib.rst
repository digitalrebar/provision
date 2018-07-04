
.. _rs_krib:

KRIB (Kubernetes Rebar Integrated Bootstrapping)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

This document provides information on how to use the Digital Rebar *KRIB* content add-on.  Use of this content will enable the operator to install Kubernetes in either a Live Boot (immutable infrastructure pattern) mode, or via installed to local hard disk OS mode.

KRIB uses the `kubeadm <https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/>`_ cluster deployment methodology coupled with Digital Rebar enhancements to help proctor the cluster Master election and secrets management.  With this content pack, you can install Kubernetes in a zero-touch manner.

.. note:: The `kubeadm` installation and configuration method is very new to the Kubernetes community.  It does not support High Availability of the master, it is still in it's infancy in regards to external networking (SDN, etc) integrations, etc.  The implementation provided via KRIB is intended to demonstrate a *Cluster Provisioning* workflow process.  We do not recommend using this Kubernetes configuration as a production solution.   It could become the base for a strong production deployment, but at this time, does not support production grade features.

  For production installations, you may be better served by using the Ansible ``kubespray`` playbooks for deployment of the cluster.  See the :ref:`rs_ansible` documentation for further details.

This document assumes you have your Digital Rebar Provisioning endpoint fully configured, tested, and working.  We assume that you are able to properly provision Machines in your environment as a base level requirement for use of the KRIB content add-on use.

KRIB Video References
---------------------

The following videos have been produced or presented by RackN related to the Digital Rebar KRIB solution.

* `KRIB Zero Config Kubernetes Cluster with RackN <https://youtu.be/OMm6Oz1NF6I>`_ on YouTube - RackN recorded presentation with just KRIB deployment information
* `KubeCon: Zero Configuration Pattern on Bare Metal <https://youtu.be/Psm9aOWzfWk>`_ on YouTube - RackN presentation at 2017 KubeCon/Cloud NativeCon in Austin TX

Immutable -vs- Local Install Mode
---------------------------------

The two primary deployment patterns that the Digital Rebar KRIB content pack supports are:

#. Live Boot (immutable infrastructure pattern - references [#]_ [#]_)
#. Local Install (standard install-to-disk pattern)

The *Live Boot* mode uses an in-memory Linux image based on the Digital Rebar Sledgehammer (CentOS based) image.  After each reboot of the Machine, the node is reloaded with the in-memory live boot image.  This enforces the concept of *immutable infrastructure* - every time a node is booted, deployed, or needs updating, simply reload the latest Live Boot image with appropriate fixes, patches, enhancements, etc.

The Local Install mode mimics the traditional "install-to-my-disk" method that most people are familiar with.

Ready State Infrastructure
--------------------------

This process assumes that your machines have been discovered and and are in a "ready state" for the next provisioning activities.  Typically, this can be done by setting the ``global`` profile to move discovered machines to ``sledgehammer-wait`` state.  A ``change-stage/map`` similar to:
  ::

    discover --> packet-discover --> sledgehammer-wait

Will provide the appropriate waiting (ready-state) stage to put your Machines in to prior to advancing through the KRIB workflow below.



KRIB Basics
-----------

KRIB is essentially nothing more than a Content Pack addition to Digital Rebar, which works with the Open version of Digital Rebar.  It uses the *cluster bootstrapping* support built in to the Digital Rebar solution which provides atomic guarantees.  This allows for a Kubernetes Master to be dynamically elected, forcing all other nodes that fail the *race-to-master* election to wait until the elected master is completed and bootstrapped.  Once the Kubernetes Master is bootstrapped, the Digital Rebar system facilitates the security token hand-off to the Minions that join the cluster, to allow them to join with out any operator intervention.

Elected -vs- Specified Master
-----------------------------

By default the KRIB process will dynamically elect a Master for the Kubernetes cluster.  This master simply wins the *race-to-master* election process and the rest of the cluster will coalesce around the elected master.   There is no failover mechanisms or High Availability (as *kubeadm* doesn't yet support this pattern).

If you wish to specify a specific machine to be the designated Master, you can do so by setting a *Param* in the cluster *Profile* to the specific *Machine* that will be come the master.  To do so, set the ``krib/cluster-master``  *Param* to the UUID of the machine to become master.  You may add this *Param* to the *Profile* in the below specifications, as follows:

  ::

    # JSON reference to add to the Profile Params section
    "krib/cluster-master": "<UUID>"

    # or drpcli command line option
    drpcli profiles set my-k8s-cluster param krib/cluster-master to <UUID>

The Kubernetes Master will be built on this Machine specified by the *<UUID>* value.

.. note:: This *MUST* be in the cluster profile because all machines in the cluster must be able to see this parameter.

Install KRIB
------------

KRIB is a Content Pack and is installed in the standard method as any other Contents.   We need the ``krib.json`` content as well as the ``task-library.json`` content packs to fully support KRIB and install the helper utility contents for stage changes.


CLI Install
===========


Using the Command Line (`drpcli`) utility, use this process:
  ::

    # KRIB content install
    curl -s https://qww9e4paf1.execute-api.us-west-2.amazonaws.com/main/catalog/content/krib -o /tmp/krib.json
    drpcli contents create -< /tmp/krib.json

    # task-libary helper content install
    curl -s https://qww9e4paf1.execute-api.us-west-2.amazonaws.com/main/catalog/content/task-library -o /tmp/task-library.json
    drpcli contents create -< /tmp/task-library.json


UX Install
==========

In the UX, follow this process:

#. Open your DRP Endpoint: (eg. https://127.0.0.1:8092/ )
#. Authenticate to your Endpoint
#. Login with your ```RackN Portal Login``` account (upper right)
#. Go to the left panel "Content Packages" menu
#. Select ``Kubernetes (KRIB: Kubernetes Rebar Immutable Bootstrapping)`` from the right side panel (you may need to select *Browser for more Content* or use the *Catalog* button)
#. *also* select the ``task-library`` content
#. Select the *Transfer* button for both content packs to add the content to your local Digital Rebar endpoint


Configuring KRIB
----------------

The basic outline for configuring KRIB follows the below steps:

#. create a *Profile* to hold the *Params* for the KRIB configuration
#. add a *Param* of name ``krib/cluster-profile`` to the *Profile* you created
#. add a stagemap workflow to the *Profile* you created above, to move machines through the KRIB install process
#. apply the Profile to the Machines you are going to add to the KRIB cluster
#. change the Stage on the Machines to set the starting point of the workflow
#. reboot the Machines in the KRIB cluster to initiate the installation

Configure with the CLI
======================

The configuration of the Cluster includes a *Stagemap* - and depending on which stage map you use, will determine if the cluster is built via install-to-local-disk or via an immutable pattern (live boot in-memory boot process).   Outside of the stagemap differences, all remaining configuration elements are the same.

You must create a *Profile* from YAML (or JSON if you prefer) with the stagemap and param required information. Modify the *Name* or other fields as appropriate - be sure you rename all subsequent fields appropriately.  This example uses CentOS 7 as the BootEnv for the install-to-local-disk option.

Additionally - ensure you correctly modify the ``access-keys`` Param to inject your apprpriate SSH public key half or halves appropriately.

  ::

    echo '
    ---
    Name: "my-k8s-cluster"
    Description: "Kubernetes install-to-local-disk"
    Params:
      krib/cluster-profile: "my-k8s-cluster"
      change-stage/map:
        centos-7-install: runner-service:Success
        runner-service: finish-install:Stop
        finish-install: docker-install:Success
        docker-install: krib-install:Success
        krib-install: complete:Success
        discover: sledgehammer-wait:Success
      access-keys:
        user1: ssh <user_1_key> user@krib
        user2: ssh <user_2_key> user@krib
    Meta:
      color: "purple"
      icon: "ship"
      title: "My Installed Kubernetes Cluster"
    ' > /tmp/krib-config.yaml


For an Immutable Kubernetes cluster install, use the below *Profile* with the stagemap below.
  ::

    echo '
    ---
    Name: "my-k8s-cluster"
    Description: "Kubernetes Live Boot (immutable) cluster"
    Params:
      krib/cluster-profile: "my-k8s-cluster"
      change-stage/map:
        access-keys: mount-local-disks:Success
        mount-local-disks: docker-install:Success
        docker-install: krib-install:Success
        krib-install: sledgehammer-wait:Success
      access-keys:
        user1: ssh <user_1_key> user1@krib
        user2: ssh <user_2_key> user2@krib
    Meta:
      color: "orange"
      icon: "ship"
      title: "My Immutable Kubernetes Cluster"
    ' > /tmp/krib-config.yaml

.. note:: ONLY select one of the two above YAML profile options.

Apply/create the Profile
  ::

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
    my_machines stage ssh-access

    # for intall-to-local-disk mode:
    my_machines stage centos-7-install

    # runs example command:
    # drpcli machines stage <UUID> ssh-access
    # or
    # drpcli machines stage <UUID> centos-7-install

    # if fails, try below for each UUID - there is a potential "stage" change bug in CLI
    # drpcli machines update <UUID> '{ "Stage": "ssh-access" }'
    # or
    # drpcli machines update <UUID> '{ "Stage": "centos-7-install" }'


For the *install-to-local-disk* mode, you now need to reboot the Machines you modified above.

.. note:: You can do this through your own tooling or power control methods.  For example, via IPMI protocol, Console access and rebooting, physically power cycling the machine, or other methods.

  Digital Rebar Provision does support installing Plugin Providers that implement IPMI control (power on/off/reboot) actions.  Some of these are available for free as a Registered user, some of these are Paid pieces.   Please see your UX ``Contents`` menu for the status of each plugin provider.

If you are using the RackN `IPMI` plugin provider (free or paid piece), you can do this with the following commands:
  ::

    my_machines action powercycle

    # runs example command:
    # drpcli machines action <UUID> powercycle

Configure with the UX
=====================

The below example outlines the process for the UX.

RackN assumes the use of CentOS 7 BootEnv during this process.  However, it should theoretically work on most of the BootEnvs.  We have not tested it, and your mileage will absolutely vary...

1. create a *Profile* for the Kubernetes Cluster (e.g. ``my-k8s-cluster``)
2. add a *Param* to that *Profile*: ``krib/cluster-profile`` = ``my-k8s-cluster``
3. add the following workflow to the ``my-k8s-cluster`` *Profile*.

  for install-to-local-disk mode:

  a. ``centos-7-install -> runner-service:Success``
  b. ``runner-service -> finish-install:Stop``
  c. ``finish-install -> docker-install:Success``
  d. ``docker-install -> krib-install:Success``
  e. ``krib-install-> complete:Success``
  f. ``discover->sledgehammer-wait:Success``

  OR

  for Immutable Kubernetes/Live Boot mode:

  a. ``ssh-access`` -> ``mount-local-disks:Success``
  b. ``mount-local-disks`` -> ``docker-install:Success``
  c. ``docker-install`` -> ``krib-install:Success``
  d. ``krib-install`` -> ``sledgehammer-wait:Success``

  The last entry is to handle discovery if you reimage the servers.

4. Add the *Profile* (eg ``my-k8s-cluster``) to all the machines you want in the cluster.
5. Change stage on all the machines to ``centos-7-install`` for install-to-local-disk, or to ``ssh-access`` for the Live Boot/Immutable Kubernetes mode
6. Reboot all the machines in your cluster if you are using the *install-to-local-disk* mode.

.. note:: You can do this through your own tooling or power control methods.  For example, via IPMI protocol, Console access and rebooting, physically power cycling the machine, or other methods.

  Digital Rebar Provision does support installing Plugin Providers that implement IPMI control (power on/off/reboot) actions.  Some of these are available for free as a Registered user, some of these are Paid pieces.   Please see your UX ``Contents`` menu for the status of each plugin provider.

Then wait for them to complete.  You can watch the Stage transitions via the Bulk Actions panel (which requires RackN Portal authentication to view).

.. note:: The reason the *Immutable Kubernetes/Live Boot* mode does not need a reboot is because they are already running *Sledgehammer* and will start installing upon the stage change.

Operating KRIB
--------------

Who is my Master?
=================

If you have not specified who the Kubernetes Master should be; and the master was chosen by election - you will need to determine which Machine is the cluster Master.
  ::

    # returns the Kubernetes cluster Machine UUID
    drpcli profiles show my-k8s-cluster | jq -r '.Params."krib/cluster-master"'

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


