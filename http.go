package rebel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (client *Client) SendMessage(channel string, message *Message) (*Message, error) {
	route := fmt.Sprintf("/channels/%s/messages", channel)

	if message.Nonce == "" {
		nonce := GenerateNonce()

		message.Nonce = nonce
	}

	marshalled, err := json.Marshal(message)

	if err != nil {
		return nil, err
	}

	resp, err := client.Request("POST", route, marshalled)

	if err != nil {
		return nil, err
	}

	var res Message
	json.Unmarshal(resp, &message)

	return &res, nil
}

func (client *Client) SendMessageContent(channel string, message string) (*Message, error) {
	return client.SendMessage(channel, &Message{
		Content: message,
	})
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
