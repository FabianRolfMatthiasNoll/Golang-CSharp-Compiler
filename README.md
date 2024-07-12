
# C# Compiler in Go

This project aims to implement a basic C# compiler written in Go. The primary motivation behind this project is to deepen understanding of both C# and Go languages by building a compiler from scratch. The project will start with implementing fundamental features and may be extended based on interest and motivation.

## Project Goals

- Implement basic arithmetic operations
- Implement classes with fields, constructors, and methods
- Support access modifiers
- Implement object instantiation
- Implement simple control flow procedures

## Current Status

- Basic lexer and pratt parser implemented
- Ongoing work on AST generation and type checking

## Project Checklist

### Lexer and Parser

- [x] Basic lexer
- [x] Basic parser

### Semantic Analysis

- [ ] Type checking
- [ ] Variable declaration and usage checks
- [ ] Function declaration and calls
- [ ] Control flow analysis

### Intermediate Representation (IR)

- [ ] Generate IR from typed AST
- [ ] Implement IR interpreter for testing

### Code Generation

- [ ] Generate Common Intermediate Language (CIL) from IR
- [ ] Handle classes and methods in CIL
- [ ] Implement variable and type management in CIL
- [ ] Implement control flow in CIL

### Execution

- [ ] Save CIL to file (e.g., .dll or .exe)
- [ ] Load and execute the CIL file using .NET runtime

## Project Structure

```
/src
│
├── /lexer
│   ├── lexer.go
│   └── tokens.go
│
├── /parser
│   ├── parser.go
│   ├── expr.go
│   ├── stmt.go
│   ├── types.go
│   └── lookups.go
│
├── /ast
│   ├── ast.go
│   └── prettyprint.go
│
├── main.go
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.16+
- .NET SDK

### Running the Project

1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/csharp-compiler-in-go.git
    cd csharp-compiler-in-go
    ```

2. Run the compiler as far as implemented:

    ```sh
    go run src/main.go
    ```

## License

This project is licensed under the MIT License.
