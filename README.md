# GDScript Toolkit

Go GDToolkit is a port of the Python [gdtoolkit](https://github.com/Scony/godot-gdscript-toolkit) to Go. It provides tools for linting and formatting GDScript code, ensuring consistent code style and quality.

This project is a fork of the original excellent [gd-toolkit) [https://github.com/Scony/godot-gdscript-toolkit]

However i was frustrated with the amount of time it'd take to format and lint files
(often several seconds) so i tasked Roo Code to port the original repo (excluding gd2py and gdradon) to Go

This repo is fully AI generated code


This implementation follows a hexagonal architecture pattern to ensure clean separation of concerns, and uses a hand-written recursive descent parser for GDScript.

## Installation
TODO


## Project Structure

```
gogdtoolkit/
├── cmd/                    # Command-line applications
│   ├── gdlint/             # GDScript linter
│   └── gdformat/           # GDScript formatter
├── internal/               # Internal packages
│   ├── core/               # Core domain logic
│   │   ├── ast/            # Abstract Syntax Tree
│   │   ├── parser/         # GDScript parser
│   │   ├── linter/         # Linting rules
│   │   └── formatter/      # Formatting logic
│   ├── ports/              # Interfaces for the core domain
│   │   ├── primary/        # Primary ports (used by adapters)
│   │   └── secondary/      # Secondary ports (implemented by adapters)
│   └── adapters/           # Implementations of ports
│       ├── primary/        # Primary adapters (CLI, API)
│       └── secondary/      # Secondary adapters (file system, config)
├── pkg/                    # Public packages
│   └── gdscript/           # GDScript utilities
└── test/                   # Test files
    ├── integration/        # Integration tests
    ├── fixtures/           # Test fixtures
    └── acceptance/         # Acceptance tests
```

## Current Implementation Status

- [x] Basic project structure
- [x] Core AST structures
- [x] Lexer for GDScript
- [x] Basic parser implementation
- [ ] Complete parser implementation
- [ ] Linter rules
- [ ] Formatter implementation
- [ ] CLI tools

## Building and Running

### Prerequisites

- Go 1.21 or higher

### Building

```bash
# Build the linter
go build -o gdlint ./cmd/gdlint

# Build the formatter
go build -o gdformat ./cmd/gdformat
```

### Running

```bash
# Run the linter
./gdlint path/to/your/script.gd

# Run the formatter
./gdformat path/to/your/script.gd
```

## Testing

```bash
# Run all tests
go test ./...

# Run specific tests
go test ./internal/core/parser
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.