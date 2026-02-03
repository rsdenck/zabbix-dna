package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	zabbixRed   = lipgloss.Color("#D20000")
	zabbixBlack = lipgloss.Color("#000000")
	zabbixWhite = lipgloss.Color("#FFFFFF")

	titleStyle = lipgloss.NewStyle().
			Foreground(zabbixWhite).
			Background(zabbixRed).
			Padding(0, 1).
			Bold(true)

	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(zabbixRed)
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.title
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		return fmt.Sprintf("\n  Executando: %s\n\n", m.choice)
	}
	if m.quitting {
		return "\n  AtÃ© logo!\n\n"
	}
	return "\n" + m.list.View()
}

func Start() (string, error) {
	items := []list.Item{
		item{title: "host list", desc: "Listar todos os hosts"},
		item{title: "proxy list", desc: "Listar todos os proxies"},
		item{title: "template list", desc: "Listar todos os templates"},
		item{title: "hostgroup list", desc: "Listar todos os grupos de hosts"},
		item{title: "backup", desc: "Realizar backup das configuraÃ§Ãµes"},
		item{title: "exporter metrics", desc: "Iniciar exportador de mÃ©tricas OTLP"},
		item{title: "exporter traces", desc: "Iniciar exportador de traces OTLP"},
		item{title: "wizard", desc: "Abrir assistente de configuraÃ§Ã£o"},
	}

	const defaultWidth = 20

	l := list.New(items, list.NewDefaultDelegate(), defaultWidth, 14)
	l.Title = "ZABBIX-DNA CLI"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	finalModel, err := tea.NewProgram(m).Run()
	if err != nil {
		return "", err
	}

	res := finalModel.(model)
	return res.choice, nil
}
