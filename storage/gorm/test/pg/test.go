/*
@Time : 2019-08-27 14:30
@Author : zr
*/
package mysql

import (
	"fmt"
	gorm2 "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func TestDbPG() *gorm2.DB {

	user := "sa"
	password := "asdf*123"
	host := "192.168.0.3"
	port := "3306"
	databaseName := "auth_test"
	charset := "utf8mb4"
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		user, password, host, port, databaseName, charset)

	client, err := gorm2.Open("mysql", connStr)
	if err != nil {
		panic(err)
	}

	return client
}
