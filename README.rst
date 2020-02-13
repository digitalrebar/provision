.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. DigitalRebar Provision documentation under Digital Rebar master license
..

.. _rs_welcome:

Digital Rebar Provision Ecosystem
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

*simple, fast and open-ecosystem server provisioning with strong Infrastructure as Code design.*

`Digital Rebar Platform (DRP) Ecosystem <https://rebar.digital>`_ is a suite of `APLv2 <https://raw.githubusercontent.com/digitalrebar/provision/v4/LICENSE>`_ components used with
the RackN DRP v4 Platform to provides a simple yet complete API-driven DHCP/PXE/TFTP provisioning and workflow system.

DRP Platform and Ecosystem is designed to be a complete data center provisioning, content scaffolding and infrastructure workflow platform with a cloud native architecture that completely replaces Cobbler, Foreman, MaaS, Ironic or similar technologies. DRP offers a single golang binary (less than 30MB) with no dependencies capable of installation on a laptop, RPi or switch supporting both bare metal and virtualized infrastructure.

Key Features:

  Platform Capabilities:
    * API-driven infrastructure-as-code automation
    * Multi-boot workflows using composable and reusable building blocks
    * Event driven actions via Websockets API
    * Extensible Plug-in Model for public, vendor and internal enhancements
    * Dynamic Workflow Contexts (allows using APIs when agents cannot be run)
    * Distributed Multi-Site Management
    * Supports ALL orchestration tools including Chef, Puppet, Ansible, SaltStack, Bosh, Terraform, etc

  Open Ecosystem Plugins:
    * RAID, IPMI, and BIOS Configuration (via commercial plugins)
    * Cloud-like pooling capabilities
    * Classification engine for automated workflow
    * High Availability modes (via commercial plugin)

.. figure::  doc/images/dr_provision.png
   :align:   right
   :width: 200 px
   :alt: Digital Rebar Provision
   :target: https://github.com/digitalrebar/provision

.. _rs_community:

Community Resources from https://rebar.digital
----------------------------------------------

.. image:: https://travis-ci.org/digitalrebar/provision.svg?branch=v4
  :target: https://travis-ci.org/digitalrebar/provision
  :alt: Build Status

.. image:: https://codecov.io/gh/digitalrebar/provision/branch/v4/graph/badge.svg
  :target: https://codecov.io/gh/digitalrebar/provision
  :alt: codecov

.. image:: https://goreportcard.com/badge/github.com/digitalrebar/provision/v4
  :target: https://goreportcard.com/report/github.com/digitalrebar/provision/v4
  :alt: Go Report Card

.. image:: https://godoc.org/github.com/digitalrebar/provision/v4?status.svg
  :target: https://godoc.org/github.com/digitalrebar/provision/v4
  :alt: GoDoc

.. image:: https://readthedocs.org/projects/provision/badge/?version=latest
  :target: http://provision.readthedocs.io/en/latest/?badge=latest
  :alt: Documentation Latest Status


* Chat/messaging via the Digital Rebar ``#community`` channel is our preferred communication method.  If you do not have a Slack invite to our channel, you can `Request a Slack Invite <http://www.rackn.com/support/slack/>`_
* `Issues and Features <https://github.com/digitalrebar/provision/issues>`_
* Full `Documentation <http://provision.readthedocs.io/en/latest/>`_ (Github `/doc <https://github.com/digitalrebar/provision/tree/v4/doc>`_ sources are updatable via pull request).
* Videos on the `DR Provision Playlist <https://www.youtube.com/playlist?list=PLXPBeIrpXjfilUi7Qj1Sl0UhjxNRSC7nx>`_ provide both specific and general background information.


.. _rs_quick:

Install & Quick Start
---------------------

.. note::  We HIGHLY recommend using the ``latest`` version of the documentation, as it contains the most up to date information.  Use the version selector in the lower right corner of your browser.

Our `Quick Start <http://provision.readthedocs.io/en/latest/doc/quickstart.html>`_ has fast play-with-it steps.  Don't worry, they are very simple and take 10 to 20 minutes.  You can choose from stable or tip.  Tip is the very bleeding edge of development.

Regular `Install <http://provision.readthedocs.io/en/latest/doc/install.html>`_ for more details on the install steps.  These include production options. (`Previous Version Docs <http://provision.readthedocs.io/en/latest/doc/quickstart.html>`_)

Components & Extensions
-----------------------

Digital Rebar Provision is composable by design.  Much of our advanced funtionality is exposed in :ref:`rs_content_packages` that are added into the system as content and plugins which have documentation embedded in the extension.

.. _rs_toc:

Table of Contents
-----------------

**Reading on Github?** Visit `Generated Docs <http://provision.readthedocs.io/>`_ for a generated ToC.

.. toctree::
   :includehidden:
   :numbered:
   :maxdepth: 1

   doc/quickstart
   doc/install
   doc/environment
   doc/features
   doc/server
   doc/configuring
   doc/release
   doc/upgrade
   doc/workflows
   doc/deployment
   doc/operation
   doc/os-support
   doc/os-support/linuxkit
   doc/ui
   doc/ux/portalux
   doc/Swagger
   doc/cli
   doc/api
   doc/dev/dev-server
   doc/dev/dev-cli
   doc/dev/dev-plugins
   doc/dev/dev-docs
   doc/faq-troubleshooting
   doc/arch
   doc/content-packages
   doc/rackn/license
   CONTRIBUTING
   Trademark
   LICENSE

.. _rs_license:

License
-------
DigitalRebar Provision code is available from multiple authors under the `Apache 2 license <https://raw.githubusercontent.com/digitalrebar/provision/v4/LICENSE>`_.

Digital Rebar Provision documentation is available from multiple authors under the `Creative Commons license <https://en.wikipedia.org/wiki/Creative_Commons_license>`_ with Attribution.

::

    Work licensed under a Creative Commons license is governed by applicable copyright law.
    This allows Creative Commons licenses to be applied to all work falling under copyright,
    including: books, plays, movies, music, articles, photographs, blogs, and websites.
    Creative Commons does not recommend the use of Creative Commons licenses for software.

    However, application of a Creative Commons license may not modify the rights allowed by
    fair use or fair dealing or exert restrictions which violate copyright exceptions.
    Furthermore, Creative Commons licenses are non-exclusive and non-revocable.
    Any work or copies of the work obtained under a Creative Commons license may continue
    to be used under that license.

    In the case of works protected by multiple Creative Common licenses,
    the user may choose either.
