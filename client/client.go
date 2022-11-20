package client

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"joystick-adapter/client/commands"
)

type ConnectionOptions struct {
	Url      string
	Protocol string
	Origin   string
}

type Connection struct {
	ws *websocket.Conn
}

// TODO: Rename Connect to ConnectWebsocket and add a ConnectTcp (Connection.ws should be a Writer or ReadWriter)
func Connect(options ConnectionOptions) (*Connection, error) {
	ws, err := websocket.Dial(options.Url, options.Protocol, options.Origin)
	if err != nil {
		return nil, err
	}
	conn := Connection{ws}
	return &conn, nil
}

type SendCommandMessage struct {
	Command string
	Data    json.RawMessage
}

func (c *Connection) SendCmd(command commands.Command) error {
	data, err := command.GetData()
	if err != nil {
		return err
	}
	msg := SendCommandMessage{
		Command: command.GetName(),
		Data:    data,
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = c.ws.Write(b)
	return err
}

func (c *Connection) Close() error {
	return c.ws.Close()
}
