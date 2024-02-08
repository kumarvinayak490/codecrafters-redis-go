package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
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
			if sc.Scan() {
				sc.Scan()
				sc.Scan()
				sc.Scan()
				fmt.Println(sc.Text())
				et, err := strconv.Atoi(sc.Text())
				if err != nil {
					fmt.Println(err)
					c.Write([]byte("+ERROR\r\n"))
					break
				}
				timer := time.After(time.Duration(et) * time.Millisecond)
				go func() {
					<-timer
					delete(store, k)
				}()
			}
			c.Write([]byte("+OK\r\n"))
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
			c.Write([]byte("$-1\r\n"))
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
