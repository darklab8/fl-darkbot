package exposer

import (
	"darkbot/app/settings/logus"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
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
			logus.Debug("sleeping to acquire server pong")
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func testQuery(url string) (string, error) {
	resp, err := http.Get("http://localhost:8080" + url)
	logus.CheckError(err, "query failed")
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return string(body), err
}

func TestHomePage(t *testing.T) {
	FixtureTestWebServer()
	body, err := testQuery("/")
	logus.CheckFatal(err, "readAll failed")

	logus.Debug(body)

	if !strings.Contains(body, "Not found") {
		t.Error("")
	}
}

func TestPlayers(t *testing.T) {
	FixtureTestWebServer()
	body, err := testQuery("/players")
	logus.CheckFatal(err, "readAll failed")

	logus.Debug(body)

	if !strings.Contains(body, "Another page") {
		t.Error("")
	}
}
