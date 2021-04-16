package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"

	"github.com/njmdk/common/network/reuseport"
)

func main() {
	e := echo.New()
	e.Any("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"ip": ctx.RealIP(),
		})
	})
	listener, err := reuseport.Listen("tcp4", ":54321")
	if err != nil {
		panic(err)
	}
	e.Listener = listener
	go func() {
		if err := e.Start(""); err != nil {
			log.Fatal("Server Shutdown", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown", err)
	}
}
