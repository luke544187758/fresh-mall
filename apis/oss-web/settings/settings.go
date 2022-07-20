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
	Name      string   `json:"name"`
	Host      string   `json:"host"`
	Mode      string   `json:"mode"`
	Version   string   `json:"version"`
	StartTime string   `json:"start_time"`
	Tags      []string `json:"tags"`

	*JWTConfig    `json:"auth"`
	*LogConfig    `json:"log"`
	*ConsulConfig `json:"consul"`
	*OssService   `json:"oss_service"`
}

type JWTConfig struct {
	JWTExpire int64  `json:"jwt_expire"`
	Secret    string `json:"secret"`
}

type LogConfig struct {
	Level      string `json:"level"`
	Filename   string `json:"filename"`
	Mode       string `json:"mode"`
	MaxSize    int    `json:"max_size"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
}

type OssService struct {
	Timeout     int64  `json:"timeout"`
	ApiKey      string `json:"api_key"`
	ApiSecret   string `json:"api_secret"`
	Host        string `json:"host"`
	CallBackUrl string `json:"call_back_url"`
	UploadDir   string `json:"upload_dir"`
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

	// 把读取到的配置信息反序列化到 Conf 变量中
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
