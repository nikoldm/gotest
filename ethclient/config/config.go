package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	InfuraURL  string
	PrivateKey string
	FromAddr   string
	ToAddr     string
}

func LoadConfig() *Config {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		// 如果找不到 .env 文件，继续使用系统环境变量
		log.Printf("警告: 未找到 .env 文件, 使用系统环境变量: %v", err)
	}

	cfg := &Config{
		InfuraURL:  os.Getenv("INFURA_URL"),
		PrivateKey: os.Getenv("PRIVATE_KEY"),
		FromAddr:   os.Getenv("FROM_ADDR"),
		ToAddr:     os.Getenv("TO_ADDR"),
	}

	// 调试输出
	logConfig(cfg)

	return cfg
}

func logConfig(cfg *Config) {
	log.Println("=== 配置信息 ===")

	if cfg.InfuraURL == "" {
		log.Println("INFURA_URL: 未设置")
	} else {
		// 隐藏部分 API Key 以保护隐私
		maskedURL := cfg.InfuraURL
		if len(maskedURL) > 30 {
			maskedURL = maskedURL[:25] + "..." + maskedURL[len(maskedURL)-5:]
		}
		log.Printf("INFURA_URL: %s", maskedURL)
	}

	if cfg.PrivateKey == "" {
		log.Println(" PRIVATE_KEY: 未设置")
	} else {
		maskedKey := cfg.PrivateKey
		if len(maskedKey) > 10 {
			maskedKey = maskedKey[:6] + "..." + maskedKey[len(maskedKey)-4:]
		}
		log.Printf("PRIVATE_KEY: %s", maskedKey)
	}

	log.Printf("FROM_ADDR: %s", cfg.FromAddr)
	log.Printf("TO_ADDR: %s", cfg.ToAddr)
	log.Println("=========.env end =========")
}
