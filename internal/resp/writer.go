package resp

import (
	"bufio"
	"strconv"
)

func WriteSimple(w *bufio.Writer, b string) error {
	w.WriteByte('+')
	w.WriteString(b)
	_, err := w.WriteString("\r\n")
	return err
}

func WriteError(w *bufio.Writer, b string) error {
	w.WriteString("-ERR ")
	w.WriteString(b)
	_, err := w.WriteString("\r\n")
	return err
}

func WriteBulk(w *bufio.Writer, b []byte) error {
	w.WriteByte('$')
	w.WriteString(strconv.Itoa(len(b)))
	w.WriteString("\r\n")
	w.Write(b)
	_, err := w.WriteString("\r\n")
	return err
}

func WriteNull(w *bufio.Writer) error {
	_, err := w.WriteString("$-1\r\n")
	return err
}

func WriteInt(w *bufio.Writer, b int) error {
	w.WriteByte(':')
	w.WriteString(strconv.Itoa(b))
	_, err := w.WriteString("\r\n")
	return err
}
