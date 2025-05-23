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

	"github.com/head-gardener/passage/pkg/bee2/belt"
)

type Config struct {
	Device struct {
		Name string
		Addr string
	}

	Metrics struct {
		Enabled bool
		Addr    string
	}

	Listener struct {
		Addr net.TCPAddr
	}

	Socket struct {
		Enabled bool
		Path    string
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

func New() (conf *Config) {
	conf = new(Config)

	conf.Device.Addr = ""
	conf.Device.Name = "tun1"

	conf.Metrics.Enabled = false
	conf.Metrics.Addr = "0.0.0.0:9031"

	conf.Socket.Enabled = false
	conf.Socket.Path = "/var/run/passage.sock"

	conf.Listener.Addr = *net.TCPAddrFromAddrPort(netip.MustParseAddrPort("0.0.0.0:53475"))

	conf.Log.Level = slog.LevelInfo

	conf.Peers = []Peer{}

	conf.Secret = ""
	conf.SecretPath = ""

	return conf
}

type Command int

const (
	CommandServe Command = iota
	CommandStatus
)

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
		return errors.New(`one of "secret" or "secretPath" must be set`)
	}

	if conf.Device.Addr == "" {
		return errors.New(`"device.addr" must be set`)
	}

	return nil
}

func ReadConfig() (cmd Command, conf *Config, err error) {
	var (
		confPath  string
		file, env Config
	)
	conf = New()

	// FIXME: why is -quickcheks here?
	filler := flagsfiller.New(flagsfiller.WithEnv("PASSAGE"))
	err = filler.Fill(flag.CommandLine, &env)
	if err != nil {
		return
	}
	flag.StringVar(&confPath, "config", "./config.yml", "config file path")

	flag.Parse()

	if flag.NArg() != 0 {
		switch flag.Arg(0) {
		case "serve":
			cmd = CommandServe
		case "status":
			cmd = CommandStatus
		default:
			return cmd, conf, fmt.Errorf("unknown command: %s", flag.Arg(0))
		}
	}

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

	mergo.MergeWithOverwrite(conf, file)
	mergo.MergeWithOverwrite(conf, env)

	if err := verifyConfig(conf); err != nil {
		return cmd, conf, fmt.Errorf("error verifying config: %w", err)
	}

	if conf.Secret == "" {
		sec, err := os.ReadFile(conf.SecretPath)
		if err != nil {
			return cmd, nil, err
		}
		conf.Secret = string(sec)
	}
	key := belt.Key{}
	belt.Hash(key[:], []byte(conf.Secret), nil)
	conf.Secret = string(key[:])

	return
}

func (conf *Config) GetSecret() []byte {
	return []byte(conf.Secret)
}
