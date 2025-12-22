package cli

import (
	"github.com/spf13/cobra"
	"maily/internal/updater"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update maily to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		return updater.Update()
	},
}
