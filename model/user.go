package model

import "time"

type UserStatus int

const (
	UserStatusBan UserStatus = -1
	_             UserStatus = iota
	UserStatusNormal
)

type User struct {
	Id         int    `gorm:"type:char(100)"`
	UserCode   string `gorm:"type:varchar(100)"`
	Name       string `gorm:"type:varchar(100)"`
	Password   string `gorm:"type:varchar(100)"`
	Status     UserStatus
	Memo       string `gorm:"type:varchar(500)"`
	EditTime   *time.Time
	CreateTime *time.Time
	DeleteAt   *time.Time
}

func (User) TableName() string {
	return "t_user"
}

type UserRole struct {
	UserId     int
	Role       string
	EditTime   *time.Time
	CreateTime *time.Time
}

func (UserRole) TableName() string {
	return "t_user_role"
}
