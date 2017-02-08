package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/rackn/rocket-skates/client/machines"

	httptransport "github.com/go-openapi/runtime/client"
	apiclient "github.com/rackn/rocket-skates/client"
)

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	hc := &http.Client{Transport: tr}

	transport := httptransport.NewWithClient("127.0.0.1:8092", "/api/v3", []string{"https"}, hc)
	ac := apiclient.New(transport, strfmt.Default)
	basicAuth := httptransport.BasicAuth("rebar", "rebar1")

	// make the request to get all items
	resp, err := ac.Machines.ListMachines(machines.NewListMachinesParams(), basicAuth)
	if err != nil {
		log.Fatalf("list machines error: %v\n", err)
	}
	log.Printf("%#v\n", resp.Payload)
}
