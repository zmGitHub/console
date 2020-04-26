package main

import (
	"flag"
	"fmt"
	"time"

	"bitbucket.org/forfd/custm-chat/webim/conf"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/dto"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

func main() {
	pwd := flag.String("pwd", "", "db pwd")
	ip := flag.String("ip", "", "db ip")
	entID := flag.String("entid", "", "ent id")
	roleID := flag.String("roleid", "", "role id")
	flag.Parse()

	mysqlconf := &conf.MySQLConfig{
		Dsn:             fmt.Sprintf("root:%s@tcp(%s:3306)/custmchat?charset=utf8mb4&parseTime=True&loc=Local", *pwd, *ip),
		MaxConn:         100,
		MaxIdle:         20,
		ConnMaxLifeTime: conf.Duration{Duration: 14400 * time.Second},
	}

	if err := db.InitMysql(mysqlconf); err != nil {
		fmt.Println("InitMysql: ", err)
		return
	}

	update(*entID, *roleID)
}

func update(entID, roleID string) {
	if entID == "" || roleID == "" {
		return
	}

	tx, err := db.Mysql.Begin()
	if err != nil {
		fmt.Println("begin tx: ", err)
		return
	}

	var dbErr error
	defer func() {
		if dbErr != nil {
			tx.Rollback()
			return
		}

		tx.Commit()
	}()

	if _, dbErr = tx.Exec("DELETE FROM custmchat.ent_app WHERE ent_id = ?", entID); dbErr != nil {
		return
	}

	if _, dbErr = tx.Exec("DELETE FROM custmchat.perm WHERE ent_id = ?", entID); dbErr != nil {
		return
	}

	var apps []string
	for app := range dto.APPPerms {
		apps = append(apps, app)
	}
	if dbErr = models.BulkCreateEntApps(tx, entID, apps); dbErr != nil {
		return
	}

	var permIDs []string
	permIDs, dbErr = models.BulkCreateEntPerms(tx, entID, dto.APPPerms)
	if dbErr != nil {
		return
	}

	if dbErr = models.UpdateRolePerms(tx, roleID, permIDs); dbErr != nil {
		return
	}
	return
}
