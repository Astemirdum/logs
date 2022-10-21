package handler

import (
	"errors"
	errr "github.com/Astemirdum/logs/internal/errs"
	"github.com/Astemirdum/logs/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (h *logHandler) CreateLog(c echo.Context) error {
	var raw models.CreateLogRequest
	if err := c.Bind(&raw); err != nil {
		return err
	}

	if err := c.Validate(raw); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	id, err := h.svc.CreateLog(ctx, raw.Raw)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, models.WriteLogResponse{ID: id})
}

func (h *logHandler) GetLog(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	log, err := h.svc.GetLog(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, errr.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, models.ReadLogResponse(log))
}

func (h *logHandler) ListLogs(c echo.Context) error {
	logs, err := h.svc.ListLogs(c.Request().Context())
	if err != nil {
		if errors.Is(err, errr.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, logsToLogsResponse(logs))
}

func logsToLogsResponse(logs []models.Log) []models.ReadLogResponse {
	resp := make([]models.ReadLogResponse, 0, len(logs))
	for _, log := range logs {
		resp = append(resp, models.ReadLogResponse(log))
	}
	return resp
}
