package internal

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParseHost(headers []byte) string {
	i := strings.LastIndex(string(headers), "Host:")
	i = i + len("Host:") + 1
	j := strings.Index(string(headers)[i:], "\r")
	host := string(headers)[i : i+j]
	return host
}

func ParseSecure(headers []byte) bool {
	i := strings.LastIndex(string(headers), "CONNECT")
	if i == -1 {
		return false
	}
	return true
}

func ParseURL(headers []byte) string {
	i := strings.LastIndex(string(headers), "http")
	j := strings.LastIndex(string(headers), "HTTP/")
	if i == -1 {
		return ""
	}
	host := string(headers)[i : j-1]
	return host
}

func ParsePort(url string) string {
	re := regexp.MustCompile(`(?m):\d+`)
	res := re.FindAllString(url, -1)
	if res == nil {
		return ""
	}
	return res[0]
}

func ParseLength(headers []byte) int {
	fmt.Print("Received headers:\n", string(headers))
	i := strings.LastIndex(string(headers), "Content-Length:")
	i = i + len("Content-Length:") + 1
	j := strings.Index(string(headers)[i:], "\r")
	l, _ := strconv.Atoi(string(headers)[i : i+j])
	return l
}

func GetRequest(conn net.Conn) *Req {
	req := &Req{
		Port: ":80",
	}
	fullMsg := make([]byte, 0, 10)
	var bmessage []byte
	var bbody []byte

	for {
		bytes := make([]byte, 10, 10)
		n, err := conn.Read(bytes)
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

	req.Secure = ParseSecure(bmessage)

	l := ParseLength(bmessage)
	host := ParseHost(bmessage)

	fmt.Print("start receiving body\n")
	for {
		if l == 0 {
			break
		}
		bytes := make([]byte, l, l)
		n, err := conn.Read(bytes)
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

func ReturnResponse(conn net.Conn, answer []byte) {
	size := len(answer)
	sentBytes := 0
	for sentBytes < size {
		n, err := conn.Write(answer)
		if err != nil {
			fmt.Println("some error in sending data", n)
		}
		sentBytes += n
	}
}

func ProxyRequest(conn net.Conn, msg []byte) []byte {
	size := len(msg)
	sentBytes := 0
	for sentBytes < size {
		n, err := conn.Write(msg)
		if err != nil {
			fmt.Println("some error in sending data", n)
		}
		sentBytes += n
	}
	answer := GetRequest(conn)
	return answer.FullMsg
}

func TlsReadMessage(conn net.Conn) ([]byte, error) {
	msg := []byte("")
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	for {
		bytes := make([]byte, 1024, 1024)
		n, err := conn.Read(bytes)
		if n == 0 {
			break
		}
		if err != nil {
			return nil, err
		}
		msg = append(msg, bytes...)
	}
	msg = append(msg, []byte("\r\n\r\n")...)
	return msg, nil
}

func TlsSendMessage(conn net.Conn, msg []byte) error {
	bytesSent := 0
	for bytesSent < len(msg) {
		n, err := conn.Write(msg)
		if err != nil {
			return err
		}
		bytesSent += n
	}
	return nil
}
