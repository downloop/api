package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	e := Error{
		Code:    code,
		Message: http.StatusText(code),
	}
	c.Logger().Error(err)
	c.JSON(code, e)
}
