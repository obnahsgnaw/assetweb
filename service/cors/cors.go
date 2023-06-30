package cors

import "strings"

var (
	DefCorsAllowMethods = []string{"OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE"}
	DefCorsAllowHeader  = []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization", "Language", "Request-Origin", "X-App-Id", "X-Security-Sign", "X-Security-Iv"}
	DefCorsExposeHeader = []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Cache-Control", "Content-Language", "Content-Type", "X-App-Id", "X-Security-Sign", "X-Security-Iv"}
)

type Config struct {
	AllowOrigin      string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
}

func (c Config) Origin() string {
	return c.AllowOrigin
}
func (c Config) Methods() string {
	var m []string
	if len(c.AllowMethods) == 0 {
		m = DefCorsAllowMethods
	} else {
		m = c.AllowMethods
	}
	return strings.Join(m, ",")
}
func (c Config) ReqHeader() string {
	var m []string
	if len(c.AllowHeaders) == 0 {
		m = DefCorsAllowHeader
	} else {
		m = c.AllowHeaders
	}
	return strings.Join(m, ",")
}
func (c Config) RespHeader() string {
	var m []string
	if len(c.ExposeHeaders) == 0 {
		m = DefCorsExposeHeader
	} else {
		m = c.ExposeHeaders
	}
	return strings.Join(m, ",")
}
func (c Config) Credential() string {
	if c.AllowCredentials {
		return "true"
	}
	return "false"
}
