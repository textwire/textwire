# Templating language for Go

## Features

This is a list of features that are implemented and planned to be implemented in the future.

- Statements
    - [ ] If statements (ex. `{{ if x == 1 }}`)
    - [ ] Else statements (ex. `{{ else }}`)
    - [ ] Else-if statements (ex. `{{ else if x == 1 }}`)
    - [ ] For statements (ex. `{{ for i, name := range names }}`)
- Expressions
    - [ ] Ternary expressions (ex. `x ? y : z`)
    - [ ] Prefix expressions (ex. `!x`)
    - [ ] Infix expressions (ex. `x * y`)
- Literals
    - [ ] String literals (ex. `"Hello, World!"`)
    - [ ] Integer literals (ex. `123`)
    - [ ] Float literals (ex. `123.456`)
    - [ ] Boolean literals (ex. `true`)
    - [ ] Slice literals (ex. `[]int{1, 2, 3}`)

## Installation

```bash
go get github.com/go-temp/go-temp
```