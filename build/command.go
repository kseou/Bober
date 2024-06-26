package build

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
)

// executes a shell command
func runCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// ensures that the build directory exists, creating it if necessary.
func ensureBuildDir(logger *log.Logger) error {
	if err := os.MkdirAll("build", 0755); err != nil {
		logger.Error("Failed to create build directory", "error", err)
		return err
	}
	return nil
}

// serveBuiltFiles starts a local HTTP server to serve files from the build directory
func serveBuiltFiles(logger *log.Logger) error {
	buildDir := "build"
	port := "8080"

	logger.Info("Starting server", "directory", buildDir, "port", port)

	http.Handle("/", http.FileServer(http.Dir(buildDir)))

	logger.Info(fmt.Sprintf("Server is running at http://localhost:%s", port))
	logger.Info("Press Ctrl+C to stop the server")

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logger.Error("Failed to start server", "error", err)
		return err
	}

	return nil
}

// runs the `make` command with the provided arguments.
func RunMake(args []string, logger *log.Logger) error {
	if len(args) > 0 {
		switch args[0] {
		case "build":
			if err := ensureBuildDir(logger); err != nil {
				return err
			}
		case "run":
			if err := ensureBuildDir(logger); err != nil {
				return err
			}

			// Check for --html5 flag
			html5Mode := false
			filteredArgs := []string{}
			for _, arg := range args[1:] {
				if arg == "--html5" {
					html5Mode = true
				} else {
					filteredArgs = append(filteredArgs, arg)
				}
			}

			if html5Mode {
				return serveBuiltFiles(logger)
			}

			// Update args to remove --html5 if it was present
			args = append([]string{"run"}, filteredArgs...)

		case "clean":
			if err := cleanBuildDir(logger); err != nil {
				return err
			}
			logger.Info("Build directory cleaned successfully")
			return nil
		}
	}

	command := "make " + strings.Join(args, " ")
	output, err := runCommand(command)
	if err != nil {
		logger.Error("Running make", "error", err)
	} else {
		fmt.Print(output)
	}

	return err
}

// cleanBuildDir removes the build directory
func cleanBuildDir(logger *log.Logger) error {
	err := os.RemoveAll("build")
	if err != nil {
		logger.Error("Failed to remove build directory", "error", err)
		return err
	}
	return nil
}
