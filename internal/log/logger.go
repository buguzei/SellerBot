package log

// Fields is an additional info for logs
type Fields map[string]interface{}

type Logger interface {
	Named(string) Logger

	Debug(message string, args Fields)

	Info(message string, args Fields)

	Warn(message string, args Fields)

	Error(message string, args Fields)

	Fatal(message string, args Fields)
}
