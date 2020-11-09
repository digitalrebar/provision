.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Security FAQ

.. _rs_security_faq:

Security Frequently Asked Questions (FAQ)
=========================================

The following questions from customer security reviews that may be generally helpful in understanding Digital Rebar security.  The questions are generally organized into categories.
 
This FAQ page is constantly evolving.  If you cannot find the answer to your question here, please let us know and we’ll add to the page.


.. _rs_security_authentication:

Authentication Questions
------------------------

What are the authentication methods for admin accounts?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

There are two internal methods.  

#. Username/Password over HTTPS connection (Basic Auth)
#. Token over HTTPS connection (Bearer Auth)
 
Tokens can grant a subset of a user’s abilities to restrict access and control.   We strongly recommend using token security as much as possible.  The preferred pattern is to use username/password authentication to create a token.  The token should be used for all subsequent requests because it is more secure and performant.
 
The Single Sign On (SSO) plugin allows Digital Rebar to delegate authentication to external authentication services such as LDAP or Active Directory.  Roles returned from that service will be mapped back into Digital Rebar roles.  No user accounts need to be created in advance.
 
With the addition of SSO capabilities, RackN chose to delegate advanced authentication features to the SSO system rather than re-implementing them in the internal authentication system.  For that reason, Digital Rebar user authentication options are kept minimal.

What are the authentication methods for user accounts?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Same as the admin accounts.  We use roles to determine privilege for accounts.

Is there an MFA for an admin account?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

There is no Multi-Factor Authentication (MFA) for the internal authentication system.  If MFA is required, consider using an SSO with MFA.

.. _rs_security_authorization:

Authorization Questions
-----------------------

Can the “rocketskates” super admin be disabled?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Yes.  This is an option in install.  The account is a default and not a hard coded requirement.

How does DRP authorization work?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

There is a role system integrated with the user and token systems to restrict access to subsets of objects and restricting the actions taken on them.

Does DRP support Multi-Tenant allocations?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Yes.  Tenants are primary objects in DRP.  Users and machines can be assigned into tenants to limit access to specific machines.
 
Tenants are maintained in a flat model.  DRP does not support nesting tenants, users or machines.

Does DRP have built in Roles?  Can they be overridden?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Yes and Yes.  Digital Rebar includes a default “superuser” role that has unrestricted rights claims (`Scope: *, Action: *, Specific: *`) in the system.  It is possible to define new roles with similar claims and then remove the superuser role.

.. _rs_security_logging:

Logging and Tracking Questions
------------------------------

Is there an accounting mechanism?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

All API actions generate events and these  are logged to files or to external event listeners.
 
The events contain what user took the action.

Are authentication activities logged?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Yes in multiple ways
 
Most simply, they are logged stdout by default.
 
In addition, all authentication activities, including token creation, also generate events against the User Model that can be forwarded.

Are admin activities logged?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Yes.  Admin is not a distinguished class in the system.  All actions are evented.

Are the logs timestamped?
~~~~~~~~~~~~~~~~~~~~~~~~~

Yes.  Events and logs are timestamped.

How long are the logs stored in the solution?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Operators define their own log rotation scheme.
 
DRP is usually run under systemd with interaction into many log capture/rotation systems.

Can the logs be sent to Splunk and/or other solutions?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Yes. We have a plugin currently that will integrate with Filebeat to send logs and events into ELK stacks.

Plugins could be written for Splunk.

Are accesses to sensitive data logged?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

We do not distinguish between access to regular data or protected data.

Is sensitive data logged?
~~~~~~~~~~~~~~~~~~~~~~~~~

We try to filter all sensitive data out of the server logs.  Sensitive data may be included at higher log levels (debug or trace) so production systems should never run at elevated log levels for prolonged periods.
 
Job logs, which are often operator created content, may contain sensitive information.  They are maintained separately so they can be quickly purged or managed independently of server logs.

.. _rs_security_condidentiality:

Confidentiality Questions
-------------------------

