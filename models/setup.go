package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// URL := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, DBNAME)
	// root:sample-password@tcp(db:3306)/jwt_auth
	db, err := gorm.Open(mysql.Open("root:gobankpassword@tcp(localhost:3306)/gobank"))

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Account{})
	DB = db
}
