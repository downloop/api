package v1

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/labstack/echo/v4"
)

func (dc DownloopContext) GetUsers(c echo.Context, params GetUsersParams) error {
	limit := 10
	offset := 0
	err := echo.QueryParamsBinder(c).
		Int("limit", &limit).
		Int("offset", &offset).
		BindError()
	if err != nil {
		return err
	}

	var users []UserModel
	res := dc.Database.Limit(limit).Offset(offset).Find(&users)
	if res.Error != nil {
		return res.Error
	}

	resp := UserResponseList{
		Data: []User{},
	}
	for _, user := range users {
		resp.Data = append(resp.Data,  User{
			Username: user.Username,
			Id:       user.ID,
		})
	}

	return c.JSON(200, resp)
}

func (dc DownloopContext) PostUsers(c echo.Context) error {
	var user UserPost
	if err := c.Bind(&user); err != nil {
		return err
	}

	model := UserModel{
		Username: user.Username,
	}
	res := dc.Database.Create(&model)
	if res.Error != nil {
		if pgError := res.Error.(*pgconn.PgError); errors.Is(res.Error, pgError) {
			var httpStatus int
			switch pgError.Code {
			case "23505":
				httpStatus = 409
			}

			return &echo.HTTPError{
				Code: httpStatus,
			}
		}
		return res.Error
	}

	resp := UserResponse{
		Data: User{
			Id:       model.ID,
			Username: model.Username,
		},
	}
	return c.JSON(201, resp)
}
