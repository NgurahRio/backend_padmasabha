package models

import "time"

type AdminSession struct {
	ID         uint       `gorm:"primaryKey"`
	UserID     uint       `gorm:"column:user_id;index;not null"`
	TokenHash  string     `gorm:"column:token_hash;size:64;uniqueIndex;not null"`
	UserAgent  string     `gorm:"column:user_agent;type:text"`
	IPAddress  string     `gorm:"column:ip_address;size:45"`
	ExpiresAt  time.Time  `gorm:"column:expires_at;index;not null"`
	LastUsedAt time.Time  `gorm:"column:last_used_at;not null"`
	RevokedAt  *time.Time `gorm:"column:revoked_at;index"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (AdminSession) TableName() string { return "admin_sessions" }
