package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Leela0o5/LeeGo/config"
	"github.com/Leela0o5/LeeGo/engine"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func TestEngineIntegration(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(echoHandler))
	defer s.Close()
	url := "ws" + strings.TrimPrefix(s.URL, "http")

	cfg := config.Config{
		URL:        url,
		NumWorkers: 2,
		Duration:   200 * time.Millisecond,
	}

	stats := engine.Run(cfg)
	if stats.TotalRequests == 0 {
		t.Error("Test failed: Engine recorded 0 total requests")
	}

	if stats.SuccessCount == 0 {
		t.Error("Test failed: Engine recorded 0 successful requests")
	}

	if len(stats.Latencies) == 0 {
		t.Error("Test failed: No latencies were recorded")
	}
}
