package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/Leela0o5/WebSocket-Load-Tester/config"
	"github.com/Leela0o5/WebSocket-Load-Tester/engine"
	"github.com/Leela0o5/WebSocket-Load-Tester/reporter"
	"github.com/Leela0o5/WebSocket-Load-Tester/tui"
)

var (
	cfgPath    string
	outputPath string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start a WebSocket load test",
	Long:  `Executes a load test against the configured WebSocket URL using parameters from the config file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := cfgPath
		if path == "" {
			if _, err := os.Stat("config.yaml"); err == nil {
				path = "config.yaml"
			}
		}

		cfg, err := config.LoadConfig(path)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		stats, done := engine.RunAsync(*cfg)
		m := tui.New(*cfg, stats, done)
		p := tea.NewProgram(m, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("tui error: %w", err)
		}

		<-done

		reporter.PrintSummary(stats)

		if outputPath != "" {
			if err := reporter.SaveJSON(stats, outputPath); err != nil {
				return fmt.Errorf("failed to save report: %w", err)
			}
			fmt.Printf("\nReport saved to %s\n", outputPath)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&cfgPath, "config", "c", "", "path to YAML config file")
	runCmd.Flags().StringVarP(&outputPath, "output", "o", "", "path to save JSON report")
}

