/*
@Time : 2019-08-27 14:11
@Author : zr
*/
package auth

import (
	"context"
	"fmt"
	"github.com/ZR233/auth/errors"
	"github.com/ZR233/auth/model"
	"path"
	"strings"
	"time"
)

type Core struct {
	serviceTree map[string]*model.Service
	roles       map[string]*model.Role

	model.Memory
	ctx    context.Context
	cancel context.CancelFunc
}

func NewCore(storage model.Storage) (c *Core, err error) {
	c = &Core{
		Memory: model.NewMemory(storage),
	}

	c.ctx, c.cancel = context.WithCancel(context.Background())

	c.Sync()
	if c.SyncErr != nil {
		return
	}

	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			case <-time.After(time.Second):
				c.Sync()
			}
		}
	}()

	return
}

func (c *Core) roleCannotUseServiceOneLayer(role *model.Role, serviceToBeCheck string) (r bool) {
	r = true
	for _, serviceRoleHas := range role.Services {
		if serviceRoleHas.Path == serviceToBeCheck {
			r = false
			return
		}
	}
	return
}

func (c *Core) roleCanUseAllServicesInServiceTrace(role *model.Role, serviceTrace []string) (err error) {
	if role.Status == model.StatusOff {
		err = fmt.Errorf("角色[%s]未启用\n%w", role.Name, errors.ErrPermissionDenied)
		return
	}

	serviceStrToBeCheck := ""
	for _, serviceOneLayer := range serviceTrace {
		serviceStrToBeCheck = path.Join(serviceStrToBeCheck, serviceOneLayer)
		if service, ok := c.Services[serviceStrToBeCheck]; ok {
			if service.Status == model.StatusOn {
				if c.roleCannotUseServiceOneLayer(role, serviceStrToBeCheck) {
					err = fmt.Errorf("[%s]没有权限执行[%s]\n%w", role.Name, serviceStrToBeCheck, errors.ErrPermissionDenied)
					return
				}
			}
		}
	}

	return
}

func (c *Core) RoleCanUseService(roleName, servicePath string) (err error) {
	c.Lock()
	defer c.Unlock()
	if c.SyncErr != nil {
		err = c.SyncErr
		return
	}
	role, ok := c.Memory.Roles[roleName]
	if !ok {
		err = fmt.Errorf("角色[%s]\n%w", roleName, errors.ErrRecordNotExist)
		return
	}

	serviceTrace := strings.Split(servicePath, "/")

	err = c.roleCanUseAllServicesInServiceTrace(role, serviceTrace)

	return
}
