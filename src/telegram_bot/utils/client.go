// Copyright 2020 aaaaaaaalesha <sks2311211@mail.ru>

package utils

import (
	"io"
	"net/http"
	"time"
)

var DefaultServiceUrl = "http://todoapp:8080/"

var client = http.Client{
	Timeout: time.Second * 15,
}

func Put(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func Delete(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
