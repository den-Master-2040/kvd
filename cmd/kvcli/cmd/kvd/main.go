package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:6380")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	in := bufio.NewScanner(os.Stdin)
	reply := bufio.NewReader(conn)
	for in.Scan() {
		line := in.Text()
		if line == "exit" {
			break
		}
		if _, err := conn.Write([]byte(line + "\r\n")); err != nil {
			fmt.Println("write error:", err)
			return
		}
		resp, err := reply.ReadString('\n')
		if err != nil {
			fmt.Println("read error:", err)
			return
		}
		fmt.Print(resp)
	}
	fmt.Println("Done")
}
