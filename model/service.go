package model

import (
	"fmt"
	"path"
	"time"
)

type Service struct {
	Path        string
	Description string
	CreateTime  time.Time
	EditTime    time.Time
	Status      Status
	memory      *Memory
}

func (s *Service) SetStatus(status Status) error {
	s.memory.Lock()
	defer s.memory.Unlock()

	s.Status = status
	s.EditTime = time.Now()
	err := s.memory.storage.ServiceUpdate(s)
	if err != nil {
		err = fmt.Errorf("service update\n%w", err)
	}
	return err
}
func (s *Service) SetDescription(description string) error {
	s.memory.Lock()
	defer s.memory.Unlock()

	s.Description = description
	s.EditTime = time.Now()
	err := s.memory.storage.ServiceUpdate(s)
	if err != nil {
		err = fmt.Errorf("service update\n%w", err)
	}
	return err
}

func (s *Service) NewSubService(name, description string) (service *Service, err error) {
	service, err = s.memory.NewService(path.Join(s.Path, name), description)
	return
}

func (s *Service) SetMemory(memory *Memory) {
	s.memory = memory
}
