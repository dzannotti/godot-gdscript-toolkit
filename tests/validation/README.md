# Parser Validation Test Suite

This directory contains comprehensive parser validation tests designed to ensure 1:1 compatibility between the Go gdtoolkit parser and the Python gdtoolkit implementation.

## Test Categories

### 1. Parser Compatibility Tests (`parser_compatibility_test.go`)
- **Purpose**: Tests Go parser against all Python gdtoolkit test fixtures
- **Scope**: Validates parsing of real-world GDScript samples from the Python test suite
- **Coverage**: Valid scripts, invalid scripts, formatter input/output pairs

### 2. AST Structure Validation (`ast_structure_validation_test.go`) 
- **Purpose**: Validates that AST structures are correctly built for various GDScript constructs
- **Scope**: Tests specific language features and edge cases
- **Coverage**: Functions, classes, variables, expressions, control flow

### 3. Fixture Parsing Tests (`fixture_parsing_test.go`)
- **Purpose**: Tests parsing against specific Python gdtoolkit fixture files
- **Scope**: Categorized testing of important GDScript features
- **Coverage**: Edge cases, regression prevention, formatter compatibility

### 4. Performance Benchmarks (`parser_benchmark_test.go`)
- **Purpose**: Benchmarks parser performance and scalability
- **Scope**: Performance testing on various code samples
- **Coverage**: Memory usage, scalability, concurrency

## Current Parser Status

As of the latest test run, the validation tests have identified several critical areas where the Go parser needs improvement to achieve 1:1 compatibility:

### ✅ Working Features
- Empty script parsing
- Basic AST structure creation
- Root class creation

### ❌ Issues Identified

#### Function Parsing
- Function parameter parsing with types fails
- Default parameter values not supported
- Return type annotations not working
- Static function declarations not recognized

#### Class Parsing  
- Class definitions not being created properly
- Class inheritance (`extends`) not working
- Nested classes not supported

#### Variable Declarations
- Type inference syntax (`:=`) not supported
- Typed variable declarations failing
- Const declarations with types not working

#### Expression Parsing
- Expression statements inside functions fail
- Arithmetic expressions not parsing
- Function calls not recognized
- Attribute access chains not working
- Array and dictionary literals not supported

#### Control Flow
- If/else statements not parsing
- For loops not supported
- While loops not supported  
- Match statements not working

#### Advanced Features
- Annotations (`@export`, `@tool`) not supported
- Signal definitions not working
- Property getters/setters not implemented
- Enum definitions not supported

## Test Results Summary

```
=== Parser Validation Results ===
✓ Empty scripts parse correctly
✗ Function definitions: FAIL
✗ Class definitions: FAIL  
✗ Variable declarations: FAIL
✗ Expression parsing: FAIL
✗ Control flow: FAIL
✗ Advanced features: FAIL

Current compatibility: ~5% (basic structure only)
Target compatibility: 100%
```

## Usage

Run all validation tests:
```bash
go test ./tests/validation/ -v
```

Run specific test categories:
```bash
# AST structure validation
go test ./tests/validation/ -run TestASTStructureValidation -v

# Fixture compatibility
go test ./tests/validation/ -run TestPythonGDToolkitFixtureCompatibility -v

# Performance benchmarks
go test ./tests/validation/ -run Benchmark -bench=. -v
```

## Next Steps

The validation tests have successfully identified the areas that need work to achieve parser compatibility:

1. **Priority 1: Core Language Constructs**
   - Fix function parameter parsing with types
   - Implement class definition parsing
   - Support variable declarations with type annotations

2. **Priority 2: Expression Support**
   - Implement expression statement parsing
   - Add support for function calls and operators
   - Support array/dictionary literals

3. **Priority 3: Control Flow**
   - Implement if/else statement parsing
   - Add for/while loop support
   - Implement match statement parsing

4. **Priority 4: Advanced Features**
   - Add annotation support
   - Implement signal/enum parsing
   - Support property getters/setters

## Test-Driven Development

These validation tests provide an excellent foundation for test-driven development:

1. Run tests to see current failures
2. Fix parser issues one by one
3. Re-run tests to verify fixes
4. Continue until 100% compatibility is achieved

The tests will catch regressions and ensure that fixes don't break previously working functionality.

## Contributing

When adding new GDScript features to the parser:

1. Add corresponding validation tests first
2. Run tests to see the expected failures
3. Implement the parser feature
4. Verify tests pass
5. Add edge case tests as needed

This ensures comprehensive coverage and maintains compatibility with the Python implementation.