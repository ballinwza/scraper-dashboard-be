package domain_ai_estate_rag

import "time"

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

// Knowledge File
type FileStatus int32

const (
	FileStatusPending   FileStatus = 0
	FileStatusCompleted FileStatus = 1
	FileStatusFailed    FileStatus = 2
)

// String ช่วยแปลง FileStatus เป็น String
func (s FileStatus) String() string {
	switch s {
	case FileStatusPending:
		return "PENDING"
	case FileStatusCompleted:
		return "COMPLETED"
	case FileStatusFailed:
		return "FAILED"
	default:
		return "UNKNOWN"
	}
}

// Chunk โครงสร้างข้อมูล Chunk ของเอกสาร
type Chunk struct {
	VectorID    string `json:"vector_id"`
	ChunkIndex  int32  `json:"chunk_index"`
	TextContent string `json:"text_content"`
	PageNumber  int32  `json:"page_number"`
	TokenCount  int32  `json:"token_count"`
}

// KnowledgeFile โครงสร้างข้อมูลหลักของ KnowledgeFile Entity
type KnowledgeFile struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	ChatbotID    string     `json:"chatbot_id"`
	Filename     string     `json:"filename"`
	FileType     string     `json:"file_type"`
	FileSizeByes int64      `json:"file_size_bytes"`
	Status       FileStatus `json:"status"`
	TotalChunks  int32      `json:"total_chunks"`
	Chunks       []Chunk    `json:"chunks,omitempty"`
	TotalPage    int32      `json:"total_page"`
	TextContent  *string    `json:"text_content,omitempty"`  // optional
	ErrorMessage *string    `json:"error_message,omitempty"` // optional
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// --- Structs สำหรับ Get Knowledge File ---

type GetKnowledgeFileResponse struct {
	File KnowledgeFile `json:"file"`
}

// --- Structs สำหรับ List Knowledge Files ---

type ListKnowledgeFilesResponse struct {
	Files      []KnowledgeFile `json:"files"`
	TotalCount int32           `json:"total_count"`
}

// --- Structs สำหรับ Delete Knowledge File ---

type DeleteKnowledgeFileResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// --- Response DTOs ---

type UploadFileMultiTenantResponse struct {
	FileID      string     `json:"file_id"`
	Status      FileStatus `json:"status"`
	TotalChunks int32      `json:"total_chunks"`
	TotalBytes  int64      `json:"total_bytes"`
	Message     string     `json:"message"`
	CreatedAt   time.Time  `json:"created_at"`
}
