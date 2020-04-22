.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Documentation; Contributions

.. _rs_dev_dev:

Contributing to Digital Rebar
=============================

As an open ecosystem project, we encourage community contributions.

.. toctree::
   :maxdepth: 1

   dev-docs
   dev-cli
   dev-curl
   dev-plugins
   dev-server

.. _rs_dev_coding_practices:

Contributor Practices
=====================

Changes to Digital Rebar, including documentation, are managed through our source code management process.  This ensures robust control over the code base and enables RackN to curate the project to ensure quality.

.. _rs_dev_commit:

Commit Message Format
---------------------

The following format is expected (eventually enforced) for all commit messages.  This format helps us assemble changelog entries and release notes.

Commit Message first Line format: ``tag(area): description``

Use on of the folowing tags:

  * build: Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)
  * ci: Changes to CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs)
  * docs: Documentation only changes
  * feat: A new feature
  * fix: A bug fix
  * perf: A code change that improves performance
  * refactor: A code change that neither fixes a bug nor adds a feature
  * style: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
  * test: Adding missing tests or correcting existing tests

Examples are areas include (but are not limited to):

  * [piece]: Server changes should list specific area. (e.g. dhcp, frontend, backend, wal, ha, manager, ...)
  * api: dr-server api
  * cli: drpcli code
  * [name]: plugin module (e.g. raid, rack, bios, ipmi, ...)
  * [name]: content directory (e.g. task-library, drp-community-content, ...)
  * ux: rackn ux

The list of things within the area field can be comma separated.

Examples:
  * feat(dhcp): option 82 processing 
  * fix(cli): allow list operations to be sorted on return
  * feat(backend,frontend,wal): manager objects are embedded in core
  * feat(ux): Search of documentation part of helper bar
  * doc(vmware): fix alignment of documentation headers.

For background, please review https://medium.com/@menuka/writing-meaningful-git-commit-messages-a62756b65c81.

