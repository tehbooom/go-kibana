package kbapi

//
// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
//
// 	"github.com/elastic/elastic-transport-go/v8/elastictransport"
// 	"github.com/tehbooom/go-kibana"
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
// 		Addresses: []string{""},
// 		Username:  "",
// 		Password:  "",
// 		Logger:    logger,
// 	})
//
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	ctx := context.Background()
//
// 	req := &FleetProxiesCreateRequest{
// 		Body: FleetProxiesCreateRequestBody{
// 			Name: "test-headers",
// 			URL:  "https://127.0.0.1:902",
// 		},
// 	}
//
// 	headers := map[string]interface{}{
// 		"Authorization": "Bearer token123",
// 		"EnableFeature": true,
// 		"MaxRetries":    5,
// 	}
//
// 	err = req.Body.SetProxyHeaders(headers)
// 	if err != nil {
// 		log.Print(err)
// 	}
//
// 	currentHeaders, err := req.Body.GetProxyHeaders()
// 	if err == nil {
// 		for k, v := range currentHeaders {
// 			fmt.Printf("%s: %v\n", k, v)
// 		}
// 	}
//
// 	resp, err := client.Fleet.Proxies.Create(ctx, req)
// 	if err != nil {
// 		log.Print(err)
// 	}
//
// 	rawBody, err := PrettyPrint(resp.Body)
// 	if err != nil {
// 		log.Print(err)
// 	}
//
// 	fmt.Print(rawBody)
// }
