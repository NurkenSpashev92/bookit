package schemas

// ErrorResponse error response
type ErrorResponse struct {
	Error string `json:"error" example:"validation failed"`
}

// MessageResponse success message response
type MessageResponse struct {
	Message string `json:"message" example:"operation completed"`
}
