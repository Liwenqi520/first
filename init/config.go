package init

import (
	"first/init/myredis"
	"first/init/mysql"
	"first/internal/app/log"
	"first/internal/defines"
	"first/internal/lang"
	"fmt"
	"time"

	"github.com/Liwenqi520/errorx"
	"github.com/Liwenqi520/gormx"
	"github.com/Liwenqi520/helper"
	"github.com/Liwenqi520/i18n"
	"github.com/Liwenqi520/logger"
	"github.com/Liwenqi520/redisx"
	"github.com/Liwenqi520/timingwheel"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

var (
	// AppConfig 应用配置
	AppConfig Config
	Cron *cron.Cron
	Wheel *timingwheel.TimingWheel
)

// Config 配置信息
type Config struct {
	ServiceIP     string        `mapstructure:"service_ip" default:":127.0.0.1"`
	ServiceScheme string        `mapstructure:"service_scheme" default:":http"`
	ServicePort   string        `mapstructure:"service_port" default:"9001"`
	Gateways      []string      `mapstructure:"gateways"`
	Debug         bool          `mapstructure:"debug" default:"true"`
	DefaultLang   string        `mapstructure:"default_lang" default:"zh"`
	AllowOrigins  []string      `mapstructure:"allow_origins"`
	AllowHeaders  []string      `mapstructure:"allow_headers"`
	AllowMethods  []string      `mapstructure:"allow_methods"`
	MySQL         gormx.Config  `mapstructure:"mysql"`
	MySQLSlave    gormx.Config  `mapstructure:"mysql_slave"` // 验证主从复制
	Redis         redisx.Config `mapstructure:"redis"`
	Logger        logger.Config `mapstructure:"log"`
	UniqueID      struct {
		Type      string `default:"xxx"`
		Snowflake struct {
			Node  int64 // 节点
			Epoch int64 //时期
		}
	}
	Gateway struct {
		HTTPPort string `mapstructure:"http_port" default:":8081"`
	}
	AliSMS           helper.AliSms      `mapstructure:"ali_sms"`
	Email            helper.EmailConfig `mapstructure:"email"`
	AppList          []string           `mapstructure:"app_list"`
	TimeZone         string             `mapstructure:"time_zone"`
	NBIDeviceService struct {
		DeviceAdminURL    string   `mapstructure:"device_admin_url"`
		DeviceAdminTcpURL string   `mapstructure:"device_admin_tcp_url"`
		Username          string   `mapstructure:"username"`
		Password          string   `mapstructure:"password"`
		Type              string   `mapstructure:"type"`
		AllowIP           []string `mapstructure:"allow_ip"`
		IsGoDeviceAdmin   bool     `mapstructure:"is_go_device_admin"`
		Timeout           int64    `mapstructure:"timeout"`
	} `mapstructure:"nbi_device_service"`
	OldDeviceAdminPlatform struct {
		DeviceAdminURL string `mapstructure:"device_admin_url"`
		Username       string `mapstructure:"username"`
		Password       string `mapstructure:"password"`
	} `mapstructure:"old_device_admin_platform"`
	Log    logger.Config `mapstructure:"log"`
	DfsURL string        `mapstructure:"dfs_url"`
	CaiYun struct {
		URL string `mapstructure:"url"`
		Key string `mapstructure:"key"`
	} `mapstructure:"caiyun"`
	TestPushDevice struct {
		Domain     string   `mapstructure:"domain"`
		CompanyID  string   `mapstructure:"company_id"`
		DeviceList []string `mapstructure:"device_list"`
	} `mapstructure:"test_push_device"`
	MapConfig struct {
		GoogleMapUrl string `mapstructure:"google_map_url"`
		GoogleMapKey string `mapstructure:"google_map_key"`
		BaiduMapUrl  string `mapstructure:"baidu_map_url"`
		BaiduMapKey  string `mapstructure:"baidu_map_key"`
	}
	RoleList []struct {
		RoleName  string `mapstructure:"role_name"`
		RoleRules string `mapstructure:"role_rules"`
		RoleType  int64  `mapstructure:"role_type"`
	} `mapstructure:"role_list"`
	SystemRoleList []struct {
		RoleName  string `mapstructure:"role_name"`
		RoleRules string `mapstructure:"role_rules"`
		RoleType  int64  `mapstructure:"role_type"`
	} `mapstructure:"system_role_list"`
	OfflineHour       int64 `mapstructure:"offline_hour"`
	CheckedPermission struct {
		Checked  []string `mapstructure:"checked"`
		Disabled []string `mapstructure:"disabled"`
	} `mapstructure:"checked_permission"`
	WaterQualityDeviceConfig struct {
		TotalTimes int64 `mapstructure:"total_times"`
	} `mapstructure:"water_quality_device_config"`
	EmailConfig struct {
		SmtpHost       string `mapstructure:"smtp_host"`
		SmtpPort       int    `mapstructure:"smtp_port"`
		SenderEmail    string `mapstructure:"sender_email"`
		SenderPassword string `mapstructure:"sender_password"`
	} `mapstructure:"email_config"`
	SmsConfig struct {
		AccessKeyId     string `mapstructure:"access_key_id"`
		AccessKeySecret string `mapstructure:"access_key_secret"`
		RegionId        string `mapstructure:"region_id"`
		SignName        string `mapstructure:"sign_name"`
		TemplateCode    string `mapstructure:"template_code"`
	} `mapstructure:"sms_config"`
	IsPushTestData       int64    `mapstructure:"is_push_test_data"`
	WorkShopServerIpList []string `mapstructure:"workshop_server_ip_list"`
	DiskThreshold        int      `mapstructure:"disk_threshold"`
	IsDiskCheck          int      `mapstructure:"is_disk_check"`
	LuzhouUrl            string   `mapstructure:"luzhou_url"`
	ErLangUrl            string   `mapstructure:"erlang_url"`
	LangjiuIpList        struct {
		IsSyncMini int      `mapstructure:"is_sync_mini"`
		Luzhou     []string `mapstructure:"luzhou"`
		Erlang     []string `mapstructure:"erlang"`
	} `mapstructure:"langjiu_ip_list"`
	MpProgram struct {
		AppID     string `mapstructure:"app_id"`
		AppSecret string `mapstructure:"app_secret"`
	} `mapstructure:"mp_program"`
	AesKey  string `mapstructure:"aes_key"`
	BaseUrl string `mapstructure:"base_url"`
	IsDev   int64  `mapstructure:"is_dev"`
}

// ConfigLoad 应用配置加载
func ConfigLoad(paths ...string) {
	vp := viper.New()
	vp.SetConfigName("app") // name of config file
	vp.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	for _, path := range paths {
		fmt.Println(path)
		vp.AddConfigPath(path) // path to look for the config file in
	}
	err := vp.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	err = vp.Unmarshal(&AppConfig) // 读取app.yml文件并赋值给AppConfig
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

// ConfigInit 配置初始化
func ConfigInit(version string, paths ...string) {
	// 配置
	ConfigLoad(paths...)
	// logger初始化
	logger.Init(AppConfig.Log, logger.String("service.name", "liwenqi"), logger.String("service.version", defines.AppVersion))

	errorx.CodePrefix = defines.AppErrorsCodePrefix

	// 语言
	i18n.LoadLangPack(lang.LangErrorCode, lang.LangKey, lang.LangTemplate)

	// 初始化redis
	myredis.RedisInit(AppConfig.Redis)
	// 初始化helper中的reids
	helper.UUID{}.Init(AppConfig.Redis)

	// 初始化mysql
	mysql.MysqlInit(AppConfig.MySQL, AppConfig.Logger)

	// 日志
	e := log.InitHook(myredis.RedisCli())
	if e != nil {
		panic(fmt.Sprintf("日志数据库初始化失败:%s", e.Error()))
	}
	logLevel := log.DebugLevel
	isCaller := true
	if !AppConfig.Debug {
		logLevel = log.InfoLevel
		isCaller = false
	}
	log.Init(logLevel,
		log.WithEnableDbWrite(true),
		log.WithServerName("first"),
		log.WithVersion(defines.Version),
		log.WithCaller(isCaller),
	)
	
	// 全局定时器，触发场景中用到
	Cron = cron.New(cron.WithSeconds())
	Cron.Start()

	// 全局时间轮
	Wheel = timingwheel.NewTimingWheel(time.Second, 3600)
	Wheel.Start()
}