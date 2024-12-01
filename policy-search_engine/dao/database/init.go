package database

import (
	"PolicySearchEngine/config"
	"PolicySearchEngine/model"
	mysqlCfg "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

var myDb *gorm.DB

func InitTable() {
	// 初始化数据表
	_ = myDb.AutoMigrate(&model.Meta{}) //使用 AutoMigrate 方法来创建或更新 Meta 表。_ = 表示忽略返回值，因为 AutoMigrate 返回一个错误，但我们在这里不处理它。
	_ = myDb.AutoMigrate(&model.Content{})
	_ = myDb.AutoMigrate(&model.Department{})
	_ = myDb.AutoMigrate(&model.Province{})
	_ = myDb.AutoMigrate(&model.SmallDepartmentMap{})
	_ = myDb.AutoMigrate(&model.SmallDepartment{})
}

func Init() {
	// 数据库配置
	cfg := mysqlCfg.Config{
		User:      config.V.GetString("mysql.user"),
		Passwd:    config.V.GetString("mysql.password"),
		Net:       "tcp",
		Addr:      config.V.GetString("mysql.addr"),
		DBName:    config.V.GetString("mysql.dbname"),
		Loc:       time.Local,
		ParseTime: true,
		// 允许原生密码
		AllowNativePasswords: true,
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(cfg.FormatDSN()),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Error), //设置日志模式为 Error，即只记录错误日志。
		}) //&gorm.Config{ ... } 设置 GORM 的配置，包括日志模式。
	if err != nil {
		log.Fatal(err)
	}

	myDb = db //将连接实例赋值给全局变量 myDb，以便其他部分的代码可以使用这个连接实例进行数据库操作。
	return
}

func MyDb() *gorm.DB {
	return myDb
} //MyDb 函数提供了一个方便的接口，让其他部分的代码可以获取到已经初始化的数据库连接实例，
// 从而进行数据库操作。这样可以避免在多个地方重复初始化数据库连接，提高了代码的复用性和可维护性。
