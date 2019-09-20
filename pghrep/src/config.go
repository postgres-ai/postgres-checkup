package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

const CONFIG_PARAM_PREFIX = "CONFIG"

type CheckupProjectConfig map[string]string

type CheckupConfig []CheckupProjectConfig

func loadConfig(path string) (CheckupConfig, error) {
	var config CheckupConfig
	if !FileExists(path) {
		return config, fmt.Errorf("Config file '%s' not found.", path)
	}

	data, cerr := ioutil.ReadFile(path)
	if cerr != nil {
		return config, fmt.Errorf("Cannot read file '%s'. %s", path, cerr)
	}

	uerr := yaml.Unmarshal([]byte(data), &config)
	if uerr != nil {
		return config, fmt.Errorf("Cannot parse config file: '%s'. %s", path, uerr)
	}

	return config, nil
}

func outputConfig(config CheckupConfig) {
	if len(config) > 0 {
		// Now we support only one project config
		for key, value := range config[0] {
			fmt.Printf("%s__%s=\"%s\"\n", CONFIG_PARAM_PREFIX, strings.Replace(key, "-", "_", -1), value)
		}
	}
}
