package main

import (
	"fmt"
	"log"
	"os"

	v1 "github.com/downloop/api/pkg/api/v1"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/craigtracey/jwksmiddleware"
	oapimiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"
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

type DownloopClaims struct {
	UUID uuid.UUID `json:"https://downloop.io/uuid"`
	jwt.StandardClaims
}

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
				Value: true,
			},
			&cli.BoolFlag{
				Name:  "validate",
				Value: true,
			},
		},
		Action: func(c *cli.Context) error {
			enforceRBAC := c.Bool("rbac")
			wipe := c.Bool("wipe")
			validate := c.Bool("validate")

			db, err := initDatabase(wipe)
			if err != nil {
				panic(err)
			}

			downloopContext := v1.DownloopContext{
				Database: db,
			}
			e := echo.New()
			e.HTTPErrorHandler = v1.HTTPErrorHandler
			e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
				Format: "method=${method}, uri=${uri}, status=${status}\n",
			}))
			if validate {
				swagger, err := v1.GetSwagger()
				if err != nil {
					panic(err)
				}
				e.Use(oapimiddleware.OapiRequestValidator(swagger))
			}

			// configure JWT middleware
			e.Use(jwksmiddleware.JWTWithConfig(jwksmiddleware.JWTConfig{
				JWTConfig: middleware.JWTConfig{
					Skipper: func(c echo.Context) bool {
						return !enforceRBAC
					},
					BeforeFunc: func(c echo.Context) {
						//body, _ := ioutil.ReadAll(c.Request().Body)
						//fmt.Printf("Got request body: %+v", string(body))
						fmt.Printf("Got Request Headers: %+v", c.Request().Header["Authorization"])
					},
					SigningMethod: "RS256",
					Claims:        &DownloopClaims{},
					SuccessHandler: func(c echo.Context) {
						user := c.Get("user").(*jwt.Token)
						claims := user.Claims.(*DownloopClaims)
						c.Set("uuid", claims.UUID)
					},
				},
				JWKSURL: "https://downloop.us.auth0.com/.well-known/jwks.json",
			}))

			v1.RegisterHandlers(e, downloopContext)
			e.Logger.Fatal(e.Start("0.0.0.0:8080"))
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func initDatabase(wipe bool) (*gorm.DB, error) {
	user := os.Getenv("PG_USERNAME")
	password := os.Getenv("PG_PASSWORD")

	conn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=require", user, password, host, port, dbname)
	fmt.Println(conn)

	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if wipe {
		db.Migrator().DropTable(&v1.UserModel{}, &v1.SessionModel{})
	}

	db.AutoMigrate(v1.UserModel{}, v1.SessionModel{})
	return db, nil
}