What information does Digital Rebar send to RackN?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Digital Rebar Server does NOT send any information to RackN.  This is required for air gap operation.  All Digital Rebar functions are available via the API and CLI.

Select commands from the CLI will connect with RackN servers to retrieve the catalog and related content.

The UX in default configurations does connect with RackN for mailbox and license validation.   We also collect non-identifying information about the endpoint such as ID, machine count and entitlements.  We do NOT store anything else about your environment or access in the RackN SaaS.  The UX automatically creates a unique anonymous identifier for mailbox communications.

See below for more details about compromises related to RackN managed systems.

What private information does RackN store?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

RackN does NOT hold any confidential or identifying information from customers’ systems or deployments.

The information collected is:

* endpoint ID
* endpoint IP address
* entitlement data (machines, license data, etc)
* deployment versions
* content packs that are installed
* the IP address of the user’s browser

To obtain a RackN license, an active email address (could be an alias) is required.  For contact purposes, we also request name and phone number.  


Is all the flow between Digital Rebar and the provisioned machines secured?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

It depends on the protocol required at each stage.  When using the API, yes.

No.  DHCP and the initial boot process (TFTP and HTTP) cannot be secured due to the limitations of the protocol.  Once started, the system transitions to secured channels.

RackN has designed some alternative paths to avoid TFTP and HTTP; however, the operational impact of these alternatives may not be justified.

RackN works very hard to minimize the time using of these protocols and can be

Does the CLI use an SSH connection?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

No.  The Remote CLI does not use SSH.   We don’t use SSH anywhere in the product.
 
The CLI uses an HTTPS connection to the DRP API. 

Can I restrict the allowed ciphers for API connections?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Yes.  DRP server v4.5.2+ can restrict clients' ciphers; however, operators need to take the addional step to restrict DRP server's ability to use ciphers with a `--min-tls-version` start-up flag.

Determine the current and available ciphers using `--tls-cipher-list` and `--tls-ciphers-available`.


Is the admin password strongly encrypted?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. note:: CHANGE YOUR ADMIN PASSWORD FROM THE DEFAULT

The password is only saved as a one-way calculated secret hash (scrypt).  This is stored on the user object.  It is possible to perform this encryption outside of the system and store the passwords as hashed data.
 
Parameters that have been flagged as Secure are stored in encrypted format.
 
Versions prior to 4.2 stored data as json files on the Server's disk.  Older versions are not recommended for production.
 
Digital Rebar does not have any external database.

How are the users IDs (login/pass) stored? Are they encrypted?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The Passwords for users (the same as admin) are stored as one-way hashes for comparison.  We do not store the user passwords on disk on their cryptographic hash.
 
Digital Rebar does not store passwords when SSO is enabled.

Does a full disk encryption feature exist or can we implement it?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Not yet.  We are in the process of exploring and implementing a LUKS process for encryption/decryption of machines during boot.  If this is interesting to you, we should talk about it.  

.. _rs_security_availability:

Service and Availability Questions
----------------------------------

What are the most likely causes of disruption or downtime?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

DRP is provided a single go-binary.  This is usually run under systemd to handle restarts after catastrophic failures.  DRP starts within seconds. 
 
DRP Enterprise includes active/passive(s) high available (HA) features to automatically synchronize data between endpoints.  By design, an additional service such as Corosync Pacemaker is needed to manage automatic failover between endpoints, if that is a concern.

What strategies and safeguards does the service/product have to help avoid disruption or downtime of the service/product?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

We have a complete HA document for this purpose with a range of options.

Can I run DRP in a (docker) container?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Yes, but there are operational considerations.

Containers may need extra configuration to handle UDP protocols like DHCP or TFTP.  Consult the install documentation.

Running in a container does not work well if you are using the Multi-Site Manager to handle upgrades of the DRP binary.

Make sure that you install DRP with the destroy container, deploy new version of container.  Then back the persistent data in a volume, so you can detach/reattach that to the new container.

