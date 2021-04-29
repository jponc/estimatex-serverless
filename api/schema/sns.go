package schema

const (
	ParticipantJoined string = "ParticipantJoined"
	ParticipantVoted  string = "ParticipantVoted"
	RevealVotes       string = "RevealVotes"
	ResetVotes        string = "ResetVotes"
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

type RevealVotesMessage struct {
	RoomID string `json:"room_id"`
}

type ResetVotesMessage struct {
	RoomID string `json:"room_id"`
}
