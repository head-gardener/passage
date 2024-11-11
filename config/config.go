package config

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
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
		Addr string
	}

	Log struct {
		Level slog.Level

	}

	Peers []Peer
}

type Peer struct {
	Addr net.UDPAddr
}

func NewConfig() (conf Config) {
	conf = Config{}

	conf.Device.MTU = 1430
	conf.Device.Name = "tun1"

	conf.Listener.Addr = "127.0.0.1:53475"

	conf.Log.Level = slog.LevelInfo

	conf.Peers = []Peer{}

	return
}

func StringToUDPAddrHook() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf(net.UDPAddr{}) {
			return data, nil
		}

		return net.ResolveUDPAddr("udp", data.(string))
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

func StringToLogLevel(s string, tag reflect.StructTag) (slog.Level, error) {
	fmt.Println("HEY!!!")
	os.Exit(1)
	return slog.LevelDebug, nil
}

func ReadConfig() (conf Config, err error) {
	var (
		confPath  string
		file, env Config
	)
	conf = NewConfig()

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
		if err != nil {
			return
		}

		decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				StringToLogLevelHook(),
				StringToUDPAddrHook(),
			),
			Result: &file,
		})

		err = decoder.Decode(raw)
		if err != nil {
			return
		}
	}

	mergo.MergeWithOverwrite(&conf, file)
	mergo.MergeWithOverwrite(&conf, env)

	return
}
