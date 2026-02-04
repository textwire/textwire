# Release notes

## v3.0.0 (2026-02-05) — Major Release

📖 [Migration Guide](https://textwire.github.io/docs/v3/upgrade) | [Announcement](https://textwire.github.io/blog/2026/02/05/textwire-v3)

### 🧑‍💻 Improvements

1. Improve error handling when trying to use `@use`, `@insert`, `@reserve` or `@component` directives in simple `EvaluateString` or `EvaluateFile` function calls. These directives are only allowed inside template files with `textwire.NewTemplate`.
2. Improve memory and performance. Read about improvements [here](https://textwire.github.io/blog/2026/02/05/textwire-v3#memory-performance).
3. Improve error messages. Now they are clearer.

### 🐛 Bug Fixes

1. Fixed incorrect file path in error messages when error happens inside `@insert` directive.
2. Fixed the `contains` function for strings; `{{ !"aaa".contains("a") }}` now returns the correct result.
3. Fixed the `contains` function for arrays; `{{ ![{}, 21].contains({age: 21}) }}` now returns the correct result.
4. Now you will get a proper error when trying to access a property on a non-object type like `{{ "str".nice }}`. Before, you would get a panic.
5. Fixed issue where you couldn't write a slot directive with space after the `@slot` keyword. This `@slot ("book")` was giving an error previously.
6. Now you will get an proper error when trying to use `@each` statement on non-array type.
7. Fixed bug where you couldn't use `@component` statement inside of a layout.

### ✨ New Features

1. Added `globals` object. You can now add `GlobalData` to your configurations and access this data in your templates using the `globals` object. For example: `globals.env`.
2. Added the `defined()` global function. It returns true if the variable is defined. [Docs](https://textwire.github.io/docs/v3/functions/global#defined)
3. Now you can add custom functions to objects as well with `RegisterObjFunc`.
4. Now you can use Go's embedded package to embed Textwire template files into a final binary.

### ⚠️ BREAKING CHANGES

1. When you defined a custom function, now it returns type `any`. If you register any custom functions, make sure to change the return type to `any`.
2. Variable `global` is now reserved.
3. Fixed precedence for prefix expressions. Instead of `((!var).func())` we now have `(!(var.func()))`.
4. Changed default file extension from `.tw.html` to `.tw`. If you still want to support it, go to your configurations in `NewTemplate` or `Configure` and add field `TemplateExt: ".tw.html"` to it.
5. Minimal Go version support is version `1.25.0`.

## [Release Notes V1](.github/CHANGELOG-V1.md)

## [Release Notes V2](.github/CHANGELOG-V2.md)

## [Emojis Meaning](.github/EMOJIS.md)
