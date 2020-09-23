
.. _rs_terraform_provider:

Terraform Provider for Digital Rebar v4.4+
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The following instructions describe how to use the Terraform Provider for
Digital Rebar v4.4+.  These changes only apply to versions that include
the embedded pooling feature.  

*WARNING: Previous generations of the provider should not be used.*

Provider source: https://github.com/rackn/terraform-provider-drp

Provider Goals & Design Considerations
--------------------------------------

The goal for Terraform Provider for Digital Rebar v4.4+ is to replicate cloud-like
behavior patterns common from cloud service providers.  Specifically, the "give me a machine(s)"
and "destroy my machine(s)" behavor.  This implementation focuses on enabling that behavior
with pools to enaure Digital Rebar operators maintain control and visibility at all times.

The Terraform Provider is _not_ to be used for general Digital Rebar configuration operations.
For this reason, only minimal API elements are exposed and there is no expectation that
Terraform will manage or maintain any Digital Rebar state information.

Reliance on Workflows set in Pool
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

To ensure consistent operational control, the Terraform Provider is _not_ allowed to set
the Machine Workflow on allocate or release.  Users of the Terraform Provider who want to
attach Workflows to Allocation are required to create and manage Pools in Digital Rebar.

Blocking
~~~~~~~~

The Terrform Provider uses blocking (`wait for`) allocate requests when requesting machines
from a pool.  That means that the machine is not considered available for Terraform until
the pool allocation Workflow completes.  This is consistent with providers written for clouds
that return only when the machine has been provisioned.

Filters
~~~~~~~

Terraform users will be able to pass-through Digital Rebar filters during allocation
requests.  This may require a degree of operational knowledge of the system to identify
filterable values.

Configuration Points
~~~~~~~~~~~~~~~~~~~~

Terraform users will be able to make the following changes to a machine during allocation
(and will be automatically removed on release):

  * set a param (can be user defined) value on the machine
  * set a profile on a machine (must be pre-defined)
  * set an SSH key on the machine (helper that just sets access-keys Param)

The following items are _not_ settable, but may be enabled in the future:

  * Machine.Description (could cause operational confusion)
  * Machine.Name (could cause issues with DNS)

Additional Resource Types
~~~~~~~~~~~~~~~~~~~~~~~~~

At this time, the Provider does not additional resource types.  Pool, Param and Profile could
be added to help users identify available read-only resources.

Terraform Formatting Rules
~~~~~~~~~~~~~~~~~~~~~~~~~~

The Terraform Provider complies with the Hashicorp syntax formatting requirements.  For this
reason, capitalization within the Provider will not match Digital Rebar capitalization.

Return Values
~~~~~~~~~~~~~

To reduce the potential for Machine state confusion, only a limited amount of information about
the allocated machine(s) are returned to Terraform by the Provider.  These include:

  * Machine.Name (translates to `name`)
  * Machine.Address (translatest to `address`)
  * Machine.Uuid (translates to `id`)
  * Machine.PoolStatus (translatest to `status`)


Using the Terraform Provider for Digital Rebar
------~---------------------------------------

Prereqs
~~~~~~~

Before using the Terraform Provider, you must have a Digital Rebar Provision (DRP) server version v4.4 or later installed with at least one machine created. These machines could be VMs, Packet servers, physical servers or even empty DRP machine models.

Terraform v1.13+ must be installed on the system.

3rd Party Provider Stanza
~~~~~~~~~~~~~~~~~~~~~~~~~

The DRP provider is maintained by RackN under https://extras.rackn.io and not available via the Terraform community repositories (at this time); consequently, users _must_ include the `drp` stanza in the `terraform.required_providers` block.  By including the reference, Terraform will be able to automatically download the provider from the RackN managed repository.

The `required_providers` block is as follows:

  ::

		terraform {
		  required_version = ">= 0.13.0"
		  required_providers {
		    drp = {
		      versions = ["2.0.0"]
		      source = "extras.rackn.io/rackn/drp"
		    }
		  }
		}

DRP Access via Environment Variables (provider "drp")
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

