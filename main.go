package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

var users = make(map[string]struct{})
var items = make(map[string]map[string]string)

func handleRegister(w http.ResponseWriter, r *http.Request) {
	// Body {"username": "foo"}
	var req struct {
		Username string
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := users[req.Username]; ok {
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	} else {
		users[req.Username] = struct{}{}
		items[req.Username] = make(map[string]string)
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

func handleItem(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if user == "" {
		// Not "authorized"
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
		value, ok := items[user][itemPath]
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
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
		items[user][itemPath] = req.Value

	case http.MethodDelete:
		delete(items[user], itemPath)
	}
}

func main() {
	http.HandleFunc(`/register`, handleRegister)
	http.HandleFunc(`/items/`, handleItem)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
