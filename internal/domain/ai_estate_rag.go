package domain

type UploadPdfResponse struct {
	FileId  string `json:"file_id"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type UploadFileMultiTenantReq struct {
	UserId    string `json:"user_id"`
	ChatbotId string `json:"chatbot_id"`
	Filename  string `json:"filename"`
	FileType  string `json:"file_type"`
	FileBytes []byte `json:"file_bytes"`
}

type UploadFileMultiTenantRes struct {
	FileId      string `json:"file_id"`
	Status      string `json:"status"`
	TotalChunks int    `json:"total_chunks"`
	TotalBytes  int    `json:"total_bytes"`
	Message     string `json:"message"`
}
