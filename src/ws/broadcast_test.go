package ws

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWebSocketScoreBroadcast(t *testing.T) {
	// Reset global hub for test
	hub = NewHub()
	go hub.Run()
	time.Sleep(50 * time.Millisecond)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WebSocketHandler(w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	defer conn.Close()

	time.Sleep(100 * time.Millisecond)

	msg := map[string]interface{}{
		"event": "score_update",
		"data": map[string]interface{}{
			"teams": []map[string]interface{}{
				{"rank": 1, "name": "Team A", "score": 500},
			},
		},
	}
	msgBytes, _ := json.Marshal(msg)

	hub.Broadcast(msgBytes)

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, received, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(received, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if result["event"] != "score_update" {
		t.Errorf("expected event score_update, got %v", result["event"])
	}
}

func TestWebSocketChallengeSolvedBroadcast(t *testing.T) {
	hub = NewHub()
	go hub.Run()
	time.Sleep(50 * time.Millisecond)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WebSocketHandler(w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	defer conn.Close()

	time.Sleep(100 * time.Millisecond)

	msg := map[string]interface{}{
		"event": "challenge_solved",
		"data": map[string]interface{}{
			"team_name":      "Team B",
			"challenge_name": "Web 101",
			"points":         100,
		},
	}
	msgBytes, _ := json.Marshal(msg)

	hub.Broadcast(msgBytes)

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, received, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(received, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if result["event"] != "challenge_solved" {
		t.Errorf("expected event challenge_solved, got %v", result["event"])
	}
}

func TestWebSocketMultipleClientsBroadcast(t *testing.T) {
	hub = NewHub()
	go hub.Run()
	time.Sleep(50 * time.Millisecond)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WebSocketHandler(w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	conns := make([]*websocket.Conn, 3)
	for i := 0; i < 3; i++ {
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("dial %d failed: %v", i, err)
		}
		defer conn.Close()
		conns[i] = conn
	}

	time.Sleep(100 * time.Millisecond)

	msg := []byte(`{"event":"ping","data":{}}`)
	hub.Broadcast(msg)

	for i, conn := range conns {
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, received, err := conn.ReadMessage()
		if err != nil {
			t.Errorf("client %d read failed: %v", i, err)
			continue
		}
		if string(received) != string(msg) {
			t.Errorf("client %d got wrong message: %s", i, received)
		}
	}
}

func TestBroadcastScoreUpdate(t *testing.T) {
	hub = NewHub()
	go hub.Run()
	time.Sleep(50 * time.Millisecond)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WebSocketHandler(w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	defer conn.Close()

	time.Sleep(100 * time.Millisecond)

	data := []byte(`{"event":"score_update","data":{"teams":[]}}`)
	BroadcastScoreUpdate(data)

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, received, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	if string(received) != string(data) {
		t.Errorf("expected %s, got %s", data, received)
	}
}

func TestWebSocketClientCount(t *testing.T) {
	hub = NewHub()
	go hub.Run()
	time.Sleep(50 * time.Millisecond)

	if hub.ClientCount() != 0 {
		t.Error("should start with 0 clients")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WebSocketHandler(w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	conn1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial 1 failed: %v", err)
	}
	defer conn1.Close()

	conn2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial 2 failed: %v", err)
	}
	defer conn2.Close()

	time.Sleep(100 * time.Millisecond)

	if hub.ClientCount() != 2 {
		t.Errorf("expected 2 clients, got %d", hub.ClientCount())
	}
}
