package models

import "gorm.io/gorm"

type User struct {
	ID        uint           `gorm:"primaryKey;column:id_users" json:"id_users"`
	Username  string         `gorm:"column:username;not null;unique" json:"username"`
	Password  string         `gorm:"column:password;not null" json:"password"`
	RoleID    uint           `gorm:"column:roleId" json:"roleId"`
	Email     string         `gorm:"column:email;not null" json:"email"`
	Image     *string        `gorm:"column:image;type:longtext" json:"image,omitempty"`

	Role      role           `gorm:"foreignKey:RoleID;references:ID"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_at;index" json:"delete_at"`
}

func (User) TableName() string {
	return "users"
}
