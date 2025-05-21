package kbapi

// package main
//
// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
//
// 	"github.com/elastic/elastic-transport-go/v8/elastictransport"
// 	"github.com/tehbooom/go-kibana"
// 	"github.com/tehbooom/go-kibana/kbapi"
// )
//
// func main() {
//
// 	// Create a custom logger using elastic-transport's logger
// 	logger := &elastictransport.ColorLogger{
// 		Output:             os.Stdout,
// 		EnableRequestBody:  true,
// 		EnableResponseBody: true,
// 	}
//
// 	client, err := kibana.NewClient(kibana.Config{
// 	})
//
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	ctx := context.Background()
//
// 	kibanaCase := kbapi.CasesObjectRequest{
// 		Assignees: &[]kbapi.CasesAssignee{
// 			{UID: ""},
// 		},
// 		Category:    kbapi.StrPtr("urmom"),
// 		Owner:       "cases",
// 		Title:       "test-332",
// 		Description: "here is something for",
// 		Tags:        []string{},
// 	}
//
// 	noneConnector := kbapi.NoneConnector{
// 		Name: "none",
// 		ID:   "none",
// 		Type: ".none",
// 	}
//
// 	kibanaCase.SetNoneConnector(noneConnector)
//
// 	req := &kbapi.CasesCreateRequest{
// 		Body: kibanaCase,
// 	}
//
// 	resp, err := client.Cases.Create(ctx, req)
// 	if err != nil {
// 		log.Print(err)
// 	}
//
// 	rawBody, err := kbapi.PrettyPrint(resp.Body)
// 	if err != nil {
// 		log.Print(err)
// 	}
//
// 	fmt.Print(rawBody)
// }
