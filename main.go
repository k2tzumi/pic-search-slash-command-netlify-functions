package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

func main() {
	verificationToken := os.Getenv("VERIFICATION_TOKEN")
	googleApiKey := os.Getenv("GOOGLE_API_KEY")
	customSearchEngineId := os.Getenv("CUSTOM_SEARCH_ENGINE_ID")

	handler, err := NewHandler(verificationToken, googleApiKey, customSearchEngineId)
	if err != nil {
		log.Fatal("new handler faild.", err)
	}

	http.HandleFunc("/", handler.Handle)

	lambda.Start(httpadapter.New(http.DefaultServeMux).ProxyWithContext)
}
