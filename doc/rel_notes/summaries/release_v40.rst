.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Release v4.0
  pair: Digital Rebar Provision; Release Notes


.. _rs_release_v40:

Digital Rebar version 4.0 [Summer 2019]
---------------------------------------

Release Date: August 8, 2019

Release Themes: License Change, New Backend

TLDR: We’re inverting our open license model to empower operations collaboration on the platform while ensuring that RackN can generate sustaining revenue from production users.

See :ref:`rs_release_summaries` for a complete list of all releases.

Summary
~~~~~~~

After 6 years, Digital Rebar has evolved into a powerful data center operations platform.  The data center open source ecosystem has also evolved.  It’s time for RackN to invert our open source strategy: we are making the bulk of our commercial extensions catalog open source while moving the RackN contributions to the Digital Rebar backend implementation under a commercial license.  The open Digital Rebar platform continues as a set of robust APIs, access clients and test libraries.  Our freemium model and commercial licensing remains the same.


Inverting the Digital Rebar License
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Maintaining software for production data center is not free and the operators of the software require confidence in its continuity.  They also want an open way to collaborate and share.  At RackN, we constantly evaluate how to achieve that balance.  Looking at our history: we’ve had robust engagement on the open parts of our catalog and no one contributing to the backend code.

As the sole stewards of Digital Rebar, we constantly evaluate how to evolve our open source foundations to ensure we are both building a strong community and ensuring financial success to support it.  Since our software must be reliable over the long term to support production data centers, we are careful to distinguish between free products and open source software.  The upcoming major release of Digital Rebar requires us to adjust that balance.

RackN is proud to announce that we’re moving the vast majority of our current content and plugin libraries to an open source model.  This includes critical infrastructure management tooling around firmware configuration, immutable deployment processes, security operations, external systems integrations including IaaS platforms, and automated inventory management.  The newly opened code based reflects years of expertise RackN has built around the Digital Rebar platform.

We make this change because we recognized that Digital Rebar “Provision” has grown from an operating system provisioner into a fully functional data center automation platform.  Our operator inspired vision for the platform is driving extensions that will reach deeper into the physical layer than ever before.  It also drives the platform to scale up and out in new and powerful ways for both Edge and Enterprise data centers.

Keeping Digital Rebar APIs, tests and clients open is essential for that platform to grow and fulfill that mission; however, it has become important to distinguish between the API and our implementation of the platform.  The RackN implementation, previously intertwined with the Digital Rebar API specification, will become closed and proprietary.  This change allows us to make changes that support scale, security and reliability for production data centers.

Availability of Digital Rebar as a platform will not change.  RackN continues to maintain our public distribution of Digital Rebar as before.  We continue to enable community members to use our implementation free and anonymously for up to 20 machines and with additional features for operators who voluntarily provide identification.  In addition, we are expanding our trial, special-use and non-profit license options.

We have heard clearly that Digital Rebar platform operators want open contribution access to our content library.  We are excited to take this step to enable collaboration within our community and look forward to seeing what we can build together.
