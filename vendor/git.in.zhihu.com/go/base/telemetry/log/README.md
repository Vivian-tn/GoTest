# 日志客户端

兼容内部日志格式规范，错误自动上报 Sentry。

## 示例

```
// 按不同的 Level 打日志，没有 context 时传 context.TODO()
log.Info(ctx, "info log")
// => [I 2021-06-02 13:44:35.235 ba757ef97fce1f1d - education_web:10000] [{"logger":"default"}] info log
log.Error(ctx, "error log")
// => [E 2021-06-02 13:44:35.235 ba757ef97fce1f1d - education_web:10000] [{"logger":"default"}] error log
log.WithFields(ctx, log.Fields{"member_id": 123}).Info("info log")
// => [I 2021-06-02 13:44:35.235 ba757ef97fce1f1d - education_web:10000] [{"member_id": 123, "logger":"default"}] info log
log.WithFields(ctx, log.Fields{"member_id": 123}).WithError(errors.New("errmsg")).Error("error log")
// => [E 2021-06-02 13:44:35.235 ba757ef97fce1f1d - education_web:10000] [{"member_id": 123, "logger":"default", "error": "errmsg"}] error log
```

logger 实例

```
// 生成新的 logger 实例
logger := log.New("user_controller")
logger.Info(ctx, "info log")
// => [I 2021-06-02 13:44:35.235 ba757ef97fce1f1d - education_web:10000] [{"logger":"user_controller"}] info log

// 使用默认 logger
log.Info(ctx, "info log")
// => [I 2021-06-02 13:44:35.235 ba757ef97fce1f1d - education_web:10000] [{"logger":"default"}] info log
```

logger 配置

```
// 配置 Level
logger.SetLevel(log.WarnLevel)

logger.Info(ctx, "info log")
// => null
logger.Error(ctx, "error log")
// => [E 2021-06-02 13:44:35.235 ba757ef97fce1f1d - education_web:10000] [{"logger":"user_controller"}] error log

// 配置 FileLineNumber 的 Level
logger.SetFileLineNumberLevel(log.ErrorLevel)

logger.Warn(ctx, "warn log")
// => [W 2021-06-02 13:44:35.235 ba757ef97fce1f1d - education_web:10000] [{"logger":"user_controller"}] warn log
logger.Error(ctx, "error log")
// => [E 2021-06-02 13:44:35.235 ba757ef97fce1f1d main:56 education_web:10000] [{"logger":"user_controller"}] error log

// 配置 SentryLevel
logger.SetSentryLevel(log.WarnLevel)

logger.Warn(ctx, "warn log") // Level >= Warning 的日志都会上报 sentry

// 配置 outputs
file, _ := os.Create("/logs/log")
logger.SetOutputs(os.Stdout, file)

logger.Warn(ctx, "warn log") // 会同时在 stdout 和 file 中输出
```
