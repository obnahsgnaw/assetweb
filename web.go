package assetweb

import (
	"errors"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
	"github.com/obnahsgnaw/application"
	"github.com/obnahsgnaw/application/endtype"
	"github.com/obnahsgnaw/application/pkg/logging/logger"
	"github.com/obnahsgnaw/application/pkg/url"
	"github.com/obnahsgnaw/application/pkg/utils"
	"github.com/obnahsgnaw/application/servertype"
	"github.com/obnahsgnaw/assetweb/internal"
	"github.com/obnahsgnaw/assetweb/service/cors"
	"go.uber.org/zap"
	"io"
	"os"
)

type Server struct {
	name           string
	app            *application.Application
	port           int
	logger         *zap.Logger
	err            error
	engine         *gin.Engine
	accessWriter   io.Writer
	errorWriter    io.Writer
	corsCnf        *cors.Config
	trustedProxies []string
	staticDir      string
	staticAsset    *assetfs.AssetFS
}

func New(app *application.Application, name string, port int, option ...Option) *Server {
	if name == "" {
		name = "asset-web"
	}
	s := &Server{
		name: name,
		app:  app,
		port: port,
	}
	s.logger, s.err = logger.New(name, app.LogConfig(), app.Debugger().Debug())
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
	s.engine, s.err = internal.NewEngine(&internal.EngineConfig{
		Debug:          s.app.Debugger().Debug(),
		AccessWriter:   s.accessWriter,
		ErrWriter:      s.errorWriter,
		TrustedProxies: s.trustedProxies,
		Cors:           s.corsCnf,
	})
	if s.err != nil {
		failedCb(s.err)
		return
	}
	if !s.initStaticDir() {
		s.initAsset()
	}
	go func() {
		host := url.Host{
			Ip:   "",
			Port: s.port,
		}
		s.logger.Info(utils.ToStr(s.name, "[", host.String(), "] listen and serving..."))
		if err := s.engine.Run(host.String()); err != nil {
			failedCb(err)
		}
	}()
}

func (s *Server) initStaticDir() bool {
	if s.staticDir != "" {
		s.engine.Static("/", s.staticDir)
		return true
	}
	return false
}

func (s *Server) initAsset() {
	if s.staticAsset != nil {
		s.engine.StaticFS("/", s.staticAsset)
	}
}

func (s *Server) Release() {
	_ = s.logger.Sync()
}
