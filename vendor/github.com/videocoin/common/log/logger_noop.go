package log

func NoOp() Logger {
	return noOp{}
}

type noOp struct{}

func (noOp) Debugf(format string, args ...interface{}) {}
func (noOp) Debugln(args ...interface{})               {}
func (noOp) Infof(format string, args ...interface{})  {}
func (noOp) Infoln(args ...interface{})                {}
func (noOp) Warnf(format string, args ...interface{})  {}
func (noOp) Warnln(args ...interface{})                {}
func (noOp) Errorf(format string, args ...interface{}) {}
func (noOp) Errorln(args ...interface{})               {}
func (noOp) Fatalf(format string, args ...interface{}) {}
func (noOp) Fatalln(args ...interface{})               {}
