package view

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"zimg/config"
	"zimg/utils"
)

type Get struct {
	cursor int
	Result string
}

func NewGet() *Get {
	current, _ := config.Get(config.Current)
	cursor := utils.IndexOf(config.All, current)
	if cursor == -1 {
		cursor = 0
	}
	return &Get{cursor: cursor}
}

func (m Get) Init() tea.Cmd {
	return nil
}

func (m Get) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			src := config.All[m.cursor]
			group := config.GetGroup(config.All[m.cursor])
			builder := strings.Builder{}
			builder.WriteString(fmt.Sprintf("[%s]配置如下:\n", src))
			for i, e := range group {
				if i != len(group)-1 {
					builder.WriteString(fmt.Sprintf("  %s\n", e))
				} else {
					builder.WriteString(fmt.Sprintf("  %s", e))
				}
			}
			m.Result = builder.String()
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

func (m Get) View() string {
	inputStyle := lipgloss.NewStyle().Foreground(HotPink)
	s := strings.Builder{}
	s.WriteString("选择要查看的图床源:\n")

	for i, choice := range config.All {
		if m.cursor == i {
			s.WriteString(inputStyle.Render(fmt.Sprintf("> %s", choice)) + "\n")
		} else {
			s.WriteString(fmt.Sprintf("  %s\n", choice))
		}
	}

	return s.String()
}
