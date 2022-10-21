package handler

import (
	mw "github.com/Astemirdum/logs/internal/handler/middleware"
	"github.com/Astemirdum/logs/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

type logHandler struct {
	svc Service
	log *zap.Logger
}

func NewHandler(srv *service.Service, log *zap.Logger) *logHandler {
	return &logHandler{
		svc: srv,
		log: log,
	}
}

func (h *logHandler) newRouter() *echo.Echo {
	e := echo.New()

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 4 << 10, // 4 KB
		LogLevel:  log.ERROR,
	}))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("2M"))
	e.Validator = &CustomValidator{validator: validator.New()}
	st := mw.NewStats()
	api := e.Group("/api/v1", st.Process)
	{
		api.POST("/logs", h.CreateLog)
		api.GET("/logs/:id", h.GetLog)
		api.GET("/logs", h.ListLogs)
	}
	e.GET("/stats", st.GetStats)

	return e
}
