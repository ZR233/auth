/*
@Time : 2019-08-27 14:42
@Author : zr
*/
package gorm

import (
	"github.com/ZR233/auth/model"
)

type Role model.Role

func (Role) TableName() string {
	return "t_auth_role"
}

type Service model.Service

func (Service) TableName() string {
	return "t_auth_service"
}

type ServiceRoleRelation struct {
	ServiceName string
	RoleName    string
}

func (ServiceRoleRelation) TableName() string {
	return "t_auth_service_role_relation"
}
