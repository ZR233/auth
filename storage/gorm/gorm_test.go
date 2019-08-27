/*
@Time : 2019-08-27 14:29
@Author : zr
*/
package gorm

import (
	"testing"
)

func TestNewStorage(t *testing.T) {

	NewStorage(DbForTest())
}

func TestStorage_Sync(t *testing.T) {
	st := NewStorage(DbForTest())

	_, err := st.Sync()
	if err != nil {
		t.Error(err)
	}
}
