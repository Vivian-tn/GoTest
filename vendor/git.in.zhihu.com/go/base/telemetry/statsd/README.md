# Graphite 客户端

Statsd 协议的 Graphite 指标客户端。

## 示例

```
// 初始化包含应用名前缀的指标客户端
client, err := statsd.New("your-appname")

// 记录 count 类型的指标
client.Count("name", 1)

// 记录 gauge 类型的指标
client.Gauge("gauge", 1)
```
