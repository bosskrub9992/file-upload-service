package configs

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string
	}
	UploadFile struct {
		MaxFileSizeMBLimit string              `mapstructure:"max_file_size_mb_limit"`
		AllowFileExtension map[string][]string `mapstructure:"allow_file_extension"`
		Bucket             string
	} `mapstructure:"upload_file"`
	Email struct {
		From string
	}
	Secret struct {
		Email struct {
			Host     string
			Port     int
			Username string
			Password string
		}
	}
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./configs")    // local
	viper.AddConfigPath("../configs")   // unit test
	viper.AddConfigPath("/app/configs") // docker
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.AutomaticEnv() // to overwrite with env var, use upper+snake case
}

func New() *Config {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}
	return &cfg
}
