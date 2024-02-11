package exposer

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/darklab/fl-darkbot/app/settings/logus"
)

var (
	is_webserver_active bool = false
)

func FixtureTestWebServer() {
	if !is_webserver_active {
		is_webserver_active = true
		go NewExposer()
		for i := 0; i < 100; i++ {
			body, _ := testQuery("/ping")
			if body == "pong!" {
				break
			}
			logus.Log.Debug("sleeping to acquire server pong")
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func testQuery(url string) (string, error) {
	resp, err := http.Get("http://localhost:8080" + url)
	if logus.Log.CheckError(err, "query failed") {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return string(body), err
}

func TestHomePage(t *testing.T) {
	FixtureTestWebServer()
	body, err := testQuery("/")
	logus.Log.CheckFatal(err, "readAll failed")

	logus.Log.Debug(body)

	if !strings.Contains(body, "Not found") {
		t.Error("")
	}
}

func TestPlayers(t *testing.T) {
	FixtureTestWebServer()
	body, err := testQuery("/players")
	logus.Log.CheckFatal(err, "readAll failed")

	logus.Log.Debug(body)

	if !strings.Contains(body, "Another page") {
		t.Error("")
	}
}
