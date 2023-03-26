package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gkawamoto/gsheet-cli/commands/shared"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var Command = &cobra.Command{
	Use:   "auth",
	Short: "Checks whether your cli is authenticated correctly",
	RunE:  commandRunE,
}

func commandRunE(cmd *cobra.Command, args []string) error {
	client := cmd.Context().Value(shared.ClientContextKey)
	if client == nil {
		return fmt.Errorf("could not authenticate")
	}

	return nil
}

func GetClient(ctx context.Context, configDir string, config *oauth2.Config) (*http.Client, error) {
	tokFile := filepath.Join(configDir, "token.json")

	tok, err := tokenFromFile(tokFile)
	if err != nil {
		if tok, err = getTokenFromWeb(ctx, config); err != nil {
			return nil, fmt.Errorf("error retrieving token from web: %w", err)
		}

		if err := saveToken(tokFile, tok); err != nil {
			return nil, fmt.Errorf("error saving token: %w", err)
		}
	}
	return config.Client(ctx, tok), nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	config.RedirectURL = "http://localhost:8097/"

	server := &http.Server{
		Addr: ":8097",
	}

	authCode := ""
	const configState = "state-token"

	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msg := "You can close this window now."

		state, code := r.URL.Query().Get("state"), r.URL.Query().Get("code")
		if state != configState {
			msg = "invalid state, please retry"
		}

		if code != "" {
			authCode = code
		}

		w.Write([]byte(msg))

		go func() {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			server.Shutdown(ctx)
		}()
	})

	authURL := config.AuthCodeURL(configState, oauth2.AccessTypeOnline)

	log.Printf("trying to open your browser pointed to %s", authURL)

	if err := openBrowser(ctx, authURL); err != nil {
		log.Println(err)
	}

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		server.Shutdown(ctx)
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return nil, fmt.Errorf("error starting server: %v", err)
	}

	tok, err := config.Exchange(ctx, authCode)
	if err != nil {
		return nil, fmt.Errorf("error retrieving token from web: %v", err)
	}
	return tok, nil
}

func openBrowser(ctx context.Context, url string) error {
	openPath, err := figureOutOpenCommand(ctx)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, openPath, url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func figureOutOpenCommand(ctx context.Context) (string, error) {
	// Try to use the `open` command on macOS.
	if _, err := exec.LookPath("open"); err == nil {
		return "open", nil
	}

	// Try to use the `xdg-open` command on Linux.
	if _, err := exec.LookPath("xdg-open"); err == nil {
		return "xdg-open", nil
	}

	// Try to use the `start` command on Windows.
	if _, err := exec.LookPath("start"); err == nil {
		return "start", nil
	}

	return "", errors.New("could not find a way to open a browser")
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	return tok, json.NewDecoder(f).Decode(tok)
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("error caching oauth token: %w", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}
