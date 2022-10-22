package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"

	oapimiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"
	v1 "github.com/downloop/api/pkg/api/v1"
	"github.com/labstack/echo/v4/middleware"
	echo "github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

const (
	host   = "downloop-downloop"
	port   = 5432
	dbname = "downloop"
)

const schema = `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
DROP TABLE sessions;
CREATE TABLE IF NOT EXISTS users(id UUID NOT NULL DEFAULT uuid_generate_v1(),
				 username VARCHAR(32));
CREATE TABLE IF NOT EXISTS sessions(id UUID NOT NULL DEFAULT uuid_generate_v1(),
                                    user_id UUID, 
				    start_time TIMESTAMP, 
			            end_time TIMESTAMP,
				    CONSTRAINT session_pkey PRIMARY KEY(id));
CREATE TABLE IF NOT EXISTS coordinates(session_id UUID, latlon POINT, time TIMESTAMP);
CREATE TABLE IF NOT EXISTS altitude(session_id UUID, alt FLOAT, time TIMESTAMP); `

func main() {
	db, err := initDatabase()
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

func initDatabase() (*sqlx.DB, error) {
	user := os.Getenv("PG_USERNAME")
	password := os.Getenv("PG_PASSWORD")

	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)

	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// create schemas if they don't exist
	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)
	c.JSON(code, err)
}
