package frontend

import (
	"fmt"
	"os"
	"time"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/zachdehooge/WatchTower/internal/backend"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

type metricsMsg struct {
	mem string
	cpu string
}

type tickMsg struct{}

// 🔁 Fetch metrics from backend
func fetchMetrics() tea.Cmd {
	return func() tea.Msg {
		return metricsMsg{
			mem: backend.MemoryMetrics(),
			cpu: backend.CpuMetrics(),
		}
	}
}

// ⏱️ Tick for continuous updates
func tickCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return tea.Batch(fetchMetrics(), tickCmd())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case metricsMsg:
		m.table.SetRows([]table.Row{
			{"Virtual Memory", msg.mem + "%"},
			{"CPU Usage", msg.cpu + "%"},
		})
		return m, nil

	case tickMsg:
		return m, tea.Batch(fetchMetrics(), tickCmd())

	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Selected: %s", m.table.SelectedRow()[1]),
			)
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	return tea.NewView(
		baseStyle.Render(m.table.View()) +
			"\n  " + m.table.HelpView() + "\n",
	)
}

func TUI() {
	columns := []table.Column{
		{Title: "Metric", Width: 15},
		{Title: "Value", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithWidth(25),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t.SetStyles(s)

	m := model{table: t}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
