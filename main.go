package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"maily/internal/auth"
	"maily/internal/ui"
)

var rootCmd = &cobra.Command{
	Use:   "maily",
	Short: "A terminal email client",
	Long:  "maily - A terminal email client for Gmail",
	Run: func(cmd *cobra.Command, args []string) {
		runTUI()
	},
}

var loginCmd = &cobra.Command{
	Use:   "login [provider]",
	Short: "Add an email account",
	Long:  "Add an email account. Currently supports: gmail",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleLogin(args[0])
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout [email]",
	Short: "Remove an account",
	Long:  "Remove an email account. If no email specified, prompts for selection.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			handleLogoutAccount(args[0])
		} else {
			handleLogout()
		}
	},
}

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "List all accounts",
	Run: func(cmd *cobra.Command, args []string) {
		handleAccounts()
	},
}

var (
	searchAccount string
	searchQuery   string
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search emails with Gmail query syntax",
	Long: `Search emails using Gmail's search syntax.

Gmail search syntax examples:
  from:sender@example.com    Emails from a sender
  subject:hello              Emails with subject containing 'hello'
  has:attachment             Emails with attachments
  is:unread                  Unread emails
  older_than:30d             Emails older than 30 days
  category:promotions        Promotional emails
  larger:5M                  Emails larger than 5MB`,
	Example: `  maily search -a me@gmail.com -q "category:promotions older_than:30d"
  maily search --account=me@gmail.com --query="is:unread"`,
	Run: func(cmd *cobra.Command, args []string) {
		handleSearch()
	},
}

func init() {
	searchCmd.Flags().StringVarP(&searchAccount, "account", "a", "", "Account email to search")
	searchCmd.Flags().StringVarP(&searchQuery, "query", "q", "", "Gmail search query (uses Gmail syntax)")
	searchCmd.MarkFlagRequired("query")

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(accountsCmd)
	rootCmd.AddCommand(searchCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runTUI() {
	store, err := auth.LoadAccountStore()
	if err != nil {
		fmt.Printf("Error loading accounts: %v\n", err)
		os.Exit(1)
	}

	if len(store.Accounts) == 0 {
		fmt.Println("No accounts configured. Run:")
		fmt.Println()
		fmt.Println("  maily login gmail")
		fmt.Println()
		os.Exit(1)
	}

	p := tea.NewProgram(
		ui.NewApp(store),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

func handleLogin(provider string) {
	switch provider {
	case "gmail":
		loginGmail()
	default:
		fmt.Printf("Unknown provider: %s\n", provider)
		fmt.Println()
		fmt.Println("Available providers:")
		fmt.Println("  gmail    Login with Gmail")
		os.Exit(1)
	}
}

func loginGmail() {
	p := tea.NewProgram(
		ui.NewLoginApp("gmail"),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func handleLogoutAccount(email string) {
	store, err := auth.LoadAccountStore()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if store.RemoveAccount(email) {
		store.Save()
		fmt.Printf("Removed account %s\n", email)
	} else {
		fmt.Printf("Account not found: %s\n", email)
	}
}

func handleLogout() {
	store, err := auth.LoadAccountStore()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if len(store.Accounts) == 0 {
		fmt.Println("No accounts configured.")
		return
	}

	// If only one account, remove it
	if len(store.Accounts) == 1 {
		email := store.Accounts[0].Credentials.Email
		store.RemoveAccount(email)
		store.Save()
		fmt.Printf("Removed account %s\n", email)
		return
	}

	// Multiple accounts - show list
	fmt.Println()
	fmt.Println("  Which account to remove?")
	fmt.Println()
	for i, acc := range store.Accounts {
		fmt.Printf("  %d. %s\n", i+1, acc.Credentials.Email)
	}
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("  Enter number (or 0 to cancel): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	num, err := strconv.Atoi(input)
	if err != nil || num < 0 || num > len(store.Accounts) {
		fmt.Println("Cancelled.")
		return
	}

	if num == 0 {
		fmt.Println("Cancelled.")
		return
	}

	email := store.Accounts[num-1].Credentials.Email
	store.RemoveAccount(email)
	store.Save()
	fmt.Printf("Removed account %s\n", email)
}

func handleAccounts() {
	store, err := auth.LoadAccountStore()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if len(store.Accounts) == 0 {
		fmt.Println("No accounts configured.")
		fmt.Println()
		fmt.Println("Run: maily login gmail")
		return
	}

	fmt.Println()
	fmt.Println("  Accounts:")
	fmt.Println()
	for _, acc := range store.Accounts {
		fmt.Printf("  %s (%s)\n", acc.Credentials.Email, acc.Provider)
	}
	fmt.Println()
}

func handleSearch() {
	if searchQuery == "" {
		fmt.Println("Error: --query (-q) is required")
		os.Exit(1)
	}

	store, err := auth.LoadAccountStore()
	if err != nil {
		fmt.Printf("Error loading accounts: %v\n", err)
		os.Exit(1)
	}

	if len(store.Accounts) == 0 {
		fmt.Println("No accounts configured. Run:")
		fmt.Println()
		fmt.Println("  maily login gmail")
		fmt.Println()
		os.Exit(1)
	}

	// Find the account
	var account *auth.Account
	if searchAccount == "" {
		if len(store.Accounts) == 1 {
			account = &store.Accounts[0]
		} else {
			fmt.Println("Error: --account (-a) is required when multiple accounts are configured")
			fmt.Println()
			fmt.Println("Available accounts:")
			for _, acc := range store.Accounts {
				fmt.Printf("  - %s\n", acc.Credentials.Email)
			}
			os.Exit(1)
		}
	} else {
		account = store.GetAccount(searchAccount)
		if account == nil {
			fmt.Printf("Error: account '%s' not found\n", searchAccount)
			fmt.Println()
			fmt.Println("Available accounts:")
			for _, acc := range store.Accounts {
				fmt.Printf("  - %s\n", acc.Credentials.Email)
			}
			os.Exit(1)
		}
	}

	// Run the search TUI
	p := tea.NewProgram(
		ui.NewSearchApp(account, searchQuery),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running search: %v\n", err)
		os.Exit(1)
	}
}
