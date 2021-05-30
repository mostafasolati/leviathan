package services

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mostafasolati/leviathan/contracts"
	"github.com/mostafasolati/leviathan/models"
	"github.com/mostafasolati/leviathan/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// NewEchoServerContainer  is responsible for registering routes and run the server
func NewEchoServerContainer(
	config contracts.IConfigService,
	logger contracts.ILogger,
) contracts.IServerContainer {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch},
	}))
	e.Use(appDetectionMiddleware())

	return &serverContainer{
		groups:        make(map[string]*echo.Group),
		e:             e,
		configService: config,
		logger:        logger,
	}
}

// SecureRoutes restrict access to certain routes by allowing access only to specified roles
func (s *serverContainer) SecureRoutes(routes map[string][]string) {

	config := middleware.JWTConfig{
		Claims:     &models.UserClaims{},
		SigningKey: []byte(s.configService.String("auth.jwt_secret")),
		Skipper: func(c echo.Context) bool {
			// Skip requests with no token.
			return c.Request().Header.Get("Authorization") == ""
		},
	}

	for route, roles := range routes {
		g := s.e.Group(route, middleware.JWTWithConfig(config))

		g.Use(accessMiddleware(roles))
		s.groups[route] = g
	}
}

// Route s a method to define a route for an API endpoint
func (s *serverContainer) Route(method, path string, handler contracts.Handler) {

	h := func(c echo.Context) error {
		server := NewServer(c, s.configService)
		err := handler(server)

		if err != nil {
			c.Logger().Error(err)
			return server.JSON(http.StatusBadRequest, &models.Error{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
		}
		return nil
	}

	g := s.e.Group("")

	for sr, group := range s.groups {
		if strings.Contains(path, sr) {
			g = group
			path = strings.Replace(path, sr, "", 1)
			break
		}
	}

	switch method {
	case http.MethodPost:
		g.POST(path, h)
	case http.MethodDelete:
		g.DELETE(path, h)
	case http.MethodPut:
		g.PUT(path, h)
	case http.MethodPatch:
		g.PATCH(path, h)
	default:
		g.GET(path, h)
	}

}

// Run starts the server in given address
func (s *serverContainer) Run(address string) {
	s.e.Start(address)
}

// TestServer starts a test server
func (s *serverContainer) TestServer() *httptest.Server {
	return httptest.NewServer(s.e.Server.Handler)
}

func accessMiddleware(roles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if c.Get("user") == nil {
				return c.JSON(401, &models.Error{
					Message: "unauthorized",
					Code:    http.StatusUnauthorized,
				})
			}

			user := c.Get("user").(*jwt.Token).Claims.(*models.UserClaims)

			for _, role := range roles {
				for _, r := range user.Roles {
					if strings.EqualFold(role, r) {
						return next(c)
					}
				}
			}

			return c.JSON(http.StatusForbidden, &models.Error{
				Message: "access forbidden",
				Code:    http.StatusForbidden,
			})

		}
	}
}

// Determine client app using the `User-Agent` header.
// Android app is known to send `okhttp/X.Y.Z`.
func appDetectionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userAgent := c.Request().Header.Get("User-Agent")
			if strings.Contains(userAgent, "okhttp") {
				c.Set("app", models.AndroidApp)
			} else if utils.IsMobileWeb(userAgent) {
				c.Set("app", models.MobileWebApp)
			} else {
				c.Set("app", models.WebApp)
			}
			return next(c)
		}
	}
}

type serverContainer struct {
	configService contracts.IConfigService
	logger        contracts.ILogger
	groups        map[string]*echo.Group
	e             *echo.Echo
}

// NewServer returns a new instance of IServer
// server is a wrapper for golang http servers
// in this case echo server
func NewServer(c echo.Context, configService contracts.IConfigService) contracts.IServer {
	return &echoServer{
		c:             c,
		configService: configService,
	}
}

// SetHeader sets a http header
func (s *echoServer) SetHeader(key, value string) {
	s.c.Response().Header().Set(key, value)
}

// Response return a proper http response to the client
func (s *echoServer) Response(status int, v interface{}) error {
	return s.JSON(status, v)
}

// Query returns a query string parameter
func (s *echoServer) Query(key string) string {
	return s.c.QueryParam(key)
}

// QueryParams returns all query string parameters
func (s *echoServer) QueryParams() map[string]string {
	params := make(map[string]string)
	for key, val := range s.c.QueryParams() {
		params[key] = val[0]
	}
	return params
}

// FormParams returns all post string parameters
func (s *echoServer) FormParams() map[string]string {
	params := make(map[string]string)
	form, err := s.c.FormParams()
	if err != nil {
		return params
	}
	for key, val := range form {
		params[key] = val[0]
	}
	return params
}

// UploadedFile returns an uploaded file.
func (s *echoServer) UploadedFile(name string) (*multipart.FileHeader, error) {
	return s.c.FormFile(name)
}

// Upload handles http file upload into the server
func (s *echoServer) Upload(field string) (string, error) {
	file, err := s.c.FormFile(field)
	if err != nil {
		return "", err
	}

	imagesRoot := filepath.Join(s.configService.StorageDir(), "files")
	return utils.UploadFile(imagesRoot, file)
}

// User returns the logged in user
func (s *echoServer) User() *models.UserClaims {
	if s.c.Get("user") == nil {
		return nil
	}
	return s.c.Get("user").(*jwt.Token).Claims.(*models.UserClaims)
}

// App returns the client app sending the HTTP request
func (s *echoServer) App() models.App {
	value := s.c.Get("app")
	if app, ok := value.(models.App); ok {
		return app
	}
	return models.WebApp
}

// Param returns a url parameter
func (s *echoServer) Param(key string) string {
	return s.c.Param(key)
}

// Bind binds a http request (posted data or query params) to a struct
func (s *echoServer) Bind(in interface{}) error {
	s.c.Response()
	return s.c.Bind(in)
}

// RawBody reads the raw body of the request.
func (s *echoServer) RawBody() ([]byte, error) {
	return ioutil.ReadAll(s.c.Request().Body)
}

// String returns a string response to the client
func (s *echoServer) String(status int, output string) error {
	return s.c.String(status, output)
}

// JSON returns a json response to the client
func (s *echoServer) JSON(status int, output interface{}) error {
	return s.c.JSON(status, output)
}

// File send a file to client for download
func (s *echoServer) File(path string) error {
	return s.c.File(path)
}

// Attachment sends a file as attachment prompting for download
func (s *echoServer) Attachment(path, name string) error {
	return s.c.Attachment(path, name)
}

// HTML serves a Go HTML template
func (s *echoServer) HTML(name string, data interface{}) error {
	tmplPath := filepath.Join(
		s.configService.StaticDir(),
		fmt.Sprintf("templates/%s.gohtml", name),
	)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, data); err != nil {
		return err
	}

	return s.c.HTMLBlob(http.StatusOK, b.Bytes())
}

// Redirect redirects the client to a url
func (s *echoServer) Redirect(status int, url string) error {
	return s.c.Redirect(status, url)
}

func (s *echoServer) Request() *http.Request {
	return s.c.Request()
}

func (s *echoServer) ResponseWriter() http.ResponseWriter {
	return s.c.Response()
}

type echoServer struct {
	c             echo.Context
	configService contracts.IConfigService
}
