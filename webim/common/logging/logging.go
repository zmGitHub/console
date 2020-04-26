package log

var Logger *CompatibleLogger

func NewLogging() {
	Logger = NewCompatibleLogger(true)
}
