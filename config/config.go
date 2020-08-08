package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Type struct {
	Storage struct {
		DbPath     string
		Engine     string // default is "wiredTiger" vs "inMemory"
		WiredTiger struct {
			EngineConfig struct {
				CacheSizeGB float32
			}
		}
	}

	SystemLog struct {
		Destination     string // default is syslog if not specified, use "file"
		Path            string // required if Destination = "file"
		TimeStampFormat string // default is "iso8601-local" vs "iso8601-utc"
		LogAppend       bool   // default is false
		Verbosity       uint   // default is zero
	}

	Security struct {
		Authorization     string // default is "disabled" vs. "enabled"
		JavascriptEnabled bool   `default:"true"` // default is true
	}

	Net struct {
		Port             uint     // default: 27017, 27018 if shard member, 27019 if CSRS member
		BindIp           string   // default is "127.0.0.1"
		Ipv6             bool     // default is false
		UnixDomainSocket struct { // not valid on Windows
			Enabled bool `default:"true"` // default is true
		} `omitwindows:"true"`
	}

	ProcessManagement struct {
		Fork bool
	}

	SetParameter struct {
		AuthenticationMechanisms string // default ?
	}
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

// Write config file fname to fpath
// Create directories for dbPath and log destination
func WriteConfig(x *Type, fpath string, fname string, isWindows bool) error {
	res2 := x.ToYaml(isWindows)
	err := os.MkdirAll(fpath, 0777)
	if err != nil {
		return fmt.Errorf("MkDirAll error: %v", err)
	}
	fn := filepath.Join(fpath, fname)
	fd, err := os.Create(fn)
	if err != nil {
		return fmt.Errorf("file create error: %v", err)
	}
	_, err = io.Copy(fd, res2)
	//err = ioutil.WriteFile(filepath.Join(fpath, fname), []byte(res2), 0644)
	if err != nil {
		_ = fd.Close()
		return fmt.Errorf("file write error: %v", err)
	}
	err = fd.Close()
	if err != nil {
		return fmt.Errorf("file close error: %v", err)
	}
	err = os.Chmod(fn, 0644)
	if err != nil {
		return fmt.Errorf("Chmod error: %v", err)
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
