/*
@Time : 2019-08-27 14:11
@Author : zr
*/
package auth

import (
	"github.com/ZR233/auth/model"
	"strings"
	"time"
)

type Core struct {
	serviceTree map[string]*model.Service
	roles       []*model.Role
	storage     model.Storage
}

func NewCore(storage model.Storage) *Core {
	c := &Core{
		storage: storage,
	}
	return c
}

func (c *Core) NewRole(name string) *model.Role {
	return &model.Role{
		Name:       name,
		State:      1,
		Storage:    c.storage,
		CreateTime: time.Now(),
		EditTime:   time.Now(),
	}
}
func (c *Core) NewService(name string) *model.Service {
	return &model.Service{
		Name:        name,
		SupService:  "",
		SubService:  nil,
		Description: "",
		CreateTime:  time.Now(),
		EditTime:    time.Now(),
		Storage:     c.storage,
	}
}
func (c *Core) Sync() error {
	serviceTree, roles, err := c.storage.Sync()
	if err != nil {
		return err
	}
	c.serviceTree = serviceTree
	c.roles = roles
	return nil
}
func rolesHasName(roles []*model.Role, name string) bool {
	r := false
	for _, v := range roles {
		if v.Name == name {
			r = true
			break
		}
	}
	return r
}

func (c *Core) Check(ServiceUrl string, roleName string) (r bool) {
	r = false
	serviceNames := strings.Split(ServiceUrl, "/")
	if len(serviceNames) == 0 {
		return true
	}
	serviceName := serviceNames[0]
	service, ok := c.serviceTree[serviceName]
	if ok {
		r = rolesHasName(service.Roles, roleName)
		if !r {
			return
		}
	} else {
		return true
	}
	serviceNames = serviceNames[1:]
	for _, servceName := range serviceNames {
		service, ok = service.SubService[servceName]
		if ok {
			r = rolesHasName(service.Roles, roleName)
			if !r {
				return
			}
		} else {
			return true
		}
	}

	return
}
