package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"maily/internal/auth"
	"maily/internal/cache"
	"maily/internal/sync"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync emails from server",
	Long:  "Perform a full sync of emails from the server for all accounts",
	Run: func(cmd *cobra.Command, args []string) {
		runSync()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}

func runSync() {
	// Load accounts
	store, err := auth.LoadAccountStore()
	if err != nil {
		fmt.Println("Error loading accounts:", err)
		os.Exit(1)
	}

	if len(store.Accounts) == 0 {
		fmt.Println("No accounts configured. Run 'maily login' first.")
		os.Exit(1)
	}

	// Create cache
	c, err := cache.New()
	if err != nil {
		fmt.Println("Error creating cache:", err)
		os.Exit(1)
	}

	fmt.Println("Syncing emails...")

	for i := range store.Accounts {
		account := &store.Accounts[i]
		fmt.Printf("  Syncing %s...", account.Credentials.Email)

		syncer := sync.NewSyncer(c, account)
		if err := syncer.FullSync("INBOX"); err != nil {
			fmt.Printf(" error: %v\n", err)
		} else {
			fmt.Println(" done")
		}
	}

	fmt.Println("Sync complete")
}
