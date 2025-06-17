package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dzannotti/gdtoolkit/internal/core/formatter"
	"github.com/dzannotti/gdtoolkit/internal/core/parser"
)

func main() {
	// Parse command-line flags
	checkOnly := flag.Bool("check", false, "Check if files are formatted without modifying them")
	flag.Parse()

	// Get the file paths from the command-line arguments
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: gdformat [--check] [file.gd...]")
		os.Exit(1)
	}

	// Process each file
	hasErrors := false
	for _, path := range args {
		if err := processFile(path, *checkOnly); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", path, err)
			hasErrors = true
		}
	}

	// Exit with non-zero status if there were errors
	if hasErrors {
		os.Exit(1)
	}
}

// processFile reads, parses, and formats a GDScript file
func processFile(path string, checkOnly bool) error {
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

	// Parse the file
	ast, errors := parser.ParseFile(path, string(content))
	if len(errors) > 0 {
		fmt.Printf("Parsing %s:\n", path)
		for _, err := range errors {
			fmt.Printf("  %v\n", err)
		}
		return fmt.Errorf("%d parsing errors", len(errors))
	}

	// Format the AST
	config := formatter.DefaultConfig()
	formattedCode, err := formatter.FormatCode(ast, config)
	if err != nil {
		return fmt.Errorf("formatting error: %w", err)
	}

	// Ensure the formatted code ends with a newline
	if !strings.HasSuffix(formattedCode, "\n") {
		formattedCode += "\n"
	}

	if checkOnly {
		// Check if the file is already formatted correctly
		if string(content) == formattedCode {
			fmt.Printf("File %s is correctly formatted\n", path)
		} else {
			fmt.Printf("File %s would be reformatted\n", path)
			return fmt.Errorf("file needs formatting")
		}
	} else {
		// Write the formatted code back to the file
		err = ioutil.WriteFile(path, []byte(formattedCode), 0644)
		if err != nil {
			return fmt.Errorf("failed to write formatted file: %w", err)
		}
		fmt.Printf("Successfully formatted %s\n", path)
	}

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
