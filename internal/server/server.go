package server

import (
	"bufio"
	"fmt"
	"net"
)

func Serve(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			fmt.Println("read error:", err)
			return
		}
		fmt.Print(msg)
		conn.Write([]byte("+PONG\r\n"))
	}
}
