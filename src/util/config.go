package util

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Name   string            `yaml:"name"`
	Params map[string]string `yaml:"params"`
}

type APIConfig struct {
	Port   int    `yaml:"port"`
	Listen string `yaml:"listen"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	API      APIConfig      `yaml:"api"`
}

// 检查配置文件是否存在，如果不存在则创建默认配置文件
// @return bool 函数调用前是否存在配置文件
// @return bool 函数调用后是否创建了配置文件
// @return error 错误
func EnsureConfigFile() (bool, bool, error) {
	if !IsFileExist("config.yaml") {

		config := Config{
			Database: DatabaseConfig{
				Name: "mongodb",
				Params: map[string]string{
					"uri":      "mongodb://localhost:27017",
					"database": "qcg-center-data",
				},
			},
			API: APIConfig{
				Port:   8989,
				Listen: "0.0.0.0",
			},
		}

		data, err := yaml.Marshal(config)

		if err != nil {
			log.Println("Failed to marshal config.yaml")
			return false, false, err
		}

		err = ioutil.WriteFile("config.yaml", data, 0644)
		if err != nil {
			log.Println("Failed to create config.yaml")
			return false, false, err
		}

		return false, true, nil
	}

	return true, false, nil
}

// 加载配置文件
func LoadConfig() (*Config, error) {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Println("Failed to read config.yaml")
		return nil, err
	}

	obj := &Config{}
	err = yaml.Unmarshal(data, obj)
	if err != nil {
		log.Println("Failed to unmarshal config.yaml")
		return nil, err
	}

	return obj, nil
}

// Save 保存配置文件
func SaveConfig(obj *Config) error {
	data, err := yaml.Marshal(obj)
	if err != nil {
		log.Println("Failed to marshal config.yaml")
		return err
	}

	err = ioutil.WriteFile("config.yaml", data, 0644)
	if err != nil {
		log.Println("Failed to write config.yaml")
		return err
	}

	return nil
}
