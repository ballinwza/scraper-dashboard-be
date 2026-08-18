package dto

import "time"

// ==========================================
// Entities DTO
// ==========================================

// ChatbotBlueprintResDTO ข้อมูลรายละเอียดของ Chatbot Blueprint
type ChatbotBlueprintResDTO struct {
	ID           string    `json:"id" example:"65f1a2b3c4d5e6f7a8b9c0d1"`
	UserID       string    `json:"user_id" example:"user-999"`
	Name         string    `json:"name" example:"Customer Support Bot"`
	Description  string    `json:"description" example:"บอทสำหรับตอบคำถามและดูแลลูกค้า"`
	SystemPrompt string    `json:"system_prompt" example:"You are a helpful assistant."`
	CreatedAt    time.Time `json:"created_at" example:"2026-08-18T10:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2026-08-18T10:05:00Z"`
}

// ==========================================
// Request DTOs
// ==========================================

// CreateMultiTenantChatbotReqDTO DTO สำหรับสร้าง Chatbot ใหม่
type CreateMultiTenantChatbotReqDTO struct {
	Name         string `json:"name" binding:"required" example:"Customer Support Bot"`
	Description  string `json:"description,omitempty" example:"บอทสำหรับตอบคำถามและดูแลลูกค้า"`
	SystemPrompt string `json:"system_prompt" binding:"required" example:"You are a helpful assistant."`
}

// GetMultiTenantChatbotReqDTO DTO สำหรับดึงข้อมูล Chatbot
type GetMultiTenantChatbotReqDTO struct {
	ID string `uri:"id" binding:"required" example:"65f1a2b3c4d5e6f7a8b9c0d1"`
}

// ListMultiTenantChatbotsReqDTO DTO สำหรับดึงรายการ Chatbot แบบ Pagination
type ListMultiTenantChatbotsReqDTO struct {
	PageSize  int32 `form:"page_size,default=10" binding:"omitempty,min=1,max=100" example:"10"`
	PageToken int32 `form:"page_token,default=0" binding:"omitempty,min=0" example:"0"`
}

// UpdateMultiTenantChatbotReqDTO DTO สำหรับอัปเดตข้อมูล Chatbot
type UpdateMultiTenantChatbotReqDTO struct {
	ID           string  `json:"id" binding:"required" example:"65f1a2b3c4d5e6f7a8b9c0d1"`
	Name         *string `json:"name,omitempty" example:"Updated Bot Name"`
	Description  *string `json:"description,omitempty" example:"Updated description"`
	SystemPrompt *string `json:"system_prompt,omitempty" example:"Updated system prompt"`
	// TODO: เพิ่มตอนทำ Masking
	// UpdateMask   []string `json:"update_mask,omitempty" example:"[\"name\", \"system_prompt\"]"`
}

// DeleteMultiTenantChatbotReqDTO DTO สำหรับลบ Chatbot
type DeleteMultiTenantChatbotReqDTO struct {
	ID string `json:"id" binding:"required" example:"65f1a2b3c4d5e6f7a8b9c0d1"`
}

// ==========================================
// Response DTOs
// ==========================================

// CreateMultiTenantChatbotResDTO คำตอบกลับเมื่อสร้าง Chatbot สำเร็จ
type CreateMultiTenantChatbotResDTO struct {
	Chatbot ChatbotBlueprintResDTO `json:"chatbot"`
}

// GetMultiTenantChatbotResDTO คำตอบกลับสำหรับการดึงข้อมูล Chatbot เดี่ยว
type GetMultiTenantChatbotResDTO struct {
	Chatbot ChatbotBlueprintResDTO `json:"chatbot"`
}

// ListMultiTenantChatbotsResDTO คำตอบกลับรายการ Chatbots แบบ Pagination
type ListMultiTenantChatbotsResDTO struct {
	Chatbots      []ChatbotBlueprintResDTO `json:"chatbots"`
	NextPageToken int32                    `json:"next_page_token" example:"10"`
	TotalCount    int32                    `json:"total_count" example:"25"`
}

// UpdateMultiTenantChatbotResDTO คำตอบกลับเมื่ออัปเดต Chatbot สำเร็จ
type UpdateMultiTenantChatbotResDTO struct {
	Chatbot ChatbotBlueprintResDTO `json:"chatbot"`
}

// DeleteMultiTenantChatbotResDTO คำตอบกลับเมื่อลบ Chatbot สำเร็จ
type DeleteMultiTenantChatbotResDTO struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Chatbot deleted successfully"`
}
