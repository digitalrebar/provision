
.. _rs_websocket:

Websocket Access
~~~~~~~~~~~~~~~~

A web socket is provide to retrieve a stream of events.  This can be access by hitting
the *https://<IP>:<APIPORT>/api/v3/ws* URL.  The URL will require authentication, but once
authorized, the resulting web socket will received events.

The events desired must be registered for.  This is done by sending a register
request to the server.  By sending, `register \*.\*.\*`, this will cause all events
that the authorized user can see to be delivered to the websocket receiver.

Multiple registers and the coorespnding deregisters are allowed.

The general form is:
  * register type.action.key
  * deregister type.action.key


The fields are:
  * Type - the type of object (e.g. profiles, machines, ...)
  * Action - the action of the object (e.g. create, save, update, destroy, ...)
  * Key - the specific key of the object (e.g. machine uuid, profile name, ...)


Some example clients, python and javascript, are provided in the DRP source tree.

