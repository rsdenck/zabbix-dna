package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage Zabbix users",
	}

	cmd.AddCommand(newUserListCmd())
	cmd.AddCommand(newUserShowCmd())
	cmd.AddCommand(newUserCreateCmd())
	cmd.AddCommand(newUserDeleteCmd())

	return cmd
}

func newUserListCmd() *cobra.Command {
	var limit int
	var search string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Zabbix users",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"output": []string{"userid", "username", "name", "surname", "roleid"},
				"limit":  limit,
			}
			if search != "" {
				params["search"] = map[string]interface{}{
					"username": search,
				}
			}

			result, err := client.Call("user.get", params)
			handleError(err)

			var users []map[string]interface{}
			json.Unmarshal(result, &users)

			headers := []string{"UserID", "Username", "First Name", "Last Name", "Role ID"}
			var rows [][]string
			for _, u := range users {
				rows = append(rows, []string{
					fmt.Sprintf("%v", u["userid"]),
					fmt.Sprintf("%v", u["username"]),
					fmt.Sprintf("%v", u["name"]),
					fmt.Sprintf("%v", u["surname"]),
					fmt.Sprintf("%v", u["roleid"]),
				})
			}

			outputResult(cmd, users, headers, rows)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 100, "Limit the number of users")
	cmd.Flags().StringVarP(&search, "search", "s", "", "Search for a user by username")

	return cmd
}

func newUserShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show [username]",
		Short: "Show details of a Zabbix user",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"filter": map[string]interface{}{
					"username": args[0],
				},
				"selectUsrgrps": "extend",
				"selectMedias":  "extend",
			}

			result, err := client.Call("user.get", params)
			handleError(err)

			var users []map[string]interface{}
			json.Unmarshal(result, &users)

			if len(users) == 0 {
				fmt.Printf("User not found: %s\n", args[0])
				return
			}

			u := users[0]
			fmt.Printf("ID:         %s\n", u["userid"])
			fmt.Printf("Username:   %s\n", u["username"])
			fmt.Printf("Name:       %s %s\n", u["name"], u["surname"])
			fmt.Printf("Role ID:    %s\n", u["roleid"])

			fmt.Println("\nGroups:")
			if groups, ok := u["usrgrps"].([]interface{}); ok {
				for _, g := range groups {
					group := g.(map[string]interface{})
					fmt.Printf("- %s (%s)\n", group["name"], group["usrgrpid"])
				}
			}

			fmt.Println("\nMedia:")
			if medias, ok := u["medias"].([]interface{}); ok {
				for _, m := range medias {
					media := m.(map[string]interface{})
					fmt.Printf("- %s (Type ID: %s)\n", media["sendto"], media["mediatypeid"])
				}
			}
		},
	}
}

func newUserCreateCmd() *cobra.Command {
	var password string
	var roleID string
	var groupID string

	cmd := &cobra.Command{
		Use:   "create [username]",
		Short: "Create a new Zabbix user",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"username": args[0],
				"passwd":   password,
				"roleid":   roleID,
				"usrgrps": []map[string]string{
					{"usrgrpid": groupID},
				},
			}

			result, err := client.Call("user.create", params)
			handleError(err)

			var resp map[string]interface{}
			json.Unmarshal(result, &resp)
			userIDs := resp["userids"].([]interface{})

			fmt.Printf("User created successfully with ID: %s\n", userIDs[0])
		},
	}

	cmd.Flags().StringVarP(&password, "password", "p", "", "Password for the user")
	cmd.Flags().StringVarP(&roleID, "roleid", "r", "1", "Role ID for the user (default: 1 - User)")
	cmd.Flags().StringVarP(&groupID, "groupid", "g", "", "User group ID for the user")
	cmd.MarkFlagRequired("password")
	cmd.MarkFlagRequired("groupid")

	return cmd
}

func newUserDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [username]",
		Short: "Delete a Zabbix user",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			// First find the user ID
			params := map[string]interface{}{
				"filter": map[string]interface{}{
					"username": args[0],
				},
			}

			result, err := client.Call("user.get", params)
			handleError(err)

			var users []map[string]interface{}
			json.Unmarshal(result, &users)

			if len(users) == 0 {
				fmt.Printf("User not found: %s\n", args[0])
				return
			}

			userID := users[0]["userid"].(string)

			// Delete the user
			_, err = client.Call("user.delete", []string{userID})
			handleError(err)

			fmt.Printf("User %s (ID: %s) deleted successfully\n", args[0], userID)
		},
	}
}
