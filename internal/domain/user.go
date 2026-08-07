package domain

import (
	"time"
)

type RefreshTokenEmbedded struct {
	TokenHash string    `bson:"token_hash" json:"-"` // เก็บเป็น SHA-256 Hash ห้ามเก็บ Token ตัวเต็ม
	UserAgent string    `bson:"user_agent,omitempty" json:"user_agent"`
	ClientIP  string    `bson:"client_ip,omitempty" json:"client_ip"`
	CreatedAt time.Time `bson:"created_at,omitempty" json:"-"`
	ExpiresAt time.Time `bson:"expires_at,omitempty" json:"-"`
}

type User struct {
	ID            string                 `bson:"_id,omitempty" json:"id"`
	Username      string                 `bson:"username,omitempty" json:"username"`
	PasswordHash  string                 `bson:"password_hash,omitempty" json:"-"` // omitempty/dash เพื่อไม่ให้หลุดออก JSON API
	Name          string                 `bson:"name,omitempty" json:"name" `
	Role          string                 `bson:"role,omitempty" json:"role"`
	IsActive      bool                   `bson:"is_active,omitempty" json:"is_active"`
	RefreshTokens []RefreshTokenEmbedded `bson:"refresh_tokens" json:"-"`
	CreatedAt     time.Time              `bson:"created_at,omitempty" json:"-"`
	UpdatedAt     time.Time              `bson:"updated_at,omitempty" json:"-"`
}
