package main

import (
	"fmt"
	"net"
)

const (
	ServerHost = "localhost"
	ServerPort = "3333"
	ServerType = "tcp"
)

type User struct {
	Conn net.Conn
}

func main() {
	fmt.Printf("Starting %v server\n", ServerType)
	l, err := net.Listen(ServerType, net.JoinHostPort(ServerHost, ServerPort))
	if err != nil {
		panic(err)
	}
	defer l.Close()
	fmt.Printf("%v server started on port '%v'\n", ServerType, ServerPort)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("error accepting connection: " + err.Error())
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	u := User{Conn: conn}
	u.Conn.Write([]byte("Welcome\n"))

	buf := make([]byte, 1024)
	_, err := u.Conn.Read(buf)
	if err != nil {
		fmt.Println("error reading input:", err.Error())
	}

	u.Conn.Write([]byte("Message received.\n"))
	u.Conn.Close()
}
