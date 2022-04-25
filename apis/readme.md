# golang 微服务开发
## user 管理 
- gin web
- viper 配置
- 路由配置
- zap日志配置
- jwt继承
- grpc继承
- protobuf
- grpc调用微服务模块
- 集成consul
- release模式下端口号自动生成
- 负载均衡配置，移除了consul发现服务的逻辑，使用负载均衡来进行处理
- 集成配置中心nacos，viper只作为本地读取nacos配置，服务其它相关配置均放在nacos中