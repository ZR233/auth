/*
@Time : 2019-08-27 14:22
@Author : zr
*/
package gorm

import (
	"fmt"
	"github.com/ZR233/auth/errors"
	"github.com/ZR233/auth/model"
	gorm2 "github.com/jinzhu/gorm"
	"time"
)

var (
	initTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)
)

type Storage struct {
	db *gorm2.DB
}

func (s *Storage) Sync(memory *model.Memory) (err error) {
	var (
		services  []*Service
		roles     []*Role
		relations []*ServiceRoleRelation
	)

	err = s.db.Find(&services).Error
	if err != nil {
		return
	}
	err = s.db.Find(&roles).Error
	if err != nil {
		return
	}
	err = s.db.Find(&relations).Error
	if err != nil {
		return
	}

	for _, role := range roles {
		r := model.Role(*role)
		memory.Roles[role.Name] = &r
	}
	for _, service := range services {
		r := model.Service(*service)
		memory.Services[service.Path] = &r
	}

	for _, relation := range relations {
		if role, ok := memory.Roles[relation.RoleName]; ok {
			if service, ok := memory.Services[relation.ServiceName]; ok {
				role.Services = append(role.Services, service)
			}
		}
	}

	return
}

func (s *Storage) UpdateRelation(role *model.Role, services []*model.Service) (err error) {
	tx := s.db.Begin()
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.Delete(ServiceRoleRelation{}, "role_name = ?", role.Name).Error
	if err != nil {
		panic(err)
	}

	for _, service := range services {
		relation := ServiceRoleRelation{
			ServiceName: service.Path,
			RoleName:    role.Name,
		}

		err = tx.Create(&relation).Error
		if err != nil {
			panic(err)
		}
	}

	return
}

func NewStorage(gorm *gorm2.DB) *Storage {
	s := &Storage{}
	s.db = gorm
	s.migrate()
	return s
}

func (s *Storage) migrate() {
	if err := s.db.AutoMigrate(
		&Service{},
		&Role{},
		&ServiceRoleRelation{},
	).Error; err != nil {
		panic(err)
	}
}

func newRole(role *model.Role) (r *Role) {
	r_ := Role(*role)
	return &r_
}
func newService(service *model.Service) (r *Service) {
	r_ := Service(*service)
	return &r_
}

func (s *Storage) RoleUpdate(role *model.Role) (err error) {
	r := newRole(role)
	return s.db.Model(&Role{}).Omit("create_time").Where("name = ?", role.Name).Updates(r).Error
}
func (s *Storage) RoleCreate(role *model.Role) (err error) {
	r := newRole(role)

	var roles []*Role
	tx := s.db.Begin()
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.Where("name = ?", role.Name).Find(&roles).Error
	if err != nil {
		panic(err)
	}
	if len(roles) > 0 {
		err = fmt.Errorf("角色[%s]\n%w", role.Name, errors.ErrRecordExist)
		panic(err)
	}
	err = tx.Create(r).Error
	if err != nil {
		panic(err)
	}
	return
}
func (s *Storage) RoleDelete(name string) (err error) {
	return s.db.Where(&Role{Name: name}).Delete(&Role{Name: name}).Error
}

func (s *Storage) ServiceUpdate(service *model.Service) (err error) {
	service_ := newService(service)
	return s.db.Model(&Service{}).Omit("create_time").Where("path = ?", service.Path).Updates(service_).Error
}
func (s *Storage) ServiceCreate(service *model.Service) (err error) {
	service_ := newService(service)
	var services []*Service
	tx := s.db.Begin()
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.Where("path = ?", service_.Path).Find(&services).Error
	if err != nil {
		panic(err)
	}
	if len(services) > 0 {
		err = fmt.Errorf("服务[%s]\n%w", service_.Path, errors.ErrRecordExist)
		panic(err)
	}
	err = tx.Create(service_).Error
	if err != nil {
		panic(err)
	}
	return
}

func (s *Storage) ServiceDelete(servicePath string) (err error) {
	return s.db.Where(&Service{Path: servicePath}).Delete(&Service{Path: servicePath}).Error
}
