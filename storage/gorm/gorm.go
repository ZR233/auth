/*
@Time : 2019-08-27 14:22
@Author : zr
*/
package gorm

import (
	"github.com/ZR233/auth/model"
	gorm2 "github.com/jinzhu/gorm"
)

type Storage struct {
	db *gorm2.DB
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
	r = &Role{
		Name:        role.Name,
		State:       role.State,
		Description: role.Description,
		EditTime:    role.EditTime,
		CreateTime:  role.CreateTime,
	}

	return
}
func newService(service *model.Service) (s *Service) {
	s = &Service{
		Name:         service.Name,
		SuperiorName: service.SupService,
		Description:  service.Description,
		CreateTime:   service.CreateTime,
		EditTime:     service.EditTime,
	}

	return
}

func (s *Storage) SaveRole(role *model.Role) (err error) {
	r := newRole(role)
	return s.db.Save(r).Error
}
func (s *Storage) SaveService(service *model.Service) (err error) {
	service_ := newService(service)
	return s.db.Save(service_).Error
}

func (s *Storage) SaveRelation(service *model.Service, roles ...*model.Role) (err error) {
	serviceName := service.Name
	tx := s.db.Begin()
	defer func() {
		if err_ := recover(); err_ != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	for _, role := range roles {
		roleName := role.Name
		r := &ServiceRoleRelation{
			ServiceName: serviceName,
			RoleName:    roleName,
		}

		err = s.db.Save(r).Error
		if err != nil {
			panic(err)
		}
	}
	return
}

func (s *Storage) Sync() (serviceTree map[string]*model.Service, roles []*model.Role, err error) {
	var services []*Service
	err = s.db.Find(&services).Error
	if err != nil {
		return
	}
	var serviceO []*model.Service
	for _, service := range services {
		var relations []*ServiceRoleRelation
		err = s.db.Where(&ServiceRoleRelation{ServiceName: service.Name}).Find(&relations).Error
		if err != nil {
			return
		}
		var roles []*model.Role

		for _, relation := range relations {
			var role_ Role
			err = s.db.Where(&Role{Name: relation.RoleName, State: 1}).Find(&role_).Error
			if err != nil {
				return
			}
			role := &model.Role{
				Name:        role_.Name,
				State:       role_.State,
				Description: role_.Description,
				EditTime:    role_.EditTime,
				CreateTime:  role_.CreateTime,
				Storage:     s,
			}
			roles = append(roles, role)
		}

		serviceO = append(serviceO, &model.Service{
			Name:        service.Name,
			SupService:  service.SuperiorName,
			SubService:  nil,
			Description: service.Description,
			CreateTime:  service.CreateTime,
			EditTime:    service.EditTime,
			Roles:       roles,
			Storage:     s,
		})
	}

	serviceTree = make(map[string]*model.Service)
	for _, v := range serviceO {
		serviceTreeInsert(serviceTree, v, serviceO)
	}
	var roles_ []*Role
	err = s.db.Where(&Role{State: 1}).Find(&roles_).Error
	if err != nil {
		return
	}
	for _, role := range roles_ {

		roles = append(roles, &model.Role{
			Name:        role.Name,
			State:       role.State,
			Description: role.Description,
			EditTime:    role.EditTime,
			CreateTime:  role.CreateTime,
			Services:    nil,
			Storage:     s,
		})
	}

	return
}

func findSupService(service *model.Service, services []*model.Service) (sub *model.Service) {
	for _, v := range services {
		if v.Name == service.SupService {
			return v
		}
	}
	return
}

func serviceGetTrace(service *model.Service, services []*model.Service) (trace []*model.Service) {
	for {
		trace = append([]*model.Service{service}, trace...)
		service = findSupService(service, services)
		if service == nil || service.Name == "" {
			return
		}
	}
}

func serviceTreeInsert(serviceTree map[string]*model.Service, service *model.Service, services []*model.Service) {
	trace := serviceGetTrace(service, services)
	if len(trace) == 0 {
		return
	}

	service_, ok := serviceTree[trace[0].Name]
	if !ok {
		serviceTree[trace[0].Name] = trace[0]
		service_ = trace[0]
	}
	trace = trace[1:]

	for _, v := range trace {
		if service_.SubService == nil {
			service_.SubService = make(map[string]*model.Service)
		}

		_, ok := service_.SubService[v.Name]
		if !ok {
			service_.SubService[v.Name] = v
		}
		service_ = service_.SubService[v.Name]
	}

}
