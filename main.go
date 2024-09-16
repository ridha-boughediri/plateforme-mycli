package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Shows the version of Buckgo",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Buckgo 0.1")
		},
	}

	rootCmd := &cobra.Command{Use: "Buckgo"}

	rootCmd.AddCommand(versionCmd)
	rootCmd.Execute()
}
