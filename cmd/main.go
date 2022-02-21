package main

import (
	"bytes"
	"crypto/tls"
	"log"
	"math/rand"
	"net"
	"os/exec"
	"path/filepath"
	"proxy/internal"
	"strconv"
	"strings"
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
		ss := string(req.FullMsg)
		ss += ""

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

			tlsconnTo, err := tls.Dial("tcp", req.Host+":443", tlsCfg)
			if err != nil {
				panic(err)
			}

			err = internal.TlsSendMessage(tlsconnTo, msg)
			if err != nil {
				panic(err)
			}

			answer, err := internal.TlsReadMessage(tlsconnTo)
			if err != nil {
				panic(err)
			}

			/* To avoid the issue described in this ticket: https://github.com/curl/curl/issues/6760 */
			if strings.LastIndex(string(answer), "Transfer-Encoding: chunked") != -1 {
				answer = bytes.Replace(answer, []byte("Transfer-Encoding: chunked"), []byte(""), -1)
			}

			err = internal.TlsSendMessage(tlssrv, answer)
			if err != nil {
				panic(err)
			}

			tlsconnTo.Close()
			tlssrv.Close()

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
