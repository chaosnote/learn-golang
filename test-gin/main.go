package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"

	"idv/chris/utils"
)

var logger = GenLogger()

func GenLogger() utils.Log {
	// 外部設置
	dir := "./log"
	lv := 0
	debug := true

	var logger utils.Log
	if debug {
		logger = utils.GenComplex(dir, lv, debug)
	} else {
		logger = utils.GenFile(dir, lv)
	}
	return logger
}

func main() {

	g := gin.New()

	redis_store, _ := redis.NewStore(10, "tcp", "ip:6379", "", "", []byte("secret"))
	g.Use(sessions.Sessions("custom_session", redis_store))

	g.Use(MiddlewareLogger())
	g.Use(ginzap.RecoveryWithZap(logger.Logger(), true)) // Recovery error

	g.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server v0.0.0")
	})

	g.GET("/session", func(c *gin.Context) {
		session := sessions.Default(c)
		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()
		c.JSON(200, gin.H{"count": count})
	})

	s := &http.Server{
		Addr:    ":8080",
		Handler: g,
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		logger.Logger().Debug("localhost:8080")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	q := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(q, syscall.SIGINT, syscall.SIGTERM)
	<-q
	logger.Logger().Debug("Shutting down server...")

}
