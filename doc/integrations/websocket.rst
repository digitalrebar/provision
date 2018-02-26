
.. _rs_websocket:

Websocket Access
~~~~~~~~~~~~~~~~

A web socket is provided to retrieve a stream of events.  This can be accessed by hitting
a URL in the form of:

  *https://<IP>:<APIPORT>/api/v3/ws?Token=<_auth_token_>* 
  
The URL will require authentication, but once authorized, the resulting web socket will receive 
events.

Events can be filtered by submitting ``register`` filters to the websocket listener.  

Register for Events
-------------------
The events desired must be registered for.  This is done by sending a register
request to the server.  By sending, ``register *.*.*``, this will cause all events
that the authorized user can see to be delivered to the websocket receiver.

Multiple registers (and the coorespnding deregisters) are allowed.

The general form is:
  * register type.action.key

The fields are:
  * Type - the type of object (e.g. profiles, machines, ...)
  * Action - the action of the object (e.g. create, save, update, destroy, ...)
  * Key - the specific key of the object (e.g. machine uuid, profile name, ...)

Some simple example are provided in the 
`DRP source tree: <https://github.com/digitalrebar/provision/tree/master/integrations/websockets/>`_

Deregister Events
-----------------

If you no longer wish to receive specific events you have registered for, you may 
use the ``deregister`` command.  The command syntax is exactly like the ``register``
command. 

The general form is:
  * deregister type.action.key

Websocket Tools
---------------

Most modern languages provide websocket libraries that you can use to create 
listeners in a given language.  Some examples include (this is NOT an exhaustive
list): 

  * Golang: https://github.com/gorilla/websocket 
  * Golang: https://github.com/gobwas/ws
  * Python: https://www.willmcgugan.com/blog/tech/post/announcing-lomond/ 
  * Python: https://github.com/Lawouach/WebSocket-for-Python/tree/master/requirements
  * Python: https://github.com/websocket-client/websocket-client 
  * Javascript: https://github.com/websockets/ws
  * google is your friend ...
    
There are several exetnsions/add-ons for web browsers that will allow you to do basic
testing of websocket listening.  Here at RackN, we have used the following with some
success:

  * Chrome/Firefox: https://github.com/WangFenjin/Simple-WebSocket-Client
  * Chrome: https://chrome.google.com/webstore/detail/smart-websocket-client/omalebghpgejjiaoknljcfmglgbpocdp?hl=en

There is a simple sample Python script available in the Digital Rebar Provision
repo for reference, see the 
`Websocket Integrations: <https://github.com/digitalrebar/provision/tree/master/integrations/websockets/>`_
page for further details.

Example Information
-------------------

Here is a simple walk through of basic testing on how to use websockets with 
Digital Rebar.  Please note this is fairly basic, but it should get you started
on how to interact with and use websockets.  This example was tested, using the 
"Simple Websocket Client" in both Chrome and Firefox that is listed above. 

We assume you have the DRP endpoint installed on your localhost in these examples. 
You can adjust the IP address/hostname to point to a remote DRP Endpoint, just 
ensure you have access to Port 8092 (by default, or the API port you specify if
you changed the default).

  URL:  ``wss://127.0.0.1:8092/api/v3/ws?token=rocketskates:r0cketsk8ts``

Note that the `token...` information is a set of credentials with permissions
to view events.  This example uses the default username/password pair.  You
may also create and :ref:`specify access Tokens <rs_grant_token>` for the 
websocket client to use.

In the *Request* input box, enter your `register` filter you'd like to receive
events for.

  Request:  ``register profiles.*.*``

This example will only output websocket events related to Parameters.  Now create 
and delete a few test parameters

  ::
    
    # now create a `bar` param on the `global` profile
    drpcli profiles set global param bar to blatz

    # now remove the param from the `global` profile
    drpcli profiles remove global param bar 

...and you should see events like:

  ::

    {"Time":"2017-12-21T23:26:43.412554192Z","Type":"profiles","Action":"save","Key":"global","Object":{"Validated":true,"Available":true,"Errors":[],"ReadOnly":false,"Meta":{"color":"blue","icon":"world","title":"Digital Rebar Provision"},"Name":"global","Description":"Global profile attached automatically to all machines.","Params":{"bar":"blatz","change-stage/map":{"centos-7-install":"packet-ssh-keys:Success","discover":"packet-discover:Success","packet-discover":"centos-7-install:Reboot","packet-ssh-keys":"complete-nowait:Success"},"kernel-console":"console=ttyS1,115200"}}}
    {"Time":"2017-12-21T23:27:15.218761478Z","Type":"profiles","Action":"save","Key":"global","Object":{"Validated":true,"Available":true,"Errors":[],"ReadOnly":false,"Meta":{"color":"blue","icon":"world","title":"Digital Rebar Provision"},"Name":"global","Description":"Global profile attached automatically to all machines.","Params":{"change-stage/map":{"centos-7-install":"packet-ssh-keys:Success","discover":"packet-discover:Success","packet-discover":"centos-7-install:Reboot","packet-ssh-keys":"complete-nowait:Success"},"kernel-console":"console=ttyS1,115200"}}}


