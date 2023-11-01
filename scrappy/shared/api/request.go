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
	if logus.CheckWarn(err, "unable to get url") {
		return []byte{}, err
	}
	body, err := io.ReadAll(resp.Body)
	if logus.CheckWarn(err, "unable to read base body") {
		return []byte{}, err
	}
	return body, nil
}
