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
// 	kafkaOutput := &kbapi.KafkaOutput{
// 		Name:           "my-kafka-output",
// 		Type:           "kafka",
// 		ConnectionType: "plaintext",
// 		Hosts:          []string{"kafka.example.com:9092"},
// 		AuthType:       "none",
// 	}
//
// 	req, err := kbapi.UpdateKafkaOutputRequest("db9c0f78-11bd-41cf-af1e-42bc84ad2214", kafkaOutput)
// 	if err != nil {
// 		log.Fatalf("Failed to create request: %v", err)
// 	}
//
// 	// req, err := kbapi.Fleet(logstashOutput)
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to create request: %v", err)
// 	// }
//
// 	resp, err := client.Fleet.Outputs.Update(ctx, req)
// 	if err != nil {
// 		log.Print(err)
// 	}
//
// 	// kafkaOutputs := resp.Body.GetOutputsByType("kafka")
// 	//
// 	// for _, output := range kafkaOutputs {
// 	// 	rawBody, err := kbapi.PrettyPrint(output)
// 	// 	if err != nil {
// 	// 		log.Print(err)
// 	// 	}
// 	//
// 	// 	fmt.Println(rawBody)
// 	//
// 	// }
//
// 	rawBody, err := kbapi.PrettyPrint(resp.Body)
// 	if err != nil {
// 		log.Print(err)
// 	}
//
// 	fmt.Print(rawBody)
// }
