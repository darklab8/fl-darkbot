/*
Reusable API code to request []byte code of smth. Reusable for player and base.
*/

package api

import (
	"io"

	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"
	"github.com/darklab8/go-utils/utils/utils_http"
)

type APIrequest struct {
}

func (a APIrequest) GetData(url types.APIurl) ([]byte, error) {
	resp, err := utils_http.Get(string(url))
	if logus.Log.CheckWarn(err, "unable to get url") {
		return []byte{}, err
	}
	body, err := io.ReadAll(resp.Body)
	if logus.Log.CheckWarn(err, "unable to read base body") {
		return []byte{}, err
	}
	return body, nil
}
