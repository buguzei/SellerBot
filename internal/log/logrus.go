package log

import "github.com/sirupsen/logrus"

type Logrus struct {
	logrus *logrus.Entry
}

func NewLogrus(lvl string) Logrus {
	var logrusLvl logrus.Level

	switch lvl {
	case "error":
		logrusLvl = logrus.ErrorLevel
	case "warn":
		logrusLvl = logrus.WarnLevel
	case "info":
		logrusLvl = logrus.InfoLevel
	case "debug":
		logrusLvl = logrus.DebugLevel
	default:
		logrusLvl = logrus.InfoLevel
	}

	logger := logrus.New()
	logger.SetLevel(logrusLvl)
	logger.Formatter = new(logrus.JSONFormatter)

	return Logrus{logrus: logrus.NewEntry(logger)}
}

func (l Logrus) Error(message string, args Fields) {
	l.logrus.WithFields(logrus.Fields(args)).Errorln(message)
}

func (l Logrus) Warn(message string, args Fields) {
	l.logrus.WithFields(logrus.Fields(args)).Warnln(message)
}

func (l Logrus) Named(name string) Logger {
	return &Logrus{
		logrus: l.logrus.WithField("name", name),
	}
}

func (l Logrus) Info(message string, args Fields) {
	l.logrus.WithFields(logrus.Fields(args)).Infoln(message)
}

func (l Logrus) Debug(message string, args Fields) {
	l.logrus.WithFields(logrus.Fields(args)).Debugln(message)
}

func (l Logrus) Fatal(message string, args Fields) {
	l.logrus.WithFields(logrus.Fields(args)).Fatalln(message)
}
