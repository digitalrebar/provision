.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Rocket Skates documentation under Digital Rebar master license

Install
~~~~~~~


Building Server
---------------

* go generate server/main.go
* go build -o rocket-skates server/\*

Running Server
--------------

* cd test-data
* sudo ../rocket-skates

NOTE: I need the sudo to bind the tftp port.  This is configurable, i.e.  *--tftp-port=30000*  