If set, the provider will use the standard DRP environment variables for access: 

  * RS_ENDPOINT for the endpoint [default "https://127.0.0.1:8092"]
  * RS_KEY for authentication [default "rocketskates:r0cketsk8ts"]

You can also define specialized behavior for the Ansible inventory
 
  * RS_TOKEN if you have defined an access token

.. note:: These values can also be set in the provider block.

Basic Resource Allocation/Release (resource "drp_machine")
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

When Terraform is applied, the "drp_machine" resource block will call the pool *allocate* API as described in :ref:`rs_pooling_ops`.  The *release* action is called when the Terraform destroys the machine.

The design attempts to be "cloud-like" and assumes that Digital Rebar operators will control the allocation and release actions behind these requests using Pool definitions.

It is important to understand that the Terraform provider is _not_ creating or destroying machines in this process.  It is simply assigning machines from a pool to the Terraform user during the apply operation.  For that reason, the provider exposes relatively few operational controls to Terraform.

.. note:: Operators can manually operate pools using `drpcli pools manage ...`


Enhanced Allocation Options
~~~~~~~~~~~~~~~~~~~~~~~~~~~

During allocation, Terraform users can include several aspects of the machines being allocated.

The primary choice for users is to select the "pool" from which to reserve machines.  Digital Rebar automatically enables the "default" pool so no changes are required; however, operators may wish to provide more options by creating pools.  The Terraform provider does not offer a way to create or change pools.

The following items may be used to influence allocations:

* add_profiles: list of profiles to add to the allocated machine (profile must exist)
* add_parameters: list of parameters (with values) to add to the allocated machine
* authorized_keys: list of ssh keys to add to the allocated machine (set via access-keys Param)
* filters: list of filter instructions used to further select machines from the pool.  See :ref:`rs_api_filters`.

The provider does not provide any additional options for destroy, but it will unwind all the choices made during allocation.


Output Values
~~~~~~~~~~~~~

After allocation, the provider sets key values for provider users.  These include:

* machine_id: maps to Machine.Uuid
* machine_ip: maps to Machine.Address
* machine_name: maps to Machine.Name
* status: maps to Machine.PoolStatus


Example Terraform Plan
----------------------

The following example plan represents all the available options in the comments.

  ::

		terraform {
		  required_version = ">= 0.13.0"
		  required_providers {
		    drp = {
		      versions = ["2.0.0"]
		      source = "extras.rackn.io/rackn/drp"
		    }
		  }
		}

	  provider "drp" {
	    username = "rocketskates"
	    password = "r0cketsk8ts"
	    endpoint = "https://127.0.0.1:8092"
	    # token  = will read from RS_TOKEN if set
	    # key    = will read from RS_KEY if set
	  }

	  resource "drp_machine" "one_random_node" {

	    # Required values
	    # there are none!

	    # Settable values
	    # pool = name of an existing DRP pool (defaults to "default")
	    # timeout = time string for max wait time (default to 5m)
	    # 
	    # List of public SSH keys to be installed (written as Param.access-keys)
	    # authorized_keys = ["ssh key"]
	    # 
	    # List of profiles to apply to node (must already exist)
	    # add_profiles = ["mandy", "clause"]
	    #
	    # list of parameters to set with their string value forms
	    # add_parameters = ["param1: value1", "param2: value2"]
	    #
	    # list of filters to reduce the nodes to draw from.
	    # follows the Digital Rebar CLI command line pattern
	    # filters = ["filter1=value1","filter2=value2"]
	    #
	    # Returned values
	    # name = machine name
	    # address = machine address
	    # status = machine status (typically "InUse")
	  }

	  output "machine_ip" {
	    value       = drp_machine.one_random_node.address
	    description = "Machine.Address (the Machine's primary IP)"
	  }

	  output "machine_id" {
	    value       = drp_machine.one_random_node.id
	    description = "Machine.Uuid"
	  }

	  output "machine_name" {
	    value       = drp_machine.one_random_node.name
	    description = "Machine.Name"
	  }