package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/gkawamoto/gsheet-cli/commands/auth"
	"github.com/gkawamoto/gsheet-cli/commands/get"
	"github.com/gkawamoto/gsheet-cli/commands/shared"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2/google"
)

var rootCmd = &cobra.Command{
	Use: "gsheet",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		configDir := cmd.Flag("config-dir").Value.String()

		if err := os.MkdirAll(configDir, 0700); err != nil {
			return fmt.Errorf("error creating config directory: %w", err)
		}

		credentialsFile := cmd.Flag("credentials-file").Value.String()

		credentialsBytes, err := os.ReadFile(credentialsFile)
		if err != nil {
			return fmt.Errorf("error reading client secret file: %w", err)
		}

		conf, err := google.ConfigFromJSON(credentialsBytes, "https://www.googleapis.com/auth/spreadsheets.readonly")
		if err != nil {
			return fmt.Errorf("error parsing client secret file to config: %w", err)
		}

		ctx = context.WithValue(ctx, shared.ConfigContextKey, conf)

		client, err := auth.GetClient(ctx, configDir, conf)
		if err != nil {
			return fmt.Errorf("error retrieving HTTP client: %w", err)
		}

		ctx = context.WithValue(ctx, shared.ClientContextKey, client)

		cmd.SetContext(ctx)

		return nil
	},
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	rootCmd.PersistentFlags().StringP("config-dir", "d", filepath.Join(os.Getenv("HOME"), ".config", "gsheet"), "config directory")
	rootCmd.PersistentFlags().StringP("credentials-file", "c", filepath.Join(os.Getenv("HOME"), ".config", "gsheet", "credentials.json"), "credentials file location")

	rootCmd.AddCommand(auth.Command)
	rootCmd.AddCommand(get.Command)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}
}
