package v1

import "github.com/labstack/echo/v4"

func (dc DownloopContext) getAll(c echo.Context, model interface{}, uid string) error {
	limit := 10
	offset := 0
	err := echo.QueryParamsBinder(c).
		Int("limit", &limit).
		Int("offset", &offset).
		BindError()
	if err != nil {
		return err
	}

	res := dc.Database.Limit(limit).Offset(offset).Where("user_id = ?", uid).Find(model)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
