package commands

import (
	"fmt"

	"github.com/ridha-boughediri/plateforme-mycli/configs"
	"github.com/ridha-boughediri/plateforme-mycli/handlers"
	"github.com/spf13/cobra"
)

var SyncCmd = &cobra.Command{
	Use:   "sync [alias] [remote] [username] [password]",
	Short: "Sync BuckGo with your cloud storage",
	Run: func(cmd *cobra.Command, args []string) {
		alias := args[0]
		remote := args[1]
		username := args[2]
		password := args[3]

		aliasConfig := configs.AliasConfig{
			Alias:    alias,
			Remote:   remote,
			Username: username,
			Password: password,
		}

		if err := handlers.SaveAlias(aliasConfig); err != nil {
			fmt.Printf("Error saving alias: %v\n", err)
			return
		}

		fmt.Println("Alias saved successfully!")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 4 {
			return fmt.Errorf("sync requires 4 arguments: [alias] [remote] [username] [password]\n\nExample: bu sync myalias https://remote.url user123 pass123")
		}
		return nil
	},
}

var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show all aliases",
	Run: func(cmd *cobra.Command, args []string) {
		aliases, err := handlers.ShowAliases()
		if err != nil {
			fmt.Printf("Error showing aliases: %v\n", err)
			return
		}

		if len(aliases) == 0 {
			fmt.Println("No aliases found")
			return
		}

		for _, alias := range aliases {
			fmt.Printf("Alias: %s\n", alias.Alias)
			fmt.Printf("Remote: %s\n", alias.Remote)
			fmt.Printf("Username: %s\n", alias.Username)
			fmt.Printf("Password: %s\n", alias.Password)
			fmt.Println()
		}
	},
}
var UnsyncCmd = &cobra.Command{
	Use:   "unsync [alias]",
	Short: "Delete an alias",
	Run: func(cmd *cobra.Command, args []string) {
		alias := args[0]

		if err := handlers.DeleteAlias(alias); err != nil {
			fmt.Printf("Error deleting alias: %v\n", err)
			return
		}

		fmt.Println("Alias deleted successfully!")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("unsync requires 1 argument: [alias]\n\nExample: bu unsync myalias")
		}
		return nil
	},
}
