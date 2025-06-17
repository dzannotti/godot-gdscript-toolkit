package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dzannotti/gdtoolkit/internal/core/linter"
	"github.com/dzannotti/gdtoolkit/internal/core/linter/rules"
)

func main() {
	// Parse command-line flags
	flag.Parse()

	// Get the file paths from the command-line arguments
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: gdlint [file.gd...]")
		os.Exit(1)
	}

	// Process each file
	hasErrors := false
	for _, path := range args {
		if err := processFile(path); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", path, err)
			hasErrors = true
		}
	}

	// Exit with non-zero status if there were errors
	if hasErrors {
		os.Exit(1)
	}
}

// processFile reads and parses a GDScript file
func processFile(path string) error {
	// Check if the file exists
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	// Check if it's a directory
	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file")
	}

	// Read the file
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Create a linter with all default rules
	lint := linter.NewLinter(
		rules.GetDefaultRules(),
		linter.DefaultConfig(),
	)

	// Lint the file
	problems, err := lint.Lint(string(content))
	if err != nil {
		return fmt.Errorf("failed to lint file: %w", err)
	}

	// Print any problems found
	if len(problems) > 0 {
		fmt.Printf("Linting %s:\n", path)
		for _, p := range problems {
			fmt.Printf("  %v\n", p)
		}
		return fmt.Errorf("%d linting problems", len(problems))
	}

	fmt.Printf("Successfully linted %s (no problems found)\n", path)
	return nil
}

// findGDScriptFiles finds all .gd files in a directory
func findGDScriptFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".gd" {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
