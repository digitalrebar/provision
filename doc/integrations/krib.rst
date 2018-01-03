
.. _rs_krib:

KRIB (Kubernetes Rebar Immutable Bootstrapping)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

This document provides information on how to use the Digital Rebar *KRIB* content add-on.  Use of this content will enable the operator to install Kubernetes in either a Live Boot (immutable infrastructure pattern) mode, or via installed to local hard disk OS mode.  

KRIB uses the `kubeadm <https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/>`_ cluster deployment methodology coupled with Digital Rebar enhancements to help proctor the cluster Master election and secrets management.  With this content pack, you can install Kubernetes in a zero-touch manner.  

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

KRIB Basics
-----------

KRIB is essentially nothing more than a Content Pack addition to Digital Rebar, which works with the Open version of Digital Rebar.  It uses the *cluster bootstrapping* support built in to the Digital Rebar solution which provides atomic guarantees.  This allows for a Kubernetes Master to be dynamically elected, forcing all other nodes that fail the *race-to-master* election to wait until the elected master is completed and bootstrapped.  Once the Kubernetes Master is bootstrapped, the Digital Rebar system facilitates the security token hand-off to the Minions that join the cluster, to allow them to join with out any operator intervention.  

Install KRIB
------------

KRIB is a Content Pack and is installed in the standard method as any other Contents.  


CLI Install
===========

Using the Command Line (`drpcli`) utility, use this process:
  ::

    curl -s https://qww9e4paf1.execute-api.us-west-2.amazonaws.com/main/catalog/content/krib -o /tmp/krib.json
    drpcli contents create -< /tmp/krib.json


UX Install
==========

In the UX, follow this process:

#. Open your DRP Endpoint: (eg. https://127.0.0.1:8092/ )
#. Authenticate to your Endpoint
#. Login with your ```RackN Portal Login``` account (upper right)
#. Go to the left panel "Content Packages" menu 
#. Select `Kubernetes (KRIB: Kubernetes Rebar Immutable Bootstrapping)` from the right side panel (you may need to select *Browser for more Content* or use the *Catalog* button)
#. Select the *Transfer* button to add the content to your local Digital Rebar endpoint


Configuring KRIB
----------------

The basic outline for configuring KRIB follows the below steps:

#. create a *Profile* to hold the *Params* for the KRIB configuration
#. add a *Param* of name `krib/cluster-profile` to the *Profile* you created
#. add a stagemap workflow to the *Profile* you created above, to move machines through the KRIB install process
#. apply the Profile to the Machines you are going to add to the KRIB cluster
#. change the Stage on the Machines to set the starting point of the workflow
#. reboot the Machines in the KRIB cluster to initiate the installation

Configure with the CLI
======================

Create the YAML for the Profile with stagemap and param required - modify the *Name* or other fields as appropriate - be sure you rename all subsequent fields appropriately.
  ::

    echo '
    ---
    Name: "my-k8s-cluster"
    Description: "My Kubernetes Cluster"
    Params:
      krib/cluster-profile: "my-k8s-cluster"
      change-stage/map:
        centos-7-install: runner-service:Success
        runner-service: finish-install:Stop
        finish-install: docker-install:Success
        docker-install: krib-install:Success
        krib-install: complete:Success
        discover: sledgehammer-wait:Success
    Meta:
      color: "purple"
      icon: "ship"
      title: "My Kubernetes Cluster"
    ' > /tmp/krib-config.yaml


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

Change stage on the Machines to initiate the Workflow transition
  ::

    my_machines stage centos-7-install

    # runs example command:
    # drpcli machines stage <UUID> centos-7-install

    # if fails, try below for each UUID - there is a potential "stage" change bug in CLI
    # drpcli machines update <UUID> '{ "Stage": "centos-7-install" }'


Now you need to reboot the Machines you modified above.  You can do this through your own tooling or power control methods.  If you are using the RackN `IPMI` plugin provider (paid piece), you can do this with the following commands:
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
3. Using workflow editor, add the following workflow to the ``my-k8s-cluster`` *Profile*.

  a. ``centos-7-install -> runner-service:Success``
  b. ``runner-service -> finish-install:Stop``
  c. ``finish-install -> docker-install:Success``
  d. ``docker-install -> krib-install:Success``
  e. ``krib-install-> complete:Success``
  f. ``discover->sledgehammer-wait:Success``

  The last entry is to handle discovery if you reimage the servers.

4. Add the *Profile* (eg ``my-k8s-cluster``) to all the machines you want in the cluster.
5. Change stage on all the machines to ``centos-7-install``
6. Reboot all the machines in your cluster.


Then wait for them to complete.  You can watch the Stage transitions via the Bulk Actions panel (which requires RackN Portal authentication to view).


Operating KRIB
--------------

This section is not yet complete.

Footnotes:
----------

.. [#] Immutable Infrastructure Reference: `Making Server Deployment 10x Faster – the ROI on Immutable Infrastructure <https://www.rackn.com/2017/10/11/making-server-deployment-10x-faster-roi-immutable-infrastructure/>`_

.. [#] Immutable Infrastructure Reference: `Go CI/CD and Immutable Infrastructure for Edge Computing Management <https://www.rackn.com/2017/09/15/go-cicd-immutable-infrastructure-edge-computing-management/>`_


