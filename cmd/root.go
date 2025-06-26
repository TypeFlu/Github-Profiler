package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github-profiler/internal/ui"
)

var (
	githubToken  string
	outputFormat string
	version      = "1.0.0"
	author       = "github@Tyeflu"
)

var rootCmd = &cobra.Command{
	Use:   "github-profiler [username]",
	Short: "A beautiful CLI tool to analyze GitHub user profiles",
	Long: `GitHub Profiler is a CLI tool that provides comprehensive analysis
of any GitHub user's public profile with beautiful visualizations,
advanced statistics, and ranking information.

Use 'github-profiler demo' to see the tool in action with sample data.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runProfiler,
}

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Run a demo with sample data",
	Long:  `Demonstrates the GitHub Profiler output using mock data.`,
	Run: func(cmd *cobra.Command, args []string) {
		runTUI("demo-user")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&githubToken, "token", "t", "", "GitHub personal access token (optional for public data)")
	rootCmd.Flags().StringVarP(&outputFormat, "format", "f", "tui", "Output format: tui, json, html")

	rootCmd.AddCommand(demoCmd)
	demoCmd.Flags().StringVarP(&outputFormat, "format", "f", "tui", "Output format: tui, json, html")

	if githubToken == "" {
		githubToken = os.Getenv("GITHUB_TOKEN")
	}
}

func runProfiler(cmd *cobra.Command, args []string) {
	username := ""
	if len(args) > 0 {
		username = args[0]
	}

	// For now, always use TUI mode for the best experience
	// Later we can add back JSON/HTML output modes
	runTUI(username)
}

func runTUI(username string) {
	model := ui.NewModel(username, githubToken, outputFormat)

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
