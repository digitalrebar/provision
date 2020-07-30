
.. _rs_terraform_provider:

Terraform Provider for Digital Rebar
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The following instructions describe how to use the Terraform Provider for
Digital Rebar v4.4+.  These changes only apply to versions that include
the embedded pooling feature.  Previous generations of the provider should
not be used.

Provider source: https://github.com/rackn/terraform-provider-drp

Prereqs
-------

Before using the Terraform Provider, you must have a Digital Rebar Provision (DRP) server version v4.4 or later installed with at least one machine created. These machines could be VMs, Packet servers, physical servers or even empty DRP machine models.

Terraform v1.12+ must be installed on the system.

DRP Access via Environment Variables (provider "drp")
-----------------------------------------------------

If set, the provider will use the standard DRP environment variables for access: 

  * RS_ENDPOINT for the endpoint [default "https://127.0.0.1:8092"]
  * RS_KEY for authentication [default "rocketskates:r0cketsk8ts"]

You can also define specialized behavior for the Ansible inventory
 
  * RS_TOKEN if you have defined an access token

.. note:: These values can also be set in the provider block.

Basic Resource Allocation/Release (resource "drp_machine")
----------------------------------------------------------

When Terraform is applied, the "drp_machine" resource block will call the pool *allocate* API as described in :ref:`rs_pooling_ops`.  The *release* action is called when the Terraform destroys the machine.

The design attempts to be "cloud-like" and assumes that Digital Rebar operators will control the allocation and release actions behind these requests using Pool definitions.

It is important to understand that the Terraform provider is _not_ creating or destroying machines in this process.  It is simply assigning machines from a pool to the Terraform user during the apply operation.  For that reason, the provider exposes relatively few operational controls to Terraform.

.. note:: Operators can manually operate pools using `drpcli pools manage ...`


Enhanced Allocation Options
---------------------------

During allocation, Terraform users can include several aspects of the machines being allocated.

The primary choice for users is to select the "pool" from which to reserve machines.  Digital Rebar automatically enables the "default" pool so no changes are required; however, operators may wish to provide more options by creating pools.  The Terraform provider does not offer a way to create or change pools.

The following items may be used to influence allocations:

* add_profiles: list of profiles to add to the allocated machine (profile must exist)
* add_parameters: list of parameters (with values) to add to the allocated machine
* authorized_keys: list of ssh keys to add to the allocated machine (set via access-keys Param)
* filters: list of filter instructions used to further select machines from the pool.  See :ref:`rs_api_filters`.

The provider does not provide any additional options for destroy, but it will unwind all the choices made during allocation.


Output Values
-------------

After allocation, the provider sets key values for provider users.  These include:

* machine_id: maps to Machine.Uuid
* machine_ip: maps to Machine.Address
* machine_name: maps to Machine.Name