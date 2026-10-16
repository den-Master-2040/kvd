package server

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"strings"

	"github.com/den-Master-2040/kvd/internal/resp"
	"github.com/den-Master-2040/kvd/internal/store"
)

func Serve(l net.Listener, st *store.Store) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go handle(conn, st)
	}
}

func handle(conn net.Conn, st *store.Store) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	for {
		args, err := resp.ReadCommand(r)
		if err != nil {
			if errors.Is(err, resp.ErrProtocol) {
				resp.WriteError(w, "protocol error")
				w.Flush()
			} else if err != io.EOF {
				log.Println("read:", err)
			}
			return
		}

		switch strings.ToUpper(string(args[0])) {
		case "PING":
			resp.WriteSimple(w, "PONG")
		case "SET":
			if len(args) != 3 {
				resp.WriteError(w, "wrong number of arguments for 'set'")
				break
			}
			st.Set(string(args[1]), args[2])
			resp.WriteSimple(w, "OK")
		case "GET":
			if len(args) != 2 {
				resp.WriteError(w, "wrong number of arguments for 'get'")
				break
			}
			if val, ok := st.Get(string(args[1])); ok {
				resp.WriteBulk(w, val)
			} else {
				resp.WriteNull(w)
			}
		case "DEL":
			if len(args) != 2 {
				resp.WriteError(w, "wrong number of arguments for 'del'")
				break
			}
			n := 0
			if st.Delete(string(args[1])) {
				n = 1
			}
			resp.WriteInt(w, n)
		default:
			resp.WriteError(w, "unknown command '"+string(args[0])+"'")
		}

		if r.Buffered() == 0 {
			if err := w.Flush(); err != nil {
				return
			}
		}
	}
}
