/*
@Time : 2019-08-27 14:57
@Author : zr
*/
package model

import "time"

type Role struct {
	Name        string
	State       int16
	Description string
	EditTime    time.Time
	CreateTime  time.Time
	Services    []*Service
	Storage     Storage
}

func (r *Role) Save() error {
	return r.Storage.SaveRole(r)
}

type Service struct {
	Name        string
	SupService  string
	SubService  map[string]*Service
	Description string
	CreateTime  time.Time
	EditTime    time.Time
	Roles       []*Role
	Storage     Storage
}

func (s *Service) AddRoles(roles ...*Role) (err error) {
	return s.Storage.SaveRelation(s, roles...)
}

func (s *Service) Save() error {
	return s.Storage.SaveService(s)
}
func (s *Service) NewSubService(name string) *Service {
	service := &Service{
		Name:        name,
		SupService:  s.Name,
		SubService:  nil,
		Description: "",
		CreateTime:  time.Now(),
		EditTime:    time.Now(),
		Storage:     s.Storage,
	}
	return service
}

type Storage interface {
	Sync() (serviceTree map[string]*Service, err error)
	SaveRole(role *Role) error
	SaveService(service *Service) error
	SaveRelation(service *Service, role ...*Role) error
}
