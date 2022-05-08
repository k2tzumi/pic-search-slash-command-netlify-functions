package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/slack-go/slack"
	"golang.org/x/oauth2/google"
	customsearch "google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"
)

type handler struct {
	verificationToken    string
	cse                  *customsearch.Service
	customSearchEngineId string
}

type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type ServiceAccountKey struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
}

func NewServiceAccountKey(
	projectId string,
	privateKeyId string,
	privateKey string,
	clientEmail string,
	clientId string,
	clientX509CertUrl string,
) ServiceAccountKey {
	serviceAccountKey := ServiceAccountKey{}
	serviceAccountKey.Type = "service_account"
	serviceAccountKey.ProjectId = projectId
	serviceAccountKey.PrivateKeyId = privateKeyId
	serviceAccountKey.PrivateKey = privateKey
	serviceAccountKey.ClientEmail = clientEmail
	serviceAccountKey.ClientId = clientId
	serviceAccountKey.AuthUri = "https://accounts.google.com/o/oauth2/auth"
	serviceAccountKey.TokenUri = "https://oauth2.googleapis.com/token"
	serviceAccountKey.AuthProviderX509CertUrl = "https://www.googleapis.com/oauth2/v1/certs"
	serviceAccountKey.ClientX509CertUrl = clientX509CertUrl

	return serviceAccountKey
}

func NewHandler(verificationToken string, serviceAccountKey ServiceAccountKey, customSearchEngineId string) (Handler, error) {
	jsonKey, err := json.Marshal(serviceAccountKey)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	// config, err := google.JWTConfigFromJSON(jsonKey, "https://www.googleapis.com/auth/cse")
	// if err != nil {
	// 	return nil, err
	// }
	// client := config.Client(ctx)

	credentials, err := google.CredentialsFromJSON(ctx, jsonKey, "https://www.googleapis.com/auth/cse")
	if err != nil {
		return nil, err
	}
	// service, err := customsearch.New(client)
	service, err := customsearch.NewService(ctx, option.WithCredentials(credentials))
	if err != nil {
		return nil, err
	}
	return &handler{verificationToken, service, customSearchEngineId}, nil
}

func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("slash commnd pase error.", err)
		return
	}

	if !s.ValidateToken(h.verificationToken) {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("validate token error.")
		return
	}

	switch s.Text {
	case "", "help":
		params := &slack.Msg{
			Text:         "*Usage*\n* /ps keyword\n* /ps ksk\n* /ps help",
			ResponseType: slack.ResponseTypeEphemeral,
		}
		b, err := json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("json marshal error.", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(b); err != nil {
			log.Fatal("Response write error.", err)
		}
	default:
		// TODO: Implement keyword counter
		links, err := h.search(s.Text, 1)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("search error.", err)
			return
		}
		params := &slack.Msg{
			Text:         pickup(links),
			ResponseType: slack.ResponseTypeInChannel,
		}
		b, err := json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("json marshal error.", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(b); err != nil {
			log.Fatal("Response write error.", err)
		}
	}
}

func (h *handler) search(keyword string, repeate int) ([]string, error) {
	// Number of search results to return
	const NUM = 10
	start := int64(NUM*(repeate-1) + 1)

	// Search query
	search := h.cse.Cse.List()
	search.Cx(h.customSearchEngineId)
	search.SearchType("image")
	search.Q(keyword)
	search.Safe("active")
	search.Lr("lang_ja")
	search.Num(NUM)
	search.Start(start)

	call, err := search.Do()
	if err != nil {
		return nil, err
	}

	links := make([]string, len(call.Items))

	for index, r := range call.Items {
		links[index] = r.Link
	}

	return links, nil
}

func pickup(links []string) string {
	rand.Seed(time.Now().UnixNano())

	return links[rand.Intn(len(links))]
}
