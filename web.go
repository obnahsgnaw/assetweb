package assetweb

import (
	"embed"
	"errors"
	"github.com/gin-contrib/gzip"
	"github.com/obnahsgnaw/application"
	"github.com/obnahsgnaw/application/endtype"
	"github.com/obnahsgnaw/application/pkg/logging/logger"
	"github.com/obnahsgnaw/application/pkg/url"
	"github.com/obnahsgnaw/application/pkg/utils"
	"github.com/obnahsgnaw/application/servertype"
	"github.com/obnahsgnaw/assetweb/html"
	"github.com/obnahsgnaw/http"
	"github.com/obnahsgnaw/http/cors"
	"github.com/obnahsgnaw/http/engine"
	"go.uber.org/zap"
	"io/fs"
	http2 "net/http"
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
	staticAsset    *embed.FS
	staticRoot     string
	routeDebug     bool
	etagManager    *EtagManager
	cacheTtl       int64
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
func (s *Server) RegisterAsset(asset *embed.FS, root string) {
	s.staticAsset = asset
	s.staticRoot = root
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
	var err error
	cnf := &engine.Config{
		Name:           s.name,
		DebugMode:      s.routeDebug,
		AccessWriter:   nil,
		ErrWriter:      nil,
		TrustedProxies: s.trustedProxies,
		Cors:           s.corsCnf,
		DefFavicon:     false,
	}
	cnf.AccessWriter, err = logger.NewDefAccessWriter(s.app.LogConfig(), func() bool {
		return s.app.Debugger().Debug()
	})
	if err != nil {
		failedCb(err)
		return
	}
	cnf.ErrWriter, err = logger.NewDefErrorWriter(s.app.LogConfig(), func() bool {
		return s.app.Debugger().Debug()
	})
	if err != nil {
		failedCb(err)
		return
	}
	s.engine, s.err = http.Default(s.host.Ip, s.host.Port, cnf)
	s.engine.Engine().Use(gzip.Gzip(gzip.DefaultCompression), CacheMiddleware(s, s.cacheTtl))
	if s.err != nil {
		failedCb(s.err)
		return
	}
	if !s.initStaticDir() {
		s.initAsset()
	}
	if s.cacheTtl > 0 {
		if err = s.etagManager.Init(); err != nil {
			failedCb(err)
			return
		}
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
		s.etagManager = newEtagManagerWithDir(s.staticDir)
		s.engine.Engine().Static("/", s.staticDir)
		return true
	}
	return false
}

func (s *Server) initAsset() {
	if s.staticAsset != nil {
		s.etagManager = newEtagManagerWithFs(s.staticAsset, s.staticRoot)
		sub, _ := fs.Sub(html.FS, "www")
		s.engine.Engine().StaticFS("/", http2.FS(sub))
	}
}

func (s *Server) Release() {
	_ = s.logger.Sync()
}
