package commands

import (
	"fmt"

	"github.com/ridha-boughediri/plateforme-mycli/handlers"
	"github.com/spf13/cobra"
)

var ObjectAddCmd = &cobra.Command{
	Use:   "oa [alias/bucketName] [localPath]",
	Short: "Add an object to a bucket",
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		localPath := args[1]

		if err := handlers.AddObject(url, localPath); err != nil {
			fmt.Printf("Error adding object: %v\n", err)
			return
		}

		fmt.Println("Object added successfully!")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("oa requires 2 arguments: [alias/bucketName] [localPath]\n\nExample: bu oa myalias/myobject /path/to/local/file")
		}
		return nil
	},
}

var ObjectDownloadCmd = &cobra.Command{
	Use:   "od [alias/bucketName/objectName] [localPath]",
	Short: "Download an object from a bucket",
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		localPath := args[1]

		if err := handlers.DownloadObject(url, localPath); err != nil {
			fmt.Printf("Error downloading object: %v\n", err)
			return
		}

		fmt.Println("Object downloaded successfully!")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("od requires 2 arguments: [alias/bucketName/objectName] [localPath]\n\nExample: bu od myalias/mybucket/myobject /path/to/local/file")
		}
		return nil
	},
}

var ObjectDeleteCmd = &cobra.Command{
	Use:   "or [alias/bucketName/objectName]",
	Short: "Delete an object from a bucket",
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		if err := handlers.DeleteObject(url); err != nil {
			fmt.Printf("Error deleting object: %v\n", err)
			return
		}

		fmt.Println("Object deleted successfully!")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("or requires 1 arguments: [alias/bucketName/objectName] \n\nExample: bu or myalias/mybucket/myobject")
		}
		return nil
	},
}
