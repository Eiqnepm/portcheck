package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"time"
)

func getOutboundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "255.255.255.255:0")
	if err != nil {
		return
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = localAddr.IP.String()
	return
}

func putQbittorrentPort(u url.URL, username string, password string, port string) (err error) {
	u.Path = "api/v2/auth/login"

	jar, err := cookiejar.New(nil)
	if err != nil {
		return
	}
	client := &http.Client{
		Jar: jar,
	}

	data := url.Values{
		"username": {username},
		"password": {password},
	}
	resp, err := client.PostForm(u.String(), data)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()
	var sid string
	for _, c := range cookies {
		if c.Name == "SID" {
			sid = c.Value
			break
		}
	}
	if sid == "" {
		return
	}

	u.Path = "api/v2/app/setPreferences"
	data = url.Values{
		"json": {fmt.Sprintf("{\"listen_port\":%s}", port)},
	}
	resp, err = client.PostForm(u.String(), data)
	if err != nil {
		return
	}
	log.Printf("port changed to %s\n", port)
	defer resp.Body.Close()

	u.Path = "api/v2/auth/logout"
	resp, err = client.Post(u.String(), "", nil)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	return
}

func queryPort(ip string, port string) (err error) {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), 5*time.Second)
	if err != nil {
		return
	}

	if conn != nil {
		defer conn.Close()
		return
	}

	return errors.New("port closed")
}

func env(key string, defaultValue string) (value string) {
	value = os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return
}

func main() {
	log.SetFlags(log.LstdFlags)

	qbitPort := env("QBITTORRENT_PORT", "6881")
	qbitWebPort := env("QBITTORRENT_WEBUI_PORT", "8080")
	qbitScheme := env("QBITTORRENT_WEBUI_SCHEME", "http")
	qbitUrl := url.URL{
		Scheme: qbitScheme,
		Host:   net.JoinHostPort("localhost", qbitWebPort),
	}
	qbitUsername := env("QBITTORRENT_USERNAME", "admin")
	qbitPassword := env("QBITTORRENT_PASSWORD", "adminadmin")
	t, err := strconv.Atoi(env("TIMEOUT", "300"))
	if err != nil {
		log.Fatal(err)
	}
	timeout := time.Duration(t) * time.Second

	firstLoop := true
	for {
		if !firstLoop {
			time.Sleep(timeout)
		}
		firstLoop = false

		outboundIp, err := getOutboundIP()
		if err != nil {
			log.Println(err)
			continue
		}

		err = queryPort(outboundIp, qbitPort)
		if err == nil {
			continue
		}
		log.Println(err)

		err = putQbittorrentPort(qbitUrl, qbitUsername, qbitPassword, "0")
		if err != nil {
			log.Println(err)
			continue
		}

		err = putQbittorrentPort(qbitUrl, qbitUsername, qbitPassword, qbitPort)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
