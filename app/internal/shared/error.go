package shared

type ErrorResponse struct {
	Error string `json:"error" example:"validation failed"`
}

type MessageResponse struct {
	Message string `json:"message" example:"operation completed"`
}
