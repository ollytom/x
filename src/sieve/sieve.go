// Package sieve provides a client of the ManageSieve protocol
// specified in RFC 5804.
package sieve

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/textproto"
	"os"
	"strings"
)

const DefaultPort int = 4190

func Dial(net, addr string) (*textproto.Conn, error) {
	conn, err := textproto.Dial(net, addr)
	if err != nil {
		return nil, err
	}
	for i := 0; i <= 10; i++ {
		line, err := conn.ReadLine()
		if err != nil {
			return nil, err
		}
		fmt.Fprintln(os.Stderr, line)
		if strings.HasPrefix(line, "OK") {
			break
		}
	}
	return conn, nil
}

// Logout sends the LOGOUT command and closes conn.
// Implementations should not use conn afterwards.
func Logout(conn *textproto.Conn) error {
	id, err := conn.Cmd("LOGOUT")
	if err != nil {
		return err
	}
	conn.StartResponse(id)
	defer conn.EndResponse(id)
	line, err := conn.ReadLine()
	if err != nil {
		return err
	}
	code, msg, found := strings.Cut(line, " ")
	if code != "OK" {
		if !found {
			return fmt.Errorf("logout failed with no message")
		}
		return errors.New(msg)
	}
	return conn.Close()
}

func Noop(conn *textproto.Conn) error {
	id, err := conn.Cmd("NOOP")
	if err != nil {
		return err
	}
	conn.StartResponse(id)
	defer conn.EndResponse(id)
	line, err := conn.ReadLine()
	if err != nil {
		return err
	}
	code, msg, found := strings.Cut(line, " ")
	if code != "OK" {
		if !found {
			return fmt.Errorf("noop failed with no message")
		}
		return errors.New(msg)
	}
	return nil
}

func StartTLS(conn *textproto.Conn, config *tls.Config) error {
	id, err := conn.Cmd("STARTTLS")
	if err != nil {
		return err
	}
	conn.StartResponse(id)
	defer conn.EndResponse(id)
	line, err := conn.ReadLine()
	if err != nil {
		return err
	}
	code, msg, found := strings.Cut(line, " ")
	if code != "OK" {
		if !found {
			return fmt.Errorf("starttls failed with no message")
		}
		return errors.New(msg)
	}
	return errors.New("TODO not yet implemented")
}
