# 测试数据导入es
```
    curl -H "Content-Type: application/json" -XPOST "192.168.220.128:9200/bank/_bulk?pretty&refresh" --data-binary "@accounts.json"
    curl "192.168.220.128:9200/_cat/indices?v"
```

# 搜索
## 任何结果
- took – 查询花费时长（毫秒）
- timed_out – 请求是否超时
- _shards – 搜索了多少分片，成功、失败或者跳过了多个分片（明细）
- max_score – 最相关的文档分数
- hits.total.value - 找到的文档总数
- hits.sort - 文档排序方式 （如没有则按相关性分数排序）
- hits._score - 文档的相关性算分 (match_all 没有算分)