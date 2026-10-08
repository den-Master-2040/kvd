package server

import (
	"bufio"
	"net"
	"testing"
)

func TestPing(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	go Serve(l)

	c, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	c.Write([]byte("PING\r\n"))
	got, err := bufio.NewReader(c).ReadString('\n')
	if err != nil || got != "+PONG\r\n" {
		t.Fatalf("got %q, err %v", got, err)
	}
}
