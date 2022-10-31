package v1

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (dc DownloopContext) GetSessions(c echo.Context) error {
	var sessions []SessionModel

	//u := c.Get("user-uuid").(uuid.UUID)
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

	u := c.Get("user-uuid").(uuid.UUID)
	model := SessionModel{
		UserID:    u,
		StartTime: session.StartTime,
		EndTime:   end,
	}

	tx := dc.Database.Create(&model)
	if tx.Error != nil {
		return tx.Error
	}

	return c.JSON(201, nil)
}

func (dc DownloopContext) GetSessionId(c echo.Context, id uuid.UUID) error {
	var model SessionModel
	tx := dc.Database.Where("id = ?", id).Find(&model)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	resp := SessionResponse{
		Data: Session{
			Id:        model.ID,
			StartTime: model.StartTime,
		},
	}
	return c.JSON(200, resp)
}

func (dc DownloopContext) DeleteSessionId(c echo.Context, id uuid.UUID) error {
	tx := dc.Database.Delete(SessionModel{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return c.NoContent(204)
}
