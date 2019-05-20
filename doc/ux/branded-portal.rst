.. Copyright (c) 2019 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Branded Portal

.. _rs_branded_portal:

Branded Portal Setup Instructions
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

In some instances, a customer may request a branded FQDN host/domain name
for the Portal.  This document outlines the steps necessary to setup a
branded domain for the customer.  This allows the following current
domain:

  * https://portal.rackn.io

To be served as a "customer branded" URL reference.  An example branded
domain might look like:

  * https://portal.bearmetal.dev/

These branded domains are often referred to as "vanity domains".


Overview
========

The customer must provide the following information or perform the following steps:

  1. Identify the branded domain - but **DO NOT** create a DNS entry for it yet
  2. Provide a valid HTTPS Certificate for the branded domain
  3. Create a verification record in the base domain for Cloudfront validation of domain
  4. Create the DNS CNAME record for the branded domain - once directed by RackN

Here is a quick overview of the process.

  * identify the branded record (but DO NOT create in DNS yet)
  * import the customer Certificate in AWS Amazon Certificate Manager (ACM)
  * create a validation DNS CNAME record the customer to verify domain ownership
  * once domain ownership is verified create the Cloudfront Distribution (use "portal" as example settings)
  * customer creates the branded record in their DNS as CNAME pointing to Cloudfront distribution URL


Customer Requirements
---------------------

The customer must complete, or be capable of directing the completion, of the
following:

  1. ability to add DNS records in the requested domain for the branded URL reference
  2. can provide a valid HTTP Certificate for the branded URL reference (FQDN)
  3. monitor the HTTPS Certificate for expiration, and provide an updated certificate if/when it expires


Detailed Configuration
======================

Here is a detailed overview of the steps necessary.


Identify Branded Domain
-----------------------
>>> who: customer

The customer must define and identify the branded record.  Generally speaking, there should
be at least one branded domain, but may be two if the customer chooses to also brand the
*tip* (or *latest*) version of the Portal.  Examples are as follows:

  * portal.bearmetal.dev - *stable* branded portal
  * latest.bearmetal.dev - *tip* or *latest* branded portal

.. note:: The customer should NOT create the DNS record at this time.  The record
          will be a CNAME entry that points to the Cloudfront distribution, and
          that reference can not be created until later in the process.


Create HTTPS Certificate
------------------------
>>> who: customer

The customer will also provide an HTTPS Certificate that is created
for each of the branded domains they choose to implement.  The customer must
provide the full Certificate PEM for RackN to install for the Cloudfront
distribution to utilize and serve correctly.


Import Certificate
------------------
>>> who: RackN

RackN employees will import the customer Certificate in AWS Amazon Certificate Manager
(ACM), making them available for the Cloudfront distribution.


Validate Domain Ownership
-------------------------
>>> who: RackN provides CNAME record
>>> who: Customer creates DNS record

This step may not be necessary if importing an external certificate.  This step
is necessary if the DNS domain is not hosted via Route 53, and you are creating
a Certificate for the domain.

Should the domain need to be verified, the AWS ACM import process will privde
a validated DNS CNAME and record to be added to the parent domain that the
Certificate is.  The customer must create this CNAME record, and the ACM
site must complete the validation step to continue.


Create Cloudfront Distribution
------------------------------
>>> who: RackN

Once domain ownership is verified (if necessary), create the Cloudfront
Distribution.  Essentially copy the existing ``portal.rackn.io`` distribution
as an example.  Unfortunately AWS has not seen fit to create a simple "copy"
function, so you must view the existing settings and compare them as you
create the new distribution.

Important settings to insure you select:
  * make certain to select **redirect HTTP to HTTPS**
  * must add customer record to **CNAMES** (eg "portal.bearmetal.dev")


Create Branded DNS CNAME Record
-------------------------------
>>> who: RackN provides the CNAME mapping
>>> who: Customer creates DNS record

Once the Cloudfront distribution has been completed, RackN employees will send
to the customer, the correct CNAME record to create for the new branded
domain.

Customer creates the branded record in their DNS as CNAME pointing to Cloudfront
distribution URL.  An example

  ::

    # example BIND file syntax:
    $ORIGIN bearmetal.dev.
    portal        IN      CNAME  d15o1P8hh7J0sj.cloudfront.net.

The Cloudfront URL will be provided once the distribution is completed.
