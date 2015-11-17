# lint-naked-returns

This program finds you naked returns in functions and methods with return arguments in Go(Golang) source code.

## Installation

```bash
go get -u github.com/zimmski/lint-naked-returns
```

## Example

This function's return statement will be flagged.

```go
func foo() (hey string) {
	return
}
```

## Usage

Search current package:
```bash
lint-naked-returns .
```

Search the package `github.com/zimmski/tavor` and its subpackages:
```bash
lint-naked-returns github.com/zimmski/tavor/...
```

Search the package `github.com/zimmski/tavor` and its subpackages using the build tag `example-main`:
```bash
lint-naked-returns --tag example-main github.com/zimmski/tavor/...
```

Show the programs help:
```bash
lint-naked-returns --help
```
