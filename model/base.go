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
	storage     Storage `json:"-"`
}

func (r *Role) Save() error {
	return r.storage.SaveRole(r)
}

func (r *Role) SetStorage(storage Storage) {
	r.storage = storage
}

type Service struct {
	Name        string
	SupService  string
	SubService  map[string]*Service
	Description string
	CreateTime  time.Time
	EditTime    time.Time
	Roles       []*Role
	storage     Storage `json:"-"`
}

func (s *Service) AddRoles(roles ...*Role) (err error) {
	return s.storage.SaveRelation(s, roles...)
}

func (s *Service) Save() error {
	return s.storage.SaveService(s)
}
func (s *Service) NewSubService(name string) *Service {
	now := time.Now()

	service := &Service{
		Name:        name,
		SupService:  s.Name,
		SubService:  nil,
		Description: "",
		CreateTime:  now,
		EditTime:    now,
		storage:     s.storage,
	}
	return service
}

func (s *Service) SetStorage(storage Storage) {
	s.storage = storage
}

type Storage interface {
	Sync() (serviceTree map[string]*Service, roles []*Role, err error)
	SaveRole(role *Role) error
	SaveService(service *Service) error
	SaveRelation(service *Service, role ...*Role) error
}
