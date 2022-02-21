package main

import (
	"crypto/tls"
	"log"
	"math/rand"
	"net"
	"os/exec"
	"path/filepath"
	"proxy/internal"
	"strconv"
	"time"
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

		port := internal.ParsePort(req.Host)
		if req.Secure {
			req.Host = req.Host[:len(req.Host)-4]
			conn.Write([]byte("HTTP/1.0 200 Connection established\r\nProxy-agent: Golang-Proxy\r\n\r\n"))
			path, _ := filepath.Abs("")
			err = exec.Command(path+"/gen_cert.sh", req.Host, strconv.Itoa(rand.Int())).Run()
			if err != nil {
				panic(err)
			}

			cert, err := tls.LoadX509KeyPair(path+"/hck.crt", path+"/cert.key")
			if err != nil {
				panic(err)
			}
			tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}}
			tlssrv := tls.Server(conn, tlsCfg)
			msg, err := internal.TlsReadMessage(tlssrv)
			if err != nil {
				panic(err)
			}
			str := string(msg)
			str = str + ""
			//clientCgg := tls.LoadX509KeyPair("","")
			tlsconnTo, err := tls.Dial("tcp", req.Host+":443", tlsCfg)
			if err != nil {
				panic(err)
			}

			bytesSent := 0
			for bytesSent < len(msg) {
				n, err := tlsconnTo.Write(msg)
				if err != nil {
					panic(err)
				}
				bytesSent += n
			}

			answer := []byte("")

			tlsconnTo.SetReadDeadline(time.Now().Add(time.Second * 5))
			for {
				bytesRec := make([]byte, 512, 512)
				n, err := tlsconnTo.Read(bytesRec)
				if n == 0 {
					break
				}
				if err != nil {
					panic(err)
				}
				answer = append(answer, bytesRec...)
			}

			answStr := string(answer)
			answStr = answStr + ""

			bytesSent = 0
			for bytesSent < len(answer) {
				n, err := tlssrv.Write(answer)
				if err != nil {
					panic(err)
				}
				bytesSent += n
			}

			tlsconnTo.Close()

		} else {
			if port != "" {
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
		}

		conn.Close()
		ln.Close()
		ln = nil
		conn = nil
	}
}
