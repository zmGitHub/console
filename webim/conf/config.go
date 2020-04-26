package conf

import (
	"io/ioutil"
	"time"

	"github.com/BurntSushi/toml"
)

var (
	IMConf              *IMConfig
	DefaultPingInterval = 30 * time.Second
	MaxPingInterval     = 60 * time.Second
	MinPingInterval     = 25 * time.Second
)

type IMConfig struct {
	Debug             bool         `toml:"debug"`
	Listen            string       `toml:"listen"`
	Host              string       `toml:"host"`
	BackendHost       string       `toml:"backend_host"`
	APIToken          string       `toml:"api_token"`
	AllowOrigins      []string     `toml:"allow_origins"`
	IPGeoConf         *IPGeoConfig `toml:"ip_location_conf"`
	DefaultEntConfigs []byte
	JWTKey            string `toml:"jwt_key"`
	JWTKeyBytes       []byte
	AgentConf         *AgentConfig         `toml:"agent_conf"`
	MySQLConf         *MySQLConfig         `toml:"mysql_conf"`
	RedisConf         *RedisConfig         `toml:"redis_conf"`
	CentrifugoConf    *CentrifugoConfig    `toml:"centrifugo_conf"`
	ElasticSearchConf *ElasticSearchConfig `toml:"elasticsearch_conf"`
	EmailConf         *EmailConfig         `toml:"email_conf"`
	SubMailConf       *SubMailConfig       `toml:"submail_conf"`
	AWSS3Conf         *AWSS3Config         `toml:"aws_s3_conf"`
	TaskHandlerConf   *TaskHandlerConfig   `toml:"task_handler_conf"`
}

type IPGeoConfig struct {
	IPDBPath string `toml:"ipdb_path"`
}

// MySQLConfig MySQL config
type MySQLConfig struct {
	Dsn             string
	MaxConn         int
	MaxIdle         int
	ConnMaxLifeTime Duration
}

type RedisConfig struct {
	Addr         string   `toml:"addr"`
	ReadTimeout  Duration `toml:"read_timeout"`
	WriteTimeout Duration `toml:"write_timeout"`
}

type CentrifugoConfig struct {
	API                     string   `toml:"api"`
	AuthKey                 string   `toml:"authkey"`
	Timeout                 Duration `toml:"timeout"`
	PingInterval            Duration `toml:"ping_interval"`
	RetryTimes              int      `toml:"retry_times"`
	AgentMessage            string   `toml:"agent_message"`
	VisitorMessage          string   `toml:"visitor_message"`
	NewConvChannel          string   `toml:"new_conv_channel"`
	OnlineAgentsChannel     string   `toml:"online_agents"`
	VisitorQueueChannel     string   `toml:"visitor_queue"`
	ConversationInstruction string   `toml:"conversation_instruction"`
	InternalMessageChannel  string   `toml:"internal_msg"`
}

type ElasticSearchConfig struct {
	Endpoint string   `toml:"endpoint"`
	Username string   `toml:"username"`
	Password string   `toml:"password"`
	Region   string   `toml:"region"`
	Timeout  Duration `toml:"timeout"`
}

type AgentConfig struct {
	MaxOnlineCount                int64    `toml:"max_online_count"`
	ActivateCodeEffectiveDuration Duration `toml:"activate_code_effective_duration"`
	AgentTokenExpire              Duration `toml:"agent_token_expire"`
	AgentConnectionTokenExpire    Duration `toml:"agent_connection_token_expire"`
	AgentOnlineExpire             Duration `toml:"agent_online_expire"`
}

type EmailConfig struct {
	APIUser  string `toml:"api_user"`
	APIKey   string `toml:"api_key"`
	EndPoint string `toml:"endpoint"`
	From     string `toml:"from"`
	Subject  string `toml:"subject"`
	HtmlTmpl string `toml:"html_tmpl"`
}

type SubMailConfig struct {
	AppID    string `toml:"app_id"`
	AppKey   string `toml:"app_key"`
	SignType string `toml:"sign_type"`
}

type AWSS3Config struct {
	EndPoint        string `toml:"end_point"`
	AccessKeyID     string `toml:"access_key_id"`
	SecretAccessKey string `toml:"secret_access_key"`
	SessionToken    string `toml:"session_token"`
	UseSSL          bool   `toml:"use_ssl"`
	Region          string `toml:"region"`
	BucketName      string `toml:"bucket_name"`
}

type TaskHandlerConfig struct {
	Host                           string `toml:"host"`
	SendEntMessageEndPoint         string `toml:"send_ent_message_endpoint"`
	SendNoRespMessageEndPoint      string `toml:"send_no_resp_message_endpoint"`
	EndConversationEndPoint        string `toml:"end_conversation_endpoint"`
	OfflineEndConversationEndPoint string `toml:"offline_end_conversation_endpoint"`
	DeleteJobEndPoint              string `toml:"delete_job_endpoint"`
}

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(text []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func LoadConfig(configFile string) error {
	content, e := ioutil.ReadFile(configFile)
	if e != nil {
		return e
	}

	if _, err := toml.Decode(string(content), &IMConf); err != nil {
		return err
	}

	IMConf.JWTKeyBytes = []byte(IMConf.JWTKey)
	return nil
}

func LoadEntDefaultConfigs(configPath string) error {
	content, e := ioutil.ReadFile(configPath)
	if e != nil {
		return e
	}

	IMConf.DefaultEntConfigs = content
	return nil
}