Which ports are required for DRP?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The current list of required and optional ports is maintained under :ref:`rs_arch_ports`.

Does DRP have unauthenticated HTTP/HTTPS reads?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Yes.  This is required by provisioning process(es) because systems bootstrapping do have foreknowledge of credentials.  No unauthenticated writes are allowed.

Where possible, Digital Rebar Platform always uses TLS encrypted API, File Transfer, and Websocket communications using industry standard certificates.  User accounts are used with Role Based Access Controls (RBAC), and multi-tenant isolation capabilities.  Generally speaking, a user generates JWT based tokens by authorizing with their user/pass pair, to build a limited use token which has specific rights (claims and scope) assigned to it.  Token management is handled internally to the service.

Network based operating system installations require the integration with hardware Network Interface Card (NIC) ROM (read-only memory) based capabilities.  As such, physical device management relies on DHCP, TFTP, and HTTP protocols to bootstrap and start network based provisioning (eg PXE) services.  These protocols are fundamentally required, can not be stripped out of the NIC ROM without rewriting with new firmware, and are not encrypted.  Wherever possible, RackN utilizes a multi-step strategy that requires starting from clear text DHCP / PXE process to get boot artifacts via TFTP, then switch to HTTP or HTTPS protocols for safety and security whenever possible.

RackN limits the exposure to unauthenticated information as much as possible:

* DRP dynamically generating templates based on machine state so the amount of information available is limited determined
* DRP transitions data exchanges to the secured API as much as possible
* DRP workflow relies on per machine limited scope tokens to limit access during workflow even during secured operations.

Does RackN support UEFI Secure Boot capabilities?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Yes.  RackN supports UEFI Secure Boot capabilities.  Additional license entitlements are required.


How can I disable insecure PXE protocols like TFTP and HTTP?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

It is possible to run DRP without TFTP or HTTP enabled; however, they may be required to be enabled for your environment.

Unfortunately, core parts of the legacy PXE bootstrap use insecure protocols.  If your infrastructure requires Legacy BIOS or has other PXE dependencies then you’ll need to enable them in DRP.

RackN works hard to minimize use of these protocols.  Please consult with RackN for suggestions about reducing or eliminating their use.

Is a self-signed TLS certificate required?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

A TLS certificate is required for the DRP API which only uses HTTPS.  It does not have to be self-signed.

The self-signed certificate is generated by default for ease of use when installing DRP.  Production users should replace the self-signed certificate with a trusted certificate.

Can I run DRP without Host Root Access?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Yes.  DRP does not require root permission; however, the DRP operational account must have sufficient permissions to open ports and perform operations.  Please see the installation guide for details.

.. _rs_security_integrity:

Integrity Questions
-------------------

Is the flow between a DRP and a provisioned machine authenticated?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

There are two sets of flows for DRP to provisioned servers.

#. The first flow is for basic booting.  These files are served over tftp/http and are not secured. 
#. The second flow is for configuration; these actions are done over the secured HTTPS ports.  These actions use token-based authentication that are restricted to the machine only.


Do DRP services intercommunicate in an authenticated way?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The DRP service is self-contained go-binary.  All services talk internally through memory operations.  Plugins are run locally and use unix/domain sockets for their communication.


What information is at risk from a "man-in-the-middle" (MITM) attack?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The CLI and UX both use authenticated HTTPS API calls to control the system that requires authenticated access to control the system.  We recommend using a chain of trust certificate, instead of self-signed, certificate for production systems.

During a client session, a time limited token is granted after initial authentication.  All subsequent requests use the token.


Does DRP use encryption and hash algorithms?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

DRP does ship with a hash of its installation tarball and a hash of all the components in that tarball for validation at installation time.  It does not self test.

Are DRP services isolated from each other?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

No, DRP is one service; however, services are managed as isolated processes in the binary

Services include: DHCP, API, HTTP Files, TFTP file, Swagger UI

Is DRP scalable?
~~~~~~~~~~~~~~~~

