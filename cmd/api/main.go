package main

// Main entry point for the API server
import (
	"template-golang/server"
)

func main() {
	s, err := server.Wire()
	if err != nil {
		panic(err)
	}
	s.Start()
}
