package models

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type Room struct {
	ID           string
	Clients      []*Client
	Code         int
	Name         string
	AdminID      string
	IsPartyStart bool
}

func (room *Room) AddClient(client Client) {
	room.Clients = append(room.Clients, &client)
}
func (room *Room) RemoveClient(client *Client) {
	var updatedClient []*Client

	for _, c := range room.Clients {
		if c.ID != client.ID {
			updatedClient = append(updatedClient, c)
		}
	}
	room.Clients = updatedClient
}

func (room *Room) StartParty() {
	room.IsPartyStart = true
	println("Start file buffering...")
	file, err := os.Open("./static/song.mp3")
	if err != nil {
		panic(err)
	}
	buffer := make([]byte, 1024*32)
	for {
		// Read a chunk from the file
		n, err := file.Read(buffer)
		if err == io.EOF {
			log.Println("Finish buffering.", err)
			room.IsPartyStart = false
			break
		}

		for _, client := range room.Clients {
			if client.IsConnected {
				println(len(room.Clients), "clients in", room.Name)
				// Write the chunk to the WebSocket connection
				err = client.Connection.WriteMessage(websocket.BinaryMessage, buffer[:n])
				if err != nil {
					println("Field to Send", client.Name, "a chunk")
					break
				}
				println("Send", client.Name, "a chunk")
			} else {
				file.Close()
				break
			}
		}
		time.Sleep(500 * time.Millisecond)

	}
}
