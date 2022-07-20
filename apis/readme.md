# golang 微服务开发
## user 管理 
- 集成gin框架，来完成web服务的搭建
- viper 配置
- zap日志配置
- jwt集成，进行权限的认证
- grpc集成，来访问为服务模块获取相关信息
- protobuf
- 集成consul来实现对user服务的发现
- 负载均衡配置，移除了consul发现服务的逻辑，使用负载均衡来进行处理
- 集成配置中心nacos，viper只作为本地读取nacos配置，服务其它相关配置均放在nacos中
- 注册user-web服务到注册中心

## 商品 管理
- 集成gin框架，来完成web服务的搭建
- viper 配置
- zap日志配置
- jwt集成，进行权限的认证
- grpc集成，来访问为服务模块获取相关信息
- protobuf
- 集成consul来实现服务的发现
- 负载均衡
- 集成配置中心nacos，viper只作为本地读取nacos配置，服务其它相关配置均放在nacos中
- 注册goods-web服务到注册中心