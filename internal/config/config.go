package config

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/netip"
	"os"
	"reflect"

	"dario.cat/mergo"
	"github.com/go-viper/mapstructure/v2"
	"github.com/itzg/go-flagsfiller"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Device struct {
		Name string
		MTU  int
	}

	Listener struct {
		Addr net.TCPAddr
	}

	Log struct {
		Level slog.Level
	}

	SecretPath string
	Secret     string

	Peers []Peer
}

type Peer struct {
	Addr net.TCPAddr
}

func New() (conf Config) {
	conf = Config{}

	conf.Device.MTU = 1430
	conf.Device.Name = "tun1"

	conf.Listener.Addr = *net.TCPAddrFromAddrPort(netip.MustParseAddrPort("0.0.0.0:53475"))

	conf.Log.Level = slog.LevelInfo

	conf.Peers = []Peer{}

	conf.Secret = ""
	conf.SecretPath = ""

	return
}

func StringToTCPAddrHook() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf(net.TCPAddr{}) {
			return data, nil
		}

		addr, err := netip.ParseAddrPort(data.(string))
		if err == nil {
			return net.TCPAddrFromAddrPort(addr), nil
		}

		// fallback for when addr contains a hostname
		return net.ResolveTCPAddr("tcp", data.(string))
	}
}

func StringToLogLevelHook() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf(slog.LevelInfo) {
			return data, nil
		}

		switch data.(string) {
		case "debug":
			return slog.LevelDebug, nil
		case "info":
			return slog.LevelInfo, nil
		case "warn":
			return slog.LevelWarn, nil
		case "error":
			return slog.LevelError, nil
		default:
			return nil, fmt.Errorf("can't parse log level: %s", data.(string))
		}
	}
}

func verifyConfig(conf *Config) error {
	if (conf.Secret == "") == (conf.SecretPath == "") {
		return errors.New(`one of "secret" or "secretPath" has to be defined`)
	}

	return nil
}

func ReadConfig() (conf Config, err error) {
	var (
		confPath  string
		file, env Config
	)
	conf = New()

	// _ = flagsfiller.New()
	filler := flagsfiller.New(flagsfiller.WithEnv("PASSAGE"))
	err = filler.Fill(flag.CommandLine, &env)
	if err != nil {
		return
	}
	flag.StringVar(&confPath, "config", "./config.yml", "config file path")

	flag.Parse()

	f, err := os.ReadFile(confPath)
	if err != nil {
		return
	} else {
		var raw interface{}
		err = yaml.Unmarshal(f, &raw)
		if err == nil {
			decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				DecodeHook: mapstructure.ComposeDecodeHookFunc(
					StringToLogLevelHook(),
					StringToTCPAddrHook(),
				),
				Result: &file,
			})

			err = decoder.Decode(raw)
			if err != nil {
				return
			}
		}
	}

	mergo.MergeWithOverwrite(&conf, file)
	mergo.MergeWithOverwrite(&conf, env)

	if err := verifyConfig(&conf); err != nil {
		return conf, fmt.Errorf("error verifying config: %w", err)
	}

	return
}

func (conf *Config) GetSecret() ([]byte, error) {
	if conf.Secret != "" {
		return []byte(conf.Secret), nil
	} else {
		return os.ReadFile(conf.SecretPath)
	}
}
