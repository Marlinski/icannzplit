package zplit

import (
	"encoding/json"
	"fmt"
	"icannzplit/util"
	"icannzplit/vpn/ipvanish"
	"io/ioutil"
	"os"
	"path/filepath"
)

// config directory
const (
	IcannZplitConfigFile = "settings.json"
)

// Config the multiplexer
type Config struct {
	IPVanishSettings ipvanish.Settings `json:"ipvanish,omitempty"`    // ipvanish-specific settings
	ExcludedIPs      []string          `json:"exclude-ips,omitempty"` // exclude ips or iprange from being rerouted
}

// ConfigInit parse an existing config or create a default one
func ConfigInit(homedir string) *Config {
	configfile := filepath.Join(homedir, IcannZplitConfigFile)

	cfg, err := loadConfig(configfile)
	if err != nil {
		return newConfig()
	}
	return cfg
}

// NewConfig returns a new configuration
func newConfig() *Config {
	util.Log.Debugf("using default config")
	return &Config{
		IPVanishSettings: ipvanish.DefaultSettings(),
	}
}

// LoadConfig reads and parse a json as config
func loadConfig(configFile string) (*Config, error) {
	jsonFile, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("could not open the config file: %+v", err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("could not read the config file: %+v", err)
	}

	var zplitConfig Config
	err = json.Unmarshal(byteValue, &zplitConfig)
	if err != nil {
		return nil, fmt.Errorf("could not parse the config file: %+v", err)
	}

	util.Log.Debugf("reading config from: %#v\n", configFile)
	return &zplitConfig, nil
}

// Save the config in a json file
func (c *Config) Save(homedir string) error {
	configfile := filepath.Join(homedir, IcannZplitConfigFile)
	util.Log.Debugf("saving settings into %#v", configfile)

	// build destination directory if it doesn't exists
	err := os.MkdirAll(homedir, os.ModePerm)
	if err != nil {
		return err
	}

	// marshalling settings
	bytes, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	// write into file
	err = ioutil.WriteFile(configfile, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
