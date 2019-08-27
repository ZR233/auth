/*
@Time : 2019-08-27 14:11
@Author : zr
*/
package auth

import (
	"github.com/ZR233/auth/model"
	"time"
)

type Core struct {
	serviceTree map[string]*model.Service
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
	serviceTree, err := c.storage.Sync()
	if err != nil {
		return err
	}
	c.serviceTree = serviceTree
	return nil
}
