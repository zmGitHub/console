# ES API

## ES 创建索引

```bash
curl -X PUT "localhost:9200/twitter"
```

## ES 创建索引的mapping

```bash
curl -X PUT "localhost:9200/test" -H 'Content-Type: application/json' -d'
{
    "settings" : {
        "number_of_shards" : 1
    },
    "mappings" : {
        "_doc" : {
            "properties" : {
                "field1" : { "type" : "text" }
            }
        }
    }
}
'
```