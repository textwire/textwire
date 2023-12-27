# Textwire. A template language for Go.

Textwire is a simple yet powerful template language for Go. It is designed to easily inject variables from Go code into a template file or just a regular string. Here is a simple example of parsing a string:

```go
import (
    "fmt"
    "github.com/textwire/textwire"
)

str := `Hello, my name is {{ name }} and I am {{ age }} years old.`

vars := map[string]interface{}{
    "name": "John",
    "age":  21,
}

parsed, err := textwire.ParseStr(str, vars)

fmt.Println(parsed) // would print: Hello, my name is John and I am 21 years old.
```


[![Go](https://github.com/textwire/textwire/actions/workflows/go.yml/badge.svg)](https://github.com/textwire/textwire/actions/workflows/go.yml)

## Features

This is a list of features that are implemented and planned to be implemented in the future.

- Statements
    - [ ] If statements `{{ if x == 1 }}`
    - [ ] Else statements `{{ else }}`
    - [ ] Else-if statements `{{ else if x == 1 }}`
    - [ ] For statements `{{ for i, name := range names }}`
- Expressions
    - [ ] Ternary expressions `x ? y : z`
    - [ ] Prefix expressions `!x`
    - [ ] Infix expressions `x * y`
- Literals
    - [x] String literals `"Hello, World!"`
    - [x] Integer literals `123` or `-234`
    - [ ] Float literals `123.456`
    - [ ] Boolean literals `true`
    - [x] Nil literal `nil`
    - [ ] Slice literals `[]int{1, 2, 3}`

## Installation

```bash
go get github.com/textwire/textwire
```