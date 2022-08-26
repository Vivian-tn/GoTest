# TZone 服务端 & 客户端

知乎内部 Thrift RPC 框架。

## 示例

```
// 初始化客户端
client := tzone.NewClient("NoteInfoService", tzone.TargetName("target_name"), tzone.Timeout(200*time.Millisecond))

// 初始化服务端
tzone.NewServer(map[string]tzone.TProcessor{
    "NoteInfoService": processor,
    "NoteService":     noteProcessor,
}).Run(":8090")
```
