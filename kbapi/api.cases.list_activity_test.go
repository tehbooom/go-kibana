package kbapi

// import (
// 	"context"
// 	"encoding/json"
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
// 	req := &kbapi.CasesListActivityRequest{
// 		ID: "ec8c015e-1384-4336-afc0-5a4abd4277c0",
// 		Params: kbapi.CasesListActivityRequestParams{
// 			Types: kbapi.SliceStrPtr([]string{"comment"}),
// 		},
// 	}
//
// 	resp, err := client.Cases.ListActivity(ctx, req)
// 	if err != nil {
// 		log.Print(err)
// 	}
//
// 	actions, err := resp.Body.GetAllTypedUserActions()
// 	if err != nil {
// 		log.Print(err)
// 	}
//
// 	typeActions, err := resp.Body.GetUserActionTypes()
//
// 	for i, action := range typeActions {
// 		switch action {
// 		case "create_case":
// 			var action kbapi.CasesPayloadCreateCase
// 			json.Unmarshal(resp.Body.UserActions[i].Payload, &action)
// 			fmt.Println(action.Description)
// 		case "comment":
// 			userAction := resp.Body.UserActions[i].Payload
// 			var typeContainer struct {
// 				Type string `json:"type"`
// 			}
// 			err := json.Unmarshal(userAction, &typeContainer)
// 			if err != nil {
// 				log.Print(err)
// 			}
//
// 			if typeContainer.Type == "alert" {
// 				var action kbapi.CasesPayloadAlert
// 				err = json.Unmarshal(userAction, &action)
// 				if err != nil {
// 					log.Print(err)
// 				}
// 				fmt.Printf("Rule name for action %s is %s\n", resp.Body.UserActions[i].ID, action.Comment.Rule.Name)
// 			} else {
// 				var action kbapi.CasesPayloadComment
// 				err = json.Unmarshal(userAction, &action)
// 				if err != nil {
// 					log.Print(err)
// 				}
// 				fmt.Printf("Comment is %s\n", action.Comment.Comment)
// 			}
//
// 			var action kbapi.CasesPayloadComment
// 			json.Unmarshal(resp.Body.UserActions[i].Payload, &action)
// 		case "connector":
// 			var action kbapi.CasesPayloadConnector
// 			json.Unmarshal(resp.Body.UserActions[i].Payload, &action)
// 			fmt.Println(action.Connector.Name)
// 		}
// 	}
//
// 	for _, action := range actions {
// 		pretty, _ := kbapi.PrettyPrint(action)
// 		fmt.Println(pretty)
// 	}
//
// 	pretty, err := kbapi.PrettyPrint(resp.Body)
// 	fmt.Println(pretty)
// }
