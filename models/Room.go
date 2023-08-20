package models

type Room struct {
	ID      string
	Clients []*Client
	Code    int
	Name    string
	AdminID string
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
