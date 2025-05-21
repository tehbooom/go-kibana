package kbapi

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
// 	geminiConfig := kbapi.GeminiConfig{
// 		APIURL: "https://127.0.0.1/us-east4",
// 	}
//
// 	geminiSecrets := kbapi.GeminiSecrets{
// 		CredentialsJSON: "supersecret",
// 	}
//
// 	reqBody := kbapi.ConnectorsUpdateRequestBody{
// 		Name: "test",
// 	}
//
// 	reqBody.SetGemini(geminiConfig, geminiSecrets)
//
// 	req := &kbapi.ConnectorsUpdateRequest{
// 		ID:   "12354",
// 		Body: reqBody,
// 	}
//
// 	resp, err := client.Connectors.Update(ctx, req)
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
