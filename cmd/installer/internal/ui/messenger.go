package ui

import (
	"claude_analysis/cmd/installer/internal/logger"

	tea "github.com/charmbracelet/bubbletea"
)

// StatusMessenger provides a way to send status updates to the UI
type StatusMessenger struct {
	Program *tea.Program
}

// NewStatusMessenger creates a new status messenger
func NewStatusMessenger(program *tea.Program) *StatusMessenger {
	return &StatusMessenger{
		Program: program,
	}
}

// Implement logger.Logger interface
var _ logger.Logger = (*StatusMessenger)(nil)

// SendStatus sends a status message to the UI
func (sm *StatusMessenger) SendStatus(msgType StatusType, message string, details ...string) {
	if sm.Program != nil {
		msg := NewStatusMsg(msgType, message, details...)
		sm.Program.Send(msg)
	}
}

// SendProgress sends a progress update to the UI
func (sm *StatusMessenger) SendProgress(step, totalSteps int, currentTask string) {
	if sm.Program != nil {
		msg := NewProgressMsg(step, totalSteps, currentTask)
		sm.Program.Send(msg)
	}
}

// Info sends an info message
func (sm *StatusMessenger) Info(message string, details ...string) {
	sm.SendStatus(StatusInfo, message, details...)
}

// Success sends a success message
func (sm *StatusMessenger) Success(message string, details ...string) {
	sm.SendStatus(StatusSuccess, message, details...)
}

// Warning sends a warning message
func (sm *StatusMessenger) Warning(message string, details ...string) {
	sm.SendStatus(StatusWarning, message, details...)
}

// Error sends an error message
func (sm *StatusMessenger) Error(message string, details ...string) {
	sm.SendStatus(StatusError, message, details...)
}

// Progress sends a progress message
func (sm *StatusMessenger) Progress(message string, details ...string) {
	sm.SendStatus(StatusProgress, message, details...)
}
