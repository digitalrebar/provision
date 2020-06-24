.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Scaling
  pair: Digital Rebar Provision; Endpoint Sizing
  pair: Digital Rebar Provision; Performance

.. _rs_scaling:

Endpoint Sizing, Performance and Scaling Guide
==============================================

RackN’s Digital Rebar Provision needs resources to run and maintain systems.  These resources are defined within a single compute system.  This guide is intended to help describe the situations that need to be considered, the knobs that can be turned, and some example configurations.

Digital Rebar :ref:`rs_high_availability` configurations are active-passive and cannot be used for load sharing.

.. _rs_scaling_parallel_workflow:

Parallel vs Sequential Workflow
-------------------------------

By design, most Digital Rebar operations occur in parallel and most of the delays in a workflow are non-workflow operations like booting or file transfers.

The provisioning long pole is reboots.  During reboots, Digital Rebar is passive waiting for the system to finish posting and begin the PXE process.

Note to improve parallelism:

* A 10 minute reboot time and required reboots.  In practice, reboots may take 3 to 30 minutes depending on the model and configuration.
* Filer could be DRP, but for scale operations it is strongly recommended that a separate scalable File Server be used.  Image download time is dependent upon the number of parallel requests and size of image.  The aggregate time is network bandwidth constrained.
* RackN kexec option may be used to eliminate one or all reboots depending on the operating systems being managed and if firmware updates are required.  Specifically, Linux-to-Linux installs (for machines already under management) do not require REPOSTs.  See :ref:`rs_kexec` in the FAQ.
* Firmware updates require a reboot.  Whenever possible, the automation determines if a change is required in order to skip unneeded reboots.
* For some customers, the post-image customizations take 10 minutes because of external service integrations (these are still parallel per machine).  The tables assume 30 seconds of work.
* RackN image deploy technology is typically 3 to 5x faster than network installs such as preseed/kickstart configuration.  Without image deploys, operating system provisioning time should be estimated as 15 minutes (or an 10 additional minutes).

.. _rs_scaling_operations:

Operational Concerns
--------------------

This sections explains how critical aspects of Digital Rebar processing impact overall system performance.

.. _rs_scaling_api_rate:

API Interaction Rate
~~~~~~~~~~~~~~~~~~~~

API Interaction Rate means how often the API is called by systems interacting with Digital Rebar.  This includes machines being provisioned, administrators using the system via UX and CLI, and third party tools requesting items.  The administrators and third-party tools usually do a small set of actions that lead to machine usage.  The exception to this pattern are event listeners consuming websockets.  Listeners consume some bandwidth and a little CPU, but is handled within the context of other actions taking place.

.. note:: Intelligent interaction design can make a very big difference in resource consumption.  API users that retrieve many full objects instead of using filters and lists will place a larger burden on the API.  RackN has added many indexes to improve query performance and allows users to use :ref:`rs_api_slim` calls with de-populated models.  Using these optimizations can significantly reduce the load on the API.

Machine provisioning consumes the most API actions because it is tightly integrated with system data.  From job rendering to status updates, the Runner drives the most API usage.  These runner API actions require some memory and CPU, but little bandwidth.

Our observation of runner based API load is that the memory consumption is bursty but recoverable over time.  CPU time is consumed rendering results.  Parallel machine provisioning concurrence may be an issue at very large scale, but is self limiting because API interactions are a very small fraction of the overall work being performed.  Typically natural timing variation distributes concurrent operations sufficiently. .

The two basic variables that drive this interaction are the number of machines and the number of machine being concurrently provisioned.

RackN is constantly updating our performance benchmarks based on system testing.  Please, contact RackN latest benchmark information.

.. _rs_scaling_file_rate:

File Serving Rate and Size
~~~~~~~~~~~~~~~~~~~~~~~~~~

File serving rate and size are linked to what the administrator is doing to the machines.  If DRP is going to host the boot images and install images, then the file serving rate will be much higher and require more bandwidth for DRP.  At a minimum, it is expected that sledgehammer and basic booting files will be served from DRP even if the target images are served from another source.  These DRP specialized images are optimized to be small and easy to serve.  As install images or installation repos are added to the mix, additional disk space, network bandwidth, and CPU will be needed to support these transfers.

For scale systems, we recommend that images be served from a load-balanced secondary service that is configured for that purpose.  These systems are often shared with the CI/CD system or general purpose object storage.

.. _rs_scaling_contexts:

Endpoint Contexts (Docker)
~~~~~~~~~~~~~~~~~~~~~~~~~~

The :ref:`rs_contexts` system allows Digital Rebar to shift a Workflow to a container(s) running on the Endpoint.  The resources requirements to support for this feature depend heavily on how it will be used.  When planning for this feature, significant additional host overhead is recommended.

.. _rs_scaling_3rd_party:

Third-Party Integrations
~~~~~~~~~~~~~~~~~~~~~~~~

