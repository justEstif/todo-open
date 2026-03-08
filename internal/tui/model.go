package tui

import (
	"context"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	apiclient "github.com/justEstif/todo-open/internal/client/api"
	"github.com/justEstif/todo-open/internal/core"
)

// view is the active screen.
type view int

const (
	viewList   view = iota
	viewDetail      // split: list left + detail right
)

// --- messages ---

type msgTasksLoaded []core.Task
type msgTaskEvent apiclient.TaskEvent
type msgError struct{ err error }

// waitForEvent returns a Cmd that blocks until the next SSE event arrives.
func waitForEvent(ch <-chan apiclient.TaskEvent) tea.Cmd {
	return func() tea.Msg {
		e, ok := <-ch
		if !ok {
			return nil
		}
		return msgTaskEvent(e)
	}
}

// --- model ---

// Model is the single source of truth for all TUI state.
type Model struct {
	client     *apiclient.Client
	eventsCh   <-chan apiclient.TaskEvent
	eventsStop func()

	tasks      []core.Task
	tasksByID  map[string]core.Task

	view     view
	cursor   int    // index into tasks
	filter   string // "all", "open", "done"

	// detail dep navigation
	depCursor    int  // cursor within deps list in detail pane
	depFocused   bool // whether dep list has focus

	// create bar
	creating  bool
	createBuf string

	// layout
	width  int
	height int

	err string
}

func newModel(client *apiclient.Client) Model {
	return Model{
		client:    client,
		filter:    "all",
		tasksByID: map[string]core.Task{},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadTasks(),
		m.startEvents(),
	)
}

func (m Model) loadTasks() tea.Cmd {
	return func() tea.Msg {
		tasks, err := m.client.ListTasks()
		if err != nil {
			return msgError{err}
		}
		return msgTasksLoaded(tasks)
	}
}

func (m Model) startEvents() tea.Cmd {
	return func() tea.Msg {
		ch, cancel, err := m.client.SubscribeEvents(context.Background())
		if err != nil {
			// Non-fatal: TUI works without live updates.
			return nil
		}
		// Store channel + cancel on the model via a dedicated msg type.
		return msgEventsReady{ch: ch, cancel: cancel}
	}
}

type msgEventsReady struct {
	ch     <-chan apiclient.TaskEvent
	cancel func()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case msgTasksLoaded:
		m.setTasks([]core.Task(msg))
		return m, nil

	case msgEventsReady:
		m.eventsCh = msg.ch
		m.eventsStop = msg.cancel
		return m, waitForEvent(m.eventsCh)

	case msgTaskEvent:
		// Reload full list on any mutation — keeps logic simple.
		return m, tea.Batch(m.loadTasks(), waitForEvent(m.eventsCh))

	case msgError:
		m.err = msg.err.Error()
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m *Model) setTasks(tasks []core.Task) {
	filtered := make([]core.Task, 0, len(tasks))
	for _, t := range tasks {
		switch m.filter {
		case "open":
			if t.Status == core.TaskStatusOpen || t.Status == core.TaskStatusInProgress || t.Status == core.TaskStatusPending {
				filtered = append(filtered, t)
			}
		case "done":
			if t.Status == core.TaskStatusDone {
				filtered = append(filtered, t)
			}
		default:
			filtered = append(filtered, t)
		}
	}
	m.tasks = filtered

	// Rebuild ID index from full list for dep resolution.
	idx := make(map[string]core.Task, len(tasks))
	for _, t := range tasks {
		idx[t.ID] = t
	}
	m.tasksByID = idx

	// Clamp cursor.
	if m.cursor >= len(m.tasks) {
		m.cursor = max(0, len(m.tasks)-1)
	}
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Create bar captures all input.
	if m.creating {
		return m.handleCreateKey(msg)
	}

	switch msg.String() {
	case "q", "ctrl+c":
		if m.eventsStop != nil {
			m.eventsStop()
		}
		return m, tea.Quit

	case "up", "k":
		if m.view == viewList || m.view == viewDetail {
			if m.cursor > 0 {
				m.cursor--
				m.depCursor = 0
				m.depFocused = false
			}
		}

	case "down", "j":
		if m.view == viewList || m.view == viewDetail {
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
				m.depCursor = 0
				m.depFocused = false
			}
		}

	case "enter":
		if m.view == viewList && len(m.tasks) > 0 {
			m.view = viewDetail
			m.depCursor = 0
			m.depFocused = false
		} else if m.view == viewDetail {
			// Jump to dep under dep cursor if focused.
			if m.depFocused {
				return m.jumpToDep()
			}
		}

	case "esc":
		if m.view == viewDetail {
			m.view = viewList
			m.depFocused = false
		}

	case "n":
		if m.view == viewList {
			m.creating = true
			m.createBuf = ""
		}

	case "d":
		if len(m.tasks) > 0 {
			return m, m.toggleDone()
		}

	case "a":
		m.filter = "all"
		return m, m.loadTasks()

	case "o":
		m.filter = "open"
		return m, m.loadTasks()

	case "D":
		m.filter = "done"
		return m, m.loadTasks()

	case "tab":
		// In detail view, toggle focus between task nav and dep nav.
		if m.view == viewDetail && len(m.tasks) > 0 {
			t := m.tasks[m.cursor]
			deps := allDeps(t)
			if len(deps) > 0 {
				m.depFocused = !m.depFocused
				m.depCursor = 0
			}
		}
	}

	return m, nil
}