Yes.  DRP scales by segmenting Data centers into pieces with content packages being a common deployment sync method.

The internal data storage uses a write logging process with check points.  This allows DRP to optimize lock and write behavior even with 1,000s of concurrent operations.
 
Additionally, DRP is light-weight and has been performance tested to ensure scale.  We have a scaling document to assist in tuning DRP host environments.

Please consult :ref:`rs_scaling` for additional details.

How sensitive data are stored?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Parameters are the primary method of storing information on plugins, machines, and profiles.  These have two forms, normal and secure.  Secure parameters are maintained in a separate data store that is encrypted. 
 
In the future, these parameters could be stored in Hashicorp Vault for example.  This is a roadmap item that is awaiting prioritization. 

See :ref:`rs_data_param_secure` for additional details.


.. _rs_remote_access:

Remote Access to DRP by RackN
-----------------------------

RackN does not connect to any DRP endpoints.  Users of the UX or CLI are connecting directly to the DRP through their network connection.


Can RackN take DRP actions via an within open session?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

No. Open sessions are directly from the user's browser to the DRP endpoint only.  No inbound communication to the DRP server or user's browser is used or allowed.

All browser to DRP endpoint communication is direct between the browser and the endpoint.  All API communication is secured with TLS and uses a time based token for authentication.


Does RackN have access to my DRP Passwords?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

No.  DRP passwords are not sent or stored by RackN.  Passwords are sent directly to the DRP endpoint from the UX or CLI: they are hash checked or passed to active directory for validation.  Once validated, the password is discarded.  This is only done to the DRP endpoint.  

RackN SaaS does not store any operator passwords (or internal data) for deployed software  The only information passed to the DRP SaaS is the DRP license identity, usage counts, and usage of plugins and contents.  This is not a specific configuration.

Does RackN have access access to on-site iLO, iDrac, etc.?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

No.  No actions are exposed (or notifications of their use) to RackN.  Further, the UX does not act against those items directly.  Requests for out-of-band management are funnelled through the DRP endpoint and must be validated by DRP security.

Can the RackN leak information to attackers including: PII data or security configurations.
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

No. The UX does not send PII or security configuration information to RackN; however, it does transmit the DRP information block including version to RackN.  If operators are not keeping up with latest CVEs (see :ref:`rs_release_summaries`) then RackN tracked information could be used to exploit known issues.


Can RackN render the DRP server unavailable or modify behavior remotely
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

No.  RackN does not have any remote control over the DRP and cannot remotely disable DRP operations.

Note that an expired RackN license key will disable DRP parts of functionality.  This is not a remote operation, it is based on information contained in the locally installed license file and cannot be modified by RackN once the license is issued.  Operators should always be aware of their license entitlements and expiration dates.


What happens if the RackN.io services are unavailable?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

In practice, this does not impact active users because it not in the active control flow; however, if RackN.io is not available then operators will not be able to download the UX for new sessions.  The processes that the UX uses collect data (see above) and send information to RackN.io are non-critical to the application and do not interrupt UX operation if interrupted.

We recommend that operators install a local copy of the RackN UX as a backup.


What is at risk from a RackN insider threat or 3rd party website compromise?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

A compromise of RackN tracked information would potentially provide an attacker with information about the DRP version installed, it's internal name or IP address (depeending on customer configuration) and number of machines managed by each DRP endpoint.

Since RackN has no access or credentials, this information is only of value for an attacker who has already penetrated the customer networks and then discovered customer DRP access information.


.. _rs_security_overall:

Overall security Security
-------------------------

.. _rs_faq_cve:

Does RackN maintaining a Common Vulnerabilities and Exposures (CVE) list?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Yes, see the list below or the Release :ref:`rs_cve` section.  The reporter and RackN customers get advanced notice before public reporting (typically 30 days).  `Create a ticket <https://rackn.zendesk.com/hc/en-us/requests/new>`_ to report an issue.

