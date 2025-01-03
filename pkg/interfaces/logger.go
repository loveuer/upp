package interfaces

type Logger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
	Panic(string, ...any)
	Fatal(string, ...any)
}
