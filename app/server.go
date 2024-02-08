package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var store = make(map[string]string)

func handleConnection(c net.Conn) {
	defer c.Close()
	sc := bufio.NewScanner(c)

	for sc.Scan() {
		switch command := strings.ToLower(sc.Text()); command {
		case "ping":
			c.Write([]byte("+PONG\r\n"))
		case "echo":
			sc.Scan()
			n := sc.Text()
			sc.Scan()
			s := sc.Text()
			c.Write([]byte(fmt.Sprintf("%s\r\n%s\r\n", n, s)))
		case "set":
			sc.Scan()
			_ = sc.Text()
			sc.Scan()
			k := sc.Text()
			sc.Scan()
			_ = sc.Text()
			sc.Scan()
			v := sc.Text()
			store[k] = v
			c.Write([]byte("+VALUE SET\r\n"))
		case "get":
			sc.Scan()
			_ = sc.Text()
			sc.Scan()
			k := sc.Text()
			fmt.Println(k)
			if v, found := store[k]; found {
				fmt.Println(v)
				c.Write([]byte(fmt.Sprintf("$%s\r\n%s\r\n", fmt.Sprint(len(v)), v)))
			}
			c.Write([]byte("+NOT FOUND\r\n"))
		}
	}
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
		go handleConnection(conn)
	}

}
