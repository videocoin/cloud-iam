package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/julienschmidt/httprouter"
	"github.com/videocoin/cloud-iam/datastore"
	"github.com/videocoin/common/log"
)

type server struct {
	router *httprouter.Router
	ds     datastore.DataStore
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) routes() {
	s.router.GET("/health", s.handleHealth())
	s.router.GET("/service_accounts/v1/metadata/x509/:sub", s.handlePublicKeyGet())
}

func (s *server) handleHealth() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {}
}

func (s *server) handlePublicKeyGet() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		sub := ps.ByName("sub")
		if _, err := uuid.Parse(sub); err != nil {
			log.Debugln(err)
			s.respond(w, r, nil, http.StatusBadRequest)
		}

		keys, err := s.ds.ListUserKeys(sub)
		if err != nil {
			log.Errorln(err)
			s.respond(w, r, nil, http.StatusInternalServerError)
		}

		jwks := make(map[string]string)
		for _, key := range keys {
			jwks[key.ID] = base64.StdEncoding.EncodeToString(key.PublicKeyData)
		}
		s.respond(w, r, jwks, http.StatusOK)
	}
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