func (m Model) handleCreateKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.creating = false
		m.createBuf = ""
	case "enter":
		title := strings.TrimSpace(m.createBuf)
		m.creating = false
		m.createBuf = ""
		if title != "" {
			return m, m.createTask(title)
		}
	case "backspace", "ctrl+h":
		if len(m.createBuf) > 0 {
			m.createBuf = m.createBuf[:len(m.createBuf)-1]
		}
	default:
		if len(msg.Runes) > 0 {
			m.createBuf += string(msg.Runes)
		}
	}
	return m, nil
}

func (m Model) toggleDone() tea.Cmd {
	if len(m.tasks) == 0 {
		return nil
	}
	t := m.tasks[m.cursor]
	return func() tea.Msg {
		if t.Status == core.TaskStatusDone {
			_, err := m.client.PatchTaskStatus(t.ID, string(core.TaskStatusOpen))
			if err != nil {
				return msgError{err}
			}
		} else {
			_, err := m.client.CompleteTask(t.ID)
			if err != nil {
				return msgError{err}
			}
		}
		tasks, err := m.client.ListTasks()
		if err != nil {
			return msgError{err}
		}
		return msgTasksLoaded(tasks)
	}
}

func (m Model) createTask(title string) tea.Cmd {
	return func() tea.Msg {
		_, err := m.client.CreateTask(title)
		if err != nil {
			return msgError{err}
		}
		tasks, err := m.client.ListTasks()
		if err != nil {
			return msgError{err}
		}
		return msgTasksLoaded(tasks)
	}
}

func (m Model) jumpToDep() (tea.Model, tea.Cmd) {
	if len(m.tasks) == 0 {
		return m, nil
	}
	t := m.tasks[m.cursor]
	deps := allDeps(t)
	if m.depCursor >= len(deps) {
		return m, nil
	}
	targetID := deps[m.depCursor]
	for i, task := range m.tasks {
		if task.ID == targetID {
			m.cursor = i
			m.depFocused = false
			m.depCursor = 0
			return m, nil
		}
	}
	return m, nil
}

// allDeps returns all dep IDs for a task in display order: blockedBy, blocking, triggers.
func allDeps(t core.Task) []string {
	out := make([]string, 0, len(t.BlockedBy)+len(t.Blocking)+len(t.TriggerIDs))
	out = append(out, t.BlockedBy...)
	out = append(out, t.Blocking...)
	out = append(out, t.TriggerIDs...)
	return out
}

// View renders the full terminal screen.
func (m Model) View() string {
	if m.width == 0 {
		return "loading…\n"
	}

	innerW := m.width - borderOverhead
	innerH := m.height - verticalOverhead

	if m.err != "" {
		return styleBorder.Width(m.width - 2).Render(
			styleHeader.Render("todo.open") + "\n\n" +
				lipgloss.NewStyle().Foreground(colorCritical).Render("error: "+m.err) + "\n\n" +
				styleKeys.Render("q quit"),
		)
	}

	header := styleHeader.Render("todo.open")

	switch m.view {
	case viewDetail:
		return m.viewDetailLayout(header, innerW, innerH)
	default:
		return m.viewListLayout(header, innerW, innerH)
	}
}

func (m Model) viewListLayout(header string, innerW, innerH int) string {
	listContent := renderList(m.tasks, m.cursor, innerW, innerH, m.filter)

	var footer string
	if m.creating {
		footer = renderCreate(m.createBuf, innerW)
	} else {
		footer = renderKeyBar(viewList, innerW)
	}

	body := lipgloss.JoinVertical(lipgloss.Left,
		header,
		listContent,
		footer,
	)
	return styleBorder.Width(m.width-2).Render(body)
}

func (m Model) viewDetailLayout(header string, innerW, innerH int) string {
	listW := int(float64(innerW) * listPaneFraction)
	detailW := innerW - listW - 1 // -1 for divider

	listContent := renderList(m.tasks, m.cursor, listW, innerH, m.filter)

	var detailContent string
	if len(m.tasks) > 0 {
		detailContent = renderDetail(m.tasks[m.cursor], m.tasksByID, detailW, innerH)
	}

	divider := styleDivider.Render(strings.Repeat("│\n", innerH+1))

	split := lipgloss.JoinHorizontal(lipgloss.Top,
		listContent,
		divider,
		detailContent,
	)

	body := lipgloss.JoinVertical(lipgloss.Left,
		header,
		split,
		renderKeyBar(viewDetail, innerW),
	)
	return styleBorder.Width(m.width-2).Render(body)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
