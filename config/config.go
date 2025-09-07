package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig() (*Config, error) {
	var config Config

	// 設定 .env 檔案
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")

	// 啟用自動環境變數支援
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 將環境變數綁定到結構體
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析設定失敗: %w", err)
	}

	return &config, nil
}
