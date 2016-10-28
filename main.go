package main

import (
	"fmt"
	"net"
	"bufio"
	"os"
)

var conn []net.Conn;

func main(){
	args := os.Args[1:];
	if(len(args) < 1){
		fmt.Println("Missing port argument");
		return;
	}

	fmt.Println("Starting server...");
	l, _ := net.Listen("tcp", ":" + args[0]);
	
	for{
		c, _ := l.Accept();
		conn = append(conn, c);
		go handle(c);
	}
}

func handle(c net.Conn){
	fmt.Println("New client!")

	read := bufio.NewReader(c);
	for{
		bytes, err := read.ReadBytes('\n');
		if err != nil { break; }
		
		for _, c1 := range conn{
			if c1 != c{
				c1.Write(bytes);
			}
		}
	}
	fmt.Println("Client disconnected")
	for i, c1 := range conn{
		if c1 == c{
			conn = append(conn[:i], conn[i+1:]...);
		}
	}
}