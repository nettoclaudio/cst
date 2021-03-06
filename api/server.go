package api

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// WebServer defines actions about an usual web server. Despite this name, it
// only regulates to start and stop methods.
type WebServer interface {
	Start() error
	Shutdown() error
}

// SecureWebServer holds the settings used to run a web server.
type SecureWebServer struct {
	CertFile string
	KeyFile  string
	Port     int
	UseTLS   bool

	echo *echo.Echo
}

// Start runs the web server using TLS and HTTP/2 protocols. Error is returned
// when it can't starts web server correctly.
func (ws *SecureWebServer) Start() error {

	ws.echo = echo.New()

	ws.echo.HideBanner = true

	ws.echo.Use(middleware.Recover())
	ws.echo.Use(middleware.Logger())

	ws.echo.GET("/health", health)

	v1 := ws.echo.Group("/v1")
	v1.POST("/scan", createScan)
	v1.GET("/scan/:image", showScans)

	address := fmt.Sprintf(":%d", ws.Port)

	if ws.UseTLS {
		return ws.echo.StartTLS(address, ws.CertFile, ws.KeyFile)
	}
	return ws.echo.Start(address)
}

// Shutdown stops web server the gracefully. Error is returned when it can't
// stops web server correctly.
func (ws *SecureWebServer) Shutdown() error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return ws.echo.Shutdown(ctx)
}
