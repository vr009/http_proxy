package internal

import (
	"fmt"
	"net"
	"strings"
)

type ClientPart struct {
	conn net.Conn
}

func NewClientPart(conn net.Conn) *ClientPart {
	return &ClientPart{conn: conn}
}

func (cp *ClientPart) ClinetWork(msg []byte) []byte {
	n, err := cp.conn.Write(msg)
	if err != nil {
		fmt.Println("some error in sending data", n)
	}

	for {
		answer := make([]byte, 0, 10)
		bytes := make([]byte, 0, 0)

		_, err = cp.conn.Read(bytes)
		if err != nil {
			fmt.Print(err)
			fmt.Print(n)
			return nil
		}
		answer = append(answer, bytes...)
		if strings.Contains(string(answer), "\r\n\r\n") {
			return answer
		}
	}
}
