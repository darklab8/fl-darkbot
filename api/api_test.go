package api

import (
	"fmt"
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
		go newApi()
		for i := 0; i < 10; i++ {
			body := query("/ping")
			if body == "pong!" {
				break
			}
			fmt.Println("sleeping to acquire server pong")
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func query(url string) string {
	resp, err := http.Get("http://localhost:8080" + url)
	if err != nil {
		fmt.Println("query failed")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("readAll failed")
	}

	return string(body)
}

func TestHomePage(t *testing.T) {
	FixtureTestWebServer()
	body := query("/")

	fmt.Println(body)

	if !strings.Contains(body, "Not found") {
		t.Error("")
	}
}

func TestPlayers(t *testing.T) {
	FixtureTestWebServer()
	body := query("/players")

	fmt.Println(body)

	if !strings.Contains(body, "Another page") {
		t.Error("")
	}
}
