# 微服务模块
## 用户服务模块
- Golang
- mysql 使用sqlx进行操作
- protobuf && grpc 集成
- viper 读取本地配置，主要是nacos配置中心相关配置
- 集成nacos来获取服务相关的配置信息
- 集成consul来实现服务的注册，以及配置健康检查

## 商品服务模块
- Golang
- mysql 使用sqlx进行操作
- protobuf && grpc 集成
- viper 读取本地配置，主要是nacos配置中心相关配置
- 集成nacos来获取服务相关的配置信息
- 集成consul来实现服务的注册，以及配置健康检查
- 商品相关数据sql文件在项目中可找到，测试数据来自京东生鲜

## 库存服务模块
- Golang
- mysql 使用sqlx进行操作
- protobuf && grpc 集成
- viper 读取本地配置，主要是nacos配置中心相关配置
- 集成nacos来获取服务相关的配置信息
- 集成consul来实现服务的注册，以及配置健康检查
- 集成redis模块，引入redis分布式锁来解决库存变化引起的一系列问题

## 订单服务模块
- Golang
- mysql 使用sqlx进行操作
- protobuf && grpc 集成
- viper 读取本地配置，主要是nacos配置中心相关配置
- 集成nacos来获取服务相关的配置信息
- 集成consul来实现服务的注册，以及配置健康检查
- 添加本地事务来实现创建订单时订单表、购物车表、订单商品表数据的原子性