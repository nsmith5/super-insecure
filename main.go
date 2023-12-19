package main

import (
	"log"

	"github.com/nsmith5/super-insecure/pkg/cli"
)

func main() {
	err := cli.New().Execute()
	if err != nil {
		log.Fatal(err)
	}
}
