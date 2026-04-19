package cli

import (
	"fmt"

	"github.com/Leela0o5/LeeGo/reporter"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report [file.json]",
	Short: "Displays a summary from a saved JSON report",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		report, err := reporter.LoadJSON(path)
		if err != nil {
			return fmt.Errorf("read report: %w", err)
		}

		fmt.Printf("Loaded Report: %s\n", path)
		reporter.PrintReport(report)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}
