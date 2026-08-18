package domain_ai_estate_rag

import (
	"time"
)

// ChatbotBlueprint โครงสร้างข้อมูลหลักของ Chatbot Entity
type ChatbotBlueprint struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	SystemPrompt string    `json:"system_prompt"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
