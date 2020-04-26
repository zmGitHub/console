package db

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"bitbucket.org/forfd/custm-chat/webim/conf"
)

var (
	defaultMaxConn         = 100
	defaultMaxIdle         = 20
	defaultConnMaxLifeTime = 14400 * time.Second
)

var Conn *gorm.DB

func NewDBConn(c *conf.IMConfig) error {
	db, err := gorm.Open("mysql", c.MySQLConf.Dsn)
	if err != nil {
		return err
	}

	if c.MySQLConf.MaxConn > 0 {
		db.DB().SetMaxOpenConns(c.MySQLConf.MaxConn)
	} else {
		db.DB().SetMaxOpenConns(defaultMaxConn)
	}

	if c.MySQLConf.MaxIdle > 0 {
		db.DB().SetMaxIdleConns(c.MySQLConf.MaxIdle)
	} else {
		db.DB().SetMaxIdleConns(defaultMaxIdle)
	}

	if c.MySQLConf.ConnMaxLifeTime.Duration.Seconds() > 0 {
		db.DB().SetConnMaxLifetime(c.MySQLConf.ConnMaxLifeTime.Duration)
	} else {
		db.DB().SetConnMaxLifetime(defaultConnMaxLifeTime)
	}

	Conn = db
	return nil
}
