// Package v1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package v1

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// Session defines model for Session.
type Session struct {
	EndTime   *time.Time          `db:"end_time" json:"end_time,omitempty"`
	Id        *openapi_types.UUID `db:"id" json:"id,omitempty"`
	StartTime time.Time           `db:"start_time" json:"start_time"`
	UserId    *openapi_types.UUID `db:"user_id" json:"user_id,omitempty"`
}

// SessionList defines model for SessionList.
type SessionList = []Session

// PostSessionsJSONBody defines parameters for PostSessions.
type PostSessionsJSONBody = Session

// PostSessionsJSONRequestBody defines body for PostSessions for application/json ContentType.
type PostSessionsJSONRequestBody = PostSessionsJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// delete a session by id
	// (DELETE /session/{id})
	DeleteSessionId(ctx echo.Context, id openapi_types.UUID) error
	// set a session by id
	// (GET /session/{id})
	GetSessionId(ctx echo.Context, id openapi_types.UUID) error

	// (GET /sessions)
	GetSessions(ctx echo.Context) error

	// (POST /sessions)
	PostSessions(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// DeleteSessionId converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteSessionId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteSessionId(ctx, id)
	return err
}

// GetSessionId converts echo context to params.
func (w *ServerInterfaceWrapper) GetSessionId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetSessionId(ctx, id)
	return err
}

// GetSessions converts echo context to params.
func (w *ServerInterfaceWrapper) GetSessions(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetSessions(ctx)
	return err
}

// PostSessions converts echo context to params.
func (w *ServerInterfaceWrapper) PostSessions(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostSessions(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.DELETE(baseURL+"/session/:id", wrapper.DeleteSessionId)
	router.GET(baseURL+"/session/:id", wrapper.GetSessionId)
	router.GET(baseURL+"/sessions", wrapper.GetSessions)
	router.POST(baseURL+"/sessions", wrapper.PostSessions)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8yTT2sbMRDFv4qY9ihn3T+nPZZCCfRQ2mMIQV6NHYVdSdGM2hij715Gu15na/qHktD6",
	"YiGxP7157+kAXRhi8OiZoD0Adbc4mLr8gkQueFnGFCImdlgP0NsbdgPKehvSYBhasIZxVXc18D4itECc",
	"nN+BhodVMNGtumBxh36FD5zMis2u0uwG2hOyFA3OLsg5O/tXUGcrjtgkfmrBj6ByRyZMN0+l+wgrQk54",
	"n11CC+3V40uvZ3DY3GHHUPQxsI+OWFiOcajQlwm30MKL5hR1M+XcHEMuM8+kZPbj3c5vgwDYcS9HNnzz",
	"fQhRmehAw1dMY0Hg1cX6Yi2MENHLYQtv6paGaPi2qmhovKo5OFvqrNgj10ykXIZd8JcWWnhf9ydhl7Yy",
	"khmQMRG0VwdwcqVwQYM3EitUp09Wccqopy7/PpNSruVjisHT2PDX67ejQuqSizwOOeqdOpWHwaT9vKuM",
	"msZTm72S7DTskM+H+4D8jydby18XPKOv+kyMveuqwuaOxvd+4v9Rd8SRpVdHM2IgVpS7Dom2uf/BO0I+",
	"N67ouSpV8WTjkv8ZOSdPyqjeEauwVfJo1Pyd/qntBM/vSH2Bv3Clal64UjSIVed9+RRoqfw+I/G7YPfP",
	"E+OyaeV/a4/8vgcAAP//hjd2GbMGAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
