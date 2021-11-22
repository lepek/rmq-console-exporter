package collectors

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	configPath := "./tests_config.toml"

	payload := []string{`^.*\.dev$`, `^.*\.super\..*$`}
	WriteDummyConfig(configPath, payload)

	config, _ := NewConfig(configPath)
	config.OnConfigChange(func(e fsnotify.Event) {
		config.Init()
	})
	config.WatchConfig()

	assert.Equal(t, true, config.filterQueue("object_test.dev"))

	// Change the config file and write the new config
	payload = []string{`^.*\.super\..*$`}
	WriteDummyConfig(configPath, payload)

	assert.Equal(t, false, config.filterQueue("object_test.dev"))

	os.Remove(configPath)
}

func TestEmpty(t *testing.T) {
	config, _ := NewConfig("")
	assert.Equal(t, true, config.filterQueue("object_test.dev"))
}

func WriteDummyConfig(configPath string, payload []string) {
	viper.Set("filters", map[string][]string{})
	viper.Set("filters.queues", payload)
	viper.WriteConfigAs(configPath)
}