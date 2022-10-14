package v1

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func (DownloopContext) GetSessions(c echo.Context) error {
	return nil
}

func (DownloopContext) PostSessions(c echo.Context) error {
	var session Session
	if err := c.Bind(&session); err != nil {
		return err
	}

	fmt.Printf("Session %+v", session)

	return nil
}
