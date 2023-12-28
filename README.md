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

## Usage with templates

### Layouts

You can define a layout for your template by creating a `[name].textwire.html`. Here is an example of a layout:

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ reserve "title" }}</title>
</head>
<body>
    {{ reserve "content" }}
</body>
</html>
```

#### The "reserve" keyword
The "reserve" keyword is used to reserve a place for dynamic content that you can insert later. For example, you can reserve a place for the title of the page and then insert it later. Here is an example of inserting a title:

```html
{{ use "layouts/main" }}

{{ insert "title", "Home page" }}
```

First, we use the layout "layouts/main" so that parser knows which layout to use. Then we insert the title into the reserved place. The first argument is the name of the reserved place and the second argument is the value that we want to insert.

As an alternative, we can insert the title like this:

```html
{{ use "layouts/main" }}

{{ insert "title" }}
    Home page
{{ end }}
```

### The "use" keyword

The "use" keyword is used to specify which layout to use. Assuming that our layout is placed in the "layouts" folder and called "main.textwire.html", we can use it like this:

```html
{{ use "layouts/main" }}
```

`"layouts/main"` is the relative path to the layout file. If you have deeply nested files and don't want to always specify the relative path, you can use the set the aliases in your `main.go` file like this:

```go
textwire.SetAliases(map[string]string{
    "@layouts": "src/views/templates/layouts",
})
```

Then you can use the layout like this:

```html
{{ use "@layouts/main" }}
````

## Installation

```bash
go get github.com/textwire/textwire
```