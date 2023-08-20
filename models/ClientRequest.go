package models

const (
	CLIENT_LOGIN_REQUEST       = 1
	CLIENT_JOIN_ROOM_REQUEST   = 6
	CLIENT_LOGIN_OUT_REQUEST   = 3
	CLIENT_GET_ROOMS_REQUEST   = 5
	CLIENT_CREATE_ROOM_REQUEST = 2
	CLIENT_GET_FILE_REQUEST    = 4
)

type ClientRequest struct {
	RequestId int
	Body      string
}

type CreateRoomBody struct {
	Name string
	Code int
}
type JoinRoomBody struct {
	ID   string
	Code int
}
