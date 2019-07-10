
.. _rs_ansible:

Dynamic Ansible Inventory using Profiles
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The following instructions show how to map Ansible Playbooks via
Digital Rebar with no Ansible specific configuration required.

The instructions are generic and could be adapted to run on any Ansible run.

Prereqs
-------

Before starting this process, a Digital Rebar Provision (DRP) server is required, along with the ability to provision machines.  These machines could be VMs, Packet servers or physical servers in a data center.  DRPCLI and Ansible must also be installed on the system.

Root ssh access to the systems is required for the script to work.  Make sure that the correct SSH keys have been installed on the target systems.  Review :ref:`rs_add_ssh` for details.

At this time, testing is on Centos 7 only using root as the login.  This documentation assumes provisioning has completed and the machines are ready for installation - there is no workflow automation to move from discovery or sledgehammer to the target o/s documented here.

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

.. note:: You may want to set `export ANSIBLE_HOST_KEY_CHECKING=False` to bypass the SSH key validation

Use non-root Login
------------------

By default, the internal `ansible_user` will be set to `root`.  You can override this by setting the `RS_ANSIBLE_USER` value.

  ::

    RS_ANSIBLE_USER="username"


Use Alternative Host Address
----------------------------

By default, the internal Machine.Address value is used.  If this address does not work (e.g. Cloud IPs) then you can specify a parameter as the source of the IP address using the `RS_HOST_ADDRESS` value.

  ::

    export RS_HOST_ADDRESS="cloud/public-ipv4"

Summary
-------

Now that these steps are completed, the Digital Rebar Provision dynamic inventory script can be used in any number of ways.
