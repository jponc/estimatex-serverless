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

type FindRoomResponse struct {
	types.Room
}

type JoinRoomRequest struct {
	RoomID string `json:"room_id"`
	Name   string `json:"name"`
}

type JoinRoomResponse struct {
	AccessToken string `json:"access_token"`
}

type CastVoteRequest struct {
	Vote string `json:"vote"`
}

type CastVoteResponse struct {
}
