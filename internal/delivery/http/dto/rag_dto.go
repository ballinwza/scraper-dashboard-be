package dto

// ==========================================
// Entities DTOs
// ==========================================

type MetadataVectorRecordDTO struct {
	UserID      string `json:"user_id" example:"user-999"`
	ChatbotID   string `json:"chatbot_id" example:"bot-888"`
	FileID      string `json:"file_id" example:"file-123"`
	ChunkIndex  int32  `json:"chunk_index" example:"1"`
	TextContent string `json:"text_content" example:"เนื้อหาในเอกสารที่ค้นพบ..."`
	PageNumber  int32  `json:"page_number" example:"5"`
	Filename    string `json:"filename" example:"document.pdf"`
}

type VectorRecordDTO struct {
	ID       string                  `json:"id" example:"vector-uuid-123"`
	Values   []float32               `json:"values,omitempty"`
	Metadata MetadataVectorRecordDTO `json:"metadata"`
}

type SearchVectorRecordItemDTO struct {
	Score  float32         `json:"score" example:"0.895"`
	Record VectorRecordDTO `json:"record"`
}

// ==========================================
// Request DTO
// ==========================================

// SearchSimilarReqDTO DTO สำหรับส่งคำค้นหา RAG
type SearchSimilarReqDTO struct {
	ChatbotID       string  `json:"chatbot_id" binding:"required" example:"bot-888"`
	QueryText       string  `json:"query_text" binding:"required" example:"สอบถามข้อมูลโปรโมชั่น"`
	TopK            *int32  `json:"top_k,omitempty" example:"5"`
	KnowledgeFileID *string `json:"knowledge_file_id,omitempty" example:"file-123"`
}

// ==========================================
// Response DTO
// ==========================================

// SearchSimilarResDTO คำตอบกลับและเอกสารอ้างอิงจากการค้นหา
type SearchSimilarResDTO struct {
	AnswerMessage string                      `json:"answer_message" example:"นี่คือคำตอบจากระบบ..."`
	Sources       []SearchVectorRecordItemDTO `json:"sources"`
}
