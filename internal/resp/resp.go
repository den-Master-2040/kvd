package resp

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
)

const (
	maxArgs   = 1024
	maxArgLen = 1 << 20 // 1 МБ
)

var ErrProtocol = errors.New("resp: protocol error")

// readHeader читает строку вида "*3\r\n" или "$5\r\n" и возвращает число.
func readHeader(r *bufio.Reader, prefix byte) (int, error) {
	line, err := r.ReadSlice('\n') //не ReadLine, чтобы остановится когда накопится слишком много символов (>4096) в хедере
	if err != nil {
		return 0, err
	}
	if len(line) < 4 || line[0] != prefix || !bytes.HasSuffix(line, []byte("\r\n")) {
		return 0, ErrProtocol
	}
	if line[1] < 0x30 || line[1] > 0x39 {
		return 0, ErrProtocol
	}

	n, err := strconv.Atoi(string(line[1 : len(line)-2]))
	if err != nil || n < 0 {
		return 0, ErrProtocol
	}
	return n, nil
}
func ReadCommand(r *bufio.Reader) ([][]byte, error) {
	n, err := readHeader(r, '*') //т.е. сначала ожидаем самый главный хэдер с количеством всех аргументов
	if err != nil {
		return nil, err
	}
	if n == 0 || n > maxArgs {
		return nil, ErrProtocol
	}

	//делаем из этого значения двумерный массив байт, где каждый элемент - это отдельный аргумент команды. Например, для команды SET k v будет массив из 3 элементов: "SET", "k", "v"
	args := make([][]byte, 0, n) // длина 0, место под n элементов

	//одна итерация цикла - это один аргумент.
	for i := 0; i < n; i++ {
		//получаем аргумент
		size, err := readHeader(r, '$')
		if err != nil {
			return nil, err
		}
		if size > maxArgLen {
			return nil, ErrProtocol
		}
		//делаем массив под команду?
		buf := make([]byte, size+2) // данные и \r\n после них

		//просто читаем буфер? то есть у нас есть итератор? он читает кол-во байт размера len(buf)?
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, err
		}
		if buf[size] != '\r' || buf[size+1] != '\n' {
			return nil, ErrProtocol
		}
		//отрезаем хвостовые \r\n и добавляем в массив команд
		args = append(args, buf[:size]) // без хвостовых \r\n
	}
	return args, nil
}
