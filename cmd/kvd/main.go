package main

import (
	"log"
	"net"

	"github.com/den-Master-2040/kvd/internal/server"
)

func main() {
	listener, err := net.Listen("tcp", ":6380")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listening on", listener.Addr())
	log.Fatal(server.Serve(listener))
}
