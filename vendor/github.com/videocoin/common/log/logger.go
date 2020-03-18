package log

type Logger interface {
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})

	Infof(format string, args ...interface{})
	Infoln(args ...interface{})

	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})

	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})

	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
}

// Fields convenience type for adding multiple fields to a log statement.
type Fields map[string]interface{}
