package api

import "sync"

const (
	_404     = `404 Not Found`
	_405     = `405 Method Not Allowed`
	_500     = `500 Internal Server Error`
	TraceKey = "X-Trace-Id"
)

type Map map[string]interface{}

type Config struct {
	DisableMessagePrint bool `json:"-"`
	// Default: 4 * 1024 * 1024
	BodyLimit int64 `json:"-"`

	// if report http.ErrServerClosed as run err
	ErrServeClose bool `json:"-"`

	DisableLogger       bool `json:"-"`
	DisableRecover      bool `json:"-"`
	DisableHttpErrorLog bool `json:"-"`

	// EnableNotImplementHandler bool        `json:"-"`
	NotFoundHandler         HandlerFunc `json:"-"`
	MethodNotAllowedHandler HandlerFunc `json:"-"`
}

var defaultConfig = &Config{
	BodyLimit: 4 * 1024 * 1024,
	NotFoundHandler: func(c *Ctx) error {
		c.Set("Content-Type", MIMETextPlain)
		_, err := c.Status(404).Write([]byte(_404))
		return err
	},
	MethodNotAllowedHandler: func(c *Ctx) error {
		c.Set("Content-Type", MIMETextPlain)
		_, err := c.Status(405).Write([]byte(_405))
		return err
	},
}

func New(config ...Config) *App {
	app := &App{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},

		pool: &sync.Pool{},

		redirectTrailingSlash:  true,  // true
		redirectFixedPath:      false, // false
		handleMethodNotAllowed: true,  // false
		useRawPath:             false, // false
		unescapePathValues:     true,  // true
		removeExtraSlash:       false, // false
	}

	if len(config) > 0 {
		app.config = &config[0]

		if app.config.BodyLimit == 0 {
			app.config.BodyLimit = defaultConfig.BodyLimit
		}

		if app.config.NotFoundHandler == nil {
			app.config.NotFoundHandler = defaultConfig.NotFoundHandler
		}

		if app.config.MethodNotAllowedHandler == nil {
			app.config.MethodNotAllowedHandler = defaultConfig.MethodNotAllowedHandler
		}

	} else {
		app.config = defaultConfig
	}

	app.RouterGroup.app = app

	if !app.config.DisableLogger {
		app.Use(NewLogger())
	}

	if !app.config.DisableRecover {
		app.Use(NewRecover(true))
	}

	app.pool.New = func() any {
		return app.allocateContext()
	}

	return app
}
