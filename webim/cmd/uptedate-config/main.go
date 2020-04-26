package main

import (
	"database/sql"
	"flag"
	"fmt"
	"time"

	"github.com/go-redis/redis"

	_ "github.com/go-sql-driver/mysql"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/handler/adapter"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

func main() {
	dbpass := flag.String("dbpass", "", "mysql dsn")
	dbhost := flag.String("host", "", "mysql dsn")
	redisAddr := flag.String("redis", "", "redis host port")
	flag.Parse()

	if *dbpass == "" || *dbhost == "" || *redisAddr == "" {
		fmt.Println("dsn or redis is empty")
		return
	}

	redisClient := redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         *redisAddr,
		DB:           0,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})

	s := `root:%s@tcp(%s)/custmchat?charset=utf8mb4&parseTime=True&loc=Local`
	pool, err := sql.Open("mysql", fmt.Sprintf(s, *dbpass, *dbhost))
	if err != nil {
		fmt.Println("connect mysql error: ", err)
		return
	}

	db := pool

	entIDs, err := getAllEnts(db)
	if err != nil {
		fmt.Println("getAllEnts ids error: ", err)
		return
	}

	for _, entID := range entIDs {
		config, err := models.EntAllConfigByEntID(db, entID)
		if err != nil {
			fmt.Println("EntAllConfigByEntID error: ", err)
			continue
		}

		var configs = &adapter.EnterpriseConfigs{}

		if err = common.Unmarshal(config.ConfigContent.String, &configs); err != nil {
			fmt.Println("ent: ", entID, " Unmarshal ent config error: ", err)
			continue
		}

		if configs.WidgetSettings != nil {
			widgetConfig := configs.WidgetSettings
			widgetConfig.RemoveBrand = "open"
			widgetConfig.TicketOnly = "close"
		}

		if configs.StandaloneWindowConfig != nil {
			configs.StandaloneWindowConfig.RemoveBrand = "open"
			configs.StandaloneWindowConfig.Ring = "open"
		}

		configsContent, err := common.Marshal(configs)
		if err != nil {
			fmt.Println("ent: ", entID, " Marshal error: ", err)
			continue
		}

		config.ConfigContent.String = configsContent
		if err := models.CreateOrUpdateConfigs(db, config); err != nil {
			fmt.Println("ent: ", entID, " CreateOrUpdateConfigs error: ", err)
			continue
		}

		_, err = redisClient.HSet(fmt.Sprintf(common.EntConfigs, entID), common.EntConfigsContent, configsContent).Result()
		if err != nil {
			fmt.Println("ent: ", entID, " redisClient.HSet error: ", err)
			continue
		}

		fmt.Println("Success Update Ent: ", entID, " Configs")
	}

	fmt.Println("Success Update All Ent Configs")
}

func getAllEnts(db *sql.DB) (ids []string, err error) {
	query := "SELECT id FROM custmchat.enterprise"

	rows, err := db.Query(query)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			return
		}

		ids = append(ids, id)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}
