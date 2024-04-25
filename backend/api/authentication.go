package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type empty struct{}

func (Server) GetAuthenticate(ctx echo.Context) error {

	return ctx.JSON(http.StatusOK, empty{})
}
