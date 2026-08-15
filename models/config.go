package models

import (
	"fmt"
	"log"
	"os"
	"sublink/utils"

	"gopkg.in/yaml.v3"
)

// type Config struct {
// 	ID    int
// 	Key   string
// 	Value string
// }

// Config 配置结构体
type Config struct {
	JwtSecret             string `yaml:"jwt_secret"`              // JWT密钥
	APIEncryptionKey      string `yaml:"api_encryption_key"`      // API加密密钥
	ExpireDays            int    `yaml:"expire_days"`             // 过期天数
	Port                  int    `yaml:"port"`                    // 端口号
	DefaultSubscriptionID int    `yaml:"default_subscription_id"` // 默认分配订阅ID
	InviteRequired        bool   `yaml:"invite_required"`         // 是否启用邀请码注册
	UACheckEnabled        bool   `yaml:"ua_check_enabled"`        // 是否启用订阅User-Agent检测(仅代理软件可拉取)
	UACheckKeywords       string `yaml:"ua_check_keywords"`       // UA关键字列表(逗号分隔，命中任一即视为代理软件)
}

// DefaultUACheckKeywords 默认UA关键字(写入配置，可在后台修改，检测时只使用配置值)
const DefaultUACheckKeywords = "clash,mihomo,surge,v2ray,shadowrocket,quantumult,stash,sing-box,hiddify,nekobox,nekoray,loon,karing,outline,hysteria,trojan,sagernet,verge,subconverter,openclash,1.1.1.1,warp"

var comment string = `# jwt_secret: JWT密钥
# expire_days: token 过期天数
# port: 启动端口
# default_subscription_id: 新注册用户默认分配的订阅ID
# invite_required: 是否要求邀请码注册
# ua_check_enabled: 开启后订阅拉取仅允许代理软件(浏览器/脚本将被拦截)
# ua_check_keywords: UA关键字列表(逗号分隔)，命中任一即视为代理软件放行
`

// 初始化配置
func ConfigInit() {
	if _, err := os.Stat("./db"); os.IsNotExist(err) {
		if mkErr := os.MkdirAll("./db", os.ModePerm); mkErr != nil {
			log.Println("创建db目录失败:", mkErr)
			return
		}
	}
	// 检查配置文件是否存在
	if _, err := os.Stat("./db/config.yaml"); os.IsNotExist(err) {

		// 如果不存在则创建默认配置文件
		defaultConfig := Config{
			JwtSecret:             utils.RandString(31), // 生成随机JWT密钥
			APIEncryptionKey:      utils.RandString(31), // 生成随机API加密密钥
			ExpireDays:            14,
			Port:                  8000, // 默认端口
			DefaultSubscriptionID: 0,
			InviteRequired:        false,
			UACheckEnabled:        false,
			UACheckKeywords:       DefaultUACheckKeywords,
		}

		// 生成yaml文件
		data, err := yaml.Marshal(&defaultConfig)
		if err != nil {
			log.Println("生成默认配置文件失败:", err)
			return
		}
		data = []byte(comment + string(data)) // 添加注释
		err = os.WriteFile("./db/config.yaml", data, 0644)
		if err != nil {
			fmt.Println("写入文件失败:", err)
			return
		}
		log.Println("配置文件不存在，已创建默认配置文件")
	}
}

// 读取配置
func ReadConfig() Config {
	file, err := os.ReadFile("./db/config.yaml")
	if err != nil {
		log.Println(err)
	}
	cfg := Config{}
	yaml.Unmarshal(file, &cfg)
	return cfg
}

// 设置配置
func SetConfig(newCfg Config) {
	cfg := newCfg
	if cfg.JwtSecret == "" {
		cfg.JwtSecret = utils.RandString(31)
	}
	if cfg.APIEncryptionKey == "" {
		cfg.APIEncryptionKey = utils.RandString(31)
	}
	if cfg.ExpireDays == 0 {
		cfg.ExpireDays = 14
	}
	if cfg.Port == 0 {
		cfg.Port = 8000
	}
	// 写入文件
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		log.Println(err)
	}
	data = []byte(comment + string(data)) // 添加注释
	os.WriteFile("./db/config.yaml", data, 0644)
}
