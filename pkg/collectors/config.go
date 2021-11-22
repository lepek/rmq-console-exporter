package collectors

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	*viper.Viper
	queueFilter *Filter
}

func NewConfig(configFilePath string) (*Config, error) {
	config := &Config{viper.New(), nil}
	if configFilePath == "" {
		return config, nil
	}
	configFile := filepath.Base(configFilePath)
	config.SetConfigName(strings.TrimSuffix(configFile, filepath.Ext(configFile)))
	config.AddConfigPath(filepath.Dir(configFilePath))
	if _, err := os.Stat(configFilePath); err != nil {
		return config, err
	}
	err := config.Init()
	return config, err
}

func (c *Config) Init() error {
	var err error
	if err = c.ReadInConfig(); err != nil {
		return err
	}
	if c.queueFilter, err = NewFilter(c.GetStringSlice("filters.queues")); err != nil {
		return err
	}
	log.Infof("Config loaded from %v", c.ConfigFileUsed())
	return nil
}

func (c *Config) IsEmpty() bool {
	if c.queueFilter == nil {
		return true
	}
	if c.queueFilter.Size() > 0 {
		return false
	}
	return true
}

func (c *Config) filterQueue(name string) bool {
	if c.IsEmpty() {
		return true
	}
	return c.queueFilter.Filter(name)
}