package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	zabbixRed    = lipgloss.Color("#d64e4e")
	zabbixDark   = lipgloss.Color("#1a1a1a")
	zabbixGray   = lipgloss.Color("#333333")
	zabbixWhite  = lipgloss.Color("#FFFFFF")
	glassOverlay = lipgloss.Color("#2a2a2a")

	// Glassmorphism Styles
	windowStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(zabbixRed).
			Padding(1, 2).
			Background(zabbixDark)

	headerStyle = lipgloss.NewStyle().
			Foreground(zabbixRed).
			Bold(true).
			MarginBottom(1)

	asciiArt = `  ███████╗ █████╗ ██████╗ ██████╗ ██╗██╗  ██╗    ██████╗██╗     ██╗ 
  ╚══███╔╝██╔══██╗██╔══██╗██╔══██╗██║╚██╗██╔╝    ██╔════╝██║     ██║ 
    ███╔╝ ███████║██████╔╝██████╔╝██║ ╚███╔╝     ██║     ██║     ██║ 
   ███╔╝  ██╔══██║██╔══██╗██╔══██╗██║ ██╔██╗     ██║     ██║     ██║ 
  ███████╗██║  ██║██████╔╝██████╔╝██║██╔╝ ██╗    ╚██████╗███████╗██║ 
  ╚══════╝╚═╝  ╚═╝╚═════╝ ╚═════╝ ╚═╝╚═╝  ╚═╝     ╚═════╝╚══════╝╚═╝ v6.4`

	asciiStyle = lipgloss.NewStyle().Foreground(zabbixRed).MarginBottom(1)

	titleStyle = lipgloss.NewStyle().
			Foreground(zabbixWhite).
			Background(zabbixRed).
			Padding(0, 1).
			Bold(true).
			MarginBottom(1)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(zabbixWhite)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(0).
				Foreground(zabbixRed).
				Bold(true).
				SetString("  > ")

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true).
			PaddingLeft(4)

	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 2 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.title)

	var title, desc string
	if index == m.Index() {
		title = selectedItemStyle.Render(str)
		desc = descStyle.Foreground(lipgloss.Color("#aaaaaa")).Render(i.desc)
	} else {
		title = itemStyle.Render(str)
		desc = descStyle.Render(i.desc)
	}

	fmt.Fprintf(w, "%s\n%s", title, desc)
}

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
		return "\n  Até logo!\n\n"
	}

	header := asciiStyle.Render(asciiArt)
	content := m.list.View()

	// Wrap everything in a glassmorphism window
	return windowStyle.Render(lipgloss.JoinVertical(lipgloss.Left, header, content))
}

func Start() (string, error) {
	items := []list.Item{
		item{title: "host list", desc: "Listar todos os hosts"},
		item{title: "proxy list", desc: "Listar todos os proxies"},
		item{title: "salt ping", desc: "Pingar minions (proxies) via SaltStack"},
		item{title: "salt run", desc: "Executar módulo em minions via SaltStack"},
		item{title: "template list", desc: "Listar todos os templates"},
		item{title: "hostgroup list", desc: "Listar todos os grupos de hosts"},
		item{title: "backup", desc: "Realizar backup das configurações"},
		item{title: "exporter metrics", desc: "Iniciar exportador de métricas OTLP"},
		item{title: "exporter traces", desc: "Iniciar exportador de traces OTLP"},
		item{title: "wizard", desc: "Abrir assistente de configuração"},
	}

	const defaultWidth = 60
	const defaultHeight = 20

	l := list.New(items, itemDelegate{}, defaultWidth, defaultHeight)
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
