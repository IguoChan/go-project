package mysqlx

import (
	"errors"
	"fmt"
	"time"

	"github.com/IguoChan/go-project/pkg/util"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
)

type Client struct {
	*gorm.DB
}

type Options struct {
	Addr         string
	Username     string
	Password     string
	DBName       string
	MaxOpenConns int
	MaxIdleConns int

	// logger, you can use log.Logger or logrus.Logger, etc...
	Logger   gLogger.Writer
	LogLevel string

	// slow query threshold
	SlowThreshold time.Duration
}

func NewClient(opt *Options) (*Client, error) {
	if opt == nil {
		return nil, errors.New("options is nil")
	}

	// logger
	level := gLogger.Info
	switch opt.LogLevel {
	case "error":
		level = gLogger.Error
	case "warn":
		level = gLogger.Warn
	}
	logger := gLogger.New(
		opt.Logger,
		gLogger.Config{
			SlowThreshold: util.SetIf0(opt.SlowThreshold, time.Second),
			Colorful:      false,
			LogLevel:      level,
		},
	)

	// client
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
				opt.Username, opt.Password, opt.Addr, opt.DBName), // DSN data source name
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		}),
		&gorm.Config{
			SkipDefaultTransaction: true, // 开启提高性能，https://gorm.io/docs/transactions.html
			Logger:                 logger,
		},
	)
	if err != nil {
		return nil, err
	}

	// get sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(opt.MaxOpenConns)
	sqlDB.SetMaxOpenConns(opt.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Hour * 4)

	return &Client{db}, nil
}

func (c *Client) Close() {
	if c != nil {
		db, err := c.DB.DB()
		if err == nil {
			_ = db.Close()
		}
	}
}
