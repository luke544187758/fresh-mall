package settings

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

var (
	Conf  = new(AppConfig)
	Nacos = new(NacosConfig)
)

type AppConfig struct {
	Port      int      `json:"port"`
	MachineID int64    `json:"machine_id"`
	Host      string   `json:"host"`
	Name      string   `json:"name"`
	Mode      string   `json:"mode"`
	StartTime string   `json:"start_time"`
	Tags      []string `json:"tags"`

	*MySQLConfig   `json:"mysql"`
	*ConsulConfig  `json:"consul"`
	*UseropService `json:"userop_service"`
}

type UseropService struct {
	Name string `json:"name"`
}

type MySQLConfig struct {
	Port         int    `json:"port"`
	MaxOpenConns int    `json:"max_open_conns"`
	MaxIdleConns int    `json:"max_idle_conns"`
	Host         string `json:"host"`
	User         string `json:"user"`
	Password     string `json:"password"`
	DB           string `json:"dbname"`
}

type ConsulConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

type NacosConfig struct {
	Port      uint64 `mapstructure:"port"`
	NameSpace string `mapstructure:"name_space"`
	UserName  string `mapstructure:"user_name"`
	Password  string `mapstructure:"password"`
	DataID    string `mapstructure:"data_id"`
	Group     string `mapstructure:"group"`
	Host      string `mapstructure:"host"`
}

func Init() (err error) {

	viper.SetConfigFile("./conf/config.yaml")

	err = viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig failed, err:%v\n", err)
		return
	}

	// 把读取到的配置信息反序列化到 Nacos 变量中
	if err := viper.Unmarshal(Nacos); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
		if err := viper.Unmarshal(Nacos); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		}
	})

	// 从nacos中获取配置信息
	//至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      Nacos.Host,
			Port:        Nacos.Port,
			ContextPath: "/nacos",
			Scheme:      "http",
		},
	}
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         Nacos.NameSpace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "../tmp/nacos/log",
		CacheDir:            "../tmp/nacos/cache",
		LogLevel:            "debug",
		Username:            Nacos.UserName,
		Password:            Nacos.Password,
	}
	//创建动态配置客户端的 推荐
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(err)
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: Nacos.DataID,
		Group:  Nacos.Group})
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(content), Conf)
	if err != nil {
		fmt.Printf("json.Unmarshal failed, err:%v\n", err)
	}

	return nil
}
