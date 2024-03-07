package assetweb

import (
	"errors"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/obnahsgnaw/application"
	"github.com/obnahsgnaw/application/endtype"
	"github.com/obnahsgnaw/application/pkg/url"
	"github.com/obnahsgnaw/application/pkg/utils"
	"github.com/obnahsgnaw/application/servertype"
	"github.com/obnahsgnaw/http"
	"github.com/obnahsgnaw/http/cors"
	"github.com/obnahsgnaw/http/engine"
	"go.uber.org/zap"
	"os"
)

type Server struct {
	name           string
	app            *application.Application
	host           url.Host
	logger         *zap.Logger
	err            error
	engine         *http.Http
	corsCnf        *cors.Config
	trustedProxies []string
	staticDir      string
	staticAsset    *assetfs.AssetFS
	routeDebug     bool
}

func New(app *application.Application, name string, host url.Host, option ...Option) *Server {
	if name == "" {
		name = "asset-web"
	}
	s := &Server{
		name: name,
		app:  app,
		host: host,
	}
	s.logger = app.Logger().Named(name)
	s.With(option...)
	return s
}

func (s *Server) ID() string {
	return s.name
}

func (s *Server) Name() string {
	return s.name
}

func (s *Server) Type() servertype.ServerType {
	return ""
}

func (s *Server) EndType() endtype.EndType {
	return ""
}

// RegisterDir register a dir
func (s *Server) RegisterDir(dirPath string) error {
	if dirPath == "" {
		return errors.New("dir required")
	}
	stat, err := os.Stat(dirPath)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return errors.New(dirPath + " is not a dir")
	}
	s.staticDir = dirPath

	return nil
}

// RegisterAsset register the asset
func (s *Server) RegisterAsset(asset *assetfs.AssetFS) {
	s.staticAsset = asset
}

func (s *Server) With(option ...Option) {
	for _, o := range option {
		o(s)
	}
}

func (s *Server) Run(failedCb func(error)) {
	if s.err != nil {
		failedCb(s.err)
		return
	}
	s.engine, s.err = http.Default(s.host, &engine.Config{
		Name:           s.name,
		DebugMode:      s.routeDebug,
		LogDebug:       s.app.Debugger().Debug() || s.app.LogConfig() == nil,
		TrustedProxies: s.trustedProxies,
		Cors:           s.corsCnf,
		LogCnf:         s.app.LogConfig(),
	})
	if s.err != nil {
		failedCb(s.err)
		return
	}
	if !s.initStaticDir() {
		s.initAsset()
	}
	go func() {
		s.logger.Info(utils.ToStr(s.name, "[http://", s.host.String(), "] listen and serving..."))
		if err := s.engine.RunAndServ(); err != nil {
			failedCb(err)
		}
	}()
}

func (s *Server) initStaticDir() bool {
	if s.staticDir != "" {
		s.engine.Engine().Static("/", s.staticDir)
		return true
	}
	return false
}

func (s *Server) initAsset() {
	if s.staticAsset != nil {
		s.engine.Engine().StaticFS("/", s.staticAsset)
	}
}

func (s *Server) Release() {
	_ = s.logger.Sync()
}
