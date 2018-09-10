Installing KRIB (Kubernetes Rebar Integrated Bootstrapping)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

This is about installing KRIB on an existing DRP endpoint.  See :ref:`component_krib` for instructions on using and extending KRIB.

.. _rs_krib:

Prerequists
-----------

You need to install the system with a `discovery` workflow that includes `sledgehammer` as a minimum.  These steps are documented in :ref:`rs_quickstart`.

Install the DRP command line interface (`drpcli`) as per :ref:`rs_cli`

You also need to install some common utilities: git, curl and jq.

Setup DRP Access
----------------

We need to make sure we have access to the system via the CLI.

  ::

    # UPDATE THESE: RS_ENDPOINT is not needed if using localhost
    export RS_ENDPOINT="[endpoint URL]"
    # UPDATE THESE: RS_KEY is not needed if using defaults
    export RS_KEY="[endpoint user:password]"
    # verify credentials
    drpcli get info

Setup Packet API Integration
----------------------------

If you are using Packet.net as a reference platform.  The following steps are specific to that platform.

Make sure you set your information in the exports!

  ::

    # UPDATE THESE: Packet Project Information for Plugin
    export PACKET_API_KEY="[packet_api_key]"
    export PACKET_PROJECT_ID="[packet_project_id]"
    # download plugin provider (update for version or archtiecture)
    curl -o packet-ipmi https://s3-us-west-2.amazonaws.com/rebar-catalog/packet-ipmi/v2.4.0-0-02301d35f9f664d6c81d904c92a9c81d3fd41d2c/amd64/linux/packet-ipmi
    # install plugin provider
    drpcli plugin_providers upload packet-ipmi from packet-ipmi
    # configure plugin
    drpcli plugins create '{ "Name": "packet-ipmi",
       "Params": {
         "packet/api-key": "$PACKET_API_KEY",
         "packet/project-id": "$PACKET_PROJECT_ID"
       },
       "Provider": "packet-ipmi"
      }'
    # verify it worked - should return true
    drpcli plugins show packet-ipmi | jq .Available

.. note:: The URLs provided for plugin downloads will change overtime for newer versions


Install KRIB Components and Cert Plugin
---------------------------------------

The following steps will install the required plugins and content for KRIB

  ::

    # download cert plugin provider (installs the plugin autmatically)
    curl -o certs https://s3-us-west-2.amazonaws.com/rebar-catalog/certs/v2.4.0-0-02301d35f9f664d6c81d904c92a9c81d3fd41d2c/amd64/linux/certs
    # install plugin provider
    drpcli plugin_providers upload certs from certs
    # verify it worked - should return true
    drpcli plugins show certs | jq .Available

    # Get code
    git clone https://github.com/digitalrebar/provision-content
    cd krib

    # KRIB content install
    drpcli contents bundle krib.yaml
    drpcli contents upload krib.yaml

.. note:: This is maintained with more detail at :ref:`component_krib`.

Create Machines (skip if using Terraform)
-----------------------------------------

You do NOT need to create machines if you are using the Terraform integration to build the KRIB cluster!  To create machines in Packet without a Terraform, use the following command:

  ::

    # copy multiple times to create extra machines - names must be unique
    drpcli machines create '{ "Name": "krib-0",
       "Params": {
         "machine-plugin": "packet-ipmi"
       }
    }'

Running KRIB
------------

Continue to next steps on :ref:`component_krib`.