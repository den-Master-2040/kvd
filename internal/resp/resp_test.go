package resp

import (
	"bufio"
	"strings"
	"testing"
)

func TestReadCommand(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{"ping", "*1\r\n$4\r\nPING\r\n", []string{"PING"}, false},
		{"set", "*3\r\n$3\r\nSET\r\n$1\r\nk\r\n$1\r\nv\r\n", []string{"SET", "k", "v"}, false},
		{"crlf in value", "*1\r\n$4\r\na\r\nb\r\n", []string{"a\r\nb"}, false},
		{"truncated", "*2\r\n$3\r\nGET\r\n", nil, true},
		{"not array", "PING\r\n", nil, true},
		{"negative len", "*1\r\n$-5\r\n", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadCommand(bufio.NewReader(strings.NewReader(tt.input)))
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d args, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if string(got[i]) != tt.want[i] {
					t.Errorf("arg %d = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}