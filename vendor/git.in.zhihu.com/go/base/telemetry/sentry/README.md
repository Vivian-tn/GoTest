# Sentry 客户端

帮助你管理好 Panic。

## 示例

```
// 配置应用引用路径
sentry.SetIncludePaths([]string{"git.in.zhihu.com/bu/app"})

// 仅捕获 Panic，不上报
var err error
sentry.Recover(func() {
    resp, err = handler(ctx, req)
}, func(e error) {
    err = e
})

// 捕获 Panic 并上报到 Sentry
sentry.CapturePanic(ctx, sentry.Tags{
    "module": "do-some-dangerous-work",
}, func() {
    panic("just panic")
})
```
