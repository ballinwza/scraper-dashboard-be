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

// // --- Response DTOs ---

// type UploadFileMultiTenantResDTO struct {
// 	FileID      string                          `json:"file_id"`
// 	Status      domain_ai_estate_rag.FileStatus `json:"status"`
// 	TotalChunks int32                           `json:"total_chunks"`
// 	TotalBytes  int64                           `json:"total_bytes"`
// 	Message     string                          `json:"message"`
// 	CreatedAt   time.Time                       `json:"created_at"`
// }
