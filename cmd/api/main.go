package main

import (
	v1 "github.com/downloop/api/pkg/api/v1"
	echo "github.com/labstack/echo/v4"
)

func main() {
	var downloopContext v1.DownloopContext
	e := echo.New()
	v1.RegisterHandlers(e, downloopContext)
	e.Logger.Fatal(e.Start("0.0.0.0:8080"))
}
