# [WIP] Telemetry

> 本模块主要面向内部组件库开发者，业务方一般无需主动引入

我们将业务逻辑中执行过程分为两类 Transaction 和 Segment，我们将指标、日志、链路追踪、错误收集的集成都集中到这两类抽象中，

## Transaction

Transaction 指独立的完整执行过程，当前默认对 Transaction 进行了如下分类：

1. HTTP，http server 中处理单个请求
2. gRPC & TZone，rpc server 处理单个的远程调用
3. Worker，worker 处理单个 job
4. Exec，一些后台独立运行的任务，跟 worker 类型类似

一个简化后的示例：

```go
func Exec(jobName string, f func(ctx context.Context) Error) Error {
    // 首先初始化 txn，传入当前
	txn, ctx, e := StartTransaction(context.Background(), &Transaction{
		System: TransactionExec,
		Method: jobName,
	})
	if e != nil {
		return WrapErrWithUnknownClass(e)
	}

	var err Error
	sentry.Recover(func() {
		err = f(ctx)
	}, func(e error) {
		err = WrapErrWithUnknownClass(e)
	})
	txn.End(ctx, err)

	return err
}
```

在 Transaction 结束的时候会尝试结束当前的 opentracing span 和记录当前的 halo span 打点、错误日志。

## Segment

Segment 指在 Transaction 中触发的外部依赖调用，当前默认对 Segment 进行了如下分类：

1. DatabaseSegment，追踪数据库访问，细分类有 MySQL 和 Redis
2. HTTPSegment，追踪往外部的 HTTP 调用
3. ProducerSegment，追踪队列发布
4. RPCSegment，追踪 RPC 调用

对于每一类 Segment 都有不同的初始化参数，这些传入的参数在结束 Segment 的时候会组装 opentracing span、halo span 等，具体的区别可以直接参考代码实现和 opentelemetry 规范。

当前 telemetry 模块还额外实现了工具方法，用于快速接入：
1. WrapRoundTripper 函数，可以对传入的 http.RoundTripper 注入 HTTPSegment 的完整埋点逻辑
2. Middleware 函数，可以返回 http.Handler 用作 http server 的中间件

## 环境变量

1. `ENABLE_TELEMETRY_RECORD_ARGUMENTS` 设置为 `1` 时, telemetry 将自动记录请求参数到 log 和 span 中
