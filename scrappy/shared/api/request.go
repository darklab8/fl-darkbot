package api

import (
	"darkbot/settings/utils/logger"
	"io"
	"net/http"
)

type APIrequest struct {
}

type APIurl string

func (a APIrequest) GetData(url APIurl) ([]byte, error) {
	resp, err := http.Get(string(url))
	logger.CheckWarn(err, "unable to get url")
	if err != nil {
		return []byte{}, err
	}
	body, err := io.ReadAll(resp.Body)
	logger.CheckWarn(err, "unable to read base body")
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}
