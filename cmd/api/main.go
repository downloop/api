package main

import (
	"net/http"

	oapimiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"
	v1 "github.com/downloop/api/pkg/api/v1"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	db, err := initDatabase(false)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	downloopContext := v1.DownloopContext{
		Database: db,
	}
	e := echo.New()
	e.HTTPErrorHandler = customHTTPErrorHandler
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	swagger, err := v1.GetSwagger()
	if err != nil {
		panic(err)
	}

	e.Use(oapimiddleware.OapiRequestValidator(swagger))
	v1.RegisterHandlers(e, downloopContext)
	e.Logger.Fatal(e.Start("0.0.0.0:8080"))
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)
	c.JSON(code, err)
}
