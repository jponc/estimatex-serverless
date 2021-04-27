package schema

type SayHelloRequest struct {
	Name string `json:"name"`
}

type SayHelloResponse struct {
	Message string `json:"message"`
}
