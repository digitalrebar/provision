# rocket-skates



## Tools

We are using [go-swagger](https://github.com/go-swagger/go-swagger) to generate code from the API specification file.
You will need a binary to generate the code.  This can be obtained from [here](https://github.com/go-swagger/go-swagger#static-binary). We are currently using version 0.8.0.

We are using [swagger-ui](https://github.com/swagger-api/swagger-ui) for a quick UI to inspect the API and drive the system.  It is customizable if we need, but we are running it straight up for now.  This needs to be at the top of the directory you run *rocket-skates* from.  

* git clone https://github.com/swagger-api/swagger-ui

*TODO* Make swagger-ui dir a config option.

## Swagger.yaml

We are using an API specification file.  This files real content for the moment lives in StopLight.  The tool also for collobartive editting of the file, testing, and other things.  We will edit there and store the updated copy in github.  StopLight has a git commit style of tracking changes as well.

This is my link, but it should get you in.
https://app.stoplight.io/wk/AEknv6vzcpJa2H5ky/HGM8K52XAR5zJyhhe/f3eE5DeAt6TCSyQLd/design

There are desktop apps for this as well.  They can drive testing and validation.  I'm looking at as well.
Make edits in the app and export the file to swagger.yml in the top directory. 

## Generating the Server

To generate code:
* swagger generate server -P models.Principal -f swagger.yaml 

This generates the following directories:

* cmd
* models
* restapi

The file that we edit for the server is:

* restapi/configure_rocket_skates.go 

This is not regenerated on generate calls. To see the original contents, move the file off and regenerate.  I've tried to make minimal changes to that file and put the code in other directories, i.e. provisioner, so that we can merge easily after generates.

The swagger command has more things, but this is enough to start.


## Generating the Client

To generate the code:
* swagger generate client -P models.Principal -f swagger.yaml 

This generates the following directories:

* client
* cmd
* models
* restapi

The file that we edit for the client is:

* cmd/rocket-skates-client/main.go 

This is completely our own file.


## Test Data

There is a test-data with a cert,key,validator and directories to run a simple instance of *rocket-skates*.  Please don't put big things here.

## Building

* go build -o rocket-skates cmd/rocket-skates-server/main.go

To Run locally:

* sudo ./rocket-skates  --tls-certificate=test-data/server.crt --tls-key=test-data/server.key --tls-port=8092 --backend=directory --file-root=test-data/tftpboot --data-root=test-data/digitalrebar

NOTE: I need the sudo to bind the tftp port.  This is configurable, i.e.  *--tftp-port=30000*  
