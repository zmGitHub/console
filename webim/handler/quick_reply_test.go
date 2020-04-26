package handler

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"bitbucket.org/forfd/custm-chat/webim/models"
)

func TestParseQuickReplyFile(t *testing.T) {
	h := &IMService{}
	f, err := os.Open("template.xls")
	assert.Nil(t, err)
	defer func() {
		f.Close()
	}()

	groups, err := h.parseQuickReplyFile("entID", "agentID", f, models.QuickReplyEnterpriseCreatorType)
	for k, v := range groups {
		fmt.Printf("group: %s\n groupInfo: %+v\n", k, *v.QuickreplyGroup)
		fmt.Printf("items: \n")
		for _, item := range v.Items {
			fmt.Printf("item: %s\n info: %+v\n", item.Title, *item)
		}
	}
}

func TestTimeFormat(t *testing.T) {
	tm := time.Now()
	layout := "2006-01-02 15:04:05"
	fileName := fmt.Sprintf("快速回复-企通-%s-%s.xlsx", tm.Format(layout), "ent")
	fmt.Println(fileName)
}
