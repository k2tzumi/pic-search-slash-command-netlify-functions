package main

import (
	"encoding/json"
	"net/http"

	"github.com/slack-go/slack"
)

type handler struct {
	verificationToken string
}

type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

func NewHandler(verificationToken string) Handler {
	return &handler{verificationToken}
}

func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !s.ValidateToken(h.verificationToken) {
		w.WriteHeader(http.StatusUnauthorized)
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
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(b); err != nil {
			panic(err)
		}
	default:
		params := &slack.Msg{
			Text:         s.Text + "\nTODO\nhttps://cdn.pixabay.com/photo/2020/05/30/09/53/todo-lists-5238324_1280.jpg",
			ResponseType: slack.ResponseTypeInChannel,
		}
		b, err := json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(b); err != nil {
			panic(err)
		}
	}
}
