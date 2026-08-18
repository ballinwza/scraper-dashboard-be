package dto

import (
	"time"

	domain_ai_estate_rag "github.com/ballinwza/scraper-dashboard-be/internal/domain/ai_estate_rag"
)

// ==========================================
// Entities DTOs
// ==========================================

// ChatMessageResDTO DTO สำหรับรายละเอียดข้อความสนทนา
type ChatMessageResDTO struct {
	Role      domain_ai_estate_rag.MessageRole `json:"role" swaggertype:"string" example:"USER"`
	Content   string                           `json:"content" example:"สวัสดีครับ มีโปรโมชั่นคอนโดอะไรบ้าง"`
	CreatedAt time.Time                        `json:"created_at" example:"2026-08-18T10:00:00Z"`
}

// ChatSessionResDTO DTO สำหรับรายละเอียด Chat Session
type ChatSessionResDTO struct {
	ID           string              `json:"id" example:"65f1a2b3c4d5e6f7a8b9c0d1"`
	UserID       string              `json:"user_id" example:"user-999"`
	ChatbotID    string              `json:"chatbot_id" example:"bot-888"`
	SessionTitle string              `json:"session_title" example:"สอบถามข้อมูลโปรโมชั่น"`
	Messages     []ChatMessageResDTO `json:"messages,omitempty"`
	CreatedAt    time.Time           `json:"created_at" example:"2026-08-18T10:00:00Z"`
	UpdatedAt    time.Time           `json:"updated_at" example:"2026-08-18T10:05:00Z"`
}

// ==========================================
// Request DTOs
// ==========================================

// CreateChatSessionReqDTO DTO สำหรับสร้าง Chat Session ใหม่
type CreateChatSessionReqDTO struct {
	ChatbotID    string `json:"chatbot_id" binding:"required" example:"bot-888"`
	SessionTitle string `json:"session_title,omitempty" example:"สอบถามข้อมูลโปรโมชั่น"`
}

// GetChatSessionReqDTO DTO สำหรับดึงข้อมูล Chat Session
type GetChatSessionReqDTO struct {
	ID string `uri:"id" binding:"required" example:"65f1a2b3c4d5e6f7a8b9c0d1"`
}

// ListChatSessionsReqDTO DTO สำหรับดึงรายการ Chat Sessions แบบ Pagination
type ListChatSessionsReqDTO struct {
	ChatbotID string `form:"chatbot_id,omitempty" example:"bot-888"`
	PageSize  int32  `form:"page_size,default=20" binding:"omitempty,min=1,max=100" example:"20"`
	PageToken int32  `form:"page_token,default=0" binding:"omitempty,min=0" example:"0"`
}

// AddChatMessageReqDTO DTO สำหรับเพิ่มข้อความลงใน Session
type AddChatMessageReqDTO struct {
	SessionID string                           `json:"session_id" binding:"required" example:"65f1a2b3c4d5e6f7a8b9c0d1"`
	Role      domain_ai_estate_rag.MessageRole `json:"role" binding:"required" swaggertype:"string" example:"USER"`
	Content   string                           `json:"content" binding:"required" example:"สวัสดีครับ มีโปรโมชั่นคอนโดอะไรบ้าง"`
}

// DeleteChatSessionReqDTO DTO สำหรับลบ Chat Session
type DeleteChatSessionReqDTO struct {
	ID string `json:"id" binding:"required" example:"65f1a2b3c4d5e6f7a8b9c0d1"`
}

// ==========================================
// Response DTOs
// ==========================================

// CreateChatSessionResDTO คำตอบกลับเมื่อสร้าง Chat Session สำเร็จ
type CreateChatSessionResDTO struct {
	Session ChatSessionResDTO `json:"session"`
}

// GetChatSessionResDTO คำตอบกลับสำหรับการดึงข้อมูล Chat Session
type GetChatSessionResDTO struct {
	Session ChatSessionResDTO `json:"session"`
}

// ListChatSessionsResDTO คำตอบกลับรายการ Chat Sessions แบบ Pagination
type ListChatSessionsResDTO struct {
	Sessions   []ChatSessionResDTO `json:"sessions"`
	TotalCount int32               `json:"total_count" example:"5"`
}

// AddChatMessageResDTO คำตอบกลับเมื่อบันทึกข้อความสำเร็จ
type AddChatMessageResDTO struct {
	Message ChatMessageResDTO `json:"message"`
}
