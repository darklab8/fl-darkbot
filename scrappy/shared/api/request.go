package api

import (
	"darkbot/utils/logger"
	"io/ioutil"
	"net/http"
)

type APIrequest struct {
	url string
}

func (a APIrequest) GetData() ([]byte, error) {
	resp, err := http.Get(a.url)
	logger.CheckWarn(err, "unable to get url")
	if err != nil {
		return []byte{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	logger.CheckWarn(err, "unable to read base body")
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func (a *APIrequest) Init(url string) {
	a.url = url
}
