package handler

// apiResponse Standard Response Structure
type userResponse struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}
