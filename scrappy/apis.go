package scrappy

import (
	"darkbot/settings"
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

type BasesAPI struct {
	API
}

func (a BasesAPI) New() APIinterface {
	a.url = settings.Config.Scrappy.Base.URL
	return a
}

type PlayerAPI struct {
	API
}

func (a PlayerAPI) New() APIinterface {
	a.url = settings.Config.Scrappy.Player.URL
	return a
}
