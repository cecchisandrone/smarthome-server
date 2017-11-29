package persistence

import (
	"fmt"
	"os"

	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

func Init() *gorm.DB {

	// Load db connection parameters
	user := os.Getenv("MYSQL_SMARTHOME_USER")
	password := os.Getenv("MYSQL_SMARTHOME_PASSWORD")
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")

	//open a db connection
	var err error
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/smarthome?charset=utf8&parseTime=True&loc=Local", user, password, host, port))
	if err != nil {
		fmt.Println(err)
		panic("Failed to connect database")
	}

	// Enable log
	db.LogMode(true)

	//Migrate the schema
	db.AutoMigrate(&model.Profile{}, &model.Configuration{}, &model.Camera{})
	db.Model(&model.Profile{}).AddForeignKey("configuration_id", "configurations(id)", "CASCADE", "CASCADE")
	db.Model(&model.Camera{}).AddForeignKey("configuration_id", "configurations(id)", "CASCADE", "CASCADE")

	return db
}