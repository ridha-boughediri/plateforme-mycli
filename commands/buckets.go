package commands

import (
	"fmt"

	"github.com/ridha-boughediri/plateforme-mycli/handlers"
	"github.com/spf13/cobra"
)

var BucketAddCmd = &cobra.Command{
	Use:   "ba [alias/bucketName]",
	Short: "Create a bucket",
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		if err := handlers.AddBucket(url); err != nil {
			fmt.Printf("Error creating bucket: %v\n", err)
			return
		}

		fmt.Println("Bucket created successfully!")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("ba requires 1 argument: [alias/bucketName]\n\nExample: bu ba myalias/mybucket")
		}
		return nil
	},
}

var BucketListCmd = &cobra.Command{
	Use:   "bl [alias/bucketName]",
	Short: "List all buckets",
	Run: func(cmd *cobra.Command, args []string) {
		alias := args[0]

		err := handlers.ListBuckets(alias)
		if err != nil {
			fmt.Printf("Error listing buckets: %v\n", err)
			return
		}

	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("bl requires 1 argument: [alias]\n\nExample: bu bl myalias")
		}
		return nil
	},
}

var BucketDeleteCmd = &cobra.Command{
	Use:   "br [alias/bucketName]",
	Short: "Delete a bucket",
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		if err := handlers.DeleteBucket(url); err != nil {
			fmt.Printf("Error deleting bucket: %v\n", err)
			return
		}

		fmt.Println("Bucket deleted successfully!")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("br requires 1 argument: [alias/bucketName]\n\nExample: bu br myalias/mybucket")
		}
		return nil
	},
}
