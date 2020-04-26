package test

import (
	"fmt"
	stdLog "log"
	"sync"

	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/external/elasticsearch"
)

var (
	tables = []string{
		"enterprise",
		"agent",
		"message_beep",
		"agent_group",
		"agent_group_relation",
		"automatic_message",
		"visit_blacklist",
		"queue_config",
		"ending_conversation",
		"ending_message",
		"conversation_transfer",
		"conversation_quality",
		"login_limit",
		"send_file",
		"leave_message_config",
		"conversation",
		"evaluation",
		"prechat_form",
		"leave_message",
		"login_records",
		"message",
		"ent_app",
		"perm",
		"role",
		"role_perm",
		"quickreply_group",
		"quickreply_item",
		"allocation_rule",
		"agent_statistic",
		"visitor_statistic",
		"conversation_statistic",
		"visit",
		"visit_page",
		"visitor",
		"visitor_tag",
		"visitor_tag_relation",
	}

	IpLocation *db.IPGeo

	initOnce sync.Once
)

func InitTest() {
	initOnce.Do(func() {
		var err error
		if err := conf.LoadConfig("../conf/config_test.toml"); err != nil {
			panic(err)
		}

		if err := db.InitMysql(conf.IMConf.MySQLConf); err != nil {
			stdLog.Println("load MySQL error: ", err)
			return
		}

		if err := db.NewRedisClient(conf.IMConf.RedisConf); err != nil {
			stdLog.Println("load Redis error: ", err)
			return
		}

		if err := elasticsearch.NewESClient(conf.IMConf.ElasticSearchConf); err != nil {
			stdLog.Println("load elasticsearch client error: ", err)
		}

		conf.IMConf.IPGeoConf.IPDBPath = "../db/ipiptest.ipdb"
		IpLocation, err = db.LoadIPGeoDB(conf.IMConf.IPGeoConf)
		if err != nil {
			stdLog.Println("LoadIPGeoDB client error: ", err)
		}

		log.NewLogging()
	})
}

func Clear() {
	for _, table := range tables {
		_, err := db.Mysql.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table))
		if err != nil {
			stdLog.Println(err)
		}
	}
}
