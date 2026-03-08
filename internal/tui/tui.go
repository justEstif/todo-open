// Package tui provides a Bubble Tea terminal UI for todo.open.
// The TUI is a pure HTTP client — all mutations flow through the server API.
// Entry point: Run(addr).
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	apiclient "github.com/justEstif/todo-open/internal/client/api"
)

// Run starts the TUI connected to the todo.open server at addr.
func Run(addr string) error {
	client := apiclient.New(addr)
	m := newModel(client)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
