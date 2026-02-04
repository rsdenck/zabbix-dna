package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize Zabbix-DNA configuration",
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Zabbix API URL (e.g., http://localhost/zabbix/api_jsonrpc.php): ")
			url, _ := reader.ReadString('\n')
			url = strings.TrimSpace(url)

			fmt.Print("Zabbix API Token (optional): ")
			token, _ := reader.ReadString('\n')
			token = strings.TrimSpace(token)

			var user, password string
			if token == "" {
				fmt.Print("Zabbix Username: ")
				user, _ = reader.ReadString('\n')
				user = strings.TrimSpace(user)

				fmt.Print("Zabbix Password: ")
				password, _ = reader.ReadString('\n')
				password = strings.TrimSpace(password)
			}

			configContent := fmt.Sprintf(`[zabbix]
url = "%s"
token = "%s"
user = "%s"
password = "%s"
timeout = 30

[otlp]
endpoint = "http://localhost:4318"
protocol = "http"
service_name = "zabbix-dna"
`, url, token, user, password)

			err := os.WriteFile("zabbix-dna.toml", []byte(configContent), 0644)
			handleError(err)

			headers := []string{"File", "Status"}
			rows := [][]string{
				{"zabbix-dna.toml", "Saved Successfully"},
			}
			outputResult(cmd, map[string]string{"file": "zabbix-dna.toml", "status": "success"}, headers, rows)
		},
	}
}
