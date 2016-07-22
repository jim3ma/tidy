package utilities

import (
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// DebugHandler is a Gin MinddleWare for inserting debug info in tidy project
func DebugHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if w, ok := c.Writer.(gin.ResponseWriter); ok {
			hostname, err := os.Hostname()
			pid := os.Getpid()
			//log.Debugf("hostname: %s", hostname)
			if err == nil {
				w.Header().Add("X-Server-Hostname", hostname)
				w.Header().Add("X-Server-PID", strconv.Itoa(pid))
			}
		} else {
			log.Debug("Can not add response header: %s")
		}
		c.Writer.Header().Add("X-Powered-By", "Tidy")
		return
	}
}
