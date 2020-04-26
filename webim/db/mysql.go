package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"bitbucket.org/forfd/custm-chat/webim/conf"
)

var Mysql *sql.DB

var (
	defaultMySQLMaxIdle = 20
	defaultMySQLMaxConn = 100
)

func InitMysql(c *conf.MySQLConfig) error {
	if c.Dsn == "" {
		return fmt.Errorf("mysql dsn is empty")
	}

	pool, err := sql.Open("mysql", c.Dsn)
	if err != nil {
		return err
	}

	maxIdle := defaultMySQLMaxIdle
	if c.MaxIdle > 0 {
		maxIdle = c.MaxIdle
	}

	maxConn := defaultMySQLMaxConn
	if c.MaxConn > 0 {
		maxConn = c.MaxConn
	}

	pool.SetMaxIdleConns(maxIdle)
	pool.SetMaxOpenConns(maxConn)

	if c.ConnMaxLifeTime.Duration.Seconds() > 0 {
		pool.SetConnMaxLifetime(c.ConnMaxLifeTime.Duration)
	} else {
		pool.SetConnMaxLifetime(defaultConnMaxLifeTime)
	}
	Mysql = pool
	return nil
}
