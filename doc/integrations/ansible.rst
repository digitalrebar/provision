
.. _rs_ansible:

Dynamic Ansible Inventory using Profiles
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The following instructions show how to map Ansible Playbooks via
Digital Rebar with no Ansible specific configuration required.

The instructions are generic and could be adapted to run on any Ansible run.

Prereqs
-------

Before starting this process, a Digital Rebar Provision (DRP) server is required, along with the ability to provision machines.  These machines could be VMs, Packet servers or physical servers in a data center.  DRPCLI and Ansible must also be installed on the system.

Root ssh access to the systems is required for the script to work.  At this time, testing is on Centos 7 only using root as the login.  This documentation assumes provisioning has completed and the machines are ready for installation - there is no workflow automation to move from discovery or sledgehammer to the target o/s documented here.

Digital Rebar Provision Ansible Configuration
---------------------------------------------

The Integrations Ansible drpmachines.py script can be used to create a dynamic inventory from a Digital Rebar Endpoint.

The hosts inventory list defaults to all machines or can be restricted by setting a `ansible=[value]` Param.

Group membership managed by directly mapping Profiles into Ansible Groups.  If machines have been assigned a profile then it will be included in the Group hosts list.  Params in Profiles will be presents as Group vars.  *There is no additional mapping required.*

Note: Ansible dynamic inventory requires JSON output instead of YAML and the format is slightly different.


Optionally, parent groups can be configured by adding the `ansible/children` Param to any profile.  The Param is a simple list of groups to be listed in the Groups children.


Ansible Dynamic Inventory from Digital Rebar Provision
------------------------------------------------------

Be certain to export the `RS_ENDPOINT` and `RS_KEY` to match the DRP endpoint information because the DRP dynamic Ansible inventory script relies on these values being set.

Optionally, you may limit the machines using the `ansible=[key]` Param by set `RS_ANSIBLE` to match the [key] value assigned.  The default is to ignore this value and use all machines.

For this example, please ensure that *jq* is installed.

Download the `drpmachines.py` inventory script to the local system to a convenient location and make it executable.  You can test the script by simply running it.  The script will output JSON in the required Ansible format.

  ::

    curl -s https://raw.githubusercontent.com/digitalrebar/provision/master/integrations/ansible/drpmachines.py -o drpmachines.py
    chmod +x drpmachines.py
    ./drpmachines.py | jq

In order to test the Ansible integration, use the ping command.  If everything is working, all the machines in the system should receive and respond to the ping command.

  ::

    ansible all -i drpmachines.py -m ping



LEGACY Dynamic Ansible Inventory (w/ Kubernetes via Kubespray Ansible)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

!! THIS APPROACH HAS BEEN DEPRECATED !!

The following instructions help with Ansible Playbooks via
Digital Rebar.  The instructions are generic and could be
adapted to run on any Ansible run.

To make the documentation more specific, they use the installation
of Kubernetes on network provisioned servers using Digital Rebar Provision (https://github.com/digitalrebar/provision) and the Kubernetes Kubespray Ansible playbooks (https://github.com/kubernetes-incubator/kubespray).

Video Demo: https://youtu.be/b5himGQ1Zew.

Prereqs
-------

Before starting this process, a Digital Rebar Provision (DRP) server is required, along with the ability to provision machines.  These machines could be VMs, Packet servers or physical servers in a data center.  DRP CLI and Ansible must also be installed on the system.

You can run the deployment on a single node; however, three nodes is the recommended minimum (one masters, three nodes).

Root ssh access to the systems is required for the script to work.  At this time, testing is on Centos 7 only using root as the login.  This documentation assumes provisioning has completed and the machines are ready for installation - there is no workflow automation to move from discovery or sledgehammer to the target o/s documented here.

Digital Rebar Provision Ansible Configuration
---------------------------------------------

The Integrations Ansible Inventory.py script can be used to create a dynamic inventory from a Digital Rebar Endpoint.  The inventory list is managed by parameters on profiles.  You may have multiple independent inventory profiles.

Note: The RackN Kubespray Package containes a preconfigured
profile that can be cloned for this use.  There is also a specialized Ansible screen in the RackN UX for creating the member list.

You must create a ansible/groups profile with the following required params:

  ::

    "ansible/groups": [
      "etcd",
      "kube-master",
      "kube-node"
    ],

You must supply a ansible/groups-members list to map hosts into groups.  The RackN UX will build this for you using a graphic selection matrix.

  ::

    "ansible/groups-members": {
      "etcd": [
        "test1.unspecified.domain.local"
      ],
      "kube-master": [
        "test1.unspecified.domain.local"
      ],
      "kube-node": [
        "test2.unspecified.domain.local"
      ]
    },

Optionally, ansible/groupvars maps profile and their params into groups:vars.

  ::

    "ansible/groupvars": {
      "etcd": "kube-etcd",
      "kube-master": "kube-master",
      "kube-node": "kube-node"
    },

Optionally, ansible/hostvars maps variables into the hostvars variable lists.

  ::

    "ansible/hostvars": {
      "ansible_user": "root"
    },

Optionally, ansible/parent-groups can be used to create children groups.

  ::

    "ansible/parent-groups": {
      "kube-cluster": [
        "kube-master",
        "kube-node"
      ]
    }


Ansible Dynamic Inventory from Digital Rebar Provision
------------------------------------------------------

Be certain to export the `RS_ENDPOINT` and `RS_KEY` to match the DRP endpoint information because the DRP dynamic Ansible inventory script relies on these values being set.

You will need to set `RS_PROFILE` to match the profile that you are using as a target.  The default is `mycluster`

For this example, please ensure that *jq* is installed.

Download the inventory script to the local system to a convenient location and make it executable.  You can test the script by simply running it.  The script will output JSON in the required Ansible format.

  ::

    curl -s https://raw.githubusercontent.com/digitalrebar/provision/master/integrations/ansible/inventory.py -o inventory.py
    chmod +x inventory.py
    RS_PROFILE=mycluster ./inventory.py | jq

In order to test the Ansible integration, use the ping command.  If everything is working, all the machines in the system should receive and respond to the ping command.

  ::

    RS_PROFILE=mycluster ansible all -i inventory.py -m ping



Kubernetes Kubespray Playbook
-----------------------------

To install Kubernetes, checkout the Kubespray playboot from https://github.com/kubernetes-incubator/kubespray using git clone.

  ::

    git clone https://github.com/kubernetes-incubator/kubespray

it is important to review the Kubespray documentation and make any of the neccessary changes to the environment.  For a simple test, run the playbook without any modifications using the following command.

  ::

    RS_PROFILE=mycluster ansible-playbook -i inventory.py cluster.yml

Wait until Kubernetes complete and log into the master using `https://[kube-master]:6443`

Summary
-------

Now that these steps are completed, the Digital Rebar Provision dynamic inventory script can be used in any number of ways.
