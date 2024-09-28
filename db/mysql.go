package db

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"time"
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

func NewMysql(viper *viper.Viper, writer io.Writer) *MysqlDb {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.db"),
		viper.GetString("mysql.charset"),
	)
	newLogger := logger.New(
		log.New(writer, "\r\n", log.LstdFlags), // io newWriter
		logger.Config{
			SlowThreshold:             time.Duration(viper.GetInt("mysql.gorm.SlowThreshold")) * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.LogLevel(viper.GetInt("mysql.gorm.level")),                          // Log level
			IgnoreRecordNotFoundError: viper.GetBool("mysql.gorm.IgnoreRecordNotFoundError"),                      // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      viper.GetBool("mysql.gorm.ParameterizedQueries"),                           // Don't include params in the SQL log
			Colorful:                  viper.GetBool("mysql.gorm.Colorful"),                                       // Disable color
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
	sqlDb.SetConnMaxIdleTime(time.Duration(viper.GetInt64("mysql.MaxIdleTime")) * time.Second)
	//空闲连接池最大数量
	sqlDb.SetMaxIdleConns(viper.GetInt("mysql.MaxIdleConns"))
	//最大打开的连接数
	sqlDb.SetMaxOpenConns(viper.GetInt("mysql.MaxOpenConns"))
	//连接可复用的最长时间
	sqlDb.SetConnMaxLifetime(time.Duration(viper.GetInt64("mysql.MaxLifetime")) * time.Minute)
	return &MysqlDb{
		db: dbObj,
	}
}
