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


Tenants
-------

Watch this space for further developments.

How Authentication Works
<<<<<<<<<<<<<<<<<<<<<<<<

Ditto
