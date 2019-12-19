/*
@Time : 2019-08-27 14:53
@Author : zr
*/
package auth

import (
	"github.com/ZR233/auth/model"
	"github.com/ZR233/auth/storage/gorm"
	"github.com/ZR233/auth/storage/gorm/test/mysql"
	"testing"
)

func testNewCore() *Core {
	st := gorm.NewStorage(mysql.TestDbMysql())
	core, err := NewCore(st)
	if err != nil {
		panic(err)
	}

	return core
}

func TestCoreDelete(t *testing.T) {

	core := testNewCore()

	testClear(core)

	if len(core.Services) != 0 {
		t.Errorf("serivce len %d", len(core.Services))
		return
	}
	if len(core.Roles) != 0 {
		t.Errorf("Roles len %d", len(core.Roles))
		return
	}
}
func testClear(core *Core) {

	_ = core.DeleteRole("普通用户")
	_ = core.DeleteRole("管理员")

	_ = core.DeleteService("User")
	_ = core.DeleteService("User/Info")
	_ = core.DeleteService("User/Edit")
	_ = core.DeleteService("User/Edit/OKButton")
}

func TestCoreCreate(t *testing.T) {
	core := testNewCore()
	testClear(core)

	normal, err := core.NewRole("普通用户", "测试")
	if err != nil {
		t.Error(err)
		return
	}
	admin, err := core.NewRole("管理员", "测试")
	if err != nil {
		t.Error(err)
		return
	}

	userS, err := core.NewService("User", "测试")
	if err != nil {
		t.Error(err)
		return
	}
	infoS, err := userS.NewSubService("Info", "测试")
	if err != nil {
		t.Error(err)
		return
	}
	editS, err := userS.NewSubService("Edit", "测试")
	if err != nil {
		t.Error(err)
		return
	}
	okButtonS, err := editS.NewSubService("OKButton", "测试")
	if err != nil {
		t.Error(err)
		return
	}

	err = admin.SetServices(userS.Path, editS.Path, infoS.Path, okButtonS.Path)
	if err != nil {
		t.Error(err)
		return
	}
	err = normal.SetServices(userS.Path, okButtonS.Path)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestCore_Check(t *testing.T) {
	c := testNewCore()

	type args struct {
		ServiceUrl string
		roleName   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"1", args{"Task", "管理员"}, false},
		{"2", args{"User", "管理员"}, false},
		{"3", args{"User", "普通用户"}, false},
		{"4", args{"User/Edit", "管理员"}, false},
		{"5", args{"User/Edit", "普通用户"}, true},
		{"6", args{"User/Edit/OKButton", "管理员"}, false},
		{"7", args{"User/Edit/OKButton", "普通用户"}, true},
		{"8", args{"User/Edit/OKButton/Test", "管理员"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := c.RoleCanUseService(tt.args.roleName, tt.args.ServiceUrl); gotR != nil && !tt.wantErr {
				t.Errorf("Check() = %v, want %v", gotR, tt.wantErr)
			}
		})
	}
}

func TestCore_Set(t *testing.T) {
	core := testNewCore()
	defer func() {
		core.DeleteRole("普通用户2")
		core.DeleteService("User2")
	}()

	normal, err := core.NewRole("普通用户2", "测试")
	if err != nil {
		t.Error(err)
		return
	}
	userS, err := core.NewService("User2", "测试")
	if err != nil {
		t.Error(err)
		return
	}

	err = normal.SetServices(userS.Path)
	if err != nil {
		t.Error(err)
		return
	}

	err = normal.SetStatus(model.StatusOff)
	if err != nil {
		t.Error(err)
		return
	}

	err = core.RoleCanUseService(normal.Name, "User")
	t.Log(err)
	if err == nil {
		t.Error(err)
		return
	}
}
