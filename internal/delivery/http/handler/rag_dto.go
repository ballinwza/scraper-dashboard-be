package handler

type generateAnswerResponse struct {
	Message string `json:"message"`
}

type generateAnswerRequest struct {
	Question string `json:"question" binding:"required" default:"สรุปเอกสารย่อๆ"`
}
