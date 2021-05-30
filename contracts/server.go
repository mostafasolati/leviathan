package contracts

import (
	"mime/multipart"
	"net/http"
	"net/http/httptest"

	"github.com/mostafasolati/leviathan/models"
)

// Handler is a function to handle http requests
type Handler func(server IServer) error

// IServer is an abstraction of underneath server implementation like echo, fast http, etc.
type IServer interface {
	Request() *http.Request
	ResponseWriter() http.ResponseWriter

	// Bind converts an http request data to a struct
	Bind(in interface{}) error

	// RawBody reads the raw body of the request.
	RawBody() ([]byte, error)

	// String returns the response to the client as a string
	String(status int, output string) error

	// JSON returns the response to the client as a JSON object
	JSON(status int, output interface{}) error

	// File serves a file as a response to the client
	File(path string) error

	// Attachment sends a file as attachment prompting for download
	Attachment(path, name string) error

	// HTML serves a Go HTML template
	HTML(name string, data interface{}) error

	// Redirect redirects the client to a url
	Redirect(status int, url string) error

	// Upload uploads a file from http request
	Upload(field string) (string, error)

	// User returns logged user if present
	User() *models.UserClaims

	// App returns the client app sending the HTTP request
	App() models.App

	// Param get the url parameter from the url
	Param(key string) string

	// Query returns the query string value
	Query(key string) string

	// QueryParams returns all query parameters
	QueryParams() map[string]string

	// FormParams returns all post string parameters
	FormParams() map[string]string

	// UploadedFile returns an uploaded file.
	UploadedFile(name string) (*multipart.FileHeader, error)

	// Response send the final response to the client with a status code
	Response(status int, v interface{}) error

	// SetHeader sets http headers for response
	SetHeader(key, value string)
}

// IServerContainer is responsible to fire the http server and register the handlers
type IServerContainer interface {

	// Route registers a route to a corresponding handler
	Route(method, path string, handler Handler)

	SecureRoutes(routes map[string][]string)

	// Run starts http server
	Run(address string)

	// TestServer returns an HTTP server for testing purposes.
	TestServer() *httptest.Server
}
