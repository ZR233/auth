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
