package forumer

import (
	"io"
	"net/http"

	"github.com/darklab8/fl-darkbot/app/forumer/forum_types"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/go-utils/utils/utils_settings"

	"golang.org/x/net/html/charset"
)

type MethodType string

const (
	GET MethodType = "GET"
)

type QueryResult struct {
	content          string
	ResponseRawQuery string
	ResponseFullUrl  string
}

func (q *QueryResult) GetContent() string {
	return q.content
}

func NewQuery(method_type MethodType, url forum_types.Url) (*QueryResult, error) {
	client := &http.Client{}
	req, err := http.NewRequest(string(method_type), string(url), nil)
	if logus.Log.CheckWarn(err, "Failed to create request") {
		return nil, err
	}
	if utils_settings.Envs.UserAgent != "" {
		req.Header.Set("User-Agent", utils_settings.Envs.UserAgent)
	} else {
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	}

	// req.Header.Set("User-Agent", "curl/7.81.0")
	req.Header.Set("ACCEPT", "*/*")
	req.Header.Set("CONTENT-LENGTH", "")
	req.Header.Set("CONTENT-TYPE", "")

	resp, err := client.Do(req)
	if logus.Log.CheckWarn(err, "Failed to make query") {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 && resp.StatusCode <= 399 {
		logus.Log.Debug("this request is redirecting!")
		redirectUrl, err := resp.Location()
		if logus.Log.CheckError(err, "Error getting redirect location") {
			return nil, err
		}

		req.URL = redirectUrl
		resp, err = client.Do(req)
		if logus.Log.CheckError(err, "Error sending redirect request:") {
			return nil, err
		}

	}

	utf8Body, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(utf8Body)

	return &QueryResult{
		content:          string(bytes),
		ResponseRawQuery: resp.Request.URL.RawQuery,
		ResponseFullUrl:  resp.Request.URL.String(),
	}, err
}
