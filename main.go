package main

import (
	"fmt"
	"github.com/jD91mZM2/stdutil"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var conns []net.Conn

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Missing port argument")
		return
	}

	fmt.Println("Starting server...")
	l, err := net.Listen("tcp", ":"+args[0])
	if err != nil {
		stdutil.PrintErr("Couldn't connect", err)
		return
	}

	go func() {
		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		for _ = range c {
			fmt.Println("Shutting down...")
			for _, conn := range conns {
				err := conn.Close()
				if err != nil {
					stdutil.PrintErr("Could not close", err)
				}
			}
			l.Close()
			if err != nil {
				stdutil.PrintErr("Could not close", err)
			}
		}
	}()

	for {
		c, err := l.Accept()
		if err != nil {
			err1, ok := err.(*net.OpError)
			if ok && err1.Temporary() {
				stdutil.PrintErr("Couldn't accept", err)
				continue
			} else {
				break
			}
		}

		conns = append(conns, c)
		go handle(c)
	}
}

func handle(c net.Conn) {
	fmt.Println("New client!")

	writer := &ConnWriter{c}
	_, err := io.Copy(writer, c)
	if err != nil {
		stdutil.PrintErr("Error while copying", err)
	}

	fmt.Println("Client disconnected")
	for i, c1 := range conns {
		if c1 == c {
			conns = append(conns[:i], conns[i+1:]...)
			err := c1.Close()
			if err != nil {
				stdutil.PrintErr("Could not close", err)
			}
		}
	}
}

type ConnWriter struct {
	Conn net.Conn
}

func (w *ConnWriter) Write(bytes []byte) (int, error) {
	for _, c := range conns {
		if c == w.Conn {
			continue
		}

		n, err := c.Write(bytes)
		if err != nil {
			return n, err
		}
		if n != len(bytes) {
			return n, io.ErrShortWrite
		}
	}
	return len(bytes), nil
}