The final component that can impact the system is third-party integrations through plugins.  Plugins consume CPU, memory, and network resources.  In general, plugins use nearly no disk space, but if this changes, it should be considered too.

Most plugins are doing minor integration actions with minimal impact on system resources.  They may wait on external services, but blocking actions consume little overhead.  This load is plugin dependent; however, current plugins do not do much more than event translation or external service requests.  We do not anticipate this pattern changing.

In general, if you plugins (any number really), an additional GB of memory and an additional core should be sufficient to keep these running.

.. _rs_scaling_components:

Scaling Components
------------------

As we see from above, the biggest drivers for scale are the total number of machines, the number of machines being provisioned concurrently, and the size and location of the provisioning data.

With these elements, we can begin to build recommendations for CPU, Memory, Networking, and Disk.

.. _rs_scaling_cpu:

Processos (CPU)
~~~~~~~~~~~~~~~

For CPU, the biggest driver is concurrent machines provisioning.  This impact is skewed by the type of hardware because BIOS posting and other factors extend the boot times of servers.  In our 750 machine stress test, we used containers to remove hardware boot time and increase API load.  With real machines, boot times will stagger the API loading across the machines to allow for higher machine counts.

The other big CPU load is serving files.  A CPU core per 1.5GB/sec of data is a reasonable baseline.  If you are expecting to draw 10GB/sec of data through the system then you will need about 7 CPU cores to do that and the concurrent clients to feed that.

The recommended minimum is 4 cores for all systems.  This allows for 1 API, 1 File server core, 1 Plugin core, and 1 OS core.  These will get shared.  This should handle up to 1000 total machines and 20 machines concurrently.

For scaling purposes, consider the number of machines concurrently over a period of time with diminishing returns as you scale up.  This means that the max a system can do is about 1000 machine sets of transactions / 5 minutes with 6 cores.  This means 6000 machine sets of transactions / 30 minutes with 6 cores (without networking cores or plugin cores).   So, if we are building a system that serves 1000 machines / 5 minutes with the system serving all images on a 10Gb link, then we would need about a 16 core system (6 API cores, 7 CPU Cores, 2 Plugin Cores, and 1 OS core).   This machine would be completely consumed during peak provisioning load.  Adding more cores will help with plugin integrations and networking and general OS operations.

.. _rs_scaling_ram:

Memory (RAM)
~~~~~~~~~~~~

Memory defines how many machines can be tracked and how much data can be sent quickly.  Concurrency doesn’t matter as much as total data transfered.  Concurrency drives a buffer for bursting, but total data defines overall memory.  For example, creating 750 machines currently will cause a burst of 3GB of memory, but overall consumption is stable around 500MB.  Much of this use is around transient event marshalling.

Additionally, networking performance is improved by caching disk images into memory.  This scenario is where additional memory is useful.  At a minimum, 2 GB should be reserved for caching the common install images of sledgehammer, pxe boot infrastructure, and boot loaders.  Additional memory can be consumed to serve images, but this is only valuable if common images are being used across the platform and there is not a separate file server.

Minimum memory is around 4GB for 1000 machines.  When considering more machines, a reasonable target of 500MB per 1000 machines should be used.  So, 30000 machines would be 15GB of memory with 4GB for caching and OS would be a starting point.  Adding more memory after that addresses image caching and bursting.

.. _rs_scaling_net:

Networking
~~~~~~~~~~

The critical design decision around networking is will the DRP Endpoint be the file server or not.  If not, then 1Gb links are sufficient for basic systems of up to 200 concurrent nodes.  The implication of this design is that a 200MB sledgehammer image is delivered every 2 seconds.  This allows 200 concurrent nodes to boot over the span of 7 minutes.  The challenge to adding more nodes is getting sledgehammer images to the machines as they boot.  This means providing more bandwidth can help basic booting operations.

Bonding can also help deliver bandwidth improvements.

The primary scaling concern will always be around data per second to the machines.  At a minimum, this load is sledgehammer and boot files (200MB per boot).  These resources could also be served through load-balancers at extreme scale.

.. _rs_scaling_store:

Storage Disk Sizing
~~~~~~~~~~~~~~~~~~~

There are two parts to disk sizing.  One is the total size of images to serve.  If the system is going to be a file server, the disks should be sized for that purpose.  This is primarily a function of the number images and their average size.

The other part of selecting disk size and type is the DRP endpoint data backing store.  DRP operates a high write database. It loads data into memory and only writes it on saves, but never reads it.  The implication of this is that write speeds are important for job logs and object storage.  SSD OS drives are highly recommended.  The data stored is not very much, but is written often.  250GB SSDs is fine for small to medium sized deployments (up 4000 machines), but 500GB should be used for larger systems.

Performance can also be helped by storing Sledgehammer and the base files in that space as well.

.. note:: Automatic job and log pruning was introduced in v4.3.  Operators of prior versions must manage their own log housekeeping.
