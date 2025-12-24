package cli

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"maily/internal/auth"
	"maily/internal/cache"
	"maily/internal/sync"
)

const (
	syncInterval = 30 * time.Minute
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Manage the background sync daemon",
	Long:  "Start or stop the background sync daemon that keeps your email cache up to date",
}

var daemonStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the background sync daemon",
	Run: func(cmd *cobra.Command, args []string) {
		foreground, _ := cmd.Flags().GetBool("foreground")
		if foreground {
			runDaemon()
		} else {
			startDaemonBackground()
		}
	},
}

var daemonStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the background sync daemon",
	Run: func(cmd *cobra.Command, args []string) {
		stopDaemon()
	},
}

var daemonStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check daemon status",
	Run: func(cmd *cobra.Command, args []string) {
		checkDaemonStatus()
	},
}

func init() {
	daemonStartCmd.Flags().BoolP("foreground", "f", false, "Run in foreground (for debugging)")
	daemonCmd.AddCommand(daemonStartCmd)
	daemonCmd.AddCommand(daemonStopCmd)
	daemonCmd.AddCommand(daemonStatusCmd)
	rootCmd.AddCommand(daemonCmd)
}

func getDaemonPidFile() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "maily", "daemon.pid")
}

func startDaemonBackground() {
	// Check if already running
	pidFile := getDaemonPidFile()
	if data, err := os.ReadFile(pidFile); err == nil {
		pid, _ := strconv.Atoi(string(data))
		if pid > 0 {
			if process, err := os.FindProcess(pid); err == nil {
				if err := process.Signal(syscall.Signal(0)); err == nil {
					fmt.Println("Daemon is already running (PID:", pid, ")")
					return
				}
			}
		}
	}

	// Start daemon in background
	executable, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		os.Exit(1)
	}

	cmd := exec.Command(executable, "daemon", "start", "--foreground")
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	// Detach from parent process
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting daemon:", err)
		os.Exit(1)
	}

	// Write PID file
	if err := os.MkdirAll(filepath.Dir(pidFile), 0700); err == nil {
		os.WriteFile(pidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0600)
	}

	fmt.Println("Daemon started (PID:", cmd.Process.Pid, ")")
}

func stopDaemon() {
	pidFile := getDaemonPidFile()
	data, err := os.ReadFile(pidFile)
	if err != nil {
		fmt.Println("Daemon is not running")
		return
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil || pid <= 0 {
		fmt.Println("Invalid PID file")
		os.Remove(pidFile)
		return
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("Daemon process not found")
		os.Remove(pidFile)
		return
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		fmt.Println("Daemon is not running")
		os.Remove(pidFile)
		return
	}

	os.Remove(pidFile)
	fmt.Println("Daemon stopped")
}

func checkDaemonStatus() {
	pidFile := getDaemonPidFile()
	data, err := os.ReadFile(pidFile)
	if err != nil {
		fmt.Println("Daemon is not running")
		return
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil || pid <= 0 {
		fmt.Println("Daemon is not running (invalid PID file)")
		return
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("Daemon is not running")
		return
	}

	if err := process.Signal(syscall.Signal(0)); err != nil {
		fmt.Println("Daemon is not running")
		os.Remove(pidFile)
		return
	}

	fmt.Println("Daemon is running (PID:", pid, ")")
}

func runDaemon() {
	// Write PID file
	pidFile := getDaemonPidFile()
	if err := os.MkdirAll(filepath.Dir(pidFile), 0700); err == nil {
		os.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0600)
	}
	defer os.Remove(pidFile)

	// Load accounts
	store, err := auth.LoadAccountStore()
	if err != nil {
		fmt.Println("Error loading accounts:", err)
		os.Exit(1)
	}

	if len(store.Accounts) == 0 {
		fmt.Println("No accounts configured")
		os.Exit(1)
	}

	// Create cache
	c, err := cache.New()
	if err != nil {
		fmt.Println("Error creating cache:", err)
		os.Exit(1)
	}

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Initial sync
	syncAllAccounts(store, c)

	// Ticker for periodic sync
	ticker := time.NewTicker(syncInterval)
	defer ticker.Stop()

	fmt.Println("Daemon started, syncing every", syncInterval)

	for {
		select {
		case <-ticker.C:
			syncAllAccounts(store, c)
		case sig := <-sigChan:
			fmt.Println("Received signal:", sig)
			return
		}
	}
}

func syncAllAccounts(store *auth.AccountStore, c *cache.Cache) {
	for i := range store.Accounts {
		account := &store.Accounts[i]
		syncer := sync.NewSyncer(c, account)

		// Sync INBOX
		if err := syncer.FullSync("INBOX"); err != nil {
			fmt.Printf("Error syncing %s: %v\n", account.Credentials.Email, err)
		} else {
			fmt.Printf("Synced %s\n", account.Credentials.Email)
		}
	}
}
