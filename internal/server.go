package internal

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type ServerPart struct {
	conn net.Conn
}

func NewServerPart(conn net.Conn) *ServerPart {
	return &ServerPart{conn: conn}
}

func ParseHost(headers []byte) string {
	i := strings.LastIndex(string(headers), "Host:")
	i = i + len("Host:") + 1
	j := strings.Index(string(headers)[i:], "\r")
	host := string(headers)[i : i+j]
	return host
}

func ParseLength(headers []byte) int {
	fmt.Print("Received headers:\n", string(headers))
	i := strings.LastIndex(string(headers), "Content-Length:")
	i = i + len("Content-Length:") + 1
	j := strings.Index(string(headers)[i:], "\r")
	l, _ := strconv.Atoi(string(headers)[i : i+j])
	return l
}

func (sp *ServerPart) SrvGetReq() *Req {
	req := &Req{}
	fullMsg := make([]byte, 0, 10)
	var bmessage []byte
	var bbody []byte

	for {
		bytes := make([]byte, 10, 10)
		n, err := sp.conn.Read(bytes)
		if err != nil || n == 0 {
			fmt.Print(err)
			fmt.Print(n)
			return nil
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

	l := ParseLength(bmessage)
	host := ParseHost(bmessage)

	fmt.Print("start receiving body\n")
	for {
		bytes := make([]byte, l, l)
		n, err := sp.conn.Read(bytes)
		if err != nil || n == 0 {
			fmt.Print(err)
			fmt.Print(n)
			return nil
		}
		bbody = append(bbody, bytes...)
		if n > l-tail-1 {
			break
		}
	}
	fmt.Print("Received body:\n", string(bbody))

	fullMsg = append(fullMsg, bmessage...)
	fullMsg = append(fullMsg, bbody...)
	req.FullMsg = fullMsg
	req.Host = host

	return req
}

func (sp *ServerPart) SrvSndRsp(answer []byte) {
	sp.conn.Write(answer)
}
