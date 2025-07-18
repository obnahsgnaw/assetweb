package web

import (
	"embed"
	"errors"
	"fmt"
	"github.com/gookit/color"
	"github.com/obnahsgnaw/application"
	"github.com/obnahsgnaw/application/pkg/url"
	"github.com/obnahsgnaw/assetweb/config"
	"github.com/obnahsgnaw/goutils/runtimeutil"
	"os"
	"strings"
)

type Web struct {
	app    *application.Application
	server *Server
	cnf    *config.Config
	opts   []Option
	err    error
}

func Default(options ...Option) *Web {
	s := &Web{
		opts: options,
	}
	s.init()
	return s
}

func (w *Web) WithFS(fs embed.FS, root string) {
	if w.err != nil {
		return
	}
	if root != "" {
		w.server.RegisterAsset(&fs, root)
	}
}

func (w *Web) WithVersionProvider(p func() string) {
	if w.err != nil {
		return
	}
	w.cnf.SetVersionProvider(p)
}

func (w *Web) init() {
	if w.cnf, w.err = config.Parse(); w.err != nil {
		w.err = errors.New("config parse failed, err=" + w.err.Error())
		return
	}
	w.app = application.New(w.cnf.Http.Name,
		application.CusCluster(application.NewCluster(w.cnf.Application.Id, w.cnf.Application.Name)),
		application.Debug(func() bool {
			return w.cnf.Application.Debug
		}),
		application.Logger(w.cnf.Log),
	)

	var rp map[string]func([]byte) []byte
	for _, item := range w.cnf.Http.Replace {
		for k, v := range item.Items {
			rp[item.File] = func(b []byte) []byte {
				return []byte(strings.ReplaceAll(string(b), k, v))
			}
		}
	}
	options := append([]Option{
		Cors(w.cnf.Cors),
		TrustedProxies(w.cnf.Http.TrustedProxies),
		RouteDebug(w.cnf.Http.RouteDebug),
		CacheTtl(w.cnf.Http.CacheTtl),
		Replace(rp),
	}, w.opts...)
	w.server = New(w.app, w.cnf.Http.Name, url.New(w.cnf.Application.InternalIp, w.cnf.Http.Port), options...)

	if dir := w.cnf.Http.Directory(); dir != "" {
		if w.err = w.server.RegisterDir(dir, w.cnf.Http.DirRoot); w.err != nil {
			w.err = errors.New("dir failed, err=" + w.err.Error())
			return
		}
	}
	w.app.AddServer(w.server)
	return
}

func (w *Web) Serve() {
	if w.err != nil {
		color.Error.Println("config parse failed, err=" + w.err.Error())
		os.Exit(1)
	}
	runtimeutil.HandleRecover(func(errMsg, stack string) {
		if w.app != nil {
			w.app.Logger().Error(errMsg)
		}
	})
	defer w.app.Release()
	w.app.Run(func(err error) {
		color.Error.Println(err.Error())
	})
	if w.Config().Http.Directory() != "" {
		w.app.Logger().Info(fmt.Sprintf("root directory: %s", w.Config().Http.Directory()))
	}
	w.app.Wait()
}

func (w *Web) Config() *config.Config {
	return w.cnf
}
