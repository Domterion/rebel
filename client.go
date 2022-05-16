package rebel

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const PING_INTERVAL = 15

func New(token string) *Client {
	client := &Client{
		Token: token,
	}

	return client
}

func (client *Client) Open() error {
	var err error

	if client.websocket != nil {
		return errors.New("websocket connection already open")
	}

	gateway := fmt.Sprintf("wss://ws.revolt.chat?version=1&format=json&token=%s", client.Token)

	client.websocket, _, err = websocket.DefaultDialer.Dial(gateway, http.Header{})
	if err != nil {
		return err
	}

	_, bytes, err := client.websocket.ReadMessage()
	if err != nil {
		return err
	}

	var res ApiResponse
	json.Unmarshal(bytes, &res)

	if res.Type != "Authenticated" {
		message := fmt.Sprintf("expected Authenticated response, got %s", res.Type)
		return errors.New(message)
	}

	// Once we confirm we are authenticated we will start heartbeating and can start our event handling

	go client.ping()

	for {
		messageType, bytes, err := client.websocket.ReadMessage()
		if err != nil {
			return err
		}

		client.handleEvent(messageType, bytes)
	}

	return nil
}

func (client *Client) ping() {
	for {
		time.Sleep(PING_INTERVAL * time.Second)

		payload := Ping{
			ApiResponse: ApiResponse{
				Type: "Ping",
			},
			Data: 0,
		}

		client.websocket.WriteJSON(payload)
	}
}

func (client *Client) handleEvent(messageType int, bytes []byte) {
	var message ApiResponse
	json.Unmarshal(bytes, &message)

	fmt.Printf("handling event type=%s\n", message.Type)

	switch message.Type {
	case "Ready":
		var ready Ready
		json.Unmarshal(bytes, &ready)

		client.onReadyFunction(client, &ready)
	case "Message":
		var message Message
		json.Unmarshal(bytes, &message)

		client.onMessageFunction(client, &message)
	}
}
