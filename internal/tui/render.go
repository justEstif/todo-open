package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/justEstif/todo-open/internal/core"
)

const (
	listPaneFraction   = 0.45 // fraction of terminal width for list in split view
	detailMinWidth     = 30
	borderOverhead     = 4 // 2 borders × 2 sides
	verticalOverhead   = 6 // header + filter bar + key bar + borders
	createBarHeight    = 2
)

// taskSymbol returns the bullet symbol for a task's status.
func taskSymbol(t core.Task) string {
	switch t.Status {
	case core.TaskStatusDone:
		return "✓"
	case core.TaskStatusArchived:
		return "○"
	case core.TaskStatusInProgress:
		return "◎"
	default:
		return "●"
	}
}

// taskRowStyle returns the base style for a task row.
func taskRowStyle(t core.Task) lipgloss.Style {
	if t.Status == core.TaskStatusDone || t.Status == core.TaskStatusArchived {
		return styleDone
	}
	return styleBase
}

// renderList renders the task list pane. width/height are the inner dimensions
// (excluding border). selected is the index of the cursor row.
func renderList(tasks []core.Task, selected int, width, height int, filter string) string {
	var sb strings.Builder

	// Filter bar
	filterLine := styleStatusBar.Render(fmt.Sprintf(
		" filter: %-6s  %s %s %s %s",
		filter,
		styleKeyHint.Render("[a]"),
		styleMuted.Render("all"),
		styleKeyHint.Render("[o]"),
		styleMuted.Render("open"),
	))
	sb.WriteString(lipgloss.NewStyle().Width(width).Render(filterLine))
	sb.WriteString("\n")

	// Task rows — allow height minus filter line minus create hint
	maxRows := height - 2
	if maxRows < 1 {
		maxRows = 1
	}

	// Window the list around the cursor
	start := 0
	if selected >= maxRows {
		start = selected - maxRows + 1
	}

	for i := start; i < len(tasks) && i < start+maxRows; i++ {
		t := tasks[i]
		sym := taskSymbol(t)
		base := taskRowStyle(t)
		pStyle := priorityStyle(string(t.Priority))

		title := t.Title
		// Truncate title to fit
		maxTitle := width - 20
		if maxTitle < 8 {
			maxTitle = 8
		}
		if len(title) > maxTitle {
			title = title[:maxTitle-1] + "…"
		}

		pri := string(t.Priority)
		if pri == "" {
			pri = "normal"
		}
		status := string(t.Status)

		var row string
		if i == selected {
			cursor := styleSelected.Render("▶ ")
			symStr := styleSelected.Render(sym + " ")
			titleStr := styleSelected.Render(fmt.Sprintf("%-*s", maxTitle, title))
			priStr := pStyle.Bold(true).Render(fmt.Sprintf("%-8s", pri))
			statStr := styleSelected.Render(fmt.Sprintf("%-12s", status))
			row = cursor + symStr + titleStr + " " + priStr + " " + statStr
		} else {
			cursor := "  "
			symStr := base.Render(sym + " ")
			titleStr := base.Render(fmt.Sprintf("%-*s", maxTitle, title))
			priStr := pStyle.Render(fmt.Sprintf("%-8s", pri))
			statStr := base.Render(fmt.Sprintf("%-12s", status))
			row = cursor + symStr + titleStr + " " + priStr + " " + statStr
		}

		sb.WriteString(lipgloss.NewStyle().Width(width).Render(row))
		sb.WriteString("\n")
	}

	// Pad remaining rows
	rendered := start + maxRows
	for i := len(tasks); i < rendered; i++ {
		sb.WriteString(lipgloss.NewStyle().Width(width).Render(""))
		sb.WriteString("\n")
	}

	return sb.String()
}

