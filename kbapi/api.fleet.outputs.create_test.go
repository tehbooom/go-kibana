package kbapi

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
//
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	ctx := context.Background()
//
// 	// Create a Kafka output configuration
// 	// kafkaOutput := &kbapi.NewKafkaOutput{
// 	// 	Name:           "my-kafka-output",
// 	// 	Type:           "kafka",
// 	// 	ConnectionType: "plaintext",
// 	// 	Hosts:          []string{"kafka.example.com:9092"},
// 	// 	AuthType:       "none",
// 	// }
// 	//
// 	// req, err := kbapi.NewKafkaOutputRequest(kafkaOutput)
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to create request: %v", err)
// 	// }
//
// 	logstashOutput := &kbapi.NewRemoteElasticsearchOutput{
// 		Name:  "my-es-output",
// 		Type:  "remote_elasticsearch",
// 		Hosts: []string{"https://kafka.example.com"},
// 	}
//
// 	req, err := kbapi.NewRemoteElasticsearchOutputRequest(logstashOutput)
// 	if err != nil {
// 		log.Fatalf("Failed to create request: %v", err)
// 	}
//
// 	resp, err := client.Fleet.Outputs.Create(ctx, req)
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
// / }
