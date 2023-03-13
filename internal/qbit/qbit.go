package qbit

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	u "net/url"
)

type session struct {
	url    u.URL
	client *http.Client
}

var s session

func Login(url u.URL, username string, password string) (err error) {
	s.url = url
	s.url.Path = "api/v2/auth/login"
	jar, err := cookiejar.New(nil)
	if err != nil {
		return
	}

	s.client = &http.Client{
		Jar: jar,
	}
	data := u.Values{
		"username": {username},
		"password": {password},
	}
	resp, err := s.client.PostForm(s.url.String(), data)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode == 200 {
		return
	}

	return errors.New(resp.Status)
}

func Logout() (err error) {
	s.url.Path = "api/v2/auth/logout"
	resp, err := s.client.Post(s.url.String(), "", nil)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode == 200 {
		return
	}

	return errors.New(resp.Status)
}

func SetPreference(preference string, value any) (err error) {
	s.url.Path = "api/v2/app/setPreferences"
	j, err := json.Marshal(map[string]any{preference: value})
	if err != nil {
		return
	}

	data := u.Values{
		"json": {string(j)},
	}
	resp, err := s.client.PostForm(s.url.String(), data)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode == 200 {
        log.Println(string(j))
		return
	}

	return errors.New(resp.Status)
}
