package conf

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	//Timeout := "10s"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", Username, Password, Host, Databese)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("连接数据库失败！！" + err.Error())
		return
	}
	DB = db
}
