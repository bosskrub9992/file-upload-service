package handlers

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/bosskrub9992/file-upload-service/responses"
	"github.com/bosskrub9992/file-upload-service/services"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	serv *services.Service
}

func New(serv *services.Service) *Handler {
	return &Handler{
		serv: serv,
	}
}

func (h Handler) Upload(c echo.Context) error {
	ctx := c.Request().Context()

	file, err := c.FormFile("file")
	if err != nil {
		slog.Error(err.Error())
		resp := responses.ErrAPIFailed
		return c.JSON(resp.HTTPStatusCode, resp)
	}

	src, err := file.Open()
	if err != nil {
		slog.Error(err.Error())
		resp := responses.ErrAPIFailed
		return c.JSON(resp.HTTPStatusCode, resp)
	}
	defer src.Close()

	srcByte, err := io.ReadAll(src)
	if err != nil {
		slog.Error(err.Error())
		resp := responses.ErrAPIFailed
		return c.JSON(resp.HTTPStatusCode, resp)
	}

	if _, err := h.serv.Upload(ctx, srcByte, file.Filename); err != nil {
		resp := responses.ErrAPIFailed
		return c.JSON(resp.HTTPStatusCode, resp)
	}

	return c.JSON(http.StatusCreated, responses.PostSuccess)
}
