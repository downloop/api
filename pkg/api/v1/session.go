package v1

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (dc DownloopContext) GetSessions(c echo.Context) error {

	var sessions []SessionModel
	err := dc.getAll(c, &sessions)
	if err != nil {
		return err
	}

	resp := SessionResponseList{
		Data: []Session{},
	}

	for _, session := range sessions {
		resp.Data = append(resp.Data, Session{
			Id:        session.ID,
			StartTime: session.StartTime,
			EndTime:   &session.EndTime,
		})
	}

	return c.JSON(200, resp)
}

func (dc DownloopContext) PostSessions(c echo.Context) error {

	var session Session
	if err := c.Bind(&session); err != nil {
		return err
	}

	end := time.Time{}
	if session.EndTime != nil {
		end = *session.EndTime
	}

	u, err := uuid.Parse("b56dd059-3200-45eb-8627-9d1480ba834b")
	if err != nil {
		fmt.Println(err)
	}
	model := SessionModel{
		UserID: u,
		StartTime: session.StartTime,
		EndTime:   end,
	}

	fmt.Printf("MODEL %+v\n ", model)

	tx := dc.Database.Create(&model)
	if tx.Error != nil {
		return tx.Error
	}

	return c.JSON(201, nil)
}

func (dc DownloopContext) GetSessionId(c echo.Context, id uuid.UUID) error {
	/*	var session Session
		if err := dc.Database.QueryRowx("SELECT * FROM sessions WHERE id = $1", id).StructScan(&session); err != nil {
			if err == sql.ErrNoRows {
				return &echo.HTTPError{Code: 404}
			}
			return err
		}
	*/
	return c.JSON(200, nil)
}

func (dc DownloopContext) DeleteSessionId(c echo.Context, id uuid.UUID) error {
	/*_, err := dc.Database.Exec("DELETE FROM sessions WHERE id = $1;", id)
	if err != nil {
		return err
	}*/
	return c.JSON(204, nil)
}
