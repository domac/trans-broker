package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime"
)

var version string

func SetVersion(ver string) {
	version = ver
}

func GetVersion() string {
	return version
}

func ShowVersion() {
	fmt.Printf(`Trans-Broker %s, Compiler: %s %s`,
		version,
		runtime.Compiler,
		runtime.Version())
	fmt.Println()
}

func VersionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-TB-VERSION", version)
		c.Next()
	}
}
