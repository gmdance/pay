package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func Post(URL string, contentType string, raw []byte) ([]byte, error) {
	resp, err := http.Post(URL, contentType, bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func Get(URL string) ([]byte, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
