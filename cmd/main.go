package main

import (
	"fmt"
	"net"
	"proxy/internal"
)

func main() {
	fmt.Println("Launching server...")

	// Устанавливаем прослушивание порта
	ln, _ := net.Listen("tcp", ":8080")
	// Открываем порт
	conn, _ := ln.Accept()

	serverPart := internal.NewServerPart(conn)
	req := serverPart.SrvGetReq()

	connTo, _ := net.Dial("tcp", req.Host)
	clientPart := internal.NewClientPart(connTo)
	answer := clientPart.ClinetWork(req.FullMsg)

	serverPart.SrvSndRsp(answer)
}
