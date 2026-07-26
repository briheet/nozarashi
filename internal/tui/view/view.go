package view

import (
	"github.com/briheet/nozarashi/internal/tui/model"
	tea "github.com/charmbracelet/bubbletea"
)

type teaModel struct {
	m *model.Model
}

func (t teaModel) Init() tea.Cmd {
	return nil
}

func (t teaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return t, nil
}

func (t teaModel) View() string {
	return ""
}
