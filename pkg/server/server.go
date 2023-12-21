package server

import (
	"net/http"

	"github.com/nsmith5/super-insecure/pkg/store"
)

func New(db store.Store) http.Server {
	server := Server{
		Store: db,
	}
	http.HandleFunc(`/register`, server.handleRegister)
	http.HandleFunc(`/items/`, server.handleItem)
	return http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}
}
