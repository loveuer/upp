package interfaces

import (
	"context"
	"mime/multipart"
	"net"
)

type ApiGroup interface {
	Group(path string, handlers ...ApiHandler) ApiGroup
	GET(path string, handlers ...ApiHandler)
	POST(path string, handlers ...ApiHandler)
	PUT(path string, handlers ...ApiHandler)
	DELETE(path string, handlers ...ApiHandler)
	HEAD(path string, handlers ...ApiHandler)
	PATCH(path string, handlers ...ApiHandler)
	OPTIONS(path string, handlers ...ApiHandler)
	Handle(method, path string, handlers ...ApiHandler)
	Use(handlers ...ApiHandler)
}

type ApiEngine interface {
	ApiGroup
	Run(address string) error
	RunListener(ln net.Listener) error
}

type ApiContext interface {
	App() Upp
	// parse body, form, json
	BodyParser(out any) error
	Context() context.Context
	Cookie(string) string
	FormFile(string) (*multipart.FileHeader, error)
	FormValue(string) string
	GetHeader(string) string
	SetHeader(string, string)
	IP() string
	Json(any) error
	Locals(interface{}, ...any) any

	// get method or rewrite method
	Method(string, ...string) string
	MultipartForm() (*multipart.Form, error)
	Param(string) string

	// get path or rewrite path
	Path(string, ...string) string
	Next() error
	Query(string) string
	QueryParse(any) error
	Redirect(string, ...int) error

	// render html response
	Render(name, layout string, data any) error

	// set response status
	Status(int)
	SendStatus(int)

	Write([]byte) (int, error)
}

type ApiHandler func(c ApiContext) error
