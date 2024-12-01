package view

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"zimg/config"
	"zimg/utils"
)

const (
	githubOwner = iota
	githubRepo
	githubBucket
	githubToken
)
const (
	giteeOwner = iota
	giteeRepo
	giteeBucket
	giteeToken
)

type Set struct {
	stage    int
	cursor   int
	inputs   []textinput.Model
	selected string
	focused  int
}

func NewSet() *Set {
	current, _ := config.Get(config.Current)
	cursor := utils.IndexOf(config.All, current)
	if cursor == -1 {
		cursor = 0
	}
	return &Set{
		stage:  1,
		cursor: cursor,
	}
}

func newInput(str string, focus bool) textinput.Model {
	input := textinput.New()
	owner, _ := config.Get(str)
	input.Placeholder = owner
	if focus {
		input.Focus()
	}
	input.Prompt = ""
	input.Validate = validator
	return input
}

func validator(s string) error {
	return nil
}

func (m *Set) Init() tea.Cmd {
	return nil
}

func (m *Set) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.stage == 1 {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC:
				return m, tea.Quit
			case tea.KeyEnter:
				m.selected = config.All[m.cursor]
				m.stage = 2

				var inputs []textinput.Model
				if m.selected == config.Github {
					inputs = make([]textinput.Model, 4)
					inputs[githubOwner] = newInput(config.GithubOwner, true)
					inputs[githubRepo] = newInput(config.GithubRepo, false)
					inputs[githubBucket] = newInput(config.GithubBucket, false)
					inputs[githubToken] = newInput(config.GithubToken, false)
				} else if m.selected == config.Gitee {
					inputs = make([]textinput.Model, 4)
					inputs[giteeOwner] = newInput(config.GiteeOwner, true)
					inputs[giteeRepo] = newInput(config.GiteeRepo, false)
					inputs[giteeBucket] = newInput(config.GiteeBucket, false)
					inputs[giteeToken] = newInput(config.GiteeToken, false)
				} else {
					utils.Boom(errors.New("[500] view.set#Update"))
				}

				m.inputs = inputs
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
	} else if m.stage == 2 {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				if m.focused == len(m.inputs)-1 {
					for i := range m.inputs {
						m.inputs[i].Blur()
					}
					m.flush()
					return m, tea.Quit
				}
				m.nextInput()
			case tea.KeyCtrlC:
				return m, tea.Quit
			case tea.KeyUp:
				m.prevInput()
			case tea.KeyDown:
				m.nextInput()
			default:
				// ...
			}
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()

		for i := range m.inputs {
			m.inputs[i], _ = m.inputs[i].Update(msg)
		}
	} else {
		utils.Boom(errors.New("[500] view.set#Update"))
	}
	return m, nil
}

func (m *Set) View() string {
	inputStyle := lipgloss.NewStyle().Foreground(HotPink)
	s := strings.Builder{}

	if m.stage == 1 {
		s.WriteString("选择要配置的图床源:\n")

		for i, choice := range config.All {
			if m.cursor == i {
				s.WriteString(inputStyle.Render(fmt.Sprintf("> %s", choice)) + "\n")
			} else {
				s.WriteString(fmt.Sprintf("  %s\n", choice))
			}
		}
	} else if m.stage == 2 {
		s.WriteString(fmt.Sprintf("正在配置[%s]图床源:\n\n", m.selected))
		if m.selected == config.Github {
			s.WriteString(inputStyle.Render(config.GithubOwner) + "\n")
			s.WriteString(m.inputs[githubOwner].View() + "\n")
			s.WriteString(inputStyle.Render(config.GithubRepo) + "\n")
			s.WriteString(m.inputs[githubRepo].View() + "\n")
			s.WriteString(inputStyle.Render(config.GithubBucket) + "\n")
			s.WriteString(m.inputs[githubBucket].View() + "\n")
			s.WriteString(inputStyle.Render(config.GithubToken) + "\n")
			s.WriteString(m.inputs[githubToken].View() + "\n")
		} else if m.selected == config.Gitee {
			s.WriteString(inputStyle.Render(config.GiteeOwner) + "\n")
			s.WriteString(m.inputs[giteeOwner].View() + "\n")
			s.WriteString(inputStyle.Render(config.GiteeRepo) + "\n")
			s.WriteString(m.inputs[giteeRepo].View() + "\n")
			s.WriteString(inputStyle.Render(config.GiteeBucket) + "\n")
			s.WriteString(m.inputs[giteeBucket].View() + "\n")
			s.WriteString(inputStyle.Render(config.GiteeToken) + "\n")
			s.WriteString(m.inputs[giteeToken].View() + "\n")
		}
	} else {
		utils.Boom(errors.New("[500] view.set#View"))
	}

	return s.String()
}

// nextInput focuses the next input field
func (m *Set) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *Set) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

func (m *Set) flush() {
	if m.selected == config.Github {
		owner := m.inputs[githubOwner].Value()
		if owner != "" {
			config.Zimgrc.Set(config.GithubOwner, owner)
		}
		repo := m.inputs[githubRepo].Value()
		if repo != "" {
			config.Zimgrc.Set(config.GithubRepo, repo)
		}
		bucket := m.inputs[githubBucket].Value()
		if bucket != "" {
			config.Zimgrc.Set(config.GithubBucket, bucket)
		}
		token := m.inputs[githubToken].Value()
		if token != "" {
			config.Zimgrc.Set(config.GithubToken, token)
		}
	} else if m.selected == config.Gitee {
		owner := m.inputs[giteeOwner].Value()
		if owner != "" {
			config.Zimgrc.Set(config.GiteeOwner, owner)
		}
		repo := m.inputs[giteeRepo].Value()
		if repo != "" {
			config.Zimgrc.Set(config.GiteeRepo, repo)
		}
		bucket := m.inputs[giteeBucket].Value()
		if bucket != "" {
			config.Zimgrc.Set(config.GiteeBucket, bucket)
		}
		token := m.inputs[giteeToken].Value()
		if token != "" {
			config.Zimgrc.Set(config.GiteeToken, token)
		}
	}
	config.Flush()
}
