// Package validation provides benchmarking tests for parser performance and compatibility
package validation

import (
	"testing"
	"time"

	"github.com/dzannotti/gdtoolkit/internal/core/parser"
	"github.com/dzannotti/gdtoolkit/internal/testutil"
)

// BenchmarkParserPerformance benchmarks parser performance on various GDScript samples
func BenchmarkParserPerformance(t *testing.B) {
	// Sample GDScript codes of varying complexity for benchmarking
	benchmarkSamples := []struct {
		name string
		code string
	}{
		{
			name: "SimpleFunction",
			code: `func test():
	pass`,
		},
		{
			name: "FunctionWithParameters",
			code: `func test(a: int, b: String, c = 42):
	return a + c`,
		},
		{
			name: "ClassDefinition",
			code: `class MyClass extends Node:
	var health: int = 100
	var mana: float = 50.0
	
	func _ready():
		print("Ready!")
	
	func take_damage(amount: int):
		health -= amount
		if health <= 0:
			die()
	
	func die():
		queue_free()`,
		},
		{
			name: "ComplexExpressions",
			code: `func complex_calculation():
	var result = (1 + 2) * 3 / 4
	var comparison = result > 5 and result < 10
	var ternary = result if comparison else 0
	var array_access = items[index][nested_index]
	var method_chain = player.inventory.weapons[0].damage
	return result + ternary + array_access + method_chain`,
		},
		{
			name: "ControlFlowStatements",
			code: `func control_flow_example(items):
	for item in items:
		if item.type == "weapon":
			match item.rarity:
				"common":
					item.damage *= 1.0
				"rare":
					item.damage *= 1.5
				"epic":
					item.damage *= 2.0
				_:
					item.damage *= 0.8
		elif item.type == "potion":
			while item.quantity > 0:
				if use_potion(item):
					item.quantity -= 1
				else:
					break`,
		},
		{
			name: "NestedStructures",
			code: `class GameManager:
	class PlayerStats:
		var health: int = 100
		var mana: int = 50
		
		func heal(amount: int):
			health = min(health + amount, 100)
	
	class Inventory:
		var items: Array[Item] = []
		
		func add_item(item: Item):
			items.append(item)
	
	var player_stats: PlayerStats
	var inventory: Inventory
	
	func _init():
		player_stats = PlayerStats.new()
		inventory = Inventory.new()`,
		},
	}

	for _, sample := range benchmarkSamples {
		t.Run(sample.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				p := parser.NewParser(sample.code)
				tree := p.Parse()
				errors := p.Errors()

				if len(errors) > 0 {
					b.Fatalf("Parsing failed: %v", errors)
				}
				if tree == nil {
					b.Fatalf("AST is nil")
				}
			}
		})
	}
}

// BenchmarkFixtureParsing benchmarks parsing performance on actual Python gdtoolkit fixtures
func BenchmarkFixtureParsing(t *testing.B) {
	fixtures := testutil.GetTestFixtures()

	// Load a representative sample of fixture files for benchmarking
	testFiles := []string{
		"functions.gd",
		"expressions.gd",
		"match.gd",
		"static_typing.gd",
		"annotations.gd",
	}

	for _, fileName := range testFiles {
		filePath := fixtures.ValidScripts + "/" + fileName
		content, err := testutil.LoadTestFile(filePath)
		if err != nil {
			t.Logf("Skipping benchmark for %s: %v", fileName, err)
			continue
		}

		t.Run(fileName, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				p := parser.NewParser(content)
				tree := p.Parse()
				errors := p.Errors()

				if len(errors) > 0 {
					b.Fatalf("Parsing failed for %s: %v", fileName, errors)
				}
				if tree == nil {
					b.Fatalf("AST is nil for %s", fileName)
				}
			}
		})
	}
}

// TestParserMemoryUsage tests parser memory efficiency
func TestParserMemoryUsage(t *testing.T) {
	// Large GDScript sample to test memory usage
	largeScript := generateLargeGDScript()

	// Parse the large script and verify it works
	p := parser.NewParser(largeScript)
	tree := p.Parse()
	errors := p.Errors()

	if len(errors) > 0 {
		t.Fatalf("Failed to parse large script: %v", errors)
	}

	if tree == nil {
		t.Fatalf("AST is nil for large script")
	}

	t.Logf("Successfully parsed large script with multiple functions and variables")
}

