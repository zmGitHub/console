# 部署

## 依赖组件

* Redis

* MySQL 5.7

* Elasticsearch 6.5.2

* Centrifugo (https://github.com/centrifugal/centrifugo)
  
  -- im server 依赖redis

* timed-task-handler (https://bitbucket.org/forfd/timed-task-handler)
  
  -- 定时任务处理 http service, 依赖MySQL
  
* Aws s3
  
  文件上传
  文件导出/导入
  读: 所有权限
  put: private
  
  
## Docker 部署

* Redis
    
    https://hub.docker.com/_/redis/
    
* MySQL

    https://hub.docker.com/_/mysql

* Centrifugo(https://hub.docker.com/r/centrifugo/centrifugo/dockerfile)
    
    config.json: doc/centrifugo_config.json
  
  `docker run --ulimit nofile=65536:65536 -v /host/dir/with/config/file:/centrifugo -p 8000:8000 -p 9000:9000 centrifugo/centrifugo centrifugo -c config.json`
  
* ES

```bash
docker run -d --name elasticsearch --net somenetwork -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" elasticsearch:6.5.4
```

在docker run完之后, 进入docker

```bash
./bin/elasticsearch-plugin install https://github.com/medcl/elasticsearch-analysis-ik/releases/download/v6.5.4/elasticsearch-analysis-ik-6.5.4.zip
```

最后 restart docker 容器

* timed-task-handler

conf.toml

```toml
listen = ":8989"
time_range_in_sec = 1 # 1秒

[mysql_conf]
Dsn = "root:D3v0ps20!I@tcp(host:port)/timedtask?charset=utf8mb4&parseTime=True&loc=Local"
MaxConn = 100
MaxIdle = 20
ConnMaxLifeTime = "8h"

[call_back_conf]
host = "http://127.0.0.1:8080"
timeout = "2s"
send_auto_msg_endpoint = "/api/v1/system/send_ent_message"
send_no_resp_endpoint = "/api/v1/system/send_no_resp_message"
end_conversation_endpoint = "/api/v1/system/end_conversation"

```

```bash
make build
docker build -t timed-task-handler:v1 .
docker run -d -p 8989:8989 timed-task-handler:v1 # -config config.toml
```

* custm-chat

config.toml

```toml
debug = true
listen = ":8090"
jwt_key = "$12secret34$"

[agent_conf]
max_online_count = 10
activate_code_effective_duration = "30m"
agent_token_expire = "24h"
agent_connection_token_expire = "2m"

[mysql_conf]
Dsn = "root:@tcp(host:port)/custmchat?charset=utf8mb4&parseTime=True&loc=Local"
MaxConn = 100
MaxIdle = 20
ConnMaxLifeTime = "8h"

[ip_location_conf]
ipdb_path = "/data/apps/custm-chat/conf/ipiptest.ipdb"

[redis_conf]
addr = "localhost:6379"
read_timeout = "2s"
write_timeout = "2s"

[centrifugo_conf]
api = "http://localhost:9000/api"
auth = "apikey bef6a4f4-8d50-42fa-9a3d-4e20d32db45e"
timeout = "2s"
retry_times = 3
online_agents = "%s_online_agents"
visitor_queue = "%s_visitor_queue"
conversation_instruction = "%s_conv_instruction"
internal_msg = "%s_internal_messages"

[elasticsearch_conf]
endpoint = "http://localhost:9200"
username = ""
password = ""
region = ""
timeout = "500ms"

[email_conf]
api_user = "czy01_test_5tagRq"
api_key = ""
endpoint = "http://api.sendcloud.net/apiv2/mail/send"
from = "GQO3GPtwHREPtt5UsAwLTgMRoFmKV5lg.sendcloud.org"
from_name = "chat186"
subject = ""
html_tmpl = ""

[aws_s3_conf]
region = "ap-southeast-1"
bucket_name = "chat-im"
access_key_id = ""
secret_access_key = ""
session_token = ""
use_ssl = false


[task_handler_conf]
host = "http://127.0.0.1:8989"
send_ent_message_endpoint = "/send_auto_message"
send_no_resp_message_endpoint = "/send_no_resp_message"
end_conversation_endpoint = "/end_conversation"
delete_job_endpoint = "/delete_task"

```

```bash
make build
docker build -t custm-chat:v1 .
docker run -d -p 8090:8090 custm-chat:v1
```