package api

import (
	"io"
	"net/http"

	"github.com/darklab/fl-darkbot/app/settings/logus"
	"github.com/darklab/fl-darkbot/app/settings/types"
)

type APIrequest struct {
}

func (a APIrequest) GetData(url types.APIurl) ([]byte, error) {
	resp, err := http.Get(string(url))
	if logus.Log.CheckWarn(err, "unable to get url") {
		return []byte{}, err
	}
	body, err := io.ReadAll(resp.Body)
	if logus.Log.CheckWarn(err, "unable to read base body") {
		return []byte{}, err
	}
	return body, nil
}
