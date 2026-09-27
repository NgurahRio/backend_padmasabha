package models

import "time"

// PasswordReset only stores hashes of the OTP and reset token.
type PasswordReset struct {
	ID             uint       `gorm:"primaryKey"`
	UserID         uint       `gorm:"column:user_id;index;not null"`
	OTPHash        string     `gorm:"column:otp_hash;size:64;not null"`
	OTPExpiresAt   time.Time  `gorm:"column:otp_expires_at;not null"`
	Attempts       int        `gorm:"column:attempts;not null;default:0"`
	LastSentAt     time.Time  `gorm:"column:last_sent_at;not null"`
	ResetTokenHash *string    `gorm:"column:reset_token_hash;size:64;index"`
	TokenExpiresAt *time.Time `gorm:"column:token_expires_at"`
	UsedAt         *time.Time `gorm:"column:used_at"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (PasswordReset) TableName() string { return "password_resets" }
