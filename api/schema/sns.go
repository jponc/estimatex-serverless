package schema

const (
	ParticipantJoined string = "ParticipantJoined"
	ParticipantVoted  string = "ParticipantVoted"
)

type ParticipantJoinedMessage struct {
	RoomID          string `json:"room_id"`
	ParticipantName string `json:"participant_name"`
}

type ParticipantVotedMessage struct {
	RoomID          string `json:"room_id"`
	ParticipantName string `json:"participant_name"`
	Vote            string `json:"vote"`
}
