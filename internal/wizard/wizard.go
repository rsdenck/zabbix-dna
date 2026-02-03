package wizard

import (
	"fmt"
	"os"
	"zabbix-dna/internal/config"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pelletier/go-toml/v2"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#D20000")).
			Padding(0, 1).
			Bold(true)
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#D20000"))
	cursorStyle  = focusedStyle.Copy()
)

type model struct {
	inputs  []textinput.Model
	focused int
	err     error
	done    bool
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			if m.focused == len(m.inputs)-1 {
				m.done = true
				return m, tea.Quit
			}
			m.focused++
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				if i == m.focused {
					cmds[i] = m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}
			return m, tea.Batch(cmds...)

		case "up", "shift+tab":
			if m.focused > 0 {
				m.focused--
			}
		case "down", "tab":
			if m.focused < len(m.inputs)-1 {
				m.focused++
			}
		}

		cmds := make([]tea.Cmd, len(m.inputs))
		for i := range m.inputs {
			if i == m.focused {
				cmds[i] = m.inputs[i].Focus()
			} else {
				m.inputs[i].Blur()
			}
		}

		for i := range m.inputs {
			m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
		}
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m model) View() string {
	if m.done {
		return "\n  ConfiguraÃ§Ã£o salva com sucesso!\n\n"
	}

	s := "\n" + titleStyle.Render("WIZARD DE CONFIGURAÃ‡ÃƒO ZABBIX-DNA") + "\n\n"

	for i := range m.inputs {
		s += fmt.Sprintf(
			"  %s\n  %s\n\n",
			m.inputs[i].Placeholder,
			m.inputs[i].View(),
		)
	}

	s += "\n  (enter para prÃ³ximo, esc para sair)\n"

	return s
}

func Start() error {
	inputs := make([]textinput.Model, 5)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "URL do Zabbix (ex: https://zabbix.exemplo.com/api_jsonrpc.php)"
	inputs[0].Focus()
	inputs[0].CharLimit = 156
	inputs[0].Width = 50

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Token de API"
	inputs[1].CharLimit = 128
	inputs[1].Width = 50
	inputs[1].EchoMode = textinput.EchoPassword
	inputs[1].EchoCharacter = '*'

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "Endpoint OTLP (ex: http://localhost:4318)"
	inputs[2].CharLimit = 156
	inputs[2].Width = 50

	inputs[3] = textinput.New()
	inputs[3].Placeholder = "Caminho do arquivo de configuraÃ§Ã£o (default: zabbix-dna.toml)"
	inputs[3].CharLimit = 100
	inputs[3].Width = 50

	inputs[4] = textinput.New()
	inputs[4].Placeholder = "NÃ­vel de Log (debug, info, warn, error)"
	inputs[4].CharLimit = 20
	inputs[4].Width = 20

	m := model{inputs: inputs}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	m = finalModel.(model)
	if !m.done {
		return nil
	}

	// Extrair valores
	zabbixURL := m.inputs[0].Value()
	apiToken := m.inputs[1].Value()
	otlpEndpoint := m.inputs[2].Value()
	configPath := m.inputs[3].Value()
	logLevel := m.inputs[4].Value()

	if configPath == "" {
		configPath = "zabbix-dna.toml"
	}

	cfg := config.Config{
		Zabbix: config.ZabbixConfig{
			URL:     zabbixURL,
			Token:   apiToken,
			Timeout: 30,
		},
		OTLP: config.OTLPConfig{
			Endpoint:    otlpEndpoint,
			Protocol:    "http",
			ServiceName: "zabbix-dna",
		},
		LogLevel: logLevel,
	}

	data, err := toml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("erro ao gerar TOML: %w", err)
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar arquivo: %w", err)
	}

	fmt.Printf("\nConfiguraÃ§Ã£o salva com sucesso em: %s\n", configPath)
	return nil
}
