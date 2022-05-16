package rebel

func (client *Client) OnReady(handler func(*Client, *Ready)) {
	client.onReadyFunction = handler
}

func (client *Client) OnMessage(handler func(*Client, *Message)) {
	client.onMessageFunction = handler
}

func (client *Client) OnMessageUpdate(handler func(*Client, *MessageUpdate)) {
	client.onMessageUpdateFunction = handler
}
