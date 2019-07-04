package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func HttpPost(URL string, contentType string, rawBody []byte) ([]byte, error) {
	resp, err := http.Post(URL, contentType, bytes.NewReader(rawBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func HttpGet(URL string) ([]byte, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
