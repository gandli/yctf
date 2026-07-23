package ws

import (
	"testing"
	"time"
)

func TestNewHub(t *testing.T) {
	hub := NewHub()
	if hub == nil {
		t.Error("hub should not be nil")
	}
	if hub.clients == nil {
		t.Error("clients map should be initialized")
	}
	if hub.broadcast == nil {
		t.Error("broadcast channel should be initialized")
	}
}

func TestHubClientCount(t *testing.T) {
	hub := NewHub()
	if hub.ClientCount() != 0 {
		t.Errorf("expected 0 clients, got %d", hub.ClientCount())
	}
}

func TestHubRegisterUnregister(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	time.Sleep(10 * time.Millisecond)

	client := &Client{hub: hub, send: make(chan []byte, 256)}
	hub.Register(client)
	time.Sleep(50 * time.Millisecond)

	if hub.ClientCount() != 1 {
		t.Errorf("expected 1 client, got %d", hub.ClientCount())
	}

	hub.Unregister(client)
	time.Sleep(50 * time.Millisecond)

	if hub.ClientCount() != 0 {
		t.Errorf("expected 0 clients, got %d", hub.ClientCount())
	}
}

func TestHubBroadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	time.Sleep(10 * time.Millisecond)

	client := &Client{hub: hub, send: make(chan []byte, 256)}
	hub.Register(client)
	time.Sleep(50 * time.Millisecond)

	message := []byte(`{"event":"test"}`)
	hub.broadcast <- message
	time.Sleep(50 * time.Millisecond)

	select {
	case msg := <-client.send:
		if string(msg) != string(message) {
			t.Errorf("got wrong message: %s", msg)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("client did not receive message")
	}
}
