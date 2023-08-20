package server

import (
	"encoding/json"
	"go-ws/models"
	"go-ws/utils"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var ws = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
}
var clients []*models.Client
var rooms []*models.Room

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	connection, err := ws.Upgrade(w, r, nil)
	if err != nil {
		println("Connection Closed")
		connection.Close()
		return
	}
	currentClient := &models.Client{ID: utils.GenerateUniqueID(), Connection: connection, IsConnected: true}
	clients = append(clients, currentClient)
	defer func() {
		RemoveClient(currentClient)
	}()
	for {
		_, message, err := currentClient.Connection.ReadMessage()
		if err != nil {
			RemoveClient(currentClient)
			break
		}

		clientRequest := &models.ClientRequest{}

		err = json.Unmarshal(message, clientRequest)
		if err != nil {
			println("Unknown request")
		}
		handleClientRequest(currentClient, clientRequest)
	}

}

func handleClientRequest(client *models.Client, request *models.ClientRequest) {
	if request.RequestId == models.CLIENT_LOGIN_REQUEST {
		client.SetName(request.Body)
		println(client.Name, "Joined")
		client.Connection.WriteMessage(websocket.TextMessage, []byte("added"))
		return
	}
	if !client.IsLogin() {
		client.Connection.WriteMessage(websocket.TextMessage, []byte("400"))
		return
	}
	if request.RequestId == models.CLIENT_GET_ROOMS_REQUEST {
		data, err := json.Marshal(rooms)
		if err != nil {
			log.Println(err)
			return
		}
		client.Connection.WriteMessage(websocket.TextMessage, data)
		return
	}
	if request.RequestId == models.CLIENT_CREATE_ROOM_REQUEST {
		createRoomRequest := &models.CreateRoomBody{}
		err := json.Unmarshal([]byte(request.Body), createRoomRequest)
		if err != nil {
			log.Println(err)
			return
		}
		rooms = append(rooms, &models.Room{ID: utils.GenerateUniqueID(), Name: createRoomRequest.Name, Code: createRoomRequest.Code, AdminID: client.ID})
		client.Connection.WriteMessage(websocket.TextMessage, []byte("created"))
	}
	if request.RequestId == models.CLIENT_JOIN_ROOM_REQUEST {
		joinRoomBody := &models.JoinRoomBody{}
		err := json.Unmarshal([]byte(request.Body), joinRoomBody)
		if err != nil {
			log.Println(err)
			return
		}
		for _, room := range rooms {
			if room.ID == joinRoomBody.ID && room.Code == joinRoomBody.Code {
				println(client.Name, "Joined", room.Name)
				JoinRoom(room, client)
			}
		}
	}
	if request.RequestId == models.CLIENT_GET_FILE_REQUEST {
		println(client.Name + " requesting buffer")

		for _, room := range rooms {
			if room.AdminID == client.ID {
				println(client.Name + " is Admin " + "in " + room.Name)
				println("room clients length", len(room.Clients))
				room.StartParty()
				return
			}
		}
		return
	}

}

func JoinRoom(room *models.Room, client *models.Client) {
	client.RoomId = room.ID
	room.AddClient(*client)
	data, err := json.Marshal(room.Clients)
	if err != nil {
		println(err)
	}
	client.Connection.WriteMessage(websocket.TextMessage, data)
}
func RemoveClient(client *models.Client) {
	println(client.ID, "Dropped")
	var updatedClients []*models.Client
	client.Connection.Close()
	for _, c := range clients {
		if c.ID == client.ID {
			updatedClients = append(updatedClients, c)
			clients = updatedClients
			if c.IsJoinedRoom() {
				println(c.Name, "IsJoinedRoom", c.Name)
				for _, room := range rooms {
					if room.ID == c.RoomId {
						println("Remove", c.Name, "from", room.Name)
						room.RemoveClient(c)
						println("Room", room.Name, "has", len(room.Clients))
						break
					}
				}
			}
			break
		}
	}
}

func RunServer() {
	http.HandleFunc("/", handleWebSocket)

	println("Server is running!")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
