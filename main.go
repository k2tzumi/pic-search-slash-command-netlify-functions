package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

func main() {
	verificationToken := os.Getenv("VERIFICATION_TOKEN")
	// googleApiKey := os.Getenv("GOOGLE_API_KEY")
	customSearchEngineId := os.Getenv("CUSTOM_SEARCH_ENGINE_ID")
	projectId := os.Getenv("PROJECT_ID")
	privateKeyId := os.Getenv("PRIVATE_KEY_ID")
	privateKey := os.Getenv("PRIVATE_KEY")
	// Replace newline character
	privateKey = strings.NewReplacer("\\n", "\n").Replace(privateKey)
	fmt.Println(privateKey)
	clientEmail := os.Getenv("CLIENT_EMAIL")
	clientId := os.Getenv("CLIENT_ID")
	clientX509CertUrl := os.Getenv("CLIENT_X509_CERT_URL")

	serviceAccountKey := NewServiceAccountKey(projectId, privateKeyId, privateKey, clientEmail, clientId, clientX509CertUrl)

	handler, err := NewHandler(verificationToken, serviceAccountKey, customSearchEngineId)
	if err != nil {
		log.Fatal("new handler faild.", err)
	}

	http.HandleFunc("/", handler.Handle)

	lambda.Start(httpadapter.New(http.DefaultServeMux).ProxyWithContext)
}
