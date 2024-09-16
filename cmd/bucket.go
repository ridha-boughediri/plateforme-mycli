package cmd

import (
    "fmt"
    "plateforme-mycli/s3manager"
    "github.com/spf13/cobra"
)

// NewBucketCommand creates a new 'bucket' command
func NewBucketCommand() *cobra.Command {
    bucketCmd := &cobra.Command{
        Use:   "bucket",
        Short: "Manage S3 buckets",
    }

    bucketCmd.AddCommand(NewCreateBucketCommand())
    bucketCmd.AddCommand(NewListBucketsCommand())

    return bucketCmd
}

// NewCreateBucketCommand creates a subcommand to create a bucket
func NewCreateBucketCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "create [bucketName]",
        Short: "Create a new S3 bucket",
        Args:  cobra.ExactArgs(1), // Requires exactly one argument (the bucket name)
        RunE: func(cmd *cobra.Command, args []string) error {
            bucketName := args[0]
            if err := s3manager.CreateBucket(bucketName); err != nil {
                return fmt.Errorf("failed to create bucket: %v", err)
            }
            fmt.Printf("Bucket %s created successfully\n", bucketName)
            return nil
        },
    }
}

// NewListBucketsCommand creates a subcommand to list all buckets
func NewListBucketsCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "list",
        Short: "List all S3 buckets",
        RunE: func(cmd *cobra.Command, args []string) error {
            if err := s3manager.ListBuckets(); err != nil {
                return fmt.Errorf("failed to list buckets: %v", err)
            }
            return nil
        },
    }
}
