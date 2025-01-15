package assetweb

import (
	"github.com/obnahsgnaw/http/cors"
)

type Option func(*Server)

func Cors(c *cors.Config) Option {
	return func(s *Server) {
		s.corsCnf = c
	}
}

func RouteDebug(fg bool) Option {
	return func(s *Server) {
		s.routeDebug = fg
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

func CacheTtl(ttl int64) Option {
	return func(s *Server) {
		s.cacheTtl = ttl
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

func Replace(rp map[string]func([]byte) []byte) Option {
	return func(s *Server) {
		if rp != nil {
			s.replace = rp
		}
	}
}
