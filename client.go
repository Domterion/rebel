package rebel

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const PING_INTERVAL = 15
const API_URL = "https://api.revolt.chat/"

func New(token string) *Client {
	client := &Client{
		Token: token,
		http:  &http.Client{},
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
		if client.onReadyFunction == nil {
			return
		}

		var ready Ready
		json.Unmarshal(bytes, &ready)

		client.onReadyFunction(client, &ready)

		// bytes, _ := client.Request("GET", "/users/@me", []byte{})
		// println(string(bytes))
	case "Message":
		if client.onMessageFunction == nil {
			return
		}

		var message Message
		json.Unmarshal(bytes, &message)

		client.onMessageFunction(client, &message)
	case "MessageUpdate":
		if client.onMessageUpdateFunction == nil {
			return
		}

		var message MessageUpdate
		json.Unmarshal(bytes, &message)

		client.onMessageUpdateFunction(client, &message)
	}
}

func (client *Client) Request(method string, route string, data []byte) ([]byte, error) {
	route = fmt.Sprintf("%s%s", API_URL, route)
	reader := bytes.NewReader(data)

	req, err := http.NewRequest(method, route, reader)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Add("x-bot-token", client.Token)

	resp, err := client.http.Do(req)
	if err != nil {
		return []byte{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
