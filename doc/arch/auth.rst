.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Auth Models

Authentication Models
<<<<<<<<<<<<<<<<<<<<<

These models work together to manage authentication, authorization,
and other access control mechanisms in *dr-provision*.

.. _rs_data_user:

User
----

Users keep track of who is allowed to talk to drp-provision, and what
actions they are allowed to take in the system.  User objects contain the
following fields:

- **Name**: A unique name for the user.  It cannot be changed after the
  user is created.
- **PasswordHash**: The scrypt hashed version of the user's Password.  This
  field is always empty when accessed via the API.  Changing the Password
  will also rotate the Secret field.
- **Secret**: A random string used to generate and validate access
  tokens.  Changing this field will invalidate any existing tokens,
  and replace Secret with a new random value.
- **Roles**: A list of Role names that the User has been assigned.

In addition to the roles asigned to the User, all Users also get a
claim that allows them to get themself, change their passwords, and
get a Token for themselves.

.. _rs_data_claim:

Claim
-----

Claims grant the ability to perform specific actions against specific
objects.  Claims have the following fields:

- **Scope**: The API top-level path or object type the Claim pertains
  to.  Object IDs referenced in the **Specific** field must be unique
  in this Scope.  This field may be one of the following values:

  - `*`, which refers to all top-level Scopes.

  - A single top-level object type or API path component.

  - A comma-seperated list of top-level object types or API components.

- **Specific**: The specific object instances referred to by the
  Scope.  This field may have one of the following values:

  - `*`, which refers to all Objects in Scope.

  - A single unique ID for the Scope.  This ID must refer to the field
    by which the object is natively referred to in the API.

  - A comma-seperated list of unique IDs.

- **Action**: The action being performed.  Common actions include
  "get", "list","update", and "delete", and different object types can
  have other actions.  This field may have one of the following values:

  - `*`, which refers to all Actions valid in the Scope under consideration,

  - A single Action.

  - A comma-seperated list of Actions

  Two actions have specialized semantics:

  - `action` By itself, `action` gives permission for all valid
    plugin-provided actions.  `action:actionName` gives permission for
    the specific actionName.  If you want to give access to more than
    one plugin-provided action, you can specify multiple
    `action:actionName` instances in the comma-seperated list of
    Actions.

  - `update` By itself, `update` gives permission to update any field
    (including Params, if the final Object has them).  You may also
    give permission to update a specific field in an Object by
    specifying an `update:/Field` specifier.  The `/Field` part after
    the colon must be a valid RFC6901 JSON Pointer to the field you
    want to allow to be updated.  This will allow updated to that
    field and any subfields it might have.  This level of access
    control works on the JSON representation of the Object.

Claims are partially ordered by the access they grant, with the
superuser claim ```{Scope:"*" Action:"*" Specific:"*"}``` granting
access to everything and the empty claim ```{Scope:"" Action:""
Specific:""}``` granting access to nothing.  If you have two claims
`a` and `b`, claim `a` is said to contain claim `b` if `a` is capable
of satisfying every authentication request `b` is.


Role
----

Roles are licensed features -- to perform any interaction with a Role
besides listing them and getting them, you must have a license with
the **rbac** feature enabled.

Roles are named lists of Claims that can be assigned to a User.  Roles
have the following fields:

- **Name**: The unique name of the Role.

- **Claims**: A list of Claims that the Role provides.

Roles are also partially ordered by the access they grant based on
their Claims.  Role `a` is said to contain role `b` if every claim `b`
contains can be satisfied by a claim on `a`

dr-provision provides a default `superuser` role that contains just
the superuser claim.  By default, the rocketskates user will be
assigned this role.


Collectively, Roles and Claims control what a caller can do with the
API.

Tenants
-------

Tenants are licensed features -- to perform any interaction with a
Tenant besides listing them and getting them, you must have a license
with the **rbac** feature enabled.

Tenants control what objects a user can see via the dr-provision API.
Tenants have the following fields:

- **Name**: The unique name of the Tenant.

- **Users**: The list of Users that are in this Tenant. Users can be
  in at most one Tenant at a time.

- **Members**: The objects that are in the Tenant.  This field is
  structured as a JSON object whose keys specify the Scope of the
  objects, and whose values are lists of object indentifiers.  Access
  is only restricted if the Scope of the object is present in the
  Members field of the tenant -- objects whose Scope is not present do
  not have restricted visibility.

Object visibility restrictions based on a Tenant are processed before
Roles are processsed, which means that a Role granting access to an
object that is not allowed by the Tenant will be ignored.

By default, Users are not members of a Tenant, and can therefore
potentially see everything via the API (subject to Role based
restrictions, of course.)

How Authentication Works
<<<<<<<<<<<<<<<<<<<<<<<<

User Tokens
-----------

User tokens are created by accessing `GET
/api/v3/users/:username/token`. By default, a token created using this
method can act on behalf of user. This includes rights for the user to
get information about themself, change their password, and fetch a
token for themselves, and any access granted by additional roles the
user has.  Users can restrict access that a token generated in this
method by passing an optional comma-seperated list of Roles as a
parameter during the API request, and any requested Roles that would
increase the scope of the allowed Claims will be silently dropped.

Machine Tokens
--------------

Certain common machine usage patterns (discovery, running tasks, etc)
also need to interact with the API, and hence need a Token that
authorizes them to perform those actions.  These tokens have a fixed
set of permissions:

- Machine Discovery: This token has the ability to create and get
  Machines, and nothing else.  It is needed to allow Sledgehammer to
  create a machine for itself during initial system discovery.

- Machine Operations: This token gives a Machine the ability to modify
  itself, get stages and tasks, create events, create a reservation
  and modify a reservation for the machine's address, and create and
  manage Jobs for itself.

These machine tokens are generated as part of template expansion via
the .GenerateToken command (which generates tokens that expire
according to the unknownTokenTimeout and knownTokenTimeout
preferences), and the .GenerateInfiniteToken command, which generates
a Machine Operations token that expires in 3 years and is intended to
grant long-term access for the task runner.  These tokens cannot be
generated by any other means.

How Tokens Are Checked
----------------------

1. A request is made to the API. If the request contains
   `Authorization: Bearer`, that token is used.  If the request
   contains `Authorization: Basic`, the contained username/password is
   checked and used to create a one-use Token.

2. Claims are created based on the API path requested and the HTTP
   method.  For example, a `GET /api/v3/users` request creates a Claim
   of ```{Scope: "users",Action:"list",Specific: ""}```, a `GET
   /api/v3/users/bob` creates a Claim of ```{Scope: "users", Action:
   "get" ,Specific: "bob"}```, a `PATCH /api/v3/bootenvs/fred` that
   wants to patch OS.Name and OS.IsoName generates ```{Scope:
   "bootenvs", Action: "update:/OS/Name", Specific: "fred"}``` and
   ```{Scope: "bootenvs", Action: "update:/OS/IsoName", Specific:
   "fred"}```, and so on.

3. The token is checked to make sure it is still valid based on the
   system Secret, the user Secret, and the grantor Secret. If any of
   these have changed, or the token has expired, the API will return
   a 403.

4. The list of created Claims is tested to see if it is contained by
   any one of the Roles contained in the Token, or by any direct
   Claims contained in the Token.  If all of the created Claims are
   satisfied, the request is considered to be authorized, otherwise
   the API will return a 403.

5. The API carries out the request and returns an appropriate
   response.
