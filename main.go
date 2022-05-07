package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

func main() {
	verificationToken := os.Getenv("VERIFICATION_TOKEN")

	handler := NewHandler(verificationToken)

	http.HandleFunc("/", handler.Handle)

	lambda.Start(httpadapter.New(http.DefaultServeMux).ProxyWithContext)
}
