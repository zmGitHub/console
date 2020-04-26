# login

POST /signin

## response

```json
{
    "ent_id": "bgmnicujo0p67acb7kag",
    "user_id": "bgmnicujo0p67acb7kc0",
    "uuid": "bgmnicujo0p67acb7kc0-bgmnio6jo0p67acb7l2g"
}
```

| 参数         | 类型   | 描述   |
| ------------ | ------ | ------ |
| ent_id     | string |        |
| user_id     | string |        |
| uuid     | string | 在后面请求的请求头中加入名为UUID, 值为response uuid的变量, 用于统计在线用户数   |