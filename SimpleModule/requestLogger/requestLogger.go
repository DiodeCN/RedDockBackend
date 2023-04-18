package requestlogger

import (
	"log"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在处理请求之前，打印请求信息
		log.Println("Request Method:", c.Request.Method)
		log.Println("Request URL:", c.Request.URL.String())
		log.Println("Request Headers:", c.Request.Header)

		// 处理请求
		c.Next()

		// 在处理请求之后，打印响应信息
		log.Println("Response Status Code:", c.Writer.Status())
	}
}
