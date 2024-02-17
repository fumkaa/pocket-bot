package server

import (
	"log"
	"net/http"
	"strconv"

	types "github.com/fumkaa/pocket-bot/pkg/db"
	client "github.com/fumkaa/pocket-bot/pkg/pocket/client/Authentication"
)

type AuthorizationServer struct {
	server       *http.Server
	pocketClient *client.Client
	bd           types.TokenDb
	redirecrUri  string
}

func NewAuthorizationServer(pocketClient *client.Client, bd types.TokenDb, redirecrUri string) *AuthorizationServer {
	return &AuthorizationServer{
		pocketClient: pocketClient,
		bd:           bd,
		redirecrUri:  redirecrUri,
	}
}

func (s *AuthorizationServer) Start() error {
	s.server = &http.Server{
		Addr:    ":80",
		Handler: s,
	}
	log.Print("server start!")
	return s.server.ListenAndServe()
}

func (s *AuthorizationServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Print("[ServeHTTP] metehod is not GET")
		return
	}

	chatIdParam := r.URL.Query().Get("chat_id")
	if chatIdParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("[ServeHTTP] parameters chat id is empty")
		return
	}

	chatId, err := strconv.Atoi(chatIdParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("[ServeHTTP] can't parse chat_id in int: ", err)
		return
	}

	reqToken, err := s.bd.RequestToken(chatId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Print("[ServeHTTP] can't get request token from bd: ", err)
		return
	}

	accessToken, err := s.pocketClient.AccessToken(reqToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("[ServeHTTP] can't get access token: ", err)
		return
	}

	err = s.bd.SaveAccessToken(chatId, accessToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("[ServeHTTP] can't save access token in bd: ", err)
		return
	}
	log.Printf("request token: %s, access token: %s", reqToken, accessToken)
	w.Header().Add("Location", s.redirecrUri)
	w.WriteHeader(http.StatusMovedPermanently)
}