// renderDetail renders the detail pane for the selected task.
// tasksByID is used to resolve dependency titles.
func renderDetail(t core.Task, tasksByID map[string]core.Task, width, height int) string {
	var sb strings.Builder

	titleStr := styleTitle.Width(width).Render(t.Title)
	sb.WriteString(titleStr + "\n\n")

	row := func(label, value string) {
		sb.WriteString(styleLabel.Render(label))
		sb.WriteString(styleValue.Render(value))
		sb.WriteString("\n")
	}

	pri := string(t.Priority)
	if pri == "" {
		pri = "normal"
	}

	row("status:", string(t.Status))
	row("priority:", priorityStyle(pri).Render(pri))
	row("created:", t.CreatedAt.Format("2006-01-02"))
	if t.StartedAt != nil {
		row("started:", t.StartedAt.Format("2006-01-02"))
	}
	if t.CompletedAt != nil {
		row("done:", t.CompletedAt.Format("2006-01-02"))
	}
	if t.Assignee != "" {
		row("assignee:", t.Assignee)
	}

	// Description
	if t.Description != "" {
		sb.WriteString("\n")
		// Word-wrap description to width
		desc := lipgloss.NewStyle().Width(width).Foreground(colorNormal).Render(t.Description)
		sb.WriteString(desc)
		sb.WriteString("\n")
	}

	// Dependencies
	hasBlocking := len(t.Blocking) > 0
	hasBlockedBy := len(t.BlockedBy) > 0

	if hasBlockedBy {
		sb.WriteString("\n")
		sb.WriteString(styleMuted.Render("blocked by") + "\n")
		for _, id := range t.BlockedBy {
			sb.WriteString(depRow(id, tasksByID, width))
		}
	}

	if hasBlocking {
		sb.WriteString("\n")
		sb.WriteString(styleMuted.Render("blocking") + "\n")
		for _, id := range t.Blocking {
			sb.WriteString(depRow(id, tasksByID, width))
		}
	}

	if len(t.TriggerIDs) > 0 {
		sb.WriteString("\n")
		sb.WriteString(styleMuted.Render("triggers") + "\n")
		for _, id := range t.TriggerIDs {
			sb.WriteString(depRow(id, tasksByID, width))
		}
	}

	return sb.String()
}

func depRow(id string, tasksByID map[string]core.Task, width int) string {
	dep, ok := tasksByID[id]
	if !ok {
		return styleMuted.Render("  ? "+id) + "\n"
	}
	sym := taskSymbol(dep)
	base := taskRowStyle(dep)
	maxTitle := width - 6
	title := dep.Title
	if len(title) > maxTitle {
		title = title[:maxTitle-1] + "…"
	}
	return base.Render(fmt.Sprintf("  %s %s", sym, title)) + "\n"
}

// renderInputBar renders a single-line text input bar with a prompt label.
func renderInputBar(prompt, input string, width int) string {
	promptStr := styleKeyHint.Render("> " + prompt + ": ")
	cursor := styleSelected.Render("█")
	line := promptStr + styleValue.Render(input) + cursor
	return lipgloss.NewStyle().Width(width).Render(line) + "\n" +
		styleKeys.Render("  enter confirm   esc cancel") + "\n"
}

// renderCreate renders the inline create bar.
func renderCreate(input string, width int) string {
	return renderInputBar("new task", input, width)
}

// renderEdit renders the inline edit bar pre-filled with the task's current title.
func renderEdit(input string, width int) string {
	return renderInputBar("edit title", input, width)
}

// renderKeyBar renders the bottom key hint bar for the given view.
func renderKeyBar(v view, width int) string {
	var hints string
	switch v {
	case viewList:
		hints = keyHint("n", "new") + "  " +
			keyHint("e", "edit") + "  " +
			keyHint("enter", "detail") + "  " +
			keyHint("d", "done") + "  " +
			keyHint("a/o/D", "filter") + "  " +
			keyHint("q", "quit")
	case viewDetail:
		hints = keyHint("esc", "back") + "  " +
			keyHint("e", "edit") + "  " +
			keyHint("d", "done") + "  " +
			keyHint("tab", "deps") + "  " +
			keyHint("enter", "jump") + "  " +
			keyHint("q", "quit")
	}
	return styleKeys.Width(width).Render(" " + hints)
}

func keyHint(key, label string) string {
	return styleKeyHint.Render(key) + " " + styleMuted.Render(label)
}
