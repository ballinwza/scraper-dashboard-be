package dto

// --- Request DTOs ---

type GetKnowledgeFileReqDTO struct {
	ID string `uri:"id" binding:"required"`
}

type ListKnowledgeFilesReqDTO struct {
	ChatbotID string `json:"chatbot_id" binding:"required"`
	Limit     int    `json:"limit" default:"10"`
	Offset    int    `json:"offset" default:"0"`
}

type DeleteKnowledgeFileReqDTO struct {
	ChatbotID string `json:"chatbot_id" binding:"required"`
}

type MultiTenantUploadFileReqDTO struct {
	UserID    string `form:"user_id" binding:"required"`
	ChatbotID string `form:"chatbot_id" binding:"required"`
}

type UploadFileMultiTenantReqDTO struct {
	ChatbotId string `form:"chatbot_id" json:"chatbot_id" binding:"required"`
}
