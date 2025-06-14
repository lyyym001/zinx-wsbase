package global

import (
	"os"

	"github.com/lyyym/zinx-wsbase/zutil/zzap"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

//if Object.Mysql != nil {
//	// 程序结束前关闭数据库链接
//	db, _ := Mysql.DB()
//	defer db.Close()
//}

//- 用到mysql需要安装驱动`gorm.io/driver/mysql`

// InitGormMysql 只读数据库连接
func InitGormMysql() *gorm.DB {
	m := Object.MysqlConfig
	if m.Dbname == "" {
		panic("sql name is empty")
	}
	dsn := m.Username + ":" + m.Password + "@tcp(" + m.Path + ")/" + m.Dbname + "?" + m.Config
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig()); err != nil {
		Glog.Error("MySQL启动异常", zap.Error(err))
		os.Exit(0)
		//return nil
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		Mysql = db
		return db
	}
}

//@author: SliverHorn
//@function: gormConfig
//@description: 根据配置决定是否开启日志
//@param: mod bool
//@return: *gorm.Config

func gormConfig() *gorm.Config {
	config := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{SingularTable: true}}

	if Glog == nil {
		panic("日志库未初始化")
	}
	//zapLogger := zapgorm2.New(zap.L())
	zapLogger := zzap.New(Glog.Named("sql"))
	// optional: configure gorm to use this zapgorm.Logger for callbacks
	//zapLogger.SetAsDefault()
	switch Object.MysqlConfig.LogMode {
	case "silent", "Silent":
		zapLogger.LogLevel = logger.Silent
	case "error", "Error":
		zapLogger.LogLevel = logger.Error
	case "warn", "Warn":
		zapLogger.LogLevel = logger.Warn
	case "info", "Info":
		zapLogger.LogLevel = logger.Error
	default:
		zapLogger.LogLevel = logger.Info
	}
	config.Logger = zapLogger
	return config
}
