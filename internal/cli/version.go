package cli

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"maily/internal/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		v := version.Version
		if len(v) > 0 && v[0] == 'v' {
			v = v[1:]
		}
		fmt.Printf("maily v%s %s/%s\n", v, runtime.GOOS, runtime.GOARCH)
		if version.Commit != "unknown" {
			fmt.Printf("commit: %s\n", version.Commit)
		}
		if version.Date != "unknown" {
			fmt.Printf("built:  %s\n", version.Date)
		}
	},
}
