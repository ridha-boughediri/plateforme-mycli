package cmd

import (
    "github.com/spf13/cobra"
)

// NewRootCommand sets up the root command
func NewRootCommand() *cobra.Command {
    rootCmd := &cobra.Command{
        Use:   "plateforme-mycli",
        Short: "CLI tool to interact with S3",
    }

    // Add the bucket command
    rootCmd.AddCommand(NewBucketCommand()) // This function must be defined in bucket.go

    return rootCmd
}
