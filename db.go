package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func dbInit() {
	//open a db connection
	var err error
	db, err = gorm.Open("mysql", "root:smarthome@tcp(192.168.99.100:4306)/smarthome?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	//Migrate the schema
	db.AutoMigrate(&todoModel{})
}
