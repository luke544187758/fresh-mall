module luke544187758/inventory-srv

go 1.16

require (
	github.com/bwmarrin/snowflake v0.3.0
	github.com/fsnotify/fsnotify v1.5.4
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-sql-driver/mysql v1.6.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/nacos-group/nacos-sdk-go v1.1.1
	github.com/spf13/viper v1.12.0
	go.uber.org/zap v1.17.0
	google.golang.org/grpc v1.46.2
	google.golang.org/protobuf v1.28.0
	gopkg.in/yaml.v3 v3.0.1 // indirect
	luke544187758 v0.0.0-00010101000000-000000000000
)

replace luke544187758 => ../common
