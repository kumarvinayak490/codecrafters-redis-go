package main

import (
	"fmt"

	"net"
	"os"
)

func handleConnection(c net.Conn) error {
	defer c.Close()
	buf := make([]byte, 1024)
	_, err := c.Read(buf)
	if err != nil {
		return err
	}
	_, err = c.Write([]byte("+PONG\r\n"))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		return err
	}
	return nil
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	defer l.Close()
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			break
		}
		err = handleConnection(conn)
		if err != nil {
			break
		}
	}

}
