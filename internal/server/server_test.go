package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"testing"

	"github.com/den-Master-2040/kvd/internal/store"
)

func TestPing(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	go Serve(l, store.New())

	c, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	r := bufio.NewReader(c) // один читатель на всё соединение

	// do кодирует команду в RESP, отправляет её и читает ответ ровно той длины, что ждём.
	do := func(want string, args ...string) {
		t.Helper()
		var b strings.Builder
		fmt.Fprintf(&b, "*%d\r\n", len(args))
		for _, a := range args {
			fmt.Fprintf(&b, "$%d\r\n%s\r\n", len(a), a)
		}
		if _, err := c.Write([]byte(b.String())); err != nil {
			t.Fatal(err)
		}
		got := make([]byte, len(want))
		if _, err := io.ReadFull(r, got); err != nil {
			t.Fatalf("%v: read: %v", args, err)
		}
		if string(got) != want {
			t.Fatalf("%v: got %q, want %q", args, got, want)
		}
	}

	do("+PONG\r\n", "PING")
	do("+OK\r\n", "SET", "name", "Den")
	do("$3\r\nDen\r\n", "GET", "name")
	do(":1\r\n", "DEL", "name")
	do(":0\r\n", "DEL", "name")
	do("$-1\r\n", "GET", "name")
	do("-ERR wrong number of arguments for 'get'\r\n", "GET")
}
