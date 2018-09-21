package app

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type AppConfig struct {
	Bind                     string   `yaml:"bind"`
	Topic                    string   `yaml:"topic"`
	BulkMessagesCount        int      `yaml:"bulk_message_count"`
	KafkaBrokers             []string `yaml:"kafka_brokers"`
	BulkMessageFlushInterval int      `yaml:"bulk_message_flush_interval_ms"`
	LogPath                  string   `yaml:"log_path"`
}

//解析配置文件
func ParseConfigFile(filepath string) (*AppConfig, error) {
	pconfig := AppConfig{}
	if config, err := ioutil.ReadFile(filepath); err == nil {

		if err = yaml.Unmarshal(config, &pconfig); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	return &pconfig, nil
}
