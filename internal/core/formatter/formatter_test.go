package formatter

import (
	"strings"
	"testing"

	"github.com/dzannotti/gdtoolkit/internal/core/parser"
)

func TestBasicFormatting(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "simple_class_with_pass",
			input: `class X:
	pass`,
			expected: `class X:
	pass`,
		},
		{
			name: "class_with_function",
			input: `class X:
	func foo():
		pass`,
			expected: `class X:
	func foo():
		pass`,
		},
		{
			name: "function_with_parameters",
			input: `func foo(a,b,c):
	pass`,
			expected: `func foo(a, b, c):
	pass`,
		},
		{
			name: "function_with_typed_parameters",
			input: `func foo(a:int,b:String,c:float):
	pass`,
			expected: `func foo(a: int, b: String, c: float):
	pass`,
		},
		{
			name: "function_with_default_parameters",
			input: `func foo(a=1,b:int=2):
	pass`,
			expected: `func foo(a = 1, b: int = 2):
	pass`,
		},
		{
			name: "function_with_return_type",
			input: `func foo()->int:
	return 1`,
			expected: `func foo() -> int:
	return 1`,
		},
		{
			name: "variable_declarations",
			input: `var a=1
var b:int=2
const C=3`,
			expected: `var a = 1
var b: int = 2
const C = 3`,
		},
		{
			name: "if_statement",
			input: `if true:
	pass
else:
	pass`,
			expected: `if true:
	pass
else:
	pass`,
		},
		{
			name: "for_loop",
			input: `for i in range(10):
	print(i)`,
			expected: `for i in range(10):
	print(i)`,
		},
		{
			name: "while_loop",
			input: `while true:
	pass`,
			expected: `while true:
	pass`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse input
			ast, errors := parser.ParseFile("test.gd", tt.input)
			if len(errors) > 0 {
				t.Fatalf("Parse errors: %v", errors)
			}

			// Format
			config := DefaultConfig()
			result, err := FormatCode(ast, config)
			if err != nil {
				t.Fatalf("Format error: %v", err)
			}

			// Normalize line endings and trim
			expected := strings.TrimSpace(tt.expected)
			actual := strings.TrimSpace(result)

			if actual != expected {
				t.Errorf("Formatting mismatch:\nExpected:\n%s\n\nActual:\n%s", expected, actual)
			}
		})
	}
}

func TestTypeHintFormatting(t *testing.T) {
	input := `class SubClass:
	enum NamedEnum { A, B, C }
var a:Array[int]
var b:Array[int] = [1]
const C:Array[int]=[1]
var e:Array[SubClass.NamedEnum]

func foo(d:Array[int])->Array[int]:
	return [1]`

	expected := `class SubClass:
	enum NamedEnum { A, B, C }


var a: Array[int]
var b: Array[int] = [1]
const C: Array[int] = [1]
var e: Array[SubClass.NamedEnum]


func foo(d: Array[int]) -> Array[int]:
	return [1]`

	// Parse input
	ast, errors := parser.ParseFile("test.gd", input)
	if len(errors) > 0 {
		t.Fatalf("Parse errors: %v", errors)
	}

	// Format
	config := DefaultConfig()
	result, err := FormatCode(ast, config)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}

	// Normalize and compare
	expected = strings.TrimSpace(expected)
	actual := strings.TrimSpace(result)

	if actual != expected {
		t.Errorf("Type hint formatting mismatch:\nExpected:\n%s\n\nActual:\n%s", expected, actual)
	}
}

func TestFunctionParameterFormatting(t *testing.T) {
	input := `class Y:
	func _init(a, b, c):
		pass

class X:
	extends Y
	func foo(a):
		pass

	func bar(a,b):
		pass

	func baz (a, b=1):
		pass

	func bax (a, b:=1):
		pass

	func bac (a, b:int):
		pass

	func bav (a, b:int=1):
		pass

	func bab (a, b:int=1,c:=1) -> int:
		return 1

	func ban (a,b,c,):
		pass

class Z:
	extends Y
	func _init (a,b:=1,c:int=1):
		pass`

	expected := `class Y:
	func _init(a, b, c):
		pass


class X:
	extends Y

	func foo(a):
		pass

	func bar(a, b):
		pass

	func baz(a, b = 1):
		pass

	func bax(a, b = 1):
		pass

	func bac(a, b: int):
		pass

	func bav(a, b: int = 1):
		pass

	func bab(a, b: int = 1, c = 1) -> int:
		return 1

	func ban(
		a,
		b,
		c,
	):
		pass


class Z:
	extends Y

	func _init(a, b = 1, c: int = 1):
		pass`

	// Parse input
	ast, errors := parser.ParseFile("test.gd", input)
	if len(errors) > 0 {
		t.Fatalf("Parse errors: %v", errors)
	}

	// Format
	config := DefaultConfig()
	result, err := FormatCode(ast, config)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}

	// Normalize and compare
	expected = strings.TrimSpace(expected)
	actual := strings.TrimSpace(result)

	if actual != expected {
		t.Errorf("Function parameter formatting mismatch:\nExpected:\n%s\n\nActual:\n%s", expected, actual)
	}
}

