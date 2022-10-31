package v1

import "github.com/labstack/echo/v4"

type whereArgs struct {
	query string
	args  []interface{}
}

func (dc DownloopContext) getAll(c echo.Context, model interface{}, where ...whereArgs) error {
	limit := 10
	offset := 0
	err := echo.QueryParamsBinder(c).
		Int("limit", &limit).
		Int("offset", &offset).
		BindError()
	if err != nil {
		return err
	}

	res := dc.Database.Limit(limit).Offset(offset).Where(where).Find(model)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
