package handler

// apiResponse Standard Response Structure
type apiResponse struct {
	Status  int         `json:"status" example:"200"`
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
}

// errorResponse Standard Error Structure
type errorResponse struct {
	Error string `json:"error" example:"Invalid request payload"`
}
