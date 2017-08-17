Kubernetes via Kubespray Ansible
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The following instructions help with the installation of Kubernetes on network provisioned servers using Digital Rebar Provision (https://github.com/digitalrebar/provision) and the Kubernetes Kubespray Ansible playbooks (https://github.com/kubernetes-incubator/kubespray).

Prereqs
-------

Before starting this process, a Digital Rebar Provision (DRP) server is required, along with the ability to provision machines.  These machines could be VMs, Packet servers or physical servers in a data center.  DRP CLI and Ansible must also be installed on the system.

You can run the deployment on a single node; however, three nodes is the recommended minimum (one masters, three nodes).

Root ssh access to the systems is required for the script to work.  At this time, testing is on Centos 7 only using root as the login.  This documentation assumes provisioning has completed and the machines are ready for installation - there is no workflow automation to move from discovery or sledgehammer to the target o/s documented here.

Digital Rebar Provision Configuration
-------------------------------------

Create the `ansible-chidren` parameter with an array schema.  This will be used to create children groups in Ansible.

  ::

    ./drpcli params create '{"name":"ansible-children", "schema":{"type":"array"}}'

Create profiles in DRP to become matching groups in Ansible.  Profiles are required for `kube-master`, `kube-node`, `etcd`, and `k8s-cluster`

  ::

    ./drpcli profiles create '{"name":"kube-master"}'
    ./drpcli profiles create '{"name":"kube-node"}'
    ./drpcli profiles create '{"name":"etcd"}'
    ./drpcli profiles create '{"name":"k8s-cluster"}'

Next, set the `ansible-children` value for the `k8s-cluster` profile.

  ::

    ./drpcli profiles params k8s-cluster '{"ansible-children":["kube-node","kube-master"]}'

Finally, assign the machines in the system to the desired profiles using `./drpcli machines addprofile [machine-uuid] kube-master` as a reference.  There must be at least one `kube-master`, `kube-node`, and `etcd`.  It is unnecessary to assign any machines to `k8s-cluster`.

Ansible Dynamic Inventory from Digital Rebar Provision
------------------------------------------------------

Be certain to export the `RS_ENDPOINT` and `RS_KEY` to match the DRP endpoint information because the DRP dynamic Ansible inventory script relies on these values being set.

Download the inventory script to the local system to a convenient location and make it executable.  You can test the script by simply running it.  The script will output JSON in the required Ansible format.

  ::
    curl -o https://raw.githubusercontent.com/digitalrebar/provision/master/integrations/ansible/inventory.py -o inventory.py
    chmod +x inventory.py
    ./inventory.py | jq

In order to test the Ansible integration, use the ping command.  If everything is working, all the machines in the system should receive and respond to the ping command. 

    ansible all -i inventory.py -m ping



Kubernetes Kubespray Playbook
-----------------------------

To install Kubernetes, checkout the Kubespray playboot from https://github.com/kubernetes-incubator/kubespray using git clone.

  ::

    git clone https://github.com/kubernetes-incubator/kubespray

it is important to review the Kubespray documentation and make any of the neccessary changes to the environment.  For a simple test, run the playbook without any modifications using the following command.

  ::

    ansible-playbook -i inventory.py cluster.yml

Wait until Kubernetes complete and log into the master using `https://[kube-master]:6443`

Summary
-------

Now that these steps are completed, the Digital Rebar Provision dynamic inventory script can be used in any number of ways. 
