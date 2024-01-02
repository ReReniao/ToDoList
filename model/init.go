package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var DB *gorm.DB
var mysqllog logger.Interface

func Database(connstring string) {
	mysqllog = logger.Default.LogMode(logger.Info)

	conf := gorm.Config{Logger: mysqllog,
		NamingStrategy: schema.NamingStrategy{
			// 	TablePrefix: "f_", // 表名前缀
			SingularTable: true, // 是否单数表名
			// 	NoLowerCase: true, // 不要大写
		},
	}

	// 如果是发行版则不输出数据库日志
	if gin.Mode() == "release" {
		conf.Logger = nil
	}

	db, err := gorm.Open(mysql.Open(connstring), &conf)
	if err != nil {
		panic("Mysql数据库连接错误")
	}
	fmt.Println("数据库连接成功")
	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(20)  // 设置连接池
	sqlDb.SetMaxOpenConns(100) // 最大连接数
	sqlDb.SetConnMaxLifetime(time.Second * 30)
	DB = db
	migration()
}
