package main

import (
	"log"

	"github.com/nsmith5/super-insecure/pkg/server"
)

func main() {
	s := server.New()

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
