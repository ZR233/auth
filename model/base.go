/*
@Time : 2019-08-27 14:57
@Author : zr
*/
package model

import (
	"fmt"
	"github.com/ZR233/auth/errors"
	"sync"
	"time"
)

type Status int

const (
	StatusOff Status = -1
	_         Status = iota
	StatusOn
)

type Storage interface {
	Sync(memory *Memory) (err error)
	RoleUpdate(role *Role) error
	RoleCreate(role *Role) (err error)
	RoleDelete(name string) (err error)
	ServiceUpdate(service *Service) (err error)
	ServiceCreate(service *Service) (err error)
	ServiceDelete(name string) (err error)
	UpdateRelation(role *Role, service []*Service) error
}

type Memory struct {
	storage  Storage
	Services map[string]*Service
	Roles    map[string]*Role
	sync.Mutex
	SyncErr error
}

func NewMemory(storage Storage) Memory {
	return Memory{
		storage:  storage,
		Services: map[string]*Service{},
		Roles:    map[string]*Role{},
	}
}

func (m *Memory) NewRole(name string, description string) (r *Role, err error) {
	m.Lock()
	defer m.Unlock()

	r = &Role{
		Name:        name,
		Description: description,
		CreateTime:  time.Now(),
		EditTime:    time.Now(),
		Status:      StatusOn,
	}
	r.SetMemory(m)
	err = m.storage.RoleCreate(r)
	if err != nil {
		return
	}
	m.Roles[name] = r
	return
}

func (m *Memory) RoleUpdateServices(roleName string, servicePaths ...string) (err error) {
	m.Lock()
	defer m.Unlock()
	var (
		ok       bool
		role     *Role
		service  *Service
		services []*Service
	)

	if role, ok = m.Roles[roleName]; !ok {
		err = fmt.Errorf("角色[%s]\n%w", roleName, errors.ErrRecordNotExist)
		return
	}
	for _, servicePath := range servicePaths {
		if service, ok = m.Services[servicePath]; !ok {
			err = fmt.Errorf("服务[%s]\n%w", servicePath, errors.ErrRecordNotExist)
			return
		}

		services = append(services, service)
	}

	err = m.storage.UpdateRelation(role, services)
	if err != nil {
		err = fmt.Errorf("持久保存失败\n%w", err)
		return
	}
	role.Services = services
	return
}

func (m *Memory) NewService(servicePath, description string) (service *Service, err error) {
	m.Lock()
	defer m.Unlock()

	service = &Service{
		Path:        servicePath,
		Description: description,
		CreateTime:  time.Now(),
		EditTime:    time.Now(),
		Status:      StatusOn,
	}
	service.SetMemory(m)
	err = m.storage.ServiceCreate(service)
	if err != nil {
		return
	}
	m.Services[servicePath] = service
	return
}
func (m *Memory) Sync() {
	m.Lock()
	defer m.Unlock()

	m.SyncErr = m.storage.Sync(m)
	if m.SyncErr != nil {
		return
	}
	return
}

func (m *Memory) DeleteRole(name string) (err error) {
	m.Lock()
	defer m.Unlock()

	err = m.storage.RoleDelete(name)
	if err != nil {
		return
	}
	delete(m.Roles, name)

	return
}
func (m *Memory) DeleteService(name string) (err error) {
	m.Lock()
	defer m.Unlock()
	err = m.storage.ServiceDelete(name)
	if err != nil {
		return
	}
	delete(m.Services, name)
	return
}
