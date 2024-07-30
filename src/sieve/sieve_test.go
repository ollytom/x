package sieve

import (
	"fmt"
	"testing"
)

func TestDial(t *testing.T) {
	addr := fmt.Sprintf("imap.migadu.com:%d", DefaultPort)
	conn, err := Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()
}