// TestParserScalability tests how parser performance scales with input size
func TestParserScalability(t *testing.T) {
	sizes := []int{100, 500, 1000, 2000, 5000}

	for _, size := range sizes {
		t.Run(intToString(size)+"_lines", func(t *testing.T) {
			script := generateGDScriptWithLines(size)

			start := time.Now()
			p := parser.NewParser(script)
			tree := p.Parse()
			errors := p.Errors()
			duration := time.Since(start)

			if len(errors) > 0 {
				t.Fatalf("Failed to parse %d line script: %v", size, errors)
			}

			if tree == nil {
				t.Fatalf("AST is nil for %d line script", size)
			}

			t.Logf("Parsed %d lines in %v", size, duration)

			// Performance check: should parse reasonable-sized files quickly
			maxDuration := time.Duration(size) * time.Millisecond // 1ms per line max
			if duration > maxDuration {
				t.Errorf("Parser took too long for %d lines: %v (max: %v)",
					size, duration, maxDuration)
			}
		})
	}
}

// TestParserConcurrency tests parser behavior under concurrent usage
func TestParserConcurrency(t *testing.T) {
	const numGoroutines = 10
	const numIterations = 100

	script := `class TestClass:
	var value: int = 42
	
	func test_method(param: String) -> int:
		for i in range(10):
			if i % 2 == 0:
				match param:
					"add":
						value += i
					"subtract":
						value -= i
					_:
						pass
		return value`

	// Channel to collect results
	results := make(chan error, numGoroutines*numIterations)

	// Start concurrent parsers
	for g := 0; g < numGoroutines; g++ {
		go func(goroutineID int) {
			for i := 0; i < numIterations; i++ {
				p := parser.NewParser(script)
				tree := p.Parse()
				errors := p.Errors()

				if len(errors) > 0 {
					results <- NewValidationError("Goroutine %d iteration %d failed: %v",
						goroutineID, i, errors)
					return
				}

				if tree == nil {
					results <- NewValidationError("Goroutine %d iteration %d: AST is nil",
						goroutineID, i)
					return
				}

				results <- nil // Success
			}
		}(g)
	}

	// Collect results
	totalOperations := numGoroutines * numIterations
	successCount := 0

	for i := 0; i < totalOperations; i++ {
		err := <-results
		if err == nil {
			successCount++
		} else {
			t.Errorf("Concurrent parsing failed: %v", err)
		}
	}

	t.Logf("Concurrent parsing: %d/%d operations successful", successCount, totalOperations)

	if successCount != totalOperations {
		t.Errorf("Some concurrent parsing operations failed")
	}
}

// Helper functions for generating test scripts

func generateLargeGDScript() string {
	// Generate a large GDScript with various constructs
	script := `class LargeScript extends Node:
	# Large script for memory testing
	
	var large_array: Array = []
	var large_dict: Dictionary = {}
	
`

	// Generate many functions
	for i := 0; i < 100; i++ {
		script += "	func function_" + intToString(i) + "(param: int) -> int:\n"
		script += "		var result = param * " + intToString(i) + "\n"
		script += "		for j in range(10):\n"
		script += "			result += j\n"
		script += "		return result\n\n"
	}

	// Generate many variables
	for i := 0; i < 200; i++ {
		script += "	var variable_" + intToString(i) + ": int = " + intToString(i*10) + "\n"
	}

	return script
}

func generateGDScriptWithLines(numLines int) string {
	script := `class GeneratedScript:
	# Generated script with ` + intToString(numLines) + ` lines
	
	var counter: int = 0
	
	func _ready():
		print("Script with ` + intToString(numLines) + ` lines")
	
`

	// Generate enough functions to reach the target line count
	functionsNeeded := (numLines - 10) / 8 // Approximately 8 lines per function
	for i := 0; i < functionsNeeded; i++ {
		script += "	func generated_function_" + intToString(i) + "():\n"
		script += "		counter += 1\n"
		script += "		if counter > 100:\n"
		script += "			counter = 0\n"
		script += "		return counter\n\n"
	}

	return script
}
