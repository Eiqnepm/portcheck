package main

import (
	"log"
	"net"
	u "net/url"
	"os"
	"strconv"
	"time"

	"github.com/Eiqnepm/portcheck/internal/network"
	"github.com/Eiqnepm/portcheck/internal/qbit"
)

func env(key string, defaultValue string) (value string) {
	value = os.Getenv(key)
	if value == "" {
		value = defaultValue
	}

	return
}

func main() {
	log.SetFlags(log.LstdFlags)

	qbitPort, err := strconv.Atoi(env("QBITTORRENT_PORT", "6881"))
	if err != nil {
		log.Fatal(err)
	}
	qbitWebPort := env("QBITTORRENT_WEBUI_PORT", "8080")
	qbitScheme := env("QBITTORRENT_WEBUI_SCHEME", "http")
	qbitUrl := u.URL{
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
	t, err = strconv.Atoi(env("DIAL_TIMEOUT", "5"))
	if err != nil {
		log.Fatal(err)
	}
	dialTimeout := time.Duration(t) * time.Second

	firstLoop := true
	for {
		if !firstLoop {
			time.Sleep(timeout)
		}
		firstLoop = false

		outboundIp, err := network.GetOutboundIP()
		if err != nil {
			log.Println(err)
			continue
		}

		err = network.QueryPort(outboundIp, qbitPort, dialTimeout)
		if err == nil {
			continue
		}
		log.Println(err)

		func() {
			err = qbit.Login(qbitUrl, qbitUsername, qbitPassword)
			if err != nil {
				log.Println(err)
				return
			}
			defer func() {
				err := qbit.Logout()
				if err != nil {
					log.Println(err)
				}
			}()

			err = qbit.SetPreference("listen_port", 0)
			if err != nil {
				log.Println(err)
				return
			}

			err = qbit.SetPreference("listen_port", qbitPort)
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}
}
