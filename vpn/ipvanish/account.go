package ipvanish

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/Marlinski/icannzplit/util"

	"github.com/op/go-logging"
)

// IPVanish Const
const (
	IPVanishConfigDir = "ipvanish"
	IPVanishConfigURL = "https://www.ipvanish.com/software/configs/configs.zip"
	IPVanishCAName    = "ca.ipvanish.com.crt"
	IPVanishAuthName  = "auth.txt"
)

var (
	configNameRegexp = regexp.MustCompile("ipvanish-([^-]+?)-(.+?)-([^-]+-[^-]+?).ovpn")
)

var (
	configDir string
	certFile  string
	authFile  string
	configs   []TunnelConfig
	log       *logging.Logger
	table     int

	stopChan chan struct{}
	wg       sync.WaitGroup
)

// TunnelConfig single endpoint configuration
type TunnelConfig struct {
	FileName    string
	CountryCode string
	CountryID   string
	DeviceName  string
}

// Init initialize configurations
func Init(homedir, user, pass string, routingTable int) error {
	if _, err := os.Stat(homedir); os.IsNotExist(err) {
		os.Mkdir(homedir, os.ModePerm)
	}

	// create ipvanish config dir if it does not exists
	configdir := filepath.Join(homedir, IPVanishConfigDir)
	if _, err := os.Stat(configdir); os.IsNotExist(err) {
		os.Mkdir(configdir, os.ModePerm)
	}
	configDir = configdir

	// pull or load the openvpn configuration for ipvanish
	err := loadConfigs(configdir)
	if err != nil {
		return err
	}

	err = saveAuthFile(user, pass)
	if err != nil {
		return err
	}

	table = routingTable

	stopChan = make(chan struct{})
	wg = sync.WaitGroup{}
	return nil
}

// load all ipvanish configurations
func loadConfigs(configdir string) error {
	empty, err := util.IsEmpty(configdir)
	if err != nil {
		return fmt.Errorf("error reading the ipvanish directory: %+v", err)
	}

	// download the ipvanish configurations
	if empty {
		configszip := filepath.Join(configdir, "configs")
		util.Log.Debugf("ipvanish: pulling configurations from: %#v", IPVanishConfigURL)
		if err := util.DownloadFile(configszip, IPVanishConfigURL); err != nil {
			return fmt.Errorf("could not download the ipvanish config files: %+v", err)
		}

		util.Log.Debugf("ipvanish: unzipping configuration to %#v", configdir)
		_, err := util.Unzip(configszip, configdir)
		if err != nil {
			return fmt.Errorf("could not unzip the ipvanish config files: %+v", err)
		}
	}

	util.Log.Debugf("ipvanish: reading configurations in: %#v", configdir)
	files, err := ioutil.ReadDir(configdir)
	if err != nil {
		return fmt.Errorf("error reading the ipvanish directory: %+v", err)
	}

	// loading all the configurations
	openvpnConfigs := make([]TunnelConfig, 0)
	for _, config := range files {
		// openvpvn config
		if filepath.Ext(config.Name()) == ".ovpn" {
			match := configNameRegexp.FindStringSubmatch(config.Name())
			if len(match) == 4 {
				openvpnConfigs = append(openvpnConfigs, TunnelConfig{
					FileName:    config.Name(),
					CountryCode: match[1],
					CountryID:   match[3],
					DeviceName:  "ipvn-" + match[3],
				})
			}
		}

		// certification authority
		if config.Name() == IPVanishCAName {
			certFile = filepath.Join(configDir, config.Name())
		}
	}

	if certFile == "" {
		return fmt.Errorf("no certificate found")
	}

	configs = openvpnConfigs
	return nil
}

func saveAuthFile(user, pass string) error {
	authFile = filepath.Join(configDir, IPVanishAuthName)
	f, err := os.Create(authFile)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s\n%s\n", user, pass)
	if err != nil {
		return err
	}

	return nil
}

// DumpConfigs dumps all IPVanish configurations
func DumpConfigs() {
	log.Debugf("ipvanish: dumping configurations..")
	for _, config := range configs {
		log.Infof("[+] Country=%#v ID=%#v - Config=%#v\n", config.CountryCode, config.CountryID, config.FileName)
	}
}
