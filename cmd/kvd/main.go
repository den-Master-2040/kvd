package main

import (
	"log"
	"net"

	"github.com/den-Master-2040/kvd/internal/server"
	"github.com/den-Master-2040/kvd/internal/store"
)

func main() {
	listener, err := net.Listen("tcp", ":6380")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listening on", listener.Addr())
	st := store.New()
	log.Fatal(server.Serve(listener, st))
}
