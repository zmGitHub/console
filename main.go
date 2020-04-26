package main

import (
	"flag"
	"log"
	"time"

	"github.com/dchest/captcha"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/auth"
	logging "bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/external/awss3"
	"bitbucket.org/forfd/custm-chat/webim/external/elasticsearch"
	"bitbucket.org/forfd/custm-chat/webim/external/submail"
	"bitbucket.org/forfd/custm-chat/webim/external/timedevent"
	"bitbucket.org/forfd/custm-chat/webim/handler"
	"bitbucket.org/forfd/custm-chat/webim/handler/monitor"
	"bitbucket.org/forfd/custm-chat/webim/imclient"
	"bitbucket.org/forfd/custm-chat/webim/router"
)

func main() {
	config := flag.String("config", "webim/conf/config_test.toml", "config path")
	ipDBPath := flag.String("ipdb", "webim/db/ipiptest.ipdb", "ipdb path")
	defaultConfig := flag.String("default-config", "webim/conf/default_config.json", "default config file")
	flag.Parse()

	if *config == "" || *defaultConfig == "" {
		log.Println("config path/default-config is not set")
		return
	}

	if err := conf.LoadConfig(*config); err != nil {
		log.Fatalln("load config error: ", err)
	}
	conf.IMConf.IPGeoConf.IPDBPath = *ipDBPath

	if err := conf.LoadEntDefaultConfigs(*defaultConfig); err != nil {
		log.Println("LoadEntDefaultConfigs: ", err)
		return
	}

	logging.NewLogging()

	if err := db.InitMysql(conf.IMConf.MySQLConf); err != nil {
		log.Println("load MySQL error: ", err)
		return
	}

	if err := db.NewRedisClient(conf.IMConf.RedisConf); err != nil {
		log.Println("load Redis error: ", err)
		return
	}

	captcha.SetCustomStore(&db.RedisStore{Expire: 2 * time.Minute})

	imclient.InitClient(conf.IMConf.CentrifugoConf)
	uploadCli := awss3.NewUploadClient(conf.IMConf.AWSS3Conf)
	emailCli := submail.NewClient(conf.IMConf.SubMailConf)
	if err := elasticsearch.NewESClient(conf.IMConf.ElasticSearchConf); err != nil {
		log.Fatalln("load elasticsearch client error: ", err)
	}

	ipGeo, err := db.LoadIPGeoDB(conf.IMConf.IPGeoConf)
	if err != nil {
		log.Println("LoadIPGeoDB error: ", err)
	}

	monitor.RegisterMetrics()

	go handler.SyncStatus()

	service := handler.NewIMService(
		handler.WithAuth(auth.NewRedisAgentAuth(db.RedisClient)),
		handler.WithEmailSender(emailCli),
		handler.WithIMClient(imclient.CentriClient),
		handler.WithUploader(uploadCli),
		handler.WithIPGeoLocation(ipGeo),
		handler.WithTaskHandler(timedevent.InitHandler(conf.IMConf.TaskHandlerConf)),
	)
	server := echo.New()
	router.SetRouter(server, service)

	logging.Logger.Infof("Starting Server With Config: %+v", *conf.IMConf)
	server.Server.Addr = conf.IMConf.Listen
	server.Logger.Fatal(gracehttp.Serve(server.Server))
}
