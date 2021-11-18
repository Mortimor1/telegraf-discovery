package telegraf

import (
	"github.com/BurntSushi/toml"
	"github.com/Mortimor1/telegraf-discovery/pkg/logging"
	"log"
	"os"
)

type Influxdb_v2 struct {
	Urls         []string `toml:"urls"`
	Token        string   `toml:"token"`
	Organization string   `toml:"organization"`
	Bucket       string   `toml:"bucket"`
	Timeout      string   `toml:"timeout"`
}

type Outputs struct {
	Influxdb_v2 []Influxdb_v2 `toml:"influxdb_v2"`
}

type Telegraf struct {
	Agent   Agent   `toml:"agent"`
	Inputs  Inputs  `toml:"inputs"`
	Outputs Outputs `toml:"outputs"`
}

type Agent struct {
	Interval       string `toml:"interval"`
	Flush_interval string `toml:"flush_interval"`
}

type Inputs struct {
	Ping []Ping `toml:"ping"`
}

type Ping struct {
	Ipv6    bool     `toml:"ipv6"`
	Method  string   `toml:"method"`
	Timeout float64  `toml:"timeout"`
	Count   int      `toml:"count"`
	Urls    []string `toml:"urls"`
}

func (t *Telegraf) LoadConfig(path string) {
	logger := logging.GetLogger()
	logger.Infof("Read Telegraf Config %s", path)
	if _, err := toml.DecodeFile(path, &t); err != nil {
		logger.Fatal(err)
	}
}

func (t *Telegraf) SaveConfig(path string) {
	logger := logging.GetLogger()
	logger.Infof("Write Telegraf Config %s", path)
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		// failed to create/open the file
		log.Fatal(err)
	}

	if err := toml.NewEncoder(f).Encode(&t); err != nil {
		// failed to encode
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		// failed to close the file
		log.Fatal(err)

	}
}
