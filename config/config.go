package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Type struct {
	Storage struct {
		DbPath     string `yaml:"dbPath"`     // required
		Engine     string `yaml:",omitempty"` // default is "wiredTiger" vs "inMemory"
		WiredTiger struct {
			EngineConfig struct {
				CacheSizeGB float32 `yaml:"cacheSizeGB,omitempty"` // default is to omit this option
			} `yaml:"engineConfig,omitempty"`
		} `yaml:"wiredTiger,omitempty"`
	}

	SystemLog struct {
		Destination     string // default is syslog if not specified, use "file"
		Path            string // required if Destination = "file"
		TimeStampFormat string `yaml:"timeStampFormat,omitempty"` // default is "iso8601-local" vs "iso8601-utc"
		LogAppend       bool   `yaml:"logAppend,omitempty"`       // default is false
		Verbosity       uint   `yaml:",omitempty"`                // default is zero
	} `yaml:"systemLog"`

	Security struct {
		Authorization     string `yaml:",omitempty"`        // default is "disabled" vs. "enabled"
		JavascriptEnabled bool   `yaml:"javascriptEnabled"` // default is true
	} `yaml:",omitempty"`

	Net struct {
		Port             uint     `yaml:",omitempty"`       // default: 27017, 27018 if shard member, 27019 if CSRS member
		BindIp           string   `yaml:"bindIp,omitempty"` // default is "127.0.0.1"
		Ipv6             bool     `yaml:",omitempty"`       // default is false
		UnixDomainSocket struct { // not valid on Windows
			Enabled bool // default is true
		} `yaml:"unixDomainSocket"`
	} `yaml:",omitempty"`

	ProcessManagement struct {
		Fork bool `yaml:",omitempty"`
	} `yaml:"processManagement,omitempty"`

	SetParameter struct {
		AuthenticationMechanisms string `yaml:"authenticationMechanisms,omitempty"` // default ?
	} `yaml:"setParameter,omitempty"`
}

// Config struct with MongoDB defaults applied
var MongoDBDefaults *Type

// Config struct with our opinionated defaults enabled
var OurDefaults *Type

func init() {

	// set defaults in config as MongoDB would have them
	t1 := &Type{}
	t1.Security.JavascriptEnabled = true
	t1.Net.UnixDomainSocket.Enabled = true
	MongoDBDefaults = t1

	// set opinionated config
	t2 := &Type{}
	// leave javascript and unix domain socket disabled, as well as logAppend
	t2.Storage.Engine = "wiredTiger"
	t2.Storage.WiredTiger.EngineConfig.CacheSizeGB = 0.5
	t2.SystemLog.Destination = "file"
	t2.Security.Authorization = "enabled"
	t2.Net.BindIp = "0.0.0.0"
	OurDefaults = t2
}

func Foo() {
	var yamlString = `
storage:
  dbPath:
  engine: wiredTiger
  wiredTiger:
    engineConfig:
      cacheSizeGB: 1

processManagement:
  fork: false
`
	out := Type{}
	//out.Storage.WiredTiger.EngineConfig.CacheSizeGB = 1
	err := yaml.Unmarshal([]byte(yamlString), &out)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return
	}

	x := Type{}
	x.ProcessManagement.Fork = true
	x.Storage.WiredTiger.EngineConfig.CacheSizeGB = 1.5
	err = printConfig(&x)
	if err != nil {
		fmt.Printf("printConfig error: %v\n", err)
		return
	}
	x.Storage.WiredTiger.EngineConfig.CacheSizeGB = 0
	_ = printConfig(&x)
	x.ProcessManagement.Fork = false
	_ = printConfig(&x)
	x.Security.Authorization = "enabled"
	x.Net.Port = 27017
	x.Net.BindIp = "0.0.0.0"
	x.Net.UnixDomainSocket.Enabled = false
	_ = printConfig(&x)
	_ = printConfig(MongoDBDefaults)
	_ = printConfig(OurDefaults)
	var something Type
	something = *OurDefaults
	_ = printConfig(&something)
}

func printConfig(x *Type) error {
	res2, err := yaml.Marshal(x)
	if err != nil {
		return fmt.Errorf("Marshal error: %v\n", err)
	}
	fmt.Printf("\nYAML Config:\n%s\n", res2)
	return nil
}

// Write config file fname to fpath
// Create directories for dbPath and log destination
func WriteConfig(x *Type, fpath string, fname string) error {
	res2, err := yaml.Marshal(x)
	if err != nil {
		return fmt.Errorf("Marshal error: %v", err)
	}
	err = os.MkdirAll(fpath, 0777)
	if err != nil {
		return fmt.Errorf("MkDirAll error: %v", err)
	}
	err = ioutil.WriteFile(filepath.Join(fpath, fname), []byte(res2), 0644)
	if err != nil {
		return fmt.Errorf("WriteFile error: %v", err)
	}
	err = os.MkdirAll(x.Storage.DbPath, 0777)
	if err != nil {
		return fmt.Errorf("MkDirAll error: %v", err)
	}
	err = os.MkdirAll(filepath.Base(x.SystemLog.Path), 0777)
	if err != nil {
		return fmt.Errorf("MkDirAll error: %v", err)
	}
	return nil
}
