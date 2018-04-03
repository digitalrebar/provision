
.. _rs_terraform:

Terraform Provider for Digital Rebar Provision
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The following instructions to install the needed Digital Rebar Provision
integrations to support Terraform.  Terraform integrations require that
machines can reboot.  

The Terraform system is not synchronized with Digital Rebar, this documentation is only
intended to cover the provider.

Note: As of DRP v3.8 (workflows enabled), changes to bootenvs will automatically
cause reboots if the DRP runner/agent is allowed to keep running.  Otherwise, an IPMI plugin is required.
We recommend using both Workflow and IPMI based reboots to ensure that the systems return to a consistent state.

Source Code: https://github.com/rackn/terraform-provider-drp 

Video Demo: https://www.youtube.com/watch?v=6MLyUVgnVo4

License
-------

The Terraform Provider for Digital Rebar Provision is APLv2 licensed.  Advanced features or workflow may require capabilities that use different licenses.

Prereqs
-------

Before starting this process, a Digital Rebar Provision (DRP) server is required, along with the ability to provision machines.  These machines could be VMs, Packet servers or physical servers in a data center.

You must also have installed Terraform and secured the Digital Rebar Provision Terraform Provider.  You can build the provider from source or retrieve a compiled version from the Github releases area, https://github.com/rackn/terraform-provider-drp/releases.

Basic Operation
---------------

The DRP Terraform Provider uses a pair of Machine Parameters to create an inventory pool.  Only machines with these parameters will be available to the provider.

The `terraform/managed` parameter determines the basic inventory availability.  This flag must be set to true for Terraform to find machines.

The `terraform/allocated` parameter determines when machines have been assigned to a Terraform plan.  When true, the machine is now being managed by Terraform.  When false, the machine is available for allocation.

Using the RackN `terraform-ready` stage will automatically set these two parameters.

The Terraform Provider can read additional fields when requesting inventory.  In this way, users can request machines with specific characteristics.


DRP Machine Configuration
-------------------------

Install the `Terraform` content package from RackN.  This content provides the parameters (described above) and a stage (`terraform-ready`) to configure Terraform properties on machines.

Summary
-------

Now that these steps are completed, the Digital Rebar Provision Terraform Provider will integrate like any cloud provider.
