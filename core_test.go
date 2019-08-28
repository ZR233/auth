/*
@Time : 2019-08-27 14:53
@Author : zr
*/
package auth

import (
	"github.com/ZR233/auth/storage/gorm"
	"testing"
)

func TestGorm(t *testing.T) {

	st := gorm.NewStorage(gorm.DbForTest())

	core := NewCore(st)
	normal := core.NewRole("普通用户")
	err := normal.Save()
	if err != nil {
		t.Error(err)
	}
	admin := core.NewRole("管理员")
	err = admin.Save()
	if err != nil {
		t.Error(err)
	}

	userS := core.NewService("User")
	infoS := userS.NewSubService("Info")
	editS := userS.NewSubService("Edit")
	okButtonS := editS.NewSubService("OKButton")

	err = userS.Save()
	if err != nil {
		t.Error(err)
	}
	err = infoS.Save()
	if err != nil {
		t.Error(err)
	}
	err = editS.Save()
	if err != nil {
		t.Error(err)
	}
	err = okButtonS.Save()
	if err != nil {
		t.Error(err)
	}
	err = userS.AddRoles(admin)
	if err != nil {
		t.Error(err)
	}
	err = userS.AddRoles(normal)
	if err != nil {
		t.Error(err)
	}
	err = editS.AddRoles(admin)
	if err != nil {
		t.Error(err)
	}
}

func TestCore_Check(t *testing.T) {
	st := gorm.NewStorage(gorm.DbForTest())

	c := NewCore(st)
	err := c.Sync()
	if err != nil {
		t.Error(err)
	}

	type args struct {
		ServiceUrl string
		roleName   string
	}
	tests := []struct {
		name  string
		args  args
		wantR bool
	}{
		{"1", args{"Task", "管理员"}, true},
		{"2", args{"User", "管理员"}, true},
		{"3", args{"User/Edit", "管理员"}, true},
		{"4", args{"User/Edit", "普通用户"}, false},
		{"5", args{"User/Edit/OKButton", "管理员"}, true},
		{"6", args{"User/Edit/OKButton", "普通用户"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := c.Check(tt.args.ServiceUrl, tt.args.roleName); gotR != tt.wantR {
				t.Errorf("Check() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}
