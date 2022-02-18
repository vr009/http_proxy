package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Launching server...")

	// Устанавливаем прослушивание порта
	ln, _ := net.Listen("tcp", ":8080")
	// Открываем порт
	conn, _ := ln.Accept()

	var bmessage []byte
	var bbody []byte

	for {
		bytes := make([]byte, 10, 10)
		n, err := conn.Read(bytes)
		if err != nil || n == 0 {
			fmt.Print(err)
			fmt.Print(n)
			return
		}
		bmessage = append(bmessage, bytes...)
		if strings.Contains(string(bmessage), "\r\n\r\n") {
			i := strings.LastIndex(string(bmessage), "\r\n\r\n") + len("\r\n\r\n")
			bbody = append(bbody, bmessage[i:len(bmessage)]...)
			bmessage = bmessage[:i]
			break
		}
	}

	tail := len(bbody)

	fmt.Print("Received headers:\n", string(bmessage))
	i := strings.LastIndex(string(bmessage), "Content-Length:")
	i = i + len("Content-Length:") + 1
	j := strings.Index(string(bmessage)[i:], "\r")
	l, _ := strconv.Atoi(string(bmessage)[i : i+j])

	fmt.Print("start receiving body\n")
	for {
		bytes := make([]byte, l, l)
		n, err := conn.Read(bytes)
		if err != nil || n == 0 {
			fmt.Print(err)
			fmt.Print(n)
			return
		}
		bbody = append(bbody, bytes...)
		if n > l-tail-1 {
			break
		}
	}
	fmt.Print("Received body:\n", string(bbody))

	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
}
