package main

import (
    "log"
    "os"
    "github.com/joho/godotenv"
    "plateforme-mycli/cmd" // Ensure this is the correct path
)

func main() {
    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found, using environment variables")
    }

    rootCmd := cmd.NewRootCommand()
    if err := rootCmd.Execute(); err != nil {
        log.Fatalf("Error executing command: %v", err)
        os.Exit(1)
    }
}
