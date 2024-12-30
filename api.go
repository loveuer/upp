package upp

import (
	"net/http"

	"github.com/loveuer/upp/pkg/api"
)

func (u *upp) API() *api.App { return u.api.engine }

func (u *upp) GET(path string, handlers ...api.HandlerFunc) {
	u.HandleAPI(http.MethodGet, path, handlers...)
}

func (u *upp) POST(path string, handlers ...api.HandlerFunc) {
	u.HandleAPI(http.MethodPost, path, handlers...)
}

func (u *upp) PUT(path string, handlers ...api.HandlerFunc) {
	u.HandleAPI(http.MethodPut, path, handlers...)
}

func (u *upp) DELETE(path string, handlers ...api.HandlerFunc) {
	u.HandleAPI(http.MethodDelete, path, handlers...)
}

func (u *upp) PATCH(path string, handlers ...api.HandlerFunc) {
	u.HandleAPI(http.MethodPatch, path, handlers...)
}

func (u *upp) HEAD(path string, handlers ...api.HandlerFunc) {
	u.HandleAPI(http.MethodHead, path, handlers...)
}

func (u *upp) OPTIONS(path string, handlers ...api.HandlerFunc) {
	u.HandleAPI(http.MethodOptions, path, handlers...)
}

func (u *upp) HandleAPI(method, path string, handlers ...api.HandlerFunc) {
	u.api.engine.Handle(method, path, handlers...)
}
