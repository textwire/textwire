# Release notes

## v3.0.0 (2026-02-05) - Major Release
📖 [Migration Guide](https://textwire.github.io/docs/v3/upgrade) | [Announcement](https://textwire.github.io/blog/2026/02/05/textwire-v3)

### 🧑‍💻 Improvements
1. Improve error handling when trying to use `@use`, `@insert`, `@reserve` or `@component` directives in simple `EvaluateString` or `EvaluateFile` function calls. These directives are only allowed inside template files with `textwire.NewTemplate`.
2. Improve memory and performance. Read improvements [here](https://textwire.github.io/blog/2026/02/05/textwire-v3#memory-performance).
3. Improve error messages. Now they are more clear.
4. Textwire files that are not a part of the program are now ignored. If you have something like `book.tw` and you don't include it anywhere, it will be ignored in Textwire v3.

### 🐛 Bug Fixes
1. Fixed incorrect file path in error messages when error happens inside of `@insert` directive.
2. Fixed `contains` function for strings, `{{ !"aaa".contains("a") }}` now returns correct result.
3. Fixed `contains` function for arrays, `{{ ![{}, 21].contains({age: 21}) }}` now returns correct result.

### ✨ New Features
1. Added `globals` object. You can now add `GlobalData` to your configurations and access this data in your templates using `globals` object. For example: `globals.env`.
2. Added `defined()` global function. It returns true if variable is defined. [docs](https://textwire.github.io/docs/v3/functions/global#defined)
3. Now you can add custom functions to objects as well with `RegisterObjFunc`.

### ⚠️ BREAKING CHANGES
1. When you defined a custom function, now it returns type `any`. If you register any custom functions make sure to change return type to `any`.
2. Variable `global` is now reserved.
3. Fixed precedence for prefix expressions. Instead of `((!var).func())` we now have `(!(var.func()))`.
4. Changed default file extension from `.tw.html` to `.tw`. If you still want to support it, go to your configurations in `NewTemplate` or `Configure` and add field `TemplateExt: ".tw.html"` to it.
5. Minimal Go version support is version `1.25`.

## [Release Notes V1](.github/CHANGELOG-V1.md)
## [Release Notes V2](.github/CHANGELOG-V2.md)
## [Emojis Meaning](.github/EMOJIS.md)
