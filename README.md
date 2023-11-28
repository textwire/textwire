# Templating language for Go

[![Go](https://github.com/go-temp/go-temp/actions/workflows/go.yml/badge.svg)](https://github.com/go-temp/go-temp/actions/workflows/go.yml)

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
    - [ ] String literals `"Hello, World!"`
    - [ ] Integer literals `123`
    - [ ] Float literals `123.456`
    - [ ] Boolean literals `true`
    - [ ] Slice literals `[]int{1, 2, 3}`

## Installation

```bash
go get github.com/go-temp/go-temp
```