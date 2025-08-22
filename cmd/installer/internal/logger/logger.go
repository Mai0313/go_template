package logger

// StatusType represents different types of status messages
type StatusType int

const (
	StatusInfo StatusType = iota
	StatusSuccess
	StatusWarning
	StatusError
	StatusProgress
)

// Logger interface for sending status updates
type Logger interface {
	Info(message string, details ...string)
	Success(message string, details ...string)
	Warning(message string, details ...string)
	Error(message string, details ...string)
	Progress(message string, details ...string)
	SendProgress(step, totalSteps int, currentTask string)
}

// Global logger instance
var GlobalLogger Logger

// Helper functions that use the global logger
func Info(message string, details ...string) {
	if GlobalLogger != nil {
		GlobalLogger.Info(message, details...)
	}
}

func Success(message string, details ...string) {
	if GlobalLogger != nil {
		GlobalLogger.Success(message, details...)
	}
}

func Warning(message string, details ...string) {
	if GlobalLogger != nil {
		GlobalLogger.Warning(message, details...)
	}
}

func Error(message string, details ...string) {
	if GlobalLogger != nil {
		GlobalLogger.Error(message, details...)
	}
}

func Progress(message string, details ...string) {
	if GlobalLogger != nil {
		GlobalLogger.Progress(message, details...)
	}
}

func SendProgress(step, totalSteps int, currentTask string) {
	if GlobalLogger != nil {
		GlobalLogger.SendProgress(step, totalSteps, currentTask)
	}
}

// NoOpLogger is a logger that does nothing (fallback)
type NoOpLogger struct{}

func (NoOpLogger) Info(message string, details ...string)                {}
func (NoOpLogger) Success(message string, details ...string)             {}
func (NoOpLogger) Warning(message string, details ...string)             {}
func (NoOpLogger) Error(message string, details ...string)               {}
func (NoOpLogger) Progress(message string, details ...string)            {}
func (NoOpLogger) SendProgress(step, totalSteps int, currentTask string) {}
