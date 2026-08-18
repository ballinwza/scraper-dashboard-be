package domain_ai_estate_rag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type MessageRole int32

const (
	MessageRoleUnspecified MessageRole = iota // 0: MESSAGE_ROLE_UNSPECIFIED
	MessageRoleUser                           // 1: USER
	MessageRoleAI                             // 2: AI
	MessageRoleSystem                         // 3: SYSTEM
)

// String แปลง MessageRole เป็น String
func (r MessageRole) String() string {
	switch r {
	case MessageRoleUser:
		return "USER"
	case MessageRoleAI:
		return "AI"
	case MessageRoleSystem:
		return "SYSTEM"
	default:
		return "MESSAGE_ROLE_UNSPECIFIED"
	}
}

// ParseMessageRole แปลงข้อความ String กลับเป็น MessageRole
func ParseMessageRole(roleStr string) MessageRole {
	switch roleStr {
	case "USER":
		return MessageRoleUser
	case "AI":
		return MessageRoleAI
	case "SYSTEM":
		return MessageRoleSystem
	default:
		return MessageRoleUnspecified
	}
}

// MarshalJSON จัดการแปลงค่าตอน serialize เป็น JSON ให้แสดงผลเป็น String
func (r MessageRole) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(r.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON จัดการแปลงค่าตอน deserialize จาก JSON (รองรับทั้ง String และ Int)
func (r *MessageRole) UnmarshalJSON(b []byte) error {
	var item interface{}
	if err := json.Unmarshal(b, &item); err != nil {
		return err
	}

	switch v := item.(type) {
	case string:
		*r = ParseMessageRole(v)
	case float64:
		*r = MessageRole(int32(v))
	default:
		return fmt.Errorf("invalid type for MessageRole: %T", item)
	}
	return nil
}

// ChatMessage โครงสร้างข้อความสนทนา
type ChatMessage struct {
	Role      MessageRole `json:"role"`
	Content   string      `json:"content"`
	CreatedAt time.Time   `json:"created_at"`
}

// ChatSession โครงสร้างประวัติการสนทนา
type ChatSession struct {
	ID           string        `json:"id"`
	UserID       string        `json:"user_id"`
	ChatbotID    string        `json:"chatbot_id"`
	SessionTitle string        `json:"session_title"`
	Messages     []ChatMessage `json:"messages"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}
