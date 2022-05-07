package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

func main() {
	verificationToken := os.Getenv("VERIFICATION_TOKEN")
	googleApiKey := os.Getenv("GOOGLE_API_KEY")
	customSearchEngineId := os.Getenv("CUSTOM_SEARCH_ENGINE_ID")

	handler := NewHandler(verificationToken, googleApiKey, customSearchEngineId)

	http.HandleFunc("/", handler.Handle)

	lambda.Start(httpadapter.New(http.DefaultServeMux).ProxyWithContext)
}
