package v1

import (
	"database/sql"
	"database/sql/driver"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Any struct{}

func (a Any) Match(v driver.Value) bool {
	return true
}

type AnyUUID struct{}

func (a AnyUUID) Match(v driver.Value) bool {
	_, err := uuid.Parse(v.(string))
	return err == nil
}

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
	e    *echo.Echo
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	if err != nil {
		panic(err)
	}

	s.DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	s.e = echo.New()
	s.e.HTTPErrorHandler = HTTPErrorHandler
}

func (s *Suite) TestPostUsers() {

	userJSON := `
	{
		"username": "foo"
	}`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := s.e.NewContext(req, rec)
	h := &DownloopContext{s.DB}

	mockedRow := sqlmock.NewRows([]string{"id"})

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_models" ("created_at","updated_at","deleted_at","username","id") VALUES ($1,$2,$3,$4,$5)`)).
		WithArgs(Any{}, Any{}, Any{}, "foo", AnyUUID{}).
		WillReturnRows(mockedRow)
	s.mock.ExpectCommit()

	if assert.NoError(s.T(), h.PostUsers(c)) {
		assert.Equal(s.T(), http.StatusCreated, rec.Code)
	}
}

func (s *Suite) TestPostUsersDuplicate() {

	userJSON := `
	{
		"username": "foo"
	}`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := s.e.NewContext(req, rec)
	h := &DownloopContext{s.DB}

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_models" ("created_at","updated_at","deleted_at","username","id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(Any{}, Any{}, Any{}, "foo", AnyUUID{}).
		WillReturnError(&pgconn.PgError{Code: "23505"})
	s.mock.ExpectRollback()

	err := h.PostUsers(c)
	if assert.Error(s.T(), err) {
		he, ok := err.(*echo.HTTPError)
		if ok {
			assert.Equal(s.T(), http.StatusConflict, he.Code)
		}
	}
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}
