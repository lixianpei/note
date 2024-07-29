# 测试数据导入es
```
    curl -H "Content-Type: application/json" -XPOST "192.168.220.128:9200/bank/_bulk?pretty&refresh" --data-binary "@accounts.json"
    curl "192.168.220.128:9200/_cat/indices?v"
```
