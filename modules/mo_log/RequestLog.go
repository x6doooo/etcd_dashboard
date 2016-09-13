package mo_log

import (
    "github.com/gin-gonic/gin"
    //"time"
    "time"
)

func RequestLog() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        c.Next()

        end := time.Now()
        latency := end.Sub(start)
        Logger.Info("[%d] %s %s", c.Writer.Status(), c.Request.RequestURI, latency.String())
    }
}


