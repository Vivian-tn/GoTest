# ZAE 客户端

支持资源发现、获取环境变量、容器配额。

## 示例

```
// 资源发现
masterAddr, err := zae.DiscoveryOne(zae.ResourceRWRedis, "name", "primary")
slaveAddrs, err := zae.DiscoveryMany(zae.ResourceRWRedis, "name", "replica")
// 手动资源注册，供测试用
err := zae.Register(zae.ResourceMySQL, "name", "master", []string{
    "user:password@host:27546/database",
})

// 获取 & 判断环境变量
zae.App()
zae.Service()
zae.Region()
zae.Environment()
zae.IsDevelopEnv()
...

// 获取容器 CPU、内存配额
zae.TotalCPU()
zae.TotalMemory()

```
