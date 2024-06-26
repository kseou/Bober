package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"bober/config"

	"github.com/charmbracelet/log"
)

// getLibraryFlags retrieves compiler and linker flags for the specified libraries.
func getLibraryFlags(libraries []struct {
	Name   string `yaml:"name,omitempty"`
	Config string `yaml:"config,omitempty"`
}, logger *log.Logger) ([]string, []string) {
	var libCflags, libLibs []string
	for _, lib := range libraries {
		if lib.Config != "" {
			cflags, err := runCommand(lib.Config + " --cflags")
			if err != nil {
				logger.Errorf("Failed to get cflags for %s: %v", lib.Config, err)
			}
			libs, err := runCommand(lib.Config + " --libs")
			if err != nil {
				logger.Errorf("Failed to get libs for %s: %v", lib.Config, err)
			}
			libCflags = append(libCflags, strings.TrimSpace(cflags))
			libLibs = append(libLibs, strings.TrimSpace(libs))
		} else if lib.Name != "" {
			libLibs = append(libLibs, "-l"+lib.Name)
		}
	}
	return libCflags, libLibs
}

// generates a Makefile based on the provided configuration and writes it to the filesystem.
// It finds source and header files, determines include directories, and constructs the Makefile content.
func GenerateMakefile(config *config.Config, logger *log.Logger) error {
	sources, err := findSourceFiles(config.Sources)
	if err != nil {
		logger.Errorf("Failed to find source files: %v", err)
		return err
	}

	// ensure there are source files to process
	if len(sources) == 0 {
		logger.Error("No source files found")
		return fmt.Errorf("no source files found")
	}

	// get compiler and linker flags for libraries
	libCflags, libLibs := getLibraryFlags(config.Libraries, logger)

	// generate the Makefile content
	makefile := generateMakefileContent(config, sources, libCflags, libLibs)

	// write the Makefile to the filesystem
	return os.WriteFile("Makefile", []byte(makefile), 0644)
}

// findSourceFiles searches for source (.cpp) and header (.hpp) files based on the provided patterns.
func findSourceFiles(patterns []string) ([]string, error) {
	var sources []string
	for _, pattern := range patterns {
		baseDir := strings.TrimSuffix(pattern, "/**/*.*")
		err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".cpp" {
				sources = append(sources, path)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("error walking the path %s: %v", baseDir, err)
		}
	}
	return sources, nil
}

// generateMakefileContent constructs the content of the Makefile based on the provided configuration and file lists.
func generateMakefileContent(config *config.Config, sources, libCflags, libLibs []string) string {
	return fmt.Sprintf(`# Automatically generated Makefile for %s by Bober!

CXX = %s
CXXFLAGS = %s -std=%s %s
LDFLAGS = %s

SRCS = %s
EXECUTABLE = build/%s

build: $(EXECUTABLE)

$(EXECUTABLE):
	$(CXX) $(CXXFLAGS) $(SRCS) -o $@ $(LDFLAGS)

# For web builds, use 'bober run --html5' to start a local server
run: $(EXECUTABLE)
	./$(EXECUTABLE)

.PHONY: build run
`,
		config.Project.Name,
		config.Cpp.Compiler,
		config.Cpp.Flags,
		config.Cpp.Standard,
		strings.Join(libCflags, " "),
		strings.Join(libLibs, " "),
		strings.Join(sources, " "),
		config.Output.Executable,
	)
}
