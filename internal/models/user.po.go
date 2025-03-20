package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model           // default fields: ID, CreatedAt, UpdatedAt, DeletedAt
	UUID       uuid.UUID `gorm:"column:uuid; type:varchar(255); not null; unique; index:idx_uuid;"`
	UserName   string    `gorm:"column:user_name;"`
	IsActive   bool      `gorm:"column:is_active; type:boolean; not null; default:true;"`
}

func (u *User) TableName() string {
	return "go_db_user"
}
