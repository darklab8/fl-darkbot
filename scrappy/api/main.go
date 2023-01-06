/*
Reusable API code to request []byte code of smth. Reusable for player and base.
*/

package api

import (
	"darkbot/utils"
	"io/ioutil"
	"net/http"
)

type APIinterface interface {
	GetData() []byte
	New() APIinterface
}

type API struct {
	url string
}

func (a API) GetData() []byte {
	resp, err := http.Get(a.url)
	utils.CheckWarn(err, "unable to get url")
	body, err := ioutil.ReadAll(resp.Body)
	utils.CheckWarn(err, "unable to read base body")
	return body
}

func (a API) Init(url string) {
	a.url = url
}