.. toctree::
   :maxdepth: 1
   :glob:

   security/*


Is DRP protected against Top 10 OWASP?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

See OWSAP reference: https://owasp.org/www-project-top-ten/

Brief comments regarding the OWASP top 10 list

#. Injection: there is no SQL database in DRP. 
#. Broken Authentication: no known issues and tokens are time and scope limited. 
#. Sensitive Data Exposure: parameters can be stored securely. 
#. XML External Entities (XXE): there is no XML in DRP.
#. Broken Access Control: no known issues.
#. Security Misconfiguration: we help mitigate this issue.  DRP makes patch and upgrade of DRP easy via the API.
#. Cross-Site Scripting XSS: DRP is API driven.
#. Insecure Deserialization: do not install the DRP agent, endpoint-exec or contexts if this is a concern.
#. Using Components with Known Vulnerabilities: we maintain a list of known component and work to mitigate them when we are aware of issues.
#. Insufficient Logging & Monitoring: we have extensive logging and encourage exporting logs to tools for additional analysis.


Has the Digital Rebar solution been penetration tested?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Yes, but we have participated in customer specific penetration tests that do not create public or sharable reports.

Do you have data flow diagrams?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

RackN has many graphics about data flows and need more details to provide the correct reference material.  There are provisioning dataflows, discovery dataflows, configuration dataflows, plugin dataflows.

We are in the process of migrating this information to this documentation site.  Please contact us if you'd like access.  

Can I customize the UX based on role?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Yes.  The UX through the ux_views plugin can create custom behaviors based on user roles.  These behaviors can be created ad hoc or through the normal content system.

Is Idle Session Timeout implemented?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The RackN UX has a user settable session timeout (default is 24 hours).  The ux_views plugin must be installed to expose this feature.
 
The DRP CLI uses maintained connections with tokens that are short lived by refreshed.  Token duration is selected when the token is created.  This way if the DRP CLI the token store to speed up connection processing times out quickly (within an hour).

Are session tampering controls implemented?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The normal communication paths are over HTTPS and reduces this issue.  In addition, tokens are encrypted by the server with it’s own uniquely generated key.
 
Additionally, tokens have markers and times in the data to facilitate secondary validation.

Which kind of data are processed by the application? Stored by the application?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Our system processes inventory and state information about the machines being managed.  In general, this is NOT your specific data but information about the system.

Some data needed to deploy the system will be potentially sensitive, e.g. ipmi/password, base words for operating system, etc.  These are stored in secure parameters.

One of the niceties of the image deploy system is that DRP doesn’t have to be involved in any of that data.  Those images can reside outside of DRP and referenced.  DRP and RackN try to keep as little information about the actual work the system is doing other than what is minimally needed to provision that system. 


.. _rs_security_general:

General Questions
-----------------

How does entitlement licensing work?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

See :ref:`rackn_licensing` for how to management entitlement files.

License entitlements are enforced from the DRP server using an encryption key.  The key controls a number of entitlements for DRP including version, number of machines/contexts/pools, HA enabled, secure boot enabled and expiration dates.

The key includes the DRP server identities (aka Endpoint ID) covered by the license and that needs to be updated for each endpoint added.  There is a self-service API that allows license holders to add endpoints to their license.

The key is distributed to operators as the ‘rackn-license’ context package that includes both the key and a plain text version of the entitlements.   This allows operators to manage licenses alongside their operational content without relying on a different path to manage licenses.

How do I know that licensing is enabled or disabled?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

DRP will generate events when entitlements have been met or dates are past the expiration.

Entitlement controlled features and plugins will fail to operate after the hard expiration date.   If they have not loaded, their plugin references will be marked as `Available: false`.

Licenses can be updated (as a content package) and reset on the fly without down time.

What is the release and patch frequency?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

In general, we try to have a release a quarter (sometime monthly depending upon feature enhancements).  We attempt to maintain compatibility and only add new features.  Bugs will be triaged and force immediate releases or wait until the quarter or monthly boundary.
