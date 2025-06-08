package tgbotapp

import "fmt"

var (
	ErrLoggerNotFound = fmt.Errorf("Logger not initialized. Consider using `Logger` middleware.")
)
