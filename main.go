package main

import (
	"os"
	"time"

	"bober/build"
	"bober/config"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

// Setting up a basic charmbracelet logger style
func setupLogger() *log.Logger {
	styles := log.DefaultStyles()
	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("ERROR!!").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("204")).
		Foreground(lipgloss.Color("0"))
	styles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	styles.Values["err"] = lipgloss.NewStyle().Bold(true)

	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("INFO").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("204")).
		Foreground(lipgloss.Color("0"))

	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    false,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
	})
	logger.SetStyles(styles)
	return logger
}

func main() {
	logger := setupLogger()

	// Load the configuration from the YAML file
	cfg, err := config.Load("project.yaml")
	if err != nil {
		logger.Errorf("Loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Generate the Makefile
	err = build.GenerateMakefile(cfg, logger)
	if err != nil {
		logger.Errorf("Generating Makefile: %v\n", err)
		os.Exit(1)
	}

	// Run `make` with provided arguments if any
	if len(os.Args) > 1 {
		err = build.RunMake(os.Args[1:], logger)
		if err != nil {
			logger.Errorf("Running make: %v\n", err)
			os.Exit(1)
		}
	} else {
		logger.Info("Makefile generated successfully.")
	}
}
