package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the version of BuckGo",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("BuckGo 0.1")
	},
}
