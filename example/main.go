package main

import (
	"github.com/obnahsgnaw/application"
	"github.com/obnahsgnaw/assetweb"
	"log"
)

func main() {
	app := application.New("demo", "Demo")

	s := assetweb.New(app, "test", 8099)

	err := s.RegisterDir("example/web")
	if err != nil {
		panic(err)
	}

	app.AddServer(s)

	app.Run(func(err error) {
		panic(err)
	})

	log.Println("Exited")
}
