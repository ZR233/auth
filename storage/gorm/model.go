/*
@Time : 2019-08-27 14:42
@Author : zr
*/
package gorm

import "time"

type Role struct {
	Name        string `gorm:"primary_key;type:varchar(100)"`
	State       int16  `gorm:"type:tinyint"`
	Description string `gorm:"type:varchar(500)"`
	EditTime    time.Time
	CreateTime  time.Time
}

func (Role) TableName() string {
	return "auth_role"
}

type Service struct {
	Name         string `gorm:"primary_key;type:varchar(100)"`
	SuperiorName string `gorm:"type:varchar(100)"`
	Description  string `gorm:"type:varchar(200)"`
	CreateTime   time.Time
	EditTime     time.Time
	Services     []*Service `gorm:"-"`
}

func (Service) TableName() string {
	return "auth_service"
}

type ServiceRoleRelation struct {
	ServiceName string `gorm:"primary_key;type:varchar(100)"`
	RoleName    string `gorm:"primary_key;type:varchar(100)"`
}

func (ServiceRoleRelation) TableName() string {
	return "auth_service_role_relation"
}
