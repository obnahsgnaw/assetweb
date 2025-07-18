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
	app  *application.Application
	c    *config.Config
	fs   *embed.FS
	root string
	opts []Option
	err  error
}

func Default(options ...Option) *Web {
	s := &Web{
		opts: options,
	}
	s.init()
	return s
}

func (w *Web) WithFS(fs embed.FS, root string) {
	w.fs = &fs
	w.root = root
}

func (w *Web) init() {
	runtimeutil.HandleRecover(func(errMsg, stack string) {
		if w.app != nil {
			w.app.Logger().Error(errMsg)
		}
	})

	cnf, err := config.Parse()
	if err != nil {
		w.err = errors.New("config parse failed, err=" + err.Error())
		return
	}
	w.c = cnf
	w.app = application.New(cnf.Http.Name,
		application.CusCluster(application.NewCluster(cnf.Application.Id, cnf.Application.Name)),
		application.Debug(func() bool {
			return cnf.Application.Debug
		}),
		application.Logger(cnf.Log),
	)
	defer w.app.Release()

	var rp map[string]func([]byte) []byte
	for _, item := range cnf.Http.Replace {
		for k, v := range item.Items {
			rp[item.File] = func(b []byte) []byte {
				return []byte(strings.ReplaceAll(string(b), k, v))
			}
		}
	}

	options := append([]Option{
		Cors(cnf.Cors),
		TrustedProxies(cnf.Http.TrustedProxies),
		RouteDebug(cnf.Http.RouteDebug),
		CacheTtl(cnf.Http.CacheTtl),
		Replace(rp),
	}, w.opts...)
	s := New(w.app, cnf.Http.Name, url.New(cnf.Application.InternalIp, cnf.Http.Port), options...)

	if dir := cnf.Http.Directory(); dir != "" {
		if err = s.RegisterDir(dir, cnf.Http.DirRoot); err != nil {
			w.err = errors.New("dir failed, err=" + err.Error())
			return
		}
	}

	if w.fs != nil && w.root != "" {
		s.RegisterAsset(w.fs, w.root)
	}

	w.app.AddServer(s)
	return
}

func (w *Web) Serve() {
	if w.err != nil {
		color.Error.Println("config parse failed, err=" + w.err.Error())
		os.Exit(1)
	}
	w.app.Run(func(err error) {
		color.Error.Println(err.Error())
	})

	if w.Config().Http.Directory() != "" {
		color.Info.Println(fmt.Sprintf("asset web[%s] root directory: %s", w.Config().Http.Name, w.Config().Http.Directory()))
	}
	color.Info.Println(fmt.Sprintf("asset web[%s] serving at: http://%s:%d", w.Config().Http.Name, w.Config().Application.InternalIp, w.Config().Http.Port))
	w.app.Wait()
	color.Info.Println(fmt.Sprintf("asset web[%s] done", w.Config().Http.Name))
}

func (w *Web) Config() *config.Config {
	return w.c
}
