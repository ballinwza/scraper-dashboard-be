package domain_ai_estate_rag

// MetadataVectorRecord รายละเอียด Metadata ของแต่ละ Vector Chunk
type MetadataVectorRecord struct {
	UserID      string `json:"user_id"`
	ChatbotID   string `json:"chatbot_id"`
	FileID      string `json:"file_id"`
	ChunkIndex  int32  `json:"chunk_index"`
	TextContent string `json:"text_content"`
	PageNumber  int32  `json:"page_number"`
	Filename    string `json:"filename"`
}

// VectorRecord โครงสร้างข้อมูล Vector และ Metadata
type VectorRecord struct {
	ID       string               `json:"id"`
	Values   []float32            `json:"values"`
	Metadata MetadataVectorRecord `json:"metadata"`
}

// SearchVectorRecordItem ผลลัพธ์พร้อมคะแนน Similarity Score
type SearchVectorRecordItem struct {
	Score  float32      `json:"score"`
	Record VectorRecord `json:"record"`
}

// RagResult ผลลัพธ์หลักการค้นหาคำตอบและแหล่งอ้างอิง
type RagResult struct {
	AnswerMessage string                   `json:"answer_message"`
	Sources       []SearchVectorRecordItem `json:"sources"`
}
