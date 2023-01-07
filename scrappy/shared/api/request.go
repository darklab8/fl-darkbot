package api

import (
	"darkbot/utils"
	"io/ioutil"
	"net/http"
)

type APIrequest struct {
	url string
}

func (a APIrequest) GetData() []byte {
	resp, err := http.Get(a.url)
	utils.CheckWarn(err, "unable to get url")
	body, err := ioutil.ReadAll(resp.Body)
	utils.CheckWarn(err, "unable to read base body")
	return body
}

func (a APIrequest) Init(url string) {
	a.url = url
}
