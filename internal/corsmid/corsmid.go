package corsmid

import (
	"github.com/gin-gonic/gin"
	"github.com/obnahsgnaw/assetweb/service/cors"
	"net/http"
)

func NewCorsMid(cf *cors.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if cf.Origin() != "" {
			c.Header("Access-Control-Allow-Origin", cf.Origin()) // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", cf.Methods())
			c.Header("Access-Control-Allow-Headers", cf.ReqHeader())
			c.Header("Access-Control-Expose-Headers", cf.RespHeader())
			c.Header("Access-Control-Allow-Credentials", cf.Credential())
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
