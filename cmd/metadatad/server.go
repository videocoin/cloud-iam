package main

import (
	"encoding/base64"
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

func newServer() *server {
	s := new(server)
	s.routes()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) routes() {
	s.router.GET("/health", s.handleHealth())
	s.router.GET("/service_accounts/v1/metadata/x509/:sub", s.handlePublicKeyGet())
	s.router.GET("/methods/v1/metadata/perm/:method", s.handlePermissionGet())
	s.router.GET("/v1/roles/:name", s.handleRoleGet())
}

func (s *server) handleHealth() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		s.respond(w, r, nil, http.StatusOK)
	}
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

func (s *server) handleRoleGet() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		role, err := s.ds.GetRole(ps.ByName("name"))
		if err != nil {
			log.Errorln(err)
			s.respond(w, r, nil, http.StatusInternalServerError)
		}
		s.respond(w, r, role, http.StatusOK)
	}
}

func (s *server) handlePermissionGet() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		/*
			role, err := s.ds.GetMethodPermission(ps.ByName("method"))
			if err != nil {
				log.Errorln(err)
				s.respond(w, r, nil, http.StatusInternalServerError)
			}
			s.respond(w, r, role, http.StatusOK)
		*/
	}
}
