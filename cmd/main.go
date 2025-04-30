package main

import (
	"github.com/gookit/color"
	"github.com/obnahsgnaw/application"
	"github.com/obnahsgnaw/application/pkg/url"
	"github.com/obnahsgnaw/assetweb"
	"github.com/obnahsgnaw/assetweb/config"
	"github.com/obnahsgnaw/assetweb/html"
	"github.com/obnahsgnaw/goutils/runtimeutil"
	"log"
	"os"
	"strings"
)

func main() {
	var app *application.Application
	runtimeutil.HandleRecover(func(errMsg, stack string) {
		if app != nil {
			app.Logger().Error(errMsg)
		}
	})

	cnf, err := config.Parse()
	if err != nil {
		color.Error.Println("config parse failed, err=" + err.Error())
		os.Exit(1)
	}

	app = application.New(cnf.Http.Name,
		application.CusCluster(application.NewCluster(cnf.Application.Id, cnf.Application.Name)),
		application.Debug(func() bool {
			return cnf.Application.Debug
		}),
		application.Logger(cnf.Log),
	)
	defer app.Release()

	s := assetweb.New(app, cnf.Http.Name, url.New(cnf.Application.InternalIp, cnf.Http.Port),
		assetweb.Cors(cnf.Cors),
		assetweb.TrustedProxies(cnf.Http.TrustedProxies),
		assetweb.RouteDebug(cnf.Http.RouteDebug),
		assetweb.CacheTtl(86400),
		assetweb.Replace(map[string]func([]byte) []byte{
			"/config.json": func(b []byte) []byte {
				return []byte(strings.ReplaceAll(string(b), "127.0.0.1", cnf.Http.ApiHost))
			},
		}),
	)

	if dir := cnf.Http.Directory(); dir != "" {
		if err = s.RegisterDir(dir); err != nil {
			color.Error.Println("dir failed, err=" + err.Error())
			os.Exit(2)
		}
	}

	s.RegisterAsset(&html.FS, "www")

	app.AddServer(s)

	app.Run(func(err error) {
		color.Error.Println(err.Error())
	})

	app.Wait()
	log.Println("Server done")
}
