package api

import (
	"darkbot/settings/logus"
	"darkbot/settings/types"
	"io"
	"net/http"
)

type APIrequest struct {
}

func (a APIrequest) GetData(url types.APIurl) ([]byte, error) {
	resp, err := http.Get(string(url))
	logus.CheckWarn(err, "unable to get url")
	if err != nil {
		return []byte{}, err
	}
	body, err := io.ReadAll(resp.Body)
	logus.CheckWarn(err, "unable to read base body")
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}
