package schema

import "github.com/jponc/estimatex-serverless/internal/types"

type SayHelloRequest struct {
	Name string `json:"name"`
}

type SayHelloResponse struct {
	Message string `json:"message"`
}

type HostRoomRequest struct {
	Name string `json:"name"`
}

type HostRoomResponse struct {
	RoomID      string `json:"room_id"`
	AccessToken string `json:"access_token"`
}

type FindRoomRequest struct {
	ID string `json:"id"`
}

type FindRoomResponse struct {
	types.Room
}
