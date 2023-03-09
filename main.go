package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"time"
)

type publicIpResponse struct {
	PublicIp string `json:"public_ip"`
}

func getPublicIp(path string) (ip string, err error) {
	u, err := url.Parse(path)
	if err != nil {
		return
	}

	u.Path = "v1/publicip/ip"

	resp, err := http.Get(u.String())
	if err != nil {
		return
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	data := publicIpResponse{}
	json.Unmarshal(bytes, &data)

	ip = data.PublicIp
	if ip == "" {
		err = errors.New("no public ip")
	}
	return
}

func putQbittorrentPort(path string, username string, password string, port string) (err error) {
	u, err := url.Parse(path)
	if err != nil {
		return
	}
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

func main() {
	log.SetFlags(log.LstdFlags)

	gluetunPath := os.Getenv("GLUETUN_PATH")
	qbittorrentPort := os.Getenv("QBITTORRENT_PORT")
	qbittorrentPath := os.Getenv("QBITTORRENT_PATH")
	qbittorrentUsername := os.Getenv("QBITTORRENT_USERNAME")
	qbittorrentPassword := os.Getenv("QBITTORRENT_PASSWORD")
	t, err := strconv.Atoi(os.Getenv("TIMEOUT"))
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

		publicIp, err := getPublicIp(gluetunPath)
		if err != nil {
			log.Println(err)
			continue
		}

		err = queryPort(publicIp, qbittorrentPort)
		if err == nil {
			continue
		}
		log.Println(err)

		err = putQbittorrentPort(qbittorrentPath, qbittorrentUsername, qbittorrentPassword, "0")
		if err != nil {
			log.Println(err)
			continue
		}

		err = putQbittorrentPort(qbittorrentPath, qbittorrentUsername, qbittorrentPassword, qbittorrentPort)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
