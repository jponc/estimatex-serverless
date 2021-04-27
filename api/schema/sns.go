package schema

const (
	ParticipantJoined string = "ParticipantJoined"
)

type ParticipantJoinedMessage struct {
	RoomID          string `json:"room_id"`
	ParticipantName string `json:"participant_name"`
}
