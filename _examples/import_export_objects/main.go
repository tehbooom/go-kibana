package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/tehbooom/go-kibana"
	"github.com/tehbooom/go-kibana/kbapi"
)

func main() {
	username := flag.String("username", "elastic", "Username to authenticate to Kibana")
	password := flag.String("password", "changeme", "Password to authenticate to Kibana")
	endpoint := flag.String("endpoint", "https://localhost:5601", "Kibana endpoint")

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := kibana.NewClient(kibana.Config{
		Addresses: []string{*endpoint},
		Username:  *username,
		Password:  *password,
		Transport: transport,
	})

	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	exportReq := &kbapi.SavedObjectExportRequest{
		Body: kbapi.SavedObjectExportRequestBody{
			Objects: []kbapi.Object{
				{ID: "system-Metrics-system-overview", Type: "dashboard"},
			},
		},
	}

	exportResp, err := client.SavedObjectExport(ctx, exportReq)
	if err != nil {
		log.Print(err)
	}

	var object []byte
	for _, msg := range exportResp.Body {
		object = append(object, msg...)
	}

	importReq := &kbapi.SavedObjectImportRequest{
		Params: kbapi.PostSavedObjectImportResponseParams{
			Overwrite: kbapi.BoolPtr(true),
		},
		Body: kbapi.SavedObjectImportRequestBody{
			File: object,
		},
	}

	importResp, err := client.SavedObjectImport(ctx, importReq)
	if err != nil {
		log.Print(err)
	}

	prettyResponse, err := kbapi.PrettyPrint(importResp.Body)
	if err != nil {
		log.Print(err)
	}

	fmt.Println(prettyResponse)
}
