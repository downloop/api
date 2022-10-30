package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	sqladapter "github.com/Blank-Xu/sql-adapter"
	"github.com/casbin/casbin/v2"
	oapimiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"
	v1 "github.com/downloop/api/pkg/api/v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	casbinmiddleware "github.com/labstack/echo-contrib/casbin"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/urfave/cli/v2"
)

const (
	host   = "downloop-downloop"
	port   = 5432
	dbname = "downloop"
)

func main() {

	app := &cli.App{
		Name: "api",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "wipe",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "rbac",
				Value: false,
			},
		},
		Action: func(c *cli.Context) error {
			enforceRBAC := c.Bool("rbac")
			db, err := initDatabase()
			if err != nil {
				panic(err)
			}

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

			// configure authorization
			database, err := db.DB()
			if err != nil {
				return err
			}
			adapter, err := sqladapter.NewAdapter(database, "postrges", "casbin_rule")
			if err != nil {
				return err
			}

			enforcer, err := casbin.NewEnforcer("/etc/downloop/rbac_model.conf", adapter)
			if err != nil {
				return err
			}
			enforcer.AddRoleForUser("craig", "admin")
			enforcer.AddPolicy("admin", "/sessions", "GET")

			e.Use(casbinmiddleware.MiddlewareWithConfig(casbinmiddleware.Config{
				Skipper: func(c echo.Context) bool {
					return !enforceRBAC
				},
				Enforcer: enforcer,
				UserGetter: func(c echo.Context) (string, error) {
					return "craig", nil
				},
			}))

			e.Use(oapimiddleware.OapiRequestValidator(swagger))
			v1.RegisterHandlers(e, downloopContext)
			e.Logger.Fatal(e.Start("0.0.0.0:8080"))
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func initDatabase() (*gorm.DB, error) {
	user := os.Getenv("PG_USERNAME")
	password := os.Getenv("PG_PASSWORD")

	conn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=require", user, password, host, port, dbname)
	fmt.Println(conn)

	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(v1.UserModel{}, v1.SessionModel{})
	return db, nil
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	e := v1.Error{
		Code:    code,
		Message: http.StatusText(code),
	}
	c.Logger().Error(err)
	c.JSON(code, e)
}
