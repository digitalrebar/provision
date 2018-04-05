
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

Note: these instructions assume that you are using the RackN extension; however, it is not required for Terraform operation.

Install the `Terraform` content package from RackN.  This content provides the parameters (described above) and a stage (`terraform-ready`) to configure Terraform properties on machines.

Create a workflow that includes the `terraform-ready` stage to set the two parameters.  Once these values are set they do not have to be run again, but there is no harm in leaving the stage in place.

Before testing Terraform, it is good practice to make sure you can manually cycle machines by changing their workflow (or change stage and reboot).  These changes should result in a machine provisioning cycle.

Make sure that you know your endpoint URI, user and password.

Installing The Terraform Digital Rebar Provision Provider
---------------------------------------------------------

Please install Terraform on your system: https://www.terraform.io/intro/getting-started/install.html

Download (or build) the DRP Provider from https://github.com/rackn/terraform-provider-drp/releases.

Initialize the plugin using `terraform init`

Create a Plan using the DRP Provider Resources
----------------------------------------------

The DRP Provider exposes all the objects in Digital Rebar Project so you can script against any operation.

The critical block is the provider block which identifies the provider and login information (shown here with default values):

  ::

	  provider "drp" {
	    api_user     = "rocketskates"
	    api_password = "r0cketsk8ts"
	    api_url      = "https://127.0.0.1:8092"
	  }

Once you have the provider block, you can name resource blocks using the normal object key values.  For example, a machine resource looks like this:

  ::

	resource "drp_machine" "one_random_node" {
	  count = 1
	  bootenv     = "ubuntu-16.04-install"
	  description = "updated description"
	  name        = "greg2"
	  userdata    = "yaml cloudinit file"
	}

There are many options to set including filters, parameters and profiles.  For a full example, please look at https://github.com/rackn/terraform-provider-drp/blob/master/test.tf.example

Running Terraform
-----------------

Just use `terraform apply` and `terraform destroy` and as normal!

Extending the Features
----------------------

Using the `terraform/owner` parameter helps administrators track who is using which machines.  You may also choose to create multiple DRP users to help track activity.

It is highly recommended that you include decommissioning steps (disk scrub, bios reset, etc) and additional burn-in to validate systems during the recovery cycle.

Using IPMI to reset machines is a safer bet than relying on the DRP runner to soft reboot systems.  If you want to make sure that you have a consistent recovery process, IPMI is highly recommended.

To improve delivery time:

1. Keep the machines running
2. Use image based provisioning instead of netboot.

If you are relying on the DRP Running workflow to start allocation and recovery, make sure that you have your tokens set to never expire!

Summary
-------

Now that these steps are completed, the Digital Rebar Provision Terraform Provider will integrate like any cloud provider.
