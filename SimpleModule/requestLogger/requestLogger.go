package requestlogger

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	// 创建或打开日志文件
	logFile, err := os.OpenFile("httplog.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return func(c *gin.Context) {
		// 在处理请求之前，记录请求信息到日志文件
		fmt.Fprintf(logFile, "Request Method: %s\n", c.Request.Method)
		fmt.Fprintf(logFile, "Request URL: %s\n", c.Request.URL.String())
		fmt.Fprintf(logFile, "Request Headers: %v\n", c.Request.Header)

		// 处理请求
		c.Next()

		// 在处理请求之后，记录响应信息到日志文件
		fmt.Fprintf(logFile, "Response Status Code: %d\n\n", c.Writer.Status())
	}
}
