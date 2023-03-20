package bot

import (
	"fmt"

	"github.com/spf13/viper"
)

var viperConfig = viper.New()

func InitConfig() {
	viperConfig.SetConfigName("config")
	viperConfig.SetConfigType("yaml")
	viperConfig.AddConfigPath(".")
	if err := viperConfig.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config failed: %v", err))
	}
	// viperConfig.SetDefault("allowed_ids", []int64{})
}

var modeConfig = viper.New()

func InitChatMode() {
	modeConfig.SetConfigName("chat_modes")
	modeConfig.AddConfigPath(".")
	modeConfig.SetConfigType("yaml")
	if err := modeConfig.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config failed: %v", err))
	}
}
