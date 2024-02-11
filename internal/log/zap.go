package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Zap struct {
	l *zap.Logger
}

func NewZap() Zap {
	cfg := zap.NewProductionConfig()
	cfg.DisableStacktrace = true

	newLogger, _ := cfg.Build()
	defer newLogger.Sync()

	return Zap{l: newLogger}
}

func (z Zap) Named(name string) Logger {
	return Zap{
		l: z.l.Named(name),
	}
}

func (z Zap) Debug(message string, args Fields) {
	for key, value := range args {
		t := checkZapType(value)

		f := zap.Field{
			Key:       key,
			Type:      t,
			Interface: value,
		}

		z.l.Debug(message, f)
	}
}

func (z Zap) Warn(message string, args Fields) {
	for key, value := range args {
		f := zap.Field{
			Key:       key,
			Interface: value,
		}

		z.l.Warn(message, f)
	}
}

func (z Zap) Error(message string, args Fields) {
	f := make([]zapcore.Field, 0, len(args))

	for key, value := range args {
		t := checkZapType(value)

		f = append(f, zap.Field{
			Key:       key,
			Type:      t,
			Interface: value,
		})
	}

	z.l.Error(message, f...)
}

func (z Zap) Fatal(message string, args Fields) {
	for key, value := range args {
		f := zap.Field{
			Key:       key,
			Interface: value,
		}

		z.l.Fatal(message, f)
	}
}

func (z Zap) Info(message string, args Fields) {
	for key, value := range args {
		f := zap.Field{
			Key:       key,
			Interface: value,
		}

		z.l.Info(message, f)
	}
}

func checkZapType(v interface{}) zapcore.FieldType {
	switch v.(type) {
	case string:
		return zapcore.StringType
	case int:
		return zapcore.Int32Type
	case error:
		return zapcore.ErrorType
	default:
		return zapcore.UnknownType
	}

	return 0
}
