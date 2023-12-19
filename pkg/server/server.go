package server

import "net/http"

func New() http.Server {
	http.HandleFunc(`/register`, handleRegister)
	http.HandleFunc(`/items/`, handleItem)
	return http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}
}
