package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/escoutdoor/ecommerce/internal/server"
)

func main() {
	s := server.NewServer()

	defer func() {
		if err := s.Close(); err != nil {
			log.Printf("Error closing the server: %s", err)
		}
	}()

	fmt.Printf("server is running on port: %s\n", strings.Join(strings.Split(s.Addr, "")[1:], ""))
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
