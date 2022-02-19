package main

import (
	"log"
	"net"
	"proxy/internal"
)

func main() {

	for {
		ln, err := net.Listen("tcp", ":8080")
		if err != nil {
			log.Println(err)
		}
		conn, _ := ln.Accept()

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

		conn.Close()
		ln.Close()
		ln = nil
		conn = nil
	}
}
