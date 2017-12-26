
.. _rs_CONTRIBUTING:

Contributing to Digital Rebar Provision
---------------------------------------

Before `submitting pull requests <https://help.github.com/articles/using-pull-requests>`_, please make sure to read and understood the Apache license.  Submitting a pull is considered to be accepting the project's license terms.

Guidelines for Pull Requests
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

We follow typical Github fork/pull request processes.

-  Must be Apache 2 license
-  For bugs & minor items (20ish lines), we can accept the request at
   our discretion
-  Does not inject vendor information (Name or Product) into Digital
   Rebar except where relevant to explain utility of push (e.g.: help
   documentation & descriptions).
-  Passes code review by Digital Rebar team reviewer
-  Does not degrade the security model of the product
-  Does not reduce code coverage
-  Items requiring more scrutiny

   -  Major changes
   -  CLI/API changes, especially breaking compatability
   -  New technology

-  Pull requests should be against a defined feature branch in the
   Digital Rebar repo 

Timing
^^^^^^

-  Accept no non-bug fix push requests within 2 weeks of a release fork
-  No SLA - code accepted at PTLs discretion.  No commitment to accept
   changes.

Coding Expectations
^^^^^^^^^^^^^^^^^^^

-  Copyright & License header will be included in files that can
   tolerate headers
-  At least 1 line comments as header for all methods
-  Documentation for API calls concurrent with pull request

Testing/ Validation
^^^^^^^^^^^^^^^^^^^

-  For core functions, the push will be validated to ensure it does NOT break build,
   deploy, or our commercial products
-  For operating systems that are non-core, we will *not* validate on
   the target OS for the push.
-  We expect that a pull request will be built and
   tested in our CI system before the push can be accepted.
