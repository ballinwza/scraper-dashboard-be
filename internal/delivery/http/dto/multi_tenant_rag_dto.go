package dto

type UploadFileMultiTenantReqDTO struct {
	ChatbotId string `form:"chatbot_id" json:"chatbot_id" binding:"required"`
}
type UploadFileMultiTenantResDTO struct {
	FileId      string `json:"file_id"`
	Status      string `json:"status"`
	TotalChunks int    `json:"total_chunks"`
	TotalBytes  int    `json:"total_bytes"`
	Message     string `json:"message"`
}
