package models

import "github.com/gorilla/websocket"

type Client struct {
	ID          string
	Connection  *websocket.Conn
	IsConnected bool
	Name        string
	RoomId      string
}

func (client *Client) SetName(name string) {
	client.Name = name
}
func (client *Client) Close() {
	client.IsConnected = false
	client.Connection.Close()
}
func (client *Client) IsLogin() bool {
	return client.IsConnected && client.Name != ""
}
func (client *Client) SetRoomId(id string) {
	client.RoomId = id
}
func (client *Client) IsJoinedRoom() bool {
	return client.RoomId != ""
}
