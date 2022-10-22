package v1

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (dc DownloopContext) GetSessions(c echo.Context) error {
	rows, err := dc.Database.Queryx("SELECT * FROM sessions;")
	if err != nil {
		return err
	}

	var sessions []Session
	for rows.Next() {
		var session Session
		if err := rows.StructScan(&session); err != nil {
			return err
		}
		sessions = append(sessions, session)
	}

	return c.JSON(200, sessions)
}

func (dc DownloopContext) PostSessions(c echo.Context) error {
	var session Session
	if err := c.Bind(&session); err != nil {
		return err
	}

	sqlStatement := "INSERT INTO sessions (start_time, end_time) VALUES ($1, $2) RETURNING id;"
	err := dc.Database.QueryRowx(sqlStatement, session.StartTime, session.EndTime).Scan(&session.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, session)
}

func (dc DownloopContext) GetSessionId(c echo.Context, id uuid.UUID) error {
	var session Session
	if err := dc.Database.QueryRowx("SELECT * FROM sessions WHERE id = $1", id).StructScan(&session); err != nil {
		if err == sql.ErrNoRows {
			return &echo.HTTPError{Code: 404}
		}
		return err
	}
	return c.JSON(200, session)
}
