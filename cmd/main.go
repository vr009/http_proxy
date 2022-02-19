package main

import (
	"net"
	"proxy/internal"
)

func main() {

	ln, _ := net.Listen("tcp", ":8080")
	conn, _ := ln.Accept()
	defer conn.Close()

	/* getting request from client */
	req := internal.GetRequest(conn)

	if internal.ParsePort(req.Host) != "" {
		req.Port = ""
	}
	connTo, err := net.Dial("tcp", req.Host+req.Port)
	if err != nil {
		panic(err)
	}

	/* proxying the message from client to dest */
	answer := internal.ProxyRequest(connTo, req.FullMsg)

	/* returning response */
	internal.ReturnResponse(conn, answer)
	connTo.Close()
}
