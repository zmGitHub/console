# message format

example

```json
{
  "creator": "visitor",
  "from_type": "agent",
  "content": "test msg",
  "content_type": "text",
  "created_at": "2019-01-07T13:40:39.371276Z"
}
```

| 参数         | 类型   | 描述   |
| ------------ | ------ | ------ |
| creator     | string |     visitor/{agent_id}   |
| from_type       | string | agent/visitor       |
| content        | string |        |
| content_type        | string |   text/picture/sticker/file     |
| created_at        | string |        |