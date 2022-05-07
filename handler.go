package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/slack-go/slack"
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

func NewHandler(verificationToken string, googleApiKey string, customSearchEngineId string) (Handler, error) {
	ctx := context.Background()
	service, err := customsearch.NewService(ctx, option.WithAPIKey(googleApiKey))
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
