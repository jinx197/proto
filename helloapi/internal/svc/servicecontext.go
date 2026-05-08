// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"helloapi/internal/config"
	"helloapi/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.MysqlDB.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// 自动建表
	db.AutoMigrate(&model.Employee{})
	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
