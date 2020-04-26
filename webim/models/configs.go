package models

import (
	"database/sql"
	"fmt"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/db"
)

func CreateOrUpdateConfigs(db XODB, config *EntAllConfig) error {
	insert := `INSERT INTO custmchat.ent_all_configs(id, ent_id, config_content, create_at, update_at) ` +
		`VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE config_content=VALUES(config_content)`

	_, err := db.Exec(insert, config.ID, config.EntID, config.ConfigContent, config.CreateAt, config.UpdateAt)
	return err
}

func GetConfigsFromCache(mysql XODB, entID string) (content string, err error) {
	entKey := fmt.Sprintf(common.EntConfigs, entID)
	content, err = db.RedisClient.HGet(entKey, common.EntConfigsContent).Result()
	if err != nil {
		config, err := EntAllConfigByEntID(mysql, entID)
		if err != nil {
			if err == sql.ErrNoRows {
				return "", nil
			}
			return "", err
		}

		content = config.ConfigContent.String
		if config.ConfigContent.Valid {
			db.RedisClient.HSet(entKey, common.EntConfigsContent, content)
		}
		return content, nil
	}
	return
}
