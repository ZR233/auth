package model

import (
	"fmt"
	"time"
)

type Role struct {
	Name        string
	Description string
	EditTime    time.Time
	CreateTime  time.Time
	Status      Status
	Services    []*Service `gorm:"-"`
	memory      *Memory
}

func (r *Role) SetStatus(status Status) error {
	r.memory.Lock()
	defer r.memory.Unlock()

	r.Status = status
	r.EditTime = time.Now()
	err := r.memory.storage.RoleUpdate(r)
	if err != nil {
		err = fmt.Errorf("role update\n%w", err)
	}
	return err
}
func (r *Role) SetDescription(description string) error {
	r.memory.Lock()
	defer r.memory.Unlock()

	r.Description = description
	r.EditTime = time.Now()
	err := r.memory.storage.RoleUpdate(r)
	if err != nil {
		err = fmt.Errorf("role update\n%w", err)
	}
	return err
}

func (r *Role) SetMemory(memory *Memory) {
	r.memory = memory
}

func (r *Role) SetServices(servicePath ...string) (err error) {
	err = r.memory.RoleUpdateServices(r.Name, servicePath...)
	return
}
