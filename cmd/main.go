package main

import (
	"github.com/gookit/color"
	"github.com/obnahsgnaw/application"
	"github.com/obnahsgnaw/application/pkg/url"
	"github.com/obnahsgnaw/application/pkg/utils"
	"github.com/obnahsgnaw/assetweb"
	"github.com/obnahsgnaw/assetweb/config"
	"github.com/obnahsgnaw/assetweb/html"
	"io/fs"
	"log"
	http2 "net/http"
	"os"
)

func main() {
	var cnf *config.Config
	var err error
	var app *application.Application
	utils.RecoverHandler("asset web", func(err, stack string) {
		if app != nil {
			app.Logger().Error(err)
		}
	})
	if cnf, err = config.Parse(); err != nil {
		color.Error.Println("config parse failed, err=" + err.Error())
		os.Exit(1)
	}
	app = application.New(application.NewCluster(cnf.Application.Id, cnf.Application.Name), cnf.Http.Name,
		application.Debug(func() bool {
			return false
		}),
		application.Logger(cnf.Log),
	)
	defer app.Release()

	s := assetweb.New(app, cnf.Http.Name, url.New(cnf.Application.InternalIp, cnf.Http.Port),
		assetweb.Cors(cnf.Cors),
		assetweb.TrustedProxies(cnf.Http.TrustedProxies),
		assetweb.RouteDebug(cnf.Http.RouteDebug),
	)

	if dir := cnf.Http.Directory(); dir != "" {
		if err = s.RegisterDir(dir); err != nil {
			color.Error.Println("dir failed, err=" + err.Error())
			os.Exit(2)
		}
	}
	sub, _ := fs.Sub(html.FS, "www")
	s.RegisterAsset(http2.FS(sub))

	app.AddServer(s)

	app.Run(func(err error) {
		color.Error.Println(err.Error())
	})
	app.Wait()

	log.Println("Exited")
}
