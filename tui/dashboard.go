package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Leela0o5/LeeGo/config"
	"github.com/Leela0o5/LeeGo/metrics"
)

type tickMsg time.Time
type doneMsg struct{}

var (
	purple  = lipgloss.Color("#9B59B6")
	cyan    = lipgloss.Color("#00D9FF")
	green   = lipgloss.Color("#2ECC71")
	red     = lipgloss.Color("#E74C3C")
	white   = lipgloss.Color("#ECF0F1")
	darkBg  = lipgloss.Color("#1A1A2E")
	panelBg = lipgloss.Color("#16213E")
	subtle  = lipgloss.Color("#7F8C8D")

	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(cyan).Background(darkBg).Padding(0, 2)
	boxStyle     = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(purple).Background(panelBg).Padding(0, 2)
	labelStyle   = lipgloss.NewStyle().Foreground(subtle).Width(14)
	valueStyle   = lipgloss.NewStyle().Foreground(white).Bold(true)
	successStyle = lipgloss.NewStyle().Foreground(green).Bold(true)
	failStyle    = lipgloss.NewStyle().Foreground(red).Bold(true)
	accentStyle  = lipgloss.NewStyle().Foreground(cyan).Bold(true)
	dimStyle     = lipgloss.NewStyle().Foreground(subtle)
)

type Model struct {
	cfg      config.Config
	stats    *metrics.Stats
	done     chan struct{}
	finished bool
	elapsed  time.Duration
	start    time.Time
}

func New(cfg config.Config, stats *metrics.Stats, done chan struct{}) Model {
	return Model{
		cfg:   cfg,
		stats: stats,
		done:  done,
		start: time.Now(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tick(), waitForDone(m.done))
}

func tick() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func waitForDone(done chan struct{}) tea.Cmd {
	return func() tea.Msg {
		<-done
		return doneMsg{}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.elapsed = time.Since(m.start)
		return m, tick()
	case doneMsg:
		m.finished = true
		m.elapsed = time.Since(m.start)
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(titleStyle.Render(" ⚡ WebSocket Load Tester ") + "\n\n")
	b.WriteString(m.configPanel() + "\n\n")
	b.WriteString(m.statsPanel() + "\n\n")
	b.WriteString(m.latencyPanel() + "\n\n")
	b.WriteString(m.progressBar() + "\n\n")

	if m.finished {
		b.WriteString(successStyle.Render("  ✓ Test complete!") + "\n")
	} else {
		b.WriteString(dimStyle.Render("  press q to quit") + "\n")
	}

	return b.String()
}

func (m Model) configPanel() string {
	rows := strings.Join([]string{
		kv("Target", m.cfg.URL),
		kv("Workers", fmt.Sprintf("%d", m.cfg.NumWorkers)),
		kv("Duration", m.cfg.Duration.String()),
		kv("Rate limit", fmt.Sprintf("%d req/s", m.cfg.RateLimit)),
	}, "\n")
	return boxStyle.Render(accentStyle.Render("Config") + "\n" + rows)
}

func (m Model) statsPanel() string {
	s := m.stats

	successRate := 0.0
	if s.TotalRequests > 0 {
		successRate = float64(s.SuccessCount) / float64(s.TotalRequests) * 100
	}

	rps := 0.0
	if m.elapsed.Seconds() > 0 {
		rps = float64(s.TotalRequests) / m.elapsed.Seconds()
	}

	rows := strings.Join([]string{
		kv("Total", fmt.Sprintf("%d", s.TotalRequests)),
		kv("Success", successStyle.Render(fmt.Sprintf("%d", s.SuccessCount))),
		kv("Failed", failStyle.Render(fmt.Sprintf("%d", s.FailureCount))),
		kv("Success %", fmt.Sprintf("%.1f%%", successRate)),
		kv("Throughput", accentStyle.Render(fmt.Sprintf("%.0f req/s", rps))),
		kv("Elapsed", m.elapsed.Round(time.Millisecond).String()),
	}, "\n")

	return boxStyle.Render(accentStyle.Render("Live Stats") + "\n" + rows)
}

func (m Model) latencyPanel() string {
	s := m.stats

	if s.SuccessCount == 0 {
		return boxStyle.Render(accentStyle.Render("Latency") + "\n" + dimStyle.Render("  waiting for data..."))
	}

	rows := strings.Join([]string{
		kv("Avg", fmtDur(s.Average())),
		kv("Min", fmtDur(s.Min())),
		kv("Max", fmtDur(s.Max())),
		kv("P50", fmtDur(s.Median())),
		kv("P95", fmtDur(s.P95())),
		kv("P99", fmtDur(s.P99())),
	}, "\n")

	return boxStyle.Render(accentStyle.Render("Latency") + "\n" + rows)
}

func (m Model) progressBar() string {
	pct := m.elapsed.Seconds() / m.cfg.Duration.Seconds()
	if pct > 1 {
		pct = 1
	}

	width := 40
	filled := int(pct * float64(width))

	bar := lipgloss.NewStyle().Foreground(purple).Render(strings.Repeat("█", filled)) +
		lipgloss.NewStyle().Foreground(subtle).Render(strings.Repeat("░", width-filled))

	label := fmt.Sprintf("  %.0f%%  %s / %s", pct*100, m.elapsed.Round(time.Second), m.cfg.Duration)

	return "  [" + bar + "] " + dimStyle.Render(label)
}

func kv(label, val string) string {
	return labelStyle.Render(label+":") + " " + valueStyle.Render(val)
}

func fmtDur(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%.2f µs", float64(d.Microseconds()))
	}
	return fmt.Sprintf("%.2f ms", float64(d.Microseconds())/1000)
}
