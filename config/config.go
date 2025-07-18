package config

import (
	"errors"
	"fmt"
	"github.com/gookit/color"
	"github.com/jessevdk/go-flags"
	"github.com/obnahsgnaw/application/pkg/logging/logger"
	"github.com/obnahsgnaw/application/pkg/utils"
	"github.com/obnahsgnaw/assetweb/version"
	"github.com/obnahsgnaw/http/cors"
	"github.com/spf13/viper"
	"net"
	"os"
	"path/filepath"
)

var Conf *Config

type Config struct {
	Version     bool   `short:"v" long:"version" description:"show version"`
	IniFile     string `short:"c" long:"conf" description:"Ini file"`
	Application *Application
	Log         *logger.Config
	Cors        *cors.Config
	Http        *Http
}

type Application struct {
	Id         string `long:"cluster-id" description:"Cluster id" required:"true" default:"web"`
	Name       string `long:"cluster-name" description:"Cluster name" required:"true" default:"web"`
	InternalIp string `long:"internal-ip" description:"Server ip address" `
	Debug      bool   `long:"debug" description:"Enable debug"`
}

type Http struct {
	curDir         string
	Name           string        `long:"name" description:"http name" required:"true" default:"web"`
	TrustedProxies []string      `long:"trusted-ip" description:"Trusted proxy ip, the gateway ip, multi"`
	RouteDebug     bool          `long:"route-debug" description:"enable http engine to debug mode"`
	Port           int           `long:"port" description:"http port" required:"true" default:"8080"`
	Dir            string        `long:"dir" description:"http dir"`
	Current        bool          `long:"current" description:"http use current pwd dir"`
	DirRoot        bool          `long:"dir-root" description:"use dir as root, otherwise static asset as root and dir as fallback"`
	CacheTtl       int64         `long:"cache-ttl" description:"cache ttl" default:"86400"`
	Replace        []ReplaceItem `long:"replace" description:"replace file item"`
}

type ReplaceItem struct {
	File  string            `long:"file" description:"file to replace"`
	Items map[string]string `long:"items" description:"item to replace"`
}

func (s *Http) Directory() string {
	if s.Dir != "" {
		return s.Dir
	}
	if s.Current {
		return s.curDir
	}
	return ""
}

func configError(msg string, err error) error {
	return utils.TitledError("config error", msg, err)
}

func Parse() (*Config, error) {
	var opt Config

	if _, err := flags.Parse(&opt); err != nil {
		var flagErr *flags.Error
		ok := errors.As(err, &flagErr)
		if !ok || flagErr.Type != flags.ErrHelp {
			color.Error.Println("config parse failed, err=" + err.Error())
		}
		os.Exit(0)
	}
	Conf = &opt

	if Conf.Version {
		fmt.Println(version.Info().String())
		os.Exit(0)
	}

	if Conf.IniFile != "" {
		cc := viper.New()
		cc.SetConfigFile(Conf.IniFile)
		if err := cc.ReadInConfig(); err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if errors.As(err, &configFileNotFoundError) {
				return nil, configError("config file not found", err)
			} else {
				return nil, configError("parse config file failed", err)
			}
		}
		if err := cc.Unmarshal(Conf); err != nil {
			return nil, configError("decode config failed", err)
		}
	}

	if Conf.Application == nil {
		Conf.Application = &Application{
			Id:         "web",
			Name:       "web",
			InternalIp: "",
			Debug:      false,
		}
	}
	if Conf.Application.InternalIp == "" {
		Conf.Application.InternalIp = getLocalIp()
	}

	if Conf.Http.Name == "" {
		Conf.Http.Name = "web"
	}

	if Conf.Http.Port <= 0 {
		return nil, configError("port required", nil)
	}

	if Conf.Http.Dir != "" {
		dr, err := filepath.Abs(Conf.Http.Dir)
		if err != nil {
			return nil, configError("dir err,", nil)
		}
		Conf.Http.Dir = dr
	}
	if Conf.Http.Current {
		dr, err := os.Getwd()
		if err != nil {
			return nil, configError("current dir fetch failed", nil)
		}
		Conf.Http.curDir = dr
	}

	return Conf, nil
}

func getLocalIp() (ip string) {
	var err error
	var conn net.Conn

	if conn, err = net.Dial("udp", "8.8.8.8:80"); err != nil {
		return "127.0.0.1"
	}
	defer func() { _ = conn.Close() }()

	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