func TestClassStructureFormatting(t *testing.T) {
	input := `class X:
	pass

	func foo():
		pass

	class Y:
		func bar():
			pass

class Z:
	pass

class Q:
	class W:
		pass`

	expected := `class X:
	pass

	func foo():
		pass

	class Y:
		func bar():
			pass


class Z:
	pass


class Q:
	class W:
		pass`

	// Parse input
	ast, errors := parser.ParseFile("test.gd", input)
	if len(errors) > 0 {
		t.Fatalf("Parse errors: %v", errors)
	}

	// Format
	config := DefaultConfig()
	result, err := FormatCode(ast, config)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}

	// Normalize and compare
	expected = strings.TrimSpace(expected)
	actual := strings.TrimSpace(result)

	if actual != expected {
		t.Errorf("Class structure formatting mismatch:\nExpected:\n%s\n\nActual:\n%s", expected, actual)
	}
}

func TestExpressionFormatting(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "binary_operations",
			input:    `var result = a+b*c-d/e`,
			expected: `var result = a + b * c - d / e`,
		},
		{
			name:     "function_calls",
			input:    `print("hello",world)`,
			expected: `print("hello", world)`,
		},
		{
			name:     "array_literals",
			input:    `var arr = [1,2,3]`,
			expected: `var arr = [1, 2, 3]`,
		},
		{
			name:     "dictionary_literals",
			input:    `var dict = {"key":"value","num":123}`,
			expected: `var dict = {"key": "value", "num": 123}`,
		},
		{
			name:     "dot_notation",
			input:    `player.position.x = 100`,
			expected: `player.position.x = 100`,
		},
		{
			name:     "index_access",
			input:    `var item = array[0]`,
			expected: `var item = array[0]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse input
			ast, errors := parser.ParseFile("test.gd", tt.input)
			if len(errors) > 0 {
				t.Fatalf("Parse errors: %v", errors)
			}

			// Format
			config := DefaultConfig()
			result, err := FormatCode(ast, config)
			if err != nil {
				t.Fatalf("Format error: %v", err)
			}

			// Normalize and compare
			expected := strings.TrimSpace(tt.expected)
			actual := strings.TrimSpace(result)

			if actual != expected {
				t.Errorf("Expression formatting mismatch:\nExpected:\n%s\n\nActual:\n%s", expected, actual)
			}
		})
	}
}

func TestControlFlowFormatting(t *testing.T) {
	input := `if condition:
	do_something()
elif other_condition:
	do_other()
else:
	do_default()

for item in collection:
	process(item)

while running:
	update()

match value:
	1:
		print("one")
	2:
		print("two")
	_:
		print("other")`

	expected := `if condition:
	do_something()
elif other_condition:
	do_other()
else:
	do_default()

for item in collection:
	process(item)

while running:
	update()

match value:
	1:
		print("one")
	2:
		print("two")
	_:
		print("other")`

	// Parse input
	ast, errors := parser.ParseFile("test.gd", input)
	if len(errors) > 0 {
		t.Fatalf("Parse errors: %v", errors)
	}

	// Format
	config := DefaultConfig()
	result, err := FormatCode(ast, config)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}

	// Normalize and compare
	expected = strings.TrimSpace(expected)
	actual := strings.TrimSpace(result)

	if actual != expected {
		t.Errorf("Control flow formatting mismatch:\nExpected:\n%s\n\nActual:\n%s", expected, actual)
	}
}

func TestConfigurationOptions(t *testing.T) {
	input := `class Test:
	func foo():
		pass`

	t.Run("default_tabs", func(t *testing.T) {
		ast, errors := parser.ParseFile("test.gd", input)
		if len(errors) > 0 {
			t.Fatalf("Parse errors: %v", errors)
		}

		config := DefaultConfig()
		result, err := FormatCode(ast, config)
		if err != nil {
			t.Fatalf("Format error: %v", err)
		}

		// Should contain tabs
		if !strings.Contains(result, "\t") {
			t.Error("Expected tabs in formatted output")
		}
	})

	t.Run("spaces_for_indent", func(t *testing.T) {
		ast, errors := parser.ParseFile("test.gd", input)
		if len(errors) > 0 {
			t.Fatalf("Parse errors: %v", errors)
		}

		config := DefaultConfig()
		spaces := 4
		config.SpacesForIndent = &spaces
		config.UseSpaces = true

		result, err := FormatCode(ast, config)
		if err != nil {
			t.Fatalf("Format error: %v", err)
		}

		// Should contain 4 spaces for indentation
		lines := strings.Split(result, "\n")
		foundIndentedLine := false
		for _, line := range lines {
			if strings.HasPrefix(line, "    ") {
				foundIndentedLine = true
				break
			}
		}
		if !foundIndentedLine {
			t.Error("Expected 4-space indentation in formatted output")
		}
	})
}
