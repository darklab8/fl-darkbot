package api

import (
	"darkbot/utils/logger"
	"io/ioutil"
	"net/http"
)

type APIrequest struct {
	url string
}

func (a APIrequest) GetData() []byte {
	resp, err := http.Get(a.url)
	logger.CheckWarn(err, "unable to get url")
	body, err := ioutil.ReadAll(resp.Body)
	logger.CheckWarn(err, "unable to read base body")
	return body
}

func (a *APIrequest) Init(url string) {
	a.url = url
}
