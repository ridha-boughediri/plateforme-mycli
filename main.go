package main

import (
	"github.com/ridha-boughediri/plateforme-mycli/commands"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "bu"}

	rootCmd.AddCommand(commands.VersionCmd)
	rootCmd.AddCommand(commands.SyncCmd)
	rootCmd.AddCommand(commands.ShowCmd)
	rootCmd.AddCommand(commands.UnsyncCmd)
	rootCmd.AddCommand(commands.ObjectAddCmd)
	rootCmd.AddCommand(commands.ObjectDownloadCmd)
	rootCmd.AddCommand(commands.ObjectDeleteCmd)
	rootCmd.AddCommand(commands.BucketAddCmd)
	rootCmd.AddCommand(commands.BucketListCmd)
	rootCmd.AddCommand(commands.BucketDeleteCmd)

	rootCmd.Execute()
}
