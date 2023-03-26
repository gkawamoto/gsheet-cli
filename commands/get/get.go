package get

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gkawamoto/gsheet-cli/commands/shared"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func init() {
	Command.PersistentFlags().StringP("spreadsheet-id", "s", "", "spreadsheet ID")
}

var Command = &cobra.Command{
	Use:     "get",
	Short:   "Get data from a spreadsheet",
	PreRunE: commandPreRunE,
	RunE:    commandRunE,
}

func commandPreRunE(cmd *cobra.Command, args []string) error {
	if cmd.Flag("spreadsheet-id").Value.String() == "" {
		return fmt.Errorf("spreadsheet ID is required")
	}

	if len(args) == 0 {
		return fmt.Errorf("range is required")
	}

	return nil
}

func commandRunE(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	spreadsheetID := cmd.Flag("spreadsheet-id").Value.String()

	return get(ctx, spreadsheetID, args, cmd.OutOrStdout())
}

func get(ctx context.Context, spreadsheetID string, ranges []string, writer io.Writer) error {
	client := ctx.Value(shared.ClientContextKey).(*http.Client)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("error retrieving Sheets client: %w", err)
	}

	encoder := json.NewEncoder(writer)

	for _, r := range ranges {
		resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, r).Do()
		if err != nil {
			return fmt.Errorf("error retrieving data from sheet: %w", err)
		}

		if err := encoder.Encode(resp.Values); err != nil {
			return fmt.Errorf("error encoding JSON: %w", err)
		}
	}

	return nil
}
