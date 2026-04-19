package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/Leela0o5/LeeGo/config"
	"github.com/Leela0o5/LeeGo/engine"
	"github.com/Leela0o5/LeeGo/reporter"
	"github.com/Leela0o5/LeeGo/tui"
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

		if cmd.Flags().Changed("url") {
			cfg.URL, _ = cmd.Flags().GetString("url")
		}
		if cmd.Flags().Changed("connections") {
			cfg.NumWorkers, _ = cmd.Flags().GetInt("connections")
		}
		if cmd.Flags().Changed("duration") {
			cfg.Duration, _ = cmd.Flags().GetDuration("duration")
		}
		if cmd.Flags().Changed("rate") {
			cfg.RateLimit, _ = cmd.Flags().GetInt("rate")
		}
		if cmd.Flags().Changed("burst") {
			cfg.Burst, _ = cmd.Flags().GetInt("burst")
		}
		if cmd.Flags().Changed("message") {
			cfg.Message, _ = cmd.Flags().GetString("message")
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

	runCmd.Flags().String("url", "", "WebSocket connection URL")
	runCmd.Flags().Int("connections", 0, "Number of concurrent workers")
	runCmd.Flags().Duration("duration", 0, "Test duration (e.g., 10s, 1m)")
	runCmd.Flags().Int("rate", 0, "Rate limit per second")
	runCmd.Flags().Int("burst", 0, "Burst size")
	runCmd.Flags().String("message", "", "Message to send")
}

