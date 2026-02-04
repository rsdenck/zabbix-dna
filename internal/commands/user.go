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

	cmd.AddCommand(newUserListCmd())    // show_users -> user list
	cmd.AddCommand(newUserShowCmd())    // show_user -> user show
	cmd.AddCommand(newUserCreateCmd())  // create_user -> user create
	cmd.AddCommand(newUserUpdateCmd())  // update_user -> user update
	cmd.AddCommand(newUserDeleteCmd())  // remove_user -> user delete
	cmd.AddCommand(newUserEnableCmd())  // enable_user -> user enable
	cmd.AddCommand(newUserDisableCmd()) // disable_user -> user disable

	return cmd
}

func newUserListCmd() *cobra.Command {
	var limit int
	var search string

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"show_users"},
		Short:   "List Zabbix users",
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
		Use:     "show [username]",
		Aliases: []string{"show_user"},
		Short:   "Show details of a Zabbix user",
		Args:    cobra.ExactArgs(1),
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
			headers := []string{"Property", "Value"}
			var rows [][]string

			rows = append(rows, []string{"UserID", fmt.Sprintf("%v", u["userid"])})
			rows = append(rows, []string{"Username", fmt.Sprintf("%v", u["username"])})
			rows = append(rows, []string{"Name", fmt.Sprintf("%v %v", u["name"], u["surname"])})
			rows = append(rows, []string{"Role ID", fmt.Sprintf("%v", u["roleid"])})

			if groups, ok := u["usrgrps"].([]interface{}); ok {
				for i, g := range groups {
					group := g.(map[string]interface{})
					label := "Group"
					if i > 0 {
						label = ""
					}
					rows = append(rows, []string{label, fmt.Sprintf("%v", group["name"])})
				}
			}

			outputResult(cmd, u, headers, rows)
		},
	}
}

func newUserCreateCmd() *cobra.Command {
	var password string
	var roleID string
	var groupID string

	cmd := &cobra.Command{
		Use:     "create [username]",
		Aliases: []string{"create_user"},
		Short:   "Create a new Zabbix user",
		Args:    cobra.ExactArgs(1),
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

			headers := []string{"Username", "Action", "Status", "ID"}
			rows := [][]string{{args[0], "Create", "Success", fmt.Sprintf("%v", userIDs[0])}}
			outputResult(cmd, resp, headers, rows)
		},
	}

	cmd.Flags().StringVarP(&password, "password", "p", "", "Password for the user")
	cmd.Flags().StringVarP(&roleID, "roleid", "r", "1", "Role ID for the user (default: 1 - User)")
	cmd.Flags().StringVarP(&groupID, "groupid", "g", "", "User group ID for the user")
	cmd.MarkFlagRequired("password")
	cmd.MarkFlagRequired("groupid")

	return cmd
}

func newUserUpdateCmd() *cobra.Command {
	var name string
	var surname string
	return &cobra.Command{
		Use:     "update [username]",
		Aliases: []string{"update_user"},
		Short:   "Update a Zabbix user",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			// Find user
			res, err := client.Call("user.get", map[string]interface{}{
				"filter": map[string]interface{}{"username": args[0]},
			})
			handleError(err)
			var users []map[string]interface{}
			json.Unmarshal(res, &users)
			if len(users) == 0 {
				fmt.Println("User not found")
				return
			}
			userID := users[0]["userid"].(string)

			params := map[string]interface{}{"userid": userID}
			if name != "" {
				params["name"] = name
			}
			if surname != "" {
				params["surname"] = surname
			}

			resp, err := client.Call("user.update", params)
			handleError(err)

			var updateResp map[string]interface{}
			json.Unmarshal(resp, &updateResp)

			headers := []string{"Username", "Action", "Status"}
			rows := [][]string{{args[0], "Update", "Success"}}
			outputResult(cmd, updateResp, headers, rows)
		},
	}
}

func newUserDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete [username]",
		Aliases: []string{"remove_user"},
		Short:   "Delete a Zabbix user",
		Args:    cobra.ExactArgs(1),
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
			resp, err := client.Call("user.delete", []string{userID})
			handleError(err)

			var deleteResp map[string]interface{}
			json.Unmarshal(resp, &deleteResp)

			headers := []string{"Username", "Action", "Status"}
			rows := [][]string{{args[0], "Delete", "Success"}}
			outputResult(cmd, deleteResp, headers, rows)
		},
	}
}

func newUserEnableCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "enable [username]",
		Aliases: []string{"enable_user"},
		Short:   "Enable a Zabbix user",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			// Find user and their groups
			res, err := client.Call("user.get", map[string]interface{}{
				"filter":        map[string]interface{}{"username": args[0]},
				"selectUsrgrps": "extend",
			})
			handleError(err)
			var users []map[string]interface{}
			json.Unmarshal(res, &users)
			if len(users) == 0 {
				fmt.Println("User not found")
				return
			}
			user := users[0]
			var lastResp map[string]interface{}
			if groups, ok := user["usrgrps"].([]interface{}); ok {
				for _, g := range groups {
					group := g.(map[string]interface{})
					groupID := group["usrgrpid"].(string)

					// Enable user group (status 0 = enabled)
					resp, err := client.Call("usergroup.update", map[string]interface{}{
						"usrgrpid":     groupID,
						"users_status": "0",
					})
					handleError(err)
					json.Unmarshal(resp, &lastResp)
				}
			}

			headers := []string{"Username", "Action", "Status", "Note"}
			rows := [][]string{{args[0], "Enable User", "Success", "Enabled via user groups"}}
			outputResult(cmd, lastResp, headers, rows)
		},
	}
}

func newUserDisableCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "disable [username]",
		Aliases: []string{"disable_user"},
		Short:   "Disable a Zabbix user",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			// Find user and their groups
			res, err := client.Call("user.get", map[string]interface{}{
				"filter":        map[string]interface{}{"username": args[0]},
				"selectUsrgrps": "extend",
			})
			handleError(err)
			var users []map[string]interface{}
			json.Unmarshal(res, &users)
			if len(users) == 0 {
				fmt.Println("User not found")
				return
			}
			user := users[0]
			var lastResp map[string]interface{}
			if groups, ok := user["usrgrps"].([]interface{}); ok {
				for _, g := range groups {
					group := g.(map[string]interface{})
					groupID := group["usrgrpid"].(string)

					// Disable user group (status 1 = disabled)
					resp, err := client.Call("usergroup.update", map[string]interface{}{
						"usrgrpid":     groupID,
						"users_status": "1",
					})
					handleError(err)
					json.Unmarshal(resp, &lastResp)
				}
			}

			headers := []string{"Username", "Action", "Status", "Note"}
			rows := [][]string{{args[0], "Disable User", "Success", "Disabled via user groups"}}
			outputResult(cmd, lastResp, headers, rows)
		},
	}
}


