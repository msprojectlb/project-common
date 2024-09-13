package db

import (
	"context"
	"fmt"
	"github.com/msprojectlb/project-common/config"
	"github.com/msprojectlb/project-common/logs/writer"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type DB interface {
	Db(ctx context.Context) *gorm.DB
}
type MysqlDb struct {
	db *gorm.DB
}

func (m *MysqlDb) Db(ctx context.Context) *gorm.DB {
	return m.db.Session(&gorm.Session{Context: ctx})
}

func NewMysql(mysqlConf config.MysqlConfig, gormConf config.GormConfig, writer *writer.GormWriter) *MysqlDb {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		mysqlConf.User,
		mysqlConf.Pwd,
		mysqlConf.Ip,
		mysqlConf.Port,
		mysqlConf.Db,
		mysqlConf.CharSet,
	)
	newLogger := logger.New(
		log.New((*lumberjack.Logger)(writer), "\r\n", log.LstdFlags), // io newWriter
		logger.Config{
			SlowThreshold:             gormConf.SlowThreshold,             // Slow SQL threshold
			LogLevel:                  gormConf.LogLevel,                  // Log level
			IgnoreRecordNotFoundError: gormConf.IgnoreRecordNotFoundError, // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      gormConf.ParameterizedQueries,      // Don't include params in the SQL log
			Colorful:                  gormConf.Colorful,                  // Disable color
		},
	)
	dbObj, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("Mysql连接失败", err.Error())
	}
	sqlDb, err := dbObj.DB()
	if err != nil {
		log.Fatal("数据库获取失败", err.Error())
	}
	//最大空闲连接时间
	sqlDb.SetConnMaxIdleTime(mysqlConf.MaxIdleTime)
	//空闲连接池最大数量
	sqlDb.SetMaxIdleConns(mysqlConf.MaxIdleConns)
	//最大打开的连接数
	sqlDb.SetMaxOpenConns(mysqlConf.MaxOpenConns)
	//连接可复用的最长时间
	sqlDb.SetConnMaxLifetime(mysqlConf.MaxLifetime)
	return &MysqlDb{
		db: dbObj,
	}
}
