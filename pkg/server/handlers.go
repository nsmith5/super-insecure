package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/nsmith5/super-insecure/pkg/store"
)

type Server struct {
	Store store.Store
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	// Body {"username": "foo"}
	var req struct {
		Username string
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.Store.UserAdd(r.Context(), req.Username); err != nil {
		if errors.Is(err, store.ErrorUserAlreadyExists) {
			http.Error(w, "user already exists", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getUser(r *http.Request) string {
	s := r.Header.Get("Authorization")
	if s == "" {
		return s
	}

	s, found := strings.CutPrefix(s, "Insecure ")
	if !found {
		return s
	}

	return s
}

func (s *Server) handleItem(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if user == "" {
		// Not "authorized"
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	exists, err := s.Store.UserExists(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	itemPath, found := strings.CutPrefix(r.URL.Path, "/items/")
	if !found {
		http.Error(w, "bad path", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		value, err := s.Store.ItemGet(r.Context(), user, itemPath)
		if errors.Is(err, store.ErrorItemNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_ = json.NewEncoder(w).Encode(struct {
			Value string
		}{
			Value: value,
		})
		return

	case http.MethodPost:
		var req struct {
			Value string
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "failed to decode body", http.StatusBadRequest)
			return
		}

		err = s.Store.ItemSet(r.Context(), user, itemPath, req.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case http.MethodDelete:
		err = s.Store.ItemDelete(r.Context(), user, itemPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
