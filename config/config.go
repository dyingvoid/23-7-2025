package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port                   int      `json:"port"`
	NumMaxTasks            int      `json:"num_max_tasks"`
	NumMaxResourcesPerTask int      `json:"num_max_resources_per_task"`
	FileExtensionWhiteList []string `json:"file_extension_white_list"`
	FileDir                string   `json:"file_dir"`
}

// TODO add paths
func RequireConfig(env string) *Config {
	f, err := os.Open("configs/config." + env + ".json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var cfg Config
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}

func (cfg *Config) GetAllowedExtensions() map[string]struct{} {
	allowed := make(map[string]struct{}, len(cfg.FileExtensionWhiteList))
	for _, ext := range cfg.FileExtensionWhiteList {
		allowed[ext] = struct{}{}
	}
	return allowed
}
