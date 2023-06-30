package assetweb

import (
	"github.com/obnahsgnaw/assetweb/service/cors"
	"io"
)

type Option func(*Server)

func AccessWriter(w io.Writer) Option {
	return func(s *Server) {
		s.accessWriter = w
	}
}

func ErrorWriter(w io.Writer) Option {
	return func(s *Server) {
		s.errorWriter = w
	}
}

func Cors(c *cors.Config) Option {
	return func(s *Server) {
		s.corsCnf = c
	}
}

func CorsAll() Option {
	return func(s *Server) {
		s.corsCnf = &cors.Config{
			AllowOrigin:      "*",
			AllowCredentials: true,
		}
	}
}

func CorsOne(origin string) Option {
	return func(s *Server) {
		s.corsCnf = &cors.Config{
			AllowOrigin:      origin,
			AllowCredentials: true,
		}
	}
}

func TrustedProxies(proxies []string) Option {
	return func(s *Server) {
		s.trustedProxies = proxies
	}
}
