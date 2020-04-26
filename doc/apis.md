-   [API](#api)
    -   [注册企业 Sign Up](#注册企业-sign-up)
    -   [Sign In](#sign-in)
    -   [激活](#激活)
    -   [Get visitor connection token](#get-visitor-connection-token)
    -   [Upload](#upload)
    -   [Init Conversation](#init-conversation)
-   [访客端](#访客端)
    -   [坐席分配](#坐席分配)
    -   [检查访客是否是黑名单访客](#检查访客是否是黑名单访客)
    -   [初始化访客记录](#初始化访客记录)
    -   [Init Visit](#init-visit)
    -   [访客发送消息](#访客发送消息)
    -   [获取历史聊天会话](#获取历史聊天会话)
    -   [坐席评价](#坐席评价)
    -   [创建留言](#创建留言)
    -   [获取询前表单](#获取询前表单)
    -   [获取企业配置](#获取企业配置)
-   [坐席端](#坐席端)
    -   [登出](#登出)
    -   [登出全部设备](#登出全部设备)
    -   [重置密码](#重置密码)
    -   [更新邮件地址(登录账号)](#更新邮件地址登录账号)
    -   [调整坐席顺序](#调整坐席顺序)
    -   [检查登录坐席是否超过限制](#检查登录坐席是否超过限制)
    -   [创建或更新坐席分配规则](#创建或更新坐席分配规则)
    -   [活跃的会话](#活跃的会话)
    -   [同事的对话](#同事的对话)
    -   [获取企业信息](#获取企业信息)
    -   [更新enterprise](#更新enterprise)
    -   [获取agent列表](#获取agent列表)
    -   [获取 agent group 列表](#获取-agent-group-列表)
    -   [创建agent group](#创建agent-group)
    -   [更新配置](#更新配置)
    -   [获取安全设置](#获取安全设置)
    -   [创建安全设置](#创建安全设置)
    -   [更新安全设置](#更新安全设置)
    -   [获取访问统计信息](#获取访问统计信息)
    -   [获取对话统计信息](#获取对话统计信息)
    -   [add agent](#add-agent)
    -   [export agent list](#export-agent-list)
    -   [get current agent info](#get-current-agent-info)
    -   [get agent detail](#get-agent-detail)
    -   [update agent info](#update-agent-info)
    -   [delete agent](#delete-agent)
    -   [获取坐席连接(centrifugo)token](#获取坐席连接centrifugotoken)
    -   [创建或更新消息提醒](#创建或更新消息提醒)
    -   [获取消息提醒配置](#获取消息提醒配置)
    -   [get agent perms](#get-agent-perms)
    -   [更改坐席状态](#更改坐席状态)
    -   [在线访客](#在线访客)
    -   [更新访客信息](#更新访客信息)
    -   [给访客打tag](#给访客打tag)
    -   [搜索访客](#搜索访客)
    -   [创建快捷回复组](#创建快捷回复组)
    -   [创建快捷回复项](#创建快捷回复项)
    -   [更新quick reply group item](#更新quick-reply-group-item)
    -   [删除quick reply item](#删除quick-reply-item)
    -   [获取自动消息设置](#获取自动消息设置)
    -   [创建自动消息](#创建自动消息)
    -   [更新自动消息](#更新自动消息)
    -   [删除自动消息](#删除自动消息)
    -   [创建prechat forms](#创建prechat-forms)
    -   [Create visitor tag](#create-visitor-tag)
    -   [Get visitor tags](#get-visitor-tags)
    -   [update visitor tag](#update-visitor-tag)
    -   [delete visitor tag](#delete-visitor-tag)
    -   [add visitor to black list](#add-visitor-to-black-list)
    -   [结束会话](#结束会话)
    -   [对话转接](#对话转接)
    -   [添加小结](#添加小结)
    -   [历史对话搜索](#历史对话搜索)
    -   [坐席发送消息](#坐席发送消息)
    -   [获取会话消息](#获取会话消息)
    -   [撤回会话消息](#撤回会话消息)
    -   [增加权限](#增加权限)
    -   [增加或更新角色的权限](#增加或更新角色的权限)
    -   [获取角色的权限](#获取角色的权限)
    -   [获取租户的所有权限](#获取租户的所有权限)
    -   [获取留言列表](#获取留言列表)
    -   [更新留言状态](#更新留言状态)

API
===

注册企业 Sign Up
----------------

POST /signup

  参数       类型     描述
  ---------- -------- ------
  name       string   
  email      string   
  password   string   
  mobile     string   

Response

``` {.json}
{
    "code": 0,
    "body": {
        "ent_id": "bgonnvul0s1crs2ugt20"
    }
}
```

curl

``` {.bash}
curl -X POST \
  http://localhost:8089/signup \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "test enterprise",
    "email": "zhiruchenupc@gmail.com",
    "password": "abc123$",
    "mobile": "13356789090"
}'
```

Sign In
-------

POST /signin

  参数       类型     描述
  ---------- -------- ------
  email      string   
  password   string   

Response

``` {.json}
{
    "ent_id": "bgrg80l5jj83bqe154fg",
    "user_id": "bgrg80l5jj83bqe154h0",
    "user_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ3JnODBsNWpqODNicWUxNTRmZyIsInVzZXJfaWQiOiJiZ3JnODBsNWpqODNicWUxNTRoMCIsImV4cCI6MTU0ODA0MjA3MSwiaXNzIjoiQ3VzdG1JTSJ9.n4bDM3KE_MpZ2XyMFZl4we3etSI9tCfruxNqCNfWu00"
}
```

curl

``` {.bash}
curl -X POST \
  http://localhost:8089/signin \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "zhiruchenupc@gmail.com",
    "password": "abc123$"
}'

< Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QyMCIsInVzZXJfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QzZyIsImV4cCI6MTU0NzYxODQwNSwiaXNzIjoiQ3VzdG1JTSJ9.hIYTbF6X4U2G-EfqOliqpXJAu1bglywxHWfSgnCCpmc
< Content-Type: application/json; charset=UTF-8
< Date: Sun, 06 Jan 2019 06:00:05 GMT
< Content-Length: 118
{"ent_id":"bgonnvul0s1crs2ugt20","user_id":"bgonnvul0s1crs2ugt3g","x_user_token":"bgonnvul0s1crs2ugt3g-bgopipel0s1ehd1bepb0"}
```

激活
----

GET /api/v1/activate

  参数             类型     描述
  ---------------- -------- --------
  activate\_code   string   激活码

Response

``` {.json}
{
  "code": 0
}
```

Get visitor connection token
----------------------------

GET /api/v1/connection\_token?visit\_id=xxxxxx

Request:

  参数        类型     描述
  ----------- -------- ------
  visit\_id   string   

Response:

``` {.json}
{
    "code": 0,
    "body": {
        "token": "xxxxxxx"
    }
}
```

Upload
------

POST /api/v1/upload?ent\_id=xxxxx

Request: form file

  参数   类型   描述
  ------ ------ ------
  file   File   

Response:

``` {.json}
{
    "code": 0,
    "body": {
        "file_url": "xxxxxxx"
    }
}
```

Init Conversation
-----------------

POST /api/v1/enterprises/:ent\_id/conversations/

  参数        类型     描述
  ----------- -------- -----------
  trace\_id   string   
  agent\_id   string   
  title       string   会话title

Response:

``` {.json}
{
  "code": 0,
  "body": {
    "id": "bg0f7n2dn47nmvtpesd0",
    "ent_id": "1234",
    "trace_id": "xxxxxx",
    "agent_id": "yyyyyy",
    "agent_msg_count": 0,
    "msg_count": 0,
    "title": "this is a conversation",
    "client_msg_count": 0,
    "duration": 0,
    "summary": "",
    "created_at": "2018-11-30T08:25:00.751035Z",
    "update_at": "2018-11-30T08:25:00.751035Z",
    "agent_type": "human",
    "client_first_send_time": "",
    "first_msg_created_at": "",
    "first_response_wait_time": 0,
    "last_msg_content": "",
    "last_msg_created_at": "",
    "quality_grade": "",
    "ended_at": "",
    "ended_by": ""
  }
}
```

访客端
======

坐席分配
--------

GET /api/v1/allocate\_agent?trace\_id=xxxx&ent\_id=xxxxxx

Response:

  参数        类型     描述
  ----------- -------- ------------------
  code        string   status code
  agent\_id   string   分配的坐席Id
  rule        string   所使用的分配规则

``` {.json}
{
    "code": 0,
    "agent_id": "",
    "rule": ""
}
```

检查访客是否是黑名单访客
------------------------

GET /enterprises/:ent\_id/visitor\_allowed?trace\_id=xxxxx

Request:

  参数        类型     描述
  ----------- -------- ------
  trace\_id   string   

Response:

``` {.json}
{
    "code": 0,
    "allow": true
}
```

初始化访客记录
--------------

Init Visit
----------

POST /api/v1/enterprises/:ent\_id/visits

Request:

  参数            类型     描述
  --------------- -------- --------
  ent\_id         string   企业id
  trace\_id       string   
  title           string   
  url             string   
  referrer\_url   string   

Response:

``` {.json}
{
    "visit_id": "xxxx",
    "visitor_id": "xxxx",
    "trace_id": "xxxx"
}
```

访客发送消息
------------

POST
/api/v1/enterprises/:ent\_id/conversations/:conversation\_id/messages

请求参数:

  参数            类型     描述
  --------------- -------- ------
  trace\_id       string   
  agent\_id       string   
  content         string   
  content\_type   string   

响应:

``` {.json}
{
    "code": 0,
    "body": "msg"
}
```

获取历史聊天会话
----------------

GET /api/v1/enterprises/:ent\_id/conversations/history

请求参数:

  参数        类型     描述
  ----------- -------- ------
  trace\_id   string   
  offset      string   
  limit       string   

响应:

``` {.json}
```

坐席评价
--------

POST
/api/v1/enterprises/:ent\_id/conversations/:conversation\_id/evaluations

  参数        类型             描述
  ----------- ---------------- -------
  agent\_id   string           
  level       int(1/2/3/4/5)   5星制
  content     string           

响应:

``` {.json}
{
    "code": 0,
    "body": {"level": 1}
}
```

创建留言
--------

POST /api/v1/enterprises/:ent\_id/leave\_messages

请求参数:

  参数            类型     描述
  --------------- -------- ------
  visitor\_name   string   
  telephone       string   
  email           string   
  wechat          string   
  qq              string   
  content         string   

响应:

``` {.json}
{
    "code": 0,
    "body": {
        "id": "xxxxxxx",
        "ent_id": "xxxxxxx",
        "content": "leave message test",
        "created_at": "",
        "updated_at": ""
    }
}
```

获取询前表单
------------

GET /api/v1/enterprises/:ent\_id/chat\_forms

响应:

``` {.json}
{
    "code": 0,
    "body": {

    }
}
```

获取企业配置
------------

GET /admin/api/v1/enterprise/configs

响应:

``` {.json}
{
    "code": 0,
    "body": {
        "conversation_configs": {
            "ending_conversation_conf": {
                "ent_id": "bgpeo4l5jj879im1vipg",
                "no_message_duration": 10,
                "offline_duration": 30,
                "status": true
            },
            "ending_message_conf": {
                "ent_id": "bgpeo4l5jj879im1vipg",
                "platform": "web",
                "agent": "坐席手动结束消息",
                "system": "系统结束消息",
                "status": true,
                "prompt": true
            },
            "visitor_queue_conf": {
                "ent_id": "bgpeo4l5jj879im1vipg",
                "queue_size": 10,
                "description": "desc queue",
                "status": true
            },
            "conversation_transfer_conf": {
                "ent_id": "bgpeo4l5jj879im1vipg",
                "duration": 30,
                "transfer_target": "bgpf6cd5jj88h5usins0",
                "target_type": "agent",
                "status": true
            },
            "conversation_quality_conf": {
                "ent_id": "bgpeo4l5jj879im1vipg",
                "grade": "🏅",
                "visitor_msg_count": 100,
                "agent_msg_count": 200,
                "status": true
            }
        },
        "security_config": {
            "login_limit": {
                "ent_id": "bgpeo4l5jj879im1vipg",
                "status": false,
                "role_ids": [],
                "city_list": [],
                "allowed_ip_list": []
            },
            "send_file": {
                "ent_id": "bgpeo4l5jj879im1vipg",
                "status": false
            }
        },
        "leave_message_config": {
            "ent_id": "",
            "introduction": "",
            "show_visitor_name": false,
            "show_telephone": false,
            "show_email": false,
            "show_wechat": false,
            "show_qq": false,
            "auto_create_category": false,
            "fill_contact": "",
            "use_default_content": false,
            "default_content": ""
        }
    }
}
```

curl

``` {.bash}
curl -X GET \
  http://localhost:8089/admin/api/v1/enterprise/configs \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ3BlbzRsNWpqODc5aW0xdmlwZyIsInVzZXJfaWQiOiJiZ3BmNmNkNWpqODhoNXVzaW5zMCIsImV4cCI6MTU0NzgwMjUzNSwiaXNzIjoiQ3VzdG1JTSJ9.n2n6yhpk-Y7VU4oatMp-eU_HgelWsAyE71QFFriVcr4' \
  -H 'cache-control: no-cache'
```

坐席端
======

登出
------

POST /admin/api/v1/signout

响应:

``` {.json}
{
    "code": 0
}
```

curl

``` {.bash}
curl -X POST \
  http://localhost:8089/admin/api/v1/signout \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QyMCIsInVzZXJfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QzZyIsImV4cCI6MTU0NzYxODQwNSwiaXNzIjoiQ3VzdG1JTSJ9.hIYTbF6X4U2G-EfqOliqpXJAu1bglywxHWfSgnCCpmc' \
  -H 'Cache-Control: no-cache'
```

登出全部设备
------------

POST /admin/api/v1/signout\_all

响应:

``` {.json}
{
    "code": 0
}
```

重置密码
--------

POST /admin/api/v1/reset\_password

  参数            类型     描述
  --------------- -------- ------
  old\_password   string   
  new\_password   string   

响应:

``` {.json}
{
    "code": 0
}
```

更新邮件地址(登录账号)
----------------------

PUT /admin/api/v1/agents/email

  参数    类型     描述
  ------- -------- ------
  email   string   

响应:

``` {.json}
{
    "code": 0
}
```

调整坐席顺序
------------

PUT /admin/api/v1/agents/rankings

Request

``` {.json}
{
  "rankings": [
      {
        "agent_id": "xxxxxxxxx",
        "ranking": 0
      },
      {
        "agent_id": "yyyyyy",
        "ranking": 2
      },
      {
        "agent_id": "zzzzzz",
        "ranking": 3
      }
  ]
}
```

Response

``` {.json}
{
  "code": 0
}
```

检查登录坐席是否超过限制
------------------------

GET /admin/api/v1/check\_online\_count

响应:

``` {.json}
{
    "code": 0
}
```

curl

``` {.bash}
curl -X GET \
  http://localhost:8089/admin/api/v1/check_online_count \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QyMCIsInVzZXJfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QzZyIsImV4cCI6MTU0NzYxMjg3MiwiaXNzIjoiQ3VzdG1JTSJ9.wxl9HmQFd01IGcgTS-QmFFpUPMSXKiKmmn64GsvluuQ' \
  -H 'Cache-Control: no-cache'
```

创建或更新坐席分配规则
----------------------

POST /admin/api/v1/allocation_rules

参数:

  参数         类型     描述
  ------------ -------- --------------
  rule_type   string   分配规则类型

rule_type:

1.  order_take_turns 按客服顺序轮流
2.  order_priority 按客服顺序优先
3.  conversation_num 按会话数

Resp

``` {.json}
{
    "code": 0,
    "body": {
        "id": "bgou0cel0s1ehd1bepeg",
        "ent_id": "bgonnvul0s1crs2ugt20",
        "rule_type": "order_take_turns"
    }
}
```

curl

``` {.bash}
curl -X POST \
  http://localhost:8089/admin/api/v1/allocation_rules \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QyMCIsInVzZXJfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QzZyIsImV4cCI6MTU0NzYzMzQ1MSwiaXNzIjoiQ3VzdG1JTSJ9.n3PdgyX1uCW9aOxX6YRs1W5kdzPybiXCNTDmUdKFTjU' \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
    "rule_type": "order_take_turns"
}'
```

活跃的会话
----------

GET /admin/api/v1/active\_conversations

Response:

``` {.json}
[
    {
        "id": "xxxxx"
    }
]
```

curl

``` {.bash}
curl -X GET \
  http://localhost:8089/admin/api/v1/active_conversations \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QyMCIsInVzZXJfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QzZyIsImV4cCI6MTU0NzYzMzQ1MSwiaXNzIjoiQ3VzdG1JTSJ9.n3PdgyX1uCW9aOxX6YRs1W5kdzPybiXCNTDmUdKFTjU' \
  -H 'Cache-Control: no-cache'
```

同事的对话
----------

GET /admin/api/v1/colleague\_conversations?offset=0&limit=10

Response:

``` {.json}
[
    {
        "id": "xxxxx"
    }
]
```

curl

``` {.bash}
curl -X GET \
  'http://localhost:8089/admin/api/v1/colleague_conversations?offset=0&limit=10' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QyMCIsInVzZXJfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QzZyIsImV4cCI6MTU0NzYzMzQ1MSwiaXNzIjoiQ3VzdG1JTSJ9.n3PdgyX1uCW9aOxX6YRs1W5kdzPybiXCNTDmUdKFTjU' \
  -H 'Cache-Control: no-cache'
```

获取企业信息
------------

GET /admin/api/v1/enterprise/info

Response:

``` {.json}
{
    "code": 0,
    "body": {
        "ent_info": {
            "ent_detail": {
                "id": "bgmnicujo0p67acb7kag",
                "name": "test ent2",
                "full_name": "",
                "province": "",
                "city": "",
                "avatar": "",
                "industry": "",
                "location": "",
                "address": "",
                "website": "",
                "email": "2550418657@qq.com",
                "mobile": "18389890023",
                "description": "",
                "created_at": "2019-01-03T02:53:39.240856Z",
                "owner": "2550418657@qq.com",
                "plan": 1,
                "agent_num": 1,
                "trial_status": 2,
                "expiration_time": "2019-02-02T02:53:39.240856Z",
                "last_activated_at": "2019-01-03T02:53:39.240856Z",
                "contact_mobile": "",
                "contact_email": "",
                "contact_qq": "",
                "contact_wechat": "",
                "contact_signature": ""
            },
            "allocation_rule": null
        },
        "super_admin": {
            "id": "bgmnicujo0p67acb7kc0",
            "ent_id": "bgmnicujo0p67acb7kag",
            "group_id": "bgmnicujo0p67acb7kb0",
            "role_id": "bgmnicujo0p67acb7kbg",
            "avatar": "",
            "username": "2550418657@qq.com",
            "real_name": "",
            "nick_name": "",
            "job_number": "",
            "serve_limit": 0,
            "is_online": false,
            "ranking": 0,
            "email": "2550418657@qq.com",
            "mobile": "18389890023",
            "qq_num": "",
            "signature": "",
            "status": "offline",
            "wechat": "",
            "is_admin": 1,
            "perms_range_type": "all",
            "account_status": "valid",
            "create_at": "2019-01-03T02:53:39.241537Z",
            "update_at": "2019-01-03T02:53:57.481448Z"
        }
    }
}
```

curl

``` {.bash}
curl -X GET \
  http://localhost:8089/admin/api/v1/enterprise/info \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QyMCIsInVzZXJfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QzZyIsImV4cCI6MTU0NzYxMjg3MiwiaXNzIjoiQ3VzdG1JTSJ9.wxl9HmQFd01IGcgTS-QmFFpUPMSXKiKmmn64GsvluuQ' \
  -H 'Cache-Control: no-cache'
```

更新enterprise
--------------

PUT /admin/api/v1/enterprises/info

参数:

  参数          类型     描述
  ------------- -------- ------
  name          string   
  avatar        string   
  industry      string   
  mobile        string   
  description   string   

Response:

``` {.json}
{
    "code": 0,
    "body": {
        "id": "bgonnvul0s1crs2ugt20",
        "name": "test ent 1234",
        "avatar": "",
        "industry": "IT",
        "email": "zhiruchenupc@gmail.com",
        "mobile": "12345678900",
        "description": "ent desc",
        "created_at": "2019-01-06T11:54:39.976567+08:00",
        "owner": "zhiruchenupc@gmail.com",
        "plan": 1,
        "agent_num": 1,
        "trial_status": 2,
        "expiration_time": "2019-02-05T11:54:39.976567+08:00",
        "last_activated_at": "2019-01-06T11:54:39.976567+08:00"
    }
}
```

curl

``` {.bash}
curl -X PUT \
  http://localhost:8089/admin/api/v1/enterprise/info \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QyMCIsInVzZXJfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QzZyIsImV4cCI6MTU0NzYzMzQ1MSwiaXNzIjoiQ3VzdG1JTSJ9.n3PdgyX1uCW9aOxX6YRs1W5kdzPybiXCNTDmUdKFTjU' \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "test ent 1234",
    "industry": "IT",
    "mobile": "12345678900",
    "description": "ent desc"
}'
```

获取agent列表
-------------

GET /admin/api/v1/enterprise/agents

Response:

``` {.json}
{
    "code": 0,
    "body": [
        {
            "id": "bgonnvul0s1crs2ugt3g",
            "ent_id": "bgonnvul0s1crs2ugt20",
            "group_id": "bgonnvul0s1crs2ugt2g",
            "role_id": "bgonnvul0s1crs2ugt30",
            "avatar": "",     
            "username": "zhiruchenupc@gmail.com",
            "real_name": "",
            "job_number": "",
            "serve_limit": 1,
            "ranking": 0,
            "email": "zhiruchenupc@gmail.com",
            "mobile": "13356789090",
            "qq_num": "",
            "signature": "",
            "wechat": "",
            "is_admin": 1,
            "perms_range_type": "all",
            "account_status": "valid",
            "create_at": "2019-01-06T11:54:39.978652+08:00",
            "update_at": "2019-01-06T11:56:54.045535+08:00",
            "deleted_at": {
                "Time": "0001-01-01T00:00:00Z",
                "Valid": false
            },
            "status": "offline"
        }
    ]
}
```

curl

``` {.bash}
curl -X GET \
  http://localhost:8089/admin/api/v1/enterprise/agents \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QyMCIsInVzZXJfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QzZyIsImV4cCI6MTU0NzYzMzQ1MSwiaXNzIjoiQ3VzdG1JTSJ9.n3PdgyX1uCW9aOxX6YRs1W5kdzPybiXCNTDmUdKFTjU' \
  -H 'Cache-Control: no-cache'
```

获取 agent group 列表
---------------------

GET /admin/api/v1/enterprises/agent\_groups

Response:

``` {.json}
{
    "code": 0,
    "body": [
        {
            "id": "bgonnvul0s1crs2ugt2g",
            "ent_id": "bgonnvul0s1crs2ugt20",
            "name": "超级管理员",
            "description": "超管"
        }
    ]
}
```

curl

``` {.bash}
curl -X GET \
  http://localhost:8089/admin/api/v1/enterprise/agent_groups \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QyMCIsInVzZXJfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QzZyIsImV4cCI6MTU0NzYzMzQ1MSwiaXNzIjoiQ3VzdG1JTSJ9.n3PdgyX1uCW9aOxX6YRs1W5kdzPybiXCNTDmUdKFTjU' \
  -H 'Cache-Control: no-cache'
```

创建agent group
---------------

POST /admin/api/v1/enterprise/agent\_groups

参数

  参数          类型     描述
  ------------- -------- ------
  name          string   
  description   string   

Response:

``` {.json}
{
    "code": 0,
    "body": {
        "id": "bgp1blul0s1fbu6c0uvg",
        "ent_id": "bgonnvul0s1crs2ugt20",
        "name": "agent group1",
        "description": "agent group desc"
    }
}
```

curl

``` {.bash}
curl -X POST \
  http://localhost:8089/admin/api/v1/enterprise/agent_groups \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QyMCIsInVzZXJfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QzZyIsImV4cCI6MTU0NzYzMzQ1MSwiaXNzIjoiQ3VzdG1JTSJ9.n3PdgyX1uCW9aOxX6YRs1W5kdzPybiXCNTDmUdKFTjU' \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "agent group1",
    "description": "agent group desc"
}'
```

更新配置
--------

PUT /admin/api/v1/enterprise/conv\_configs

Response:

``` {.json}
{
    "code": 0,
    "body": {
        "conversation_configs": {
            "ending_conversation_conf": {
                "ent_id": "bgpeo4l5jj879im1vipg",
                "no_message_duration": 10,
                "offline_duration": 30,
                "status": true
            },
            "ending_message_conf": {
                "ent_id": "bgpeo4l5jj879im1vipg",
                "platform": "web",
                "agent": "坐席手动结束消息",
                "system": "系统结束消息",
                "status": true,
                "prompt": true
            },
            "visitor_queue_conf": {
                "ent_id": "bgpeo4l5jj879im1vipg",
                "queue_size": 10,
                "description": "desc queue",
                "status": true
            },
            "conversation_transfer_conf": {
                "ent_id": "bgpeo4l5jj879im1vipg",
                "duration": 30,
                "transfer_target": "bgpf6cd5jj88h5usins0",
                "target_type": "agent",
                "status": true
            },
            "conversation_quality_conf": {
                "ent_id": "bgpeo4l5jj879im1vipg",
                "grade": "🏅",
                "visitor_msg_count": 100,
                "agent_msg_count": 200,
                "status": true
            }
        }
    }
}
```

curl

``` {.bash}
curl -X PUT \
  http://localhost:8089/admin/api/v1/enterprise/conv_configs \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ3BlbzRsNWpqODc5aW0xdmlwZyIsInVzZXJfaWQiOiJiZ3BmNmNkNWpqODhoNXVzaW5zMCIsImV4cCI6MTU0NzgwMjUzNSwiaXNzIjoiQ3VzdG1JTSJ9.n2n6yhpk-Y7VU4oatMp-eU_HgelWsAyE71QFFriVcr4' \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
    "conv_configs": {
        "ending_conversation_conf": {
            "no_message_duration": 10,
            "offline_duration": 30,
            "status": true
        },
        "ending_message_conf": {
            "platform": "web",
            "agent": "坐席手动结束消息",
            "system": "系统结束消息",
            "status": true,
            "prompt": true
        },
        "visitor_queue_conf": {
            "queue_size": 10,
            "description": "desc queue",
            "status": true
        },
        "conversation_transfer_conf": {
            "duration": 30,
            "transfer_target": "bgpf6cd5jj88h5usins0",
            "target_type": "agent",
            "status": true
        },
        "conversation_quality_conf": {
            "grade": "🏅",
            "visitor_msg_count": 100,
            "agent_msg_count": 200,
            "status": true
        }
    }
}'
```

获取安全设置
------------

GET /admin/api/v1/enterprise/security\_configs

Response:

``` {.json}
{
    "code": 0,
    "body": {
        "login_limit": {
            "ent_id": "bgonnvul0s1crs2ugt20",
            "status": false,
            "role_ids": [],
            "city_list": [],
            "allowed_ip_list": []
        },
        "send_file": {
            "ent_id": "bgonnvul0s1crs2ugt20",
            "status": false
        }
    }
}
```

curl

``` {.bash}
curl -X GET \
  http://localhost:8089/admin/api/v1/enterprise/security_configs \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QyMCIsInVzZXJfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QzZyIsImV4cCI6MTU0NzYzMzQ1MSwiaXNzIjoiQ3VzdG1JTSJ9.n3PdgyX1uCW9aOxX6YRs1W5kdzPybiXCNTDmUdKFTjU' \
  -H 'Cache-Control: no-cache'
```

创建安全设置
------------

POST /admin/api/v1/enterprise/security\_configs

参数

  参数           类型     描述
  -------------- -------- ------
  login\_limit   object   
  send\_file     object   

login\_limit

  参数                类型         描述
  ------------------- ------------ ------
  ent\_id             string       
  status              bool         
  role\_ids           \[\]string   
  city\_list          \[\]string   
  allowed\_ip\_list   \[\]string   

send\_file

  参数      类型     描述
  --------- -------- ------
  ent\_id   string   
  status    bool     

Response:

``` {.json}
{
    "code": 0,
    "body": {
        "login_limit": {
            "ent_id": "bgonnvul0s1crs2ugt20",
            "status": true,
            "role_ids": null,
            "city_list": [
                "北京"
            ],
            "allowed_ip_list": null
        },
        "send_file": {
            "ent_id": "bgonnvul0s1crs2ugt20",
            "status": true
        }
    }
}
```

curl

``` {.bash}
curl -X POST \
  http://localhost:8089/admin/api/v1/enterprise/security_configs \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QyMCIsInVzZXJfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QzZyIsImV4cCI6MTU0NzYzMzQ1MSwiaXNzIjoiQ3VzdG1JTSJ9.n3PdgyX1uCW9aOxX6YRs1W5kdzPybiXCNTDmUdKFTjU' \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
    "login_limit": {
        "status": true,
        "city_list": ["北京"]
    },
    
    "send_file": {
        "status": true
    }
}'
```

更新安全设置
------------

PUT /admin/api/v1/enterprise/security\_configs

参数

  参数          类型     描述
  ------------- -------- ------
  name          string   
  description   string   

Response:

``` {.json}
{
    "code": 0,
    "body": {
        "login_limit": {
            "ent_id": "bgonnvul0s1crs2ugt20",
            "status": true,
            "role_ids": null,
            "city_list": [
                "北京",
                "兰州"
            ],
            "allowed_ip_list": null
        },
        "send_file": {
            "ent_id": "bgonnvul0s1crs2ugt20",
            "status": true
        }
    }
}
```

curl

``` {.bash}
curl -X PUT \
  http://localhost:8089/admin/api/v1/enterprise/security_configs \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QyMCIsInVzZXJfaWQiOiJiZ29ubnZ1bDBzMWNyczJ1Z3QzZyIsImV4cCI6MTU0NzYzMzQ1MSwiaXNzIjoiQ3VzdG1JTSJ9.n3PdgyX1uCW9aOxX6YRs1W5kdzPybiXCNTDmUdKFTjU' \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
    "login_limit": {
        "status": true,
        "city_list": ["北京", "兰州"]
    },
    
    "send_file": {
        "status": true
    }
}'
```

获取访问统计信息
----------------

GET /admin/api/v1/enterprise/reports/visit

获取对话统计信息
----------------

GET /admin/api/v1/enterprise/reports/conversation

add agent
---------

POST /admin/api/v1/agents

参数

  参数             类型         描述
  ---------------- ------------ ------
  email            string       
  real\_name       string       
  init\_password   string       
  job\_num         string       
  perms\_range     string()     
  perms\_groups    \[\]string   
  serve\_limit     int          
  group\_id        string       
  role\_id         string       

perms\_range

1.  personal
2.  all
3.  part

Response:

``` {.json}
{
    "code": 0,
    "body": {
        "id": "bgpf6cd5jj88h5usins0",
        "ent_id": "bgpeo4l5jj879im1vipg",
        "group_id": "bgpetvl5jj879im1vji0",
        "role_id": "bgpeo4l5jj879im1viqg",
        "avatar": "",
        "username": "2550418657@qq.com",
        "real_name": "Bob",
        "job_number": "001",
        "serve_limit": 1,
        "ranking": 0,
        "email": "2550418657@qq.com",
        "mobile": "",
        "qq_num": "",
        "signature": "",
        "status": "offline",
        "wechat": "",
        "is_admin": 0,
        "perms_range_type": "personal",
        "account_status": "created",
        "create_at": "2019-01-07T06:35:29.374514Z",
        "update_at": "2019-01-07T06:35:29.374514Z",
        "deleted_at": ""
    }
}
```

curl

``` {.bash}
curl -X POST \
  http://localhost:8089/admin/api/v1/agents/ \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ3BlbzRsNWpqODc5aW0xdmlwZyIsInVzZXJfaWQiOiJiZ3BlbzRsNWpqODc5aW0xdmlyMCIsImV4cCI6MTU0NzcwNTQ4MiwiaXNzIjoiQ3VzdG1JTSJ9.HkK55YIQzCb1uzpoBZlBX3BGVT06OVMiwfUM_WtE0F0' \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
    "email": "2550418657@qq.com",
    "real_name": "Bob",
    "init_password": "abc123$",
    "job_num": "001",
    "perms_range" : "personal",
    "serve_limit": 1,
    "group_id": "bgpetvl5jj879im1vji0",
    "role_id": "bgpeo4l5jj879im1viqg"
}'
```

export agent list
-----------------

GET /admin/api/v1/agents/export

Response:

``` {.json}
{
    "code": 0,
    "body": {
        "location": "https://chat-im.s3.ap-southeast-1.amazonaws.com/bgpeo4l5jj879im1vipg/files/tmp/bgpeo4l5jj879im1vipg-1546843457-agent-list.xlsx"
    }
}
```

curl

``` {.bash}
curl -X GET \
  http://localhost:8089/admin/api/v1/agents/export \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ3BlbzRsNWpqODc5aW0xdmlwZyIsInVzZXJfaWQiOiJiZ3BlbzRsNWpqODc5aW0xdmlyMCIsImV4cCI6MTU0NzcwNTQ4MiwiaXNzIjoiQ3VzdG1JTSJ9.HkK55YIQzCb1uzpoBZlBX3BGVT06OVMiwfUM_WtE0F0' \
  -H 'cache-control: no-cache'
```

get current agent info
----------------------

GET /admin/api/v1/agents/info

Response:

``` {.json}
{
    "code": 0,
    "body": {
        "id": "bgrg80l5jj83bqe154h0",
        "ent_id": "bgrg80l5jj83bqe154fg",
        "group_id": "bgrg80l5jj83bqe154g0",
        "role_id": "bgrg80l5jj83bqe154gg",
        "avatar": "",
        "username": "2550418657@qq.com",
        "real_name": "",
        "nick_name": "",
        "job_number": "",
        "serve_limit": 1,
        "is_online": false,
        "ranking": 0,
        "email": "2550418657@qq.com",
        "mobile": "18868905690",
        "qq_num": "",
        "signature": "",
        "status": "unavailable",
        "wechat": "",
        "is_admin": 1,
        "perms_range_type": "all",
        "account_status": "valid",
        "create_at": "2019-01-10T16:36:18.684521+08:00",
        "update_at": "2019-01-11T11:42:42.593282+08:00",
        "deleted_at": ""
    }
}
```

get agent detail
----------------

GET /admin/api/v1/agents/:agent\_id

Response:

``` {.json}
{
    "code": 0,
    "body": {
        "id": "bgrg80l5jj83bqe154h0",
        "ent_id": "bgrg80l5jj83bqe154fg",
        "group_id": "bgrg80l5jj83bqe154g0",
        "role_id": "bgrg80l5jj83bqe154gg",
        "avatar": "",
        "username": "2550418657@qq.com",
        "real_name": "",
        "nick_name": "",
        "job_number": "",
        "serve_limit": 1,
        "is_online": false,
        "ranking": 0,
        "email": "2550418657@qq.com",
        "mobile": "18868905690",
        "qq_num": "",
        "signature": "",
        "status": "unavailable",
        "wechat": "",
        "is_admin": 1,
        "perms_range_type": "all",
        "account_status": "valid",
        "create_at": "2019-01-10T16:36:18.684521+08:00",
        "update_at": "2019-01-11T11:42:42.593282+08:00",
        "deleted_at": ""
    }
}
```

curl

``` {.bash}
curl -X GET \
  http://localhost:8089/admin/api/v1/agents/bgpf6cd5jj88h5usins0 \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ3BlbzRsNWpqODc5aW0xdmlwZyIsInVzZXJfaWQiOiJiZ3BlbzRsNWpqODc5aW0xdmlyMCIsImV4cCI6MTU0NzcwNTQ4MiwiaXNzIjoiQ3VzdG1JTSJ9.HkK55YIQzCb1uzpoBZlBX3BGVT06OVMiwfUM_WtE0F0' \
  -H 'cache-control: no-cache'
```

update agent info
-----------------

PUT /admin/api/v1/agents/:agent\_id

参数:

  参数             类型         描述
  ---------------- ------------ ------
  email            string       
  real\_name       string       
  init\_password   string       
  job\_num         string       
  perms\_range     string       
  perms\_groups    \[\]string   
  serve\_limit     int          
  group\_id        string       
  role\_id         string       

Response:

``` {.json}
{
    "code": 0,
    "body": {
        "id": "bgpgqad5jj8912h5nge0",
        "ent_id": "bgpgqad5jj8912h5ngcg",
        "group_id": "bgpetvl5jj879im1vji0",
        "role_id": "bgpeo4l5jj879im1viqg",
        "avatar": "",
        "username": "gloryz@dingtalk.com",
        "real_name": "BobDD",
        "job_number": "002",
        "serve_limit": 1,
        "ranking": 0,
        "email": "gloryz@dingtalk.com",
        "mobile": "18868905690",
        "qq_num": "",
        "signature": "",
        "status": "offline",
        "wechat": "",
        "is_admin": 1,
        "perms_range_type": "all",
        "account_status": "valid",
        "create_at": "2019-01-07T16:26:17.525318+08:00",
        "update_at": "2019-01-07T09:16:55.594575Z",
        "deleted_at": ""
    }
}
```

curl

``` {.bash}
curl -X PUT \
  http://localhost:8089/admin/api/v1/agents/bgpgqad5jj8912h5nge0 \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ3BncWFkNWpqODkxMmg1bmdjZyIsInVzZXJfaWQiOiJiZ3BncWFkNWpqODkxMmg1bmdlMCIsImV4cCI6MTU0NzcxMzc5MCwiaXNzIjoiQ3VzdG1JTSJ9.WsCJ_QvA0XTFJAMM6nAUij-Pmp37Dz2la901GC1Mblo' \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
    "email": "2550418657@qq.com",
    "real_name": "BobDD",
    "job_num": "002",
    "perms_range" : "personal",
    "serve_limit": 1,
    "group_id": "bgpetvl5jj879im1vji0",
    "role_id": "bgpeo4l5jj879im1viqg"
}'
```

delete agent
------------

DELETE /admin/api/v1/agents/:agent\_id

Response:

``` {.json}
{
    "code": 0,
    "body": null
}
```

获取坐席连接(centrifugo)token
-----------------------------

GET /admin/api/v1/agents/connection\_token

Response:

``` {.json}
{
    "code": 0,
    "body": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NDY4NTMwMDksInN1YiI6ImJncGdxYWQ1amo4OTEyaDVuZ2UwIn0.KsorGodV4FTQdlHxh2Y7FRkX6k5l0_nIlnM59Di8ewI"
    }
}
```

curl

``` {.bash}
curl -X GET \
  http://localhost:8089/admin/api/v1/agents/connection_token \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ3BncWFkNWpqODkxMmg1bmdjZyIsInVzZXJfaWQiOiJiZ3BncWFkNWpqODkxMmg1bmdlMCIsImV4cCI6MTU0NzcxMzc5MCwiaXNzIjoiQ3VzdG1JTSJ9.WsCJ_QvA0XTFJAMM6nAUij-Pmp37Dz2la901GC1Mblo' \
  -H 'Postman-Token: 64c8e9ab-eef8-4dc0-8e55-cec6f6e8b492' \
  -H 'cache-control: no-cache'
```

创建或更新消息提醒
------------------

POST /admin/api/v1/agents/message\_beep

参数:

  参数                          类型     描述
  ----------------------------- -------- ------
  client\_type                  string   
  beep\_type                    string   
  new\_conversation             bool     
  new\_message                  bool     
  conversation\_transfer\_in    bool     
  conversation\_transfer\_out   bool     
  colleague\_conversation       bool     

Response:

``` {.json}
{
    "code": 0,
    "body": {

    }
}
```

获取消息提醒配置
----------------

GET /admin/api/v1/agents/message\_beep

Response:

``` {.json}
{
    "code": 0,
    "body": {

    }
}
```

get agent perms
---------------

GET /admin/api/v1/agents/:agent\_id/perms

Response

``` {.json}
{
    "code": 0,
    "body": [
        {
            "id": "bgpgqad5jj8912h5ngig",
            "ent_id": "bgpgqad5jj8912h5ngcg",
            "app_name": "ent_info",
            "name": "check_update_ent_info",
            "created_at": "2019-01-07T16:26:17.536129+08:00",
            "updated_at": "2019-01-07T16:26:17.536129+08:00"
        },
        {
            "id": "bgpgqad5jj8912h5ngj0",
            "ent_id": "bgpgqad5jj8912h5ngcg",
            "app_name": "ent_info",
            "name": "check_ent_agent_group_info",
            "created_at": "2019-01-07T16:26:17.536129+08:00",
            "updated_at": "2019-01-07T16:26:17.536129+08:00"
        }
    ]
}
```

curl

``` {.bash}
curl -X GET \
  http://localhost:8089/admin/api/v1/agents/bgpgqad5jj8912h5nge0/perms \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ3BncWFkNWpqODkxMmg1bmdjZyIsInVzZXJfaWQiOiJiZ3BncWFkNWpqODkxMmg1bmdlMCIsImV4cCI6MTU0NzcxMzc5MCwiaXNzIjoiQ3VzdG1JTSJ9.WsCJ_QvA0XTFJAMM6nAUij-Pmp37Dz2la901GC1Mblo' \
  -H 'cache-control: no-cache'
```

更改坐席状态
------------

PUT /admin/api/v1/agents/status

  参数     类型     描述
  -------- -------- ------------------
  status   string   (online/offline)

Response:

``` {.json}
{
    "code": 0,
    "body": null
}
```

在线访客
--------

GET /admin/api/v1/visitors

Response:

``` {.json}
{
    "code": 0,
    "body": {}
}
```

更新访客信息
------------

PUT /admin/api/v1/visitors

  参数        类型     描述
  ----------- -------- ------
  trace\_id   string   
  name        string   
  mobile      string   
  wechat      string   
  email       string   
  qq\_num     string   
  remark      string   

Response:

``` {.json}
{
    "code": 0,
    "body": {}
}
```

给访客打tag
-----------

POST /admin/api/v1/visitors/tags

  参数        类型     描述
  ----------- -------- ------
  trace\_id   string   
  tag\_id     string   

Response:

``` {.json}
{
    "code": 0,
    "body": null
}
```

搜索访客
--------

POST /admin/api/v1/visitors/search

参数

  参数     类型         描述
  -------- ------------ ------
  rules    \[\]\*rule   
  offset   string       
  limit    string       

Response:

``` {.json}
{
    "code": 0,
    "body": {}
}
```

创建快捷回复组
--------------

POST /admin/api/v1/quickreplies/

Response

``` {.json}
{
    "code": 0,
    "body": {
        "id": "bgpfijl5jj88jf88hofg",
        "ent_id": "bgpeo4l5jj879im1vipg",
        "title": "quick reply 001",
        "created_by": "bgpeo4l5jj879im1vir0",
        "created_at": "2019-01-07T07:01:34.924006Z"
    }
}
```

curl

``` {.bash}
curl -X POST \
  http://localhost:8089/admin/api/v1/quickreplies/ \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ3BlbzRsNWpqODc5aW0xdmlwZyIsInVzZXJfaWQiOiJiZ3BlbzRsNWpqODc5aW0xdmlyMCIsImV4cCI6MTU0NzcwNTQ4MiwiaXNzIjoiQ3VzdG1JTSJ9.HkK55YIQzCb1uzpoBZlBX3BGVT06OVMiwfUM_WtE0F0' \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
    "title": "quick reply 001"
}'
```

创建快捷回复项
--------------

POST /admin/api/v1/quickreplies/:group\_id/items

Response

``` {.json}
{
    "code": 0,
    "body": {
        "id": "bgpg3t55jj88rceoouh0",
        "quickreply_group_id": "bgpfijl5jj88jf88hofg",
        "title": "quick reply title",
        "content": "quick reply content",
        "created_by": "bgpeo4l5jj879im1vir0",
        "created_at": "2019-01-07T07:38:28.810108Z"
    }
}
```

CURL

``` {.bash}
curl -X POST \
  http://localhost:8089/admin/api/v1/quickreplies/bgpfijl5jj88jf88hofg/items \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ3BlbzRsNWpqODc5aW0xdmlwZyIsInVzZXJfaWQiOiJiZ3BlbzRsNWpqODc5aW0xdmlyMCIsImV4cCI6MTU0NzcwNTQ4MiwiaXNzIjoiQ3VzdG1JTSJ9.HkK55YIQzCb1uzpoBZlBX3BGVT06OVMiwfUM_WtE0F0' \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
    "title": "quick reply title",
    "content": "quick reply content"
}'
```

更新quick reply group item
--------------------------

PUT /admin/api/v1/quickreplies/items/:item\_id

Request:

``` {.json}
{
  "title": "",
  "content": ""
}
```

Response:

``` {.json}
{
  "id": "",
  "quickreply_group_id": "",
  "title": "item title",
  "content": "item content",
  "created_by": "",
  "created_at": ""
}
```

删除quick reply item
--------------------

DELETE /admin/api/v1/quickreplies/items/:item\_id

Response:

``` {.json}
{
  "code": 0
}
```

获取自动消息设置
----------------

GET /admin/api/v1/automessages

Response:

``` {.go}
// AutomaticMessage represents a row from 'custmchat.automatic_message'.
type AutomaticMessage struct {
    ID           string    `json:"id"`            // id
    EntID        string    `json:"ent_id"`        // ent_id
    ChannelType  string    `json:"channel_type"`  // channel_type
    MsgType      string    `json:"msg_type"`      // msg_type
    MsgContent   string    `json:"msg_content"`   // msg_content
    AfterSeconds int       `json:"after_seconds"` // after_seconds
    Enabled      bool      `json:"enabled"`       // enabled
    CreatedAt    time.Time `json:"created_at"`    // created_at

    // xo fields
    _exists, _deleted bool
}
```

``` {.json}
[
  {
    "id": "",
    "ent_id": "",
    "channel_type": "",
    "msg_type": "",
    "msg_content": "",
    "after_seconds": 1,
    "enabled": false,
    "created_at": ""
  }
]
```

创建自动消息
------------

POST /admin/api/v1/automessages

Request:

``` {.go}
type AutomaticMessage struct {
    ChannelType  string `json:"channel_type"`  // channel_type
    MsgType      string `json:"msg_type"`      // msg_type
    MsgContent   string `json:"msg_content"`   // msg_content
    AfterSeconds int    `json:"after_seconds"` // after_seconds
    Enabled      bool   `json:"enabled"`       // enabled
}
```

``` {.json}
{
  "channel_type": "",
  "msg_type": "",
  "msg_content": "",
  "after_seconds": 10,
  "enabled": true
}
```

Response:

``` {.json}
{
    "id": "",
    "ent_id": "",
    "channel_type": "",
    "msg_type": "",
    "msg_content": "",
    "after_seconds": 1,
    "enabled": true,
    "created_at": ""
}
```

更新自动消息
------------

PUT /admin/api/v1/automessages/:msg\_id

Request:

``` {.json}
{
  "channel_type": "",
  "msg_type": "",
  "msg_content": "",
  "after_seconds": 10,
  "enabled": true
}
```

Response:

``` {.json}
{
    "id": "",
    "ent_id": "",
    "channel_type": "",
    "msg_type": "",
    "msg_content": "",
    "after_seconds": 1,
    "enabled": true,
    "created_at": ""
}
```

删除自动消息
------------

DELETE /admin/api/v1/automessages/:msg\_id

Response:

``` {.json}
{
  "code": 0
}
```

创建prechat forms
-----------------

POST /admin/api/v1/prechat\_forms

Request:

``` {.json}
{
  "title": "test form",
  "inputs": {
    "status": "open",
    "description": "input form",
    "fields": [
      {
        "display_name": "telephone",
        "field_name": "telephone",
        "value_type": "text",
        "required": true
      }
    ]
  },
  "menus": {
    "status": "open",
    "description": "menu desc",
    "fields": [
      {
        "description": "support",
        "value": "support",
        "agent_type": "agent"
      }
    ]
  }
}
```

Response:

``` {.json}
{
  "code": 0,
  "body": {
    "id": "bgrc0vt5jj82a4bhs95g",
    "title": "test form",
    "inputs": {
      "status": "open",
      "description": "input form",
      "fields": [
        {
          "display_name": "telephone",
          "field_name": "telephone",
          "value_type": "text",
          "required": true
        }
      ]
    },
    "menus": {
      "status": "open",
      "description": "menu desc",
      "fields": [
        {
          "description": "support",
          "value": "support",
          "agent_type": "agent"
        }
      ]
    },
    "created_at": "2019-01-10T03:48:15.160604Z",
    "updated_at": "2019-01-10T03:48:15.160604Z"
  }
}
```

Create visitor tag
------------------

POST /admin/api/v1/visitor\_tags

Request:

``` {.json}
{
    "name": "test tag",
    "color": "red"
}
```

Response:

``` {.json}
{
    "code": 0,
    "body": {
        "id": "bgrga3l5jj83bqe1557g",
        "ent_id": "bgrg80l5jj83bqe154fg",
        "creator": "bgrg80l5jj83bqe154h0",
        "name": "test tag",
        "color": "red",
        "use_count": 0,
        "created_at": "2019-01-10T08:40:46.983393Z",
        "updated_at": "2019-01-10T08:40:46.983393Z"
    }
}
```

curl

``` {.bash}
curl -X POST \
  http://localhost:8089/admin/api/v1/visitor_tags/ \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ3JnODBsNWpqODNicWUxNTRmZyIsInVzZXJfaWQiOiJiZ3JnODBsNWpqODNicWUxNTRoMCIsImV4cCI6MTU0Nzk3MzQ5OCwiaXNzIjoiQ3VzdG1JTSJ9.dbVwIGJ9Lge3bon7hopyQFQtmnw6DivM314VpItGlow' \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: ca10dff4-ac6e-483a-9acd-cefff053240d' \
  -H 'cache-control: no-cache' \
  -d '{
    "name": "test tag",
    "color": "red"
}'
```

Get visitor tags
----------------

``` {.bash}
curl -X GET \
  http://localhost:8089/admin/api/v1/visitor_tags/ \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ3JnODBsNWpqODNicWUxNTRmZyIsInVzZXJfaWQiOiJiZ3JnODBsNWpqODNicWUxNTRoMCIsImV4cCI6MTU0Nzk3MzQ5OCwiaXNzIjoiQ3VzdG1JTSJ9.dbVwIGJ9Lge3bon7hopyQFQtmnw6DivM314VpItGlow' \
  -H 'Postman-Token: 840bf5b2-123a-441d-a7b7-ce423b92be86' \
  -H 'cache-control: no-cache'
```

Response:

``` {.json}
{
    "code": 0,
    "body": [
        {
            "id": "bgrga3l5jj83bqe1557g",
            "ent_id": "bgrg80l5jj83bqe154fg",
            "creator": "bgrg80l5jj83bqe154h0",
            "name": "test tag",
            "color": "red",
            "use_count": 0,
            "created_at": "2019-01-10T16:40:46.983393+08:00",
            "updated_at": "2019-01-10T16:40:46.983393+08:00"
        },
        {
            "id": "bgrgh9t5jj83bqe15580",
            "ent_id": "bgrg80l5jj83bqe154fg",
            "creator": "bgrg80l5jj83bqe154h0",
            "name": "test tag2",
            "color": "blue",
            "use_count": 0,
            "created_at": "2019-01-10T16:56:07.364274+08:00",
            "updated_at": "2019-01-10T16:56:07.364274+08:00"
        }
    ]
}
```

update visitor tag
------------------

``` {.bash}
curl -X PUT \
  http://localhost:8089/admin/api/v1/visitor_tags/bgrgh9t5jj83bqe15580 \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ3JnODBsNWpqODNicWUxNTRmZyIsInVzZXJfaWQiOiJiZ3JnODBsNWpqODNicWUxNTRoMCIsImV4cCI6MTU0Nzk3MzQ5OCwiaXNzIjoiQ3VzdG1JTSJ9.dbVwIGJ9Lge3bon7hopyQFQtmnw6DivM314VpItGlow' \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 92435e18-c558-4dd2-a244-01ef0500e67b' \
  -H 'cache-control: no-cache' \
  -d '{
    "name": "test tag update",
    "color": "浅蓝"
}'
```

Response:

``` {.json}
{
    "code": 0,
    "body": {
        "id": "bgrgh9t5jj83bqe15580",
        "ent_id": "bgrg80l5jj83bqe154fg",
        "creator": "bgrg80l5jj83bqe154h0",
        "name": "test tag update",
        "color": "浅蓝",
        "use_count": 0,
        "created_at": "2019-01-10T16:56:07.364274+08:00",
        "updated_at": "2019-01-10T09:00:40.680643Z"
    }
}
```

delete visitor tag
------------------

``` {.bash}
curl -X DELETE \
  http://localhost:8089/admin/api/v1/visitor_tags/bgrga3l5jj83bqe1557g \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRfaWQiOiJiZ3JnODBsNWpqODNicWUxNTRmZyIsInVzZXJfaWQiOiJiZ3JnODBsNWpqODNicWUxNTRoMCIsImV4cCI6MTU0Nzk3MzQ5OCwiaXNzIjoiQ3VzdG1JTSJ9.dbVwIGJ9Lge3bon7hopyQFQtmnw6DivM314VpItGlow' \
  -H 'Postman-Token: 7b6483cb-9fb6-465c-92c9-d94273746175' \
  -H 'cache-control: no-cache'
```

Response:

``` {.json}
{
  "code": 0
}
```

add visitor to black list
-------------------------

POST /admin/api/v1/blacklists

Request:

``` {.go}
type addBlacklistReq struct {
    EntID   string `json:"ent_id"`
    TraceID string `json:"trace_id"`
    VisitID string `json:"visit_id"`
    ConvID  string `json:"conv_id"`
}
```

``` {.json}
{
  "ent_id": "",
  "trace_id": "",
  "visit_id": "",
  "conv_id": ""
}
```

Response:

``` {.json}
```

结束会话
--------

PUT /admin/api/v1/conversations/:conversation\_id/end

Request:

``` {.json}
{
  "end_by": ""
}
```

Response:

``` {.json}
{
  "code": 0
}
```

对话转接
--------

PUT /admin/api/v1/conversations/:conversation\_id/transfer

Request:

``` {.go}
type TransferConversationReq struct {
    TraceID     string `json:"trace_id"`
    TargetAgent string `json:"target_agent"`
}
```

``` {.json}
{
  "trace_id": "",
  "target_agent": ""
}
```

添加小结
--------

PUT /admin/api/v1/conversations/:conversation\_id/summary

Request:

``` {.json}
{
 "content": ""
}
```

历史对话搜索
------------

POST /admin/api/v1/conversations/search

Request:

``` {.json}
{
  "visitor_name": "visit_name",
  "start_time_range": {
    "start_at": "2019-01-10T21:28:28.442293+08:00",
    "end_at": "2019-01-10T21:28:28.442294+08:00"
  },
  "end_time_range": {
    "start_at": "2019-01-11T07:28:28.442294+08:00",
    "end_at": "2019-01-11T17:28:28.442305+08:00"
  },
  "agent_id": "agent_id",
  "tags": [
    "tag1",
    "tag2"
  ],
  "city_name": "北京",
  "eval_level": 5,
  "quality_grade": "A",
  "source": "",
  "telephone": "18866899090",
  "message_content": "some message",
  "ip": "192.168.1.12",
  "remark": "a_good_visitor",
  "offset": 0,
  "limit": 100
}
```

Response:

``` {.json}
```

坐席发送消息
------------

POST /admin/api/v1/conversations/:conversation_id/messages

  参数            类型     描述
  --------------- -------- ------
  trace_id       string   
  agent_id       string   
  content         string   
  content_type   string
  internal       bool  

Response:

``` {.json}
```

访客发送消息
------------

POST /admin/api/v1/conversations/:conversation_id/messages

  参数            类型     描述
  trace\_id       string   
  agent\_id       string   
  content         string   
  content\_type   string   


获取会话消息
------------

GET  /admin/api/v1/conversations/:conversation_id/messages?offset=0&limit=10 

Response:

``` {.json}
```

撤回会话消息
------------

PUT /admin/api/v1/conversations/:conversation_id/messages/:message_id/revoke 

Response:

``` {.json}
{
  "code": 0
}
```

增加权限
-------------

POST /admin/api/v1/roles

Request:

```json
{
  "name": ""
}
```

Response:

```json
{
  "id": "",
  "ent_id": "",
  "name": ""
}
```

增加或更新角色的权限
-------------------

POST /admin/api/v1/roles/:role_id/perms

Request:

```json
{
  "perm_ids": []
}
```

Response:

```json

```

获取角色的权限
--------------

GET /admin/api/v1/roles/:role_id/perms

Response:

```json
{
  "code": 0,
  "body": []
}
```

获取租户的所有权限
---------------

GET /admin/api/v1/perms

Response:

```json
{
    "code": 0,
    "body": {
        "agent_settings": [
            {
                "id": "bgmnicujo0p67acb7kt0",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "agent_settings",
                "name": "check_manage_blacklist",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7l00",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "agent_settings",
                "name": "check_update_agent_allocation",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7l1g",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "agent_settings",
                "name": "check_update_agent_evaluation",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7kvg",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "agent_settings",
                "name": "check_update_auto_message",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7l0g",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "agent_settings",
                "name": "check_update_conversation_invite",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7l10",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "agent_settings",
                "name": "check_update_conversation_rule",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7ku0",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "agent_settings",
                "name": "check_update_prechat_form",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7kv0",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "agent_settings",
                "name": "check_update_quick_reply",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7kug",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "agent_settings",
                "name": "check_update_tag",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7ktg",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "agent_settings",
                "name": "check_update_visitor_queue",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            }
        ],
        "conversation": [
            {
                "id": "bgmnicujo0p67acb7kl0",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "conversation",
                "name": "end_others_conversation",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7kkg",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "conversation",
                "name": "others_transfer_conversation",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7kk0",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "conversation",
                "name": "transfer_others_conversation",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            }
        ],
        "customer": [
            {
                "id": "bgmnicujo0p67acb7kng",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "customer",
                "name": "check_update_customer_card",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7ko0",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "customer",
                "name": "check_update_customer_list",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7kn0",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "customer",
                "name": "use_customer_app",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            }
        ],
        "data_report": [
            {
                "id": "bgmnicujo0p67acb7kp0",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "data_report",
                "name": "export_data",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7kog",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "data_report",
                "name": "use_data_report",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            }
        ],
        "ent_info": [
            {
                "id": "bgmnicujo0p67acb7kq0",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "ent_info",
                "name": "check_ent_agent_group_info",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7kpg",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "ent_info",
                "name": "check_update_ent_info",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7ks0",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "ent_info",
                "name": "check_update_ent_role",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7ksg",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "ent_info",
                "name": "check_update_ent_safety",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7krg",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "ent_info",
                "name": "check_update_visitor_limit",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7kr0",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "ent_info",
                "name": "create_update_delete_agent_account",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7kqg",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "ent_info",
                "name": "create_update_delete_agent_group",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            }
        ],
        "history_conversation": [
            {
                "id": "bgmnicujo0p67acb7klg",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "history_conversation",
                "name": "check_others_conversation",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7km0",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "history_conversation",
                "name": "export_conversation",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7kmg",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "history_conversation",
                "name": "update_conversation_summary",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            }
        ],
        "integrate_settings": [
            {
                "id": "bgmnicujo0p67acb7kh0",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "integrate_settings",
                "name": "check_update_chat_link",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7kgg",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "integrate_settings",
                "name": "check_update_website_plugin",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            }
        ],
        "visitor": [
            {
                "id": "bgmnicujo0p67acb7kig",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "visitor",
                "name": "add_delete_visitor_tag",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7kj0",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "visitor",
                "name": "add_update_visitor_card",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7kjg",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "visitor",
                "name": "block_unblock_visitor",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7ki0",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "visitor",
                "name": "invite_visitor",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            },
            {
                "id": "bgmnicujo0p67acb7khg",
                "ent_id": "bgmnicujo0p67acb7kag",
                "app_name": "visitor",
                "name": "update_agent_status",
                "created_at": "2019-01-03T02:53:39.24449Z",
                "updated_at": "2019-01-03T02:53:39.24449Z"
            }
        ]
    }
}
```

获取留言列表
-------------

GET /admin/api/v1/leave_messages?status=handled/unhandled&offset=0&limit=10


Response:

```go
type LeaveMessage struct {
	ID        string    `json:"id"`         // id
	EntID     string    `json:"ent_id"`     // ent_id
	Content   string    `json:"content"`    // content
	Status    string    `json:"status"`     // status
	CreatedAt time.Time `json:"created_at"` // created_at
	UpdatedAt time.Time `json:"updated_at"` // updated_at

	// xo fields
	_exists, _deleted bool
}
```

```json
{
  "code": 0,
  "body": [
    {
     "id": "", 
     "ent_id": "", 
     "content": "", 
     "status": "",
     "created_at": "",
     "updated_at": ""
    }
  ]
}
```

更新留言状态
------------

PUT /admin/api/v1/leave_messages/:msg_id

| 名称      | 描述 |
| ----------- | ----------- |
| status      | 留言处理状态(handled/unhandled)       |

Request:

```json
{
  "status": "handled"
}
```

Response:

```json
{
  "code": 0,
  "body": null
}
```