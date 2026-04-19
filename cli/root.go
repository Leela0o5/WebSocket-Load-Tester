package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "loader",
	Short: "A high-performance WebSocket load testing tool",
	Long: `A CLI tool designed for stress testing WebSocket servers 
with precise rate limiting and detailed latency reporting.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Root-level flags
}
