package view

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"zimg/config"
	"zimg/utils"
)

type Use struct {
	cursor int
}

func NewUse() *Use {
	current, _ := config.Get(config.Current)
	cursor := utils.IndexOf(config.All, current)
	if cursor == -1 {
		cursor = 0
	}
	return &Use{cursor: cursor}
}

func (m Use) Init() tea.Cmd {
	return nil
}

func (m Use) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			m.flush()
			return m, tea.Quit
		case tea.KeyDown:
			m.cursor++
			if m.cursor >= len(config.All) {
				m.cursor = 0
			}
		case tea.KeyUp:
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(config.All) - 1
			}
		default:
			// ...
		}
	}

	return m, nil
}

func (m Use) View() string {
	inputStyle := lipgloss.NewStyle().Foreground(HotPink)
	s := strings.Builder{}
	s.WriteString("选择要使用的图床源:\n")

	for i, choice := range config.All {
		if m.cursor == i {
			s.WriteString(inputStyle.Render(fmt.Sprintf("> %s", choice)) + "\n")
		} else {
			s.WriteString(fmt.Sprintf("  %s\n", choice))
		}
	}

	return s.String()
}

func (m Use) flush() {
	config.Set(config.Current, config.All[m.cursor])
	config.Flush()
}
