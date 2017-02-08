# rocket-skates



# Notes!!!

Get Tools:
We are using go-swagger and swagger-ui.

https://github.com/go-swagger/go-swagger
File here:  https://github.com/go-swagger/go-swagger#static-binary
I'm using version 0.8.0

And this: https://github.com/swagger-api/swagger-ui
git clone https://github.com/swagger-api/swagger-ui

To edit the swagger.yaml: Use STOPLIGHT

This is my link, but it should get you in.
https://app.stoplight.io/wk/AEknv6vzcpJa2H5ky/HGM8K52XAR5zJyhhe/f3eE5DeAt6TCSyQLd/design

There are desktop apps for this as well.  They can drive testing and validation.  I'm looking at as well.
Make edits in the app and export the file to swagger.yml in the top directory. 

To generate code:
swagger generate server -f swagger.yaml
swagger generate client -f swagger.yaml

This generates:

client
cmd
models
restapi

The file that we edit is:

restapi/configure_rocket_skates.go 

This is not regenerated on generate calls. To see the original contents, move the file off and regenerate.  I've tried to make minimal changes to that file and put the code in other directories, i.e. provisioner, so that we can merge easily after generates.

The swagger command has more things, but this is enough to start.


test-data is useful for running the app locally.


To Build:
go build cmd/rocket-skates-server/main.go

To Run locally:
sudo ./main  --tls-certificate=test-data/server.crt --tls-key=test-data/server.key --tls-port=8092 --backend=directory --file-root=test-data/tftpboot --data-root=test-data/digitalrebar
