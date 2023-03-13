package network

import (
	"errors"
	"log"
	"net"
	"strconv"
	"time"
)

func GetOutboundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "255.255.255.255:0")
	if err != nil {
		return
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = localAddr.IP.String()
	return
}

func QueryPort(ip string, port int, timeout time.Duration) (err error) {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, strconv.Itoa(port)), timeout)
	if err != nil {
		return
	}

	if conn == nil {
		return errors.New("port closed")
	}

	e := conn.Close()
	if e != nil {
		log.Println(e)
	}

	return
}
