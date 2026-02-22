# Release notes

## v3.2.3 (2026-02-22)

- 🐛 Fixed file watcher stops warching files when parser returns syntax error. Now, it will keep watching files changes even after syntax error happens.

## v3.2.2 (2026-02-22)

- 🐛 Fixed styles after `@dump` directive output being modified.

## v3.2.1 (2026-02-22)

- 🐛 Fixed bug with literal pointers like `*int`, `*string`, etc. Textwire would panic when you try to pass them to Textwire template.

## v3.2.0 (2026-02-16)

- ✨ Added a file wather that wathes your file changes and refreshes Textwire AST nodes. It prevents you from restarting server everytime you want to see any changes in the browser.
- 🧑‍💻 Accessing undefined property on an object does not give error anymore. It makes it consistant with accessing array on non-existant index. This `{{ {}.prop }}` returns nil now, but this `{{ {}.prop.second }}` causes error because you are trying to call property `second` on type `NIL`.
- 🐛 Fixed global function `defined`, it was returning `true` in cases like this `{{ name = "john"; defined(name.somemethod()) }}` because it was returning `true` when any error happens inside `defined`. Now, it only checks for undefined variables and undefined properties on objects.
- 🧑‍💻 Added so that now you can use any literal value in logical OR and logical AND expressions. Before, you could only use boolean on both sides. For example, now you can do `{{ "nice" && 13 ? "Yes" : "No" }}` and it returns `Yes` becuase non-empty string is `true` and non-zero int is also `true`.
- 🧑‍💻 Performance improve ment for Go's slice convertion into Textwire array. Here are the benchmarks:
    | Size | ⚡ Speed | 💾 Memory | 📉 Allocations |
    |------|----------|-----------|----------------|
    | small | **1.38× faster** | **48.7% less** | **6.4% fewer** |
    | medium | **1.23× faster** | **41.6% less** | **1.0% fewer** |
    | large | **1.67× faster** | **65.5% less** | **0.2% fewer** |
    | huge | **2.74× faster** | **73.7% less** | **0.03% more** |

## v3.1.2 (2026-02-15)

- 🧑‍💻 Full architecture change in directory and file structure. Doesn't break public API, just internal refactorking. Will be a breaking change if you use Textwire internals like parser, lexer, etc.

## v3.1.1 (2026-02-14)

- 🐛 Renamed 1 test file because it was causing some weird printing after `go get -u ./...` command.

## v3.1.0 (2026-02-14)

- 🧑‍💻 Added tests to make sure some age cases work.
- ✨ You can now add a fallback second argument to the reserve directive like this `@reserve('title', 'My Blog')` [#66](https://github.com/textwire/textwire/issues/66). It will be used when you didn't pass any inserts to that matching reseve.
- 🧑‍💻 Made slots defined in a component file optional.
- 🐛 Fixed bug where you could define multiple `@reserve` directives in a layout file. Now, you'll get an error.

## v3.0.2 (2026-02-11)

- ♻️ Lots of refactoring and improvements to the codebase, including adding tests for edge cases.
- 🐛 Fixed broken link on error page with debug mode on.

## v3.0.1 (2026-02-08)

- 🧑‍💻 Improve error messages.
- 🧑‍💻 Potentially a breaking change, but it should not be. You will receive a clear error if you trying to use `@insert` in a template file without defining `@use` directive. Previously, `@insert` would result in empty string in your template.
- 🐛 Fixed incorrect duplicate slot counter for error message when you have multiple duplicate slots in the same comonent.

## v3.0.0 (2026-02-07) — Major Release

📖 [Migration Guide](https://textwire.github.io/docs/v3/upgrade) | [Announcement](https://textwire.github.io/blog/2026/02/05/textwire-v3)

### 🧑‍💻 Improvements

1. Improve error handling when trying to use `@use`, `@insert`, `@reserve` or `@component` directives in simple `EvaluateString` or `EvaluateFile` function calls. These directives are only allowed inside template files with `textwire.NewTemplate`.
2. Improve memory and performance. Read about improvements [here](https://textwire.github.io/blog/2026/02/05/textwire-v3#memory-performance).
3. Improve error messages. Now they are clearer.
4. Added tons of tests to make sure version 3 is stable.
5. You'll get a clear error when using 2 or more `@use` directives in the same template.

### 🐛 Bug Fixes

1. Fixed incorrect file path in error messages when error happens inside `@insert` directive.
2. Fixed the `contains` function for strings; `{{ !"aaa".contains("a") }}` now returns the correct result.
3. Fixed the `contains` function for arrays; `{{ ![{}, 21].contains({age: 21}) }}` now returns the correct result.
4. Now you will get a proper error when trying to access a property on a non-object type like `{{ "str".nice }}`. Before, you would get a panic.
5. Fixed issue where you couldn't write a slot directive with space after the `@slot` keyword. This `@slot ("book")` was giving an error previously.
6. Now you will get an proper error when trying to use `@each` directive on non-array type.
7. Fixed bug where you couldn't use `@component` directive inside of a layout and `@component` inside of other components.

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
6. Components in Textwire v2 would pass variables to their children automatially without manual passing. It was a bug. In Textwire v3 each component has its scope. You need to pass data manually `@component('user', { user })`.
7. Fixed variable leak from template to layout non-explicitly. In Textwire v2, if you had a variable in your template, it would be accessible in your layout without passing it explicitly. In Textwire v3, this is not available anymore.

## [Release Notes V1](.github/CHANGELOG-V1.md)

## [Release Notes V2](.github/CHANGELOG-V2.md)

## [Emojis Meaning](.github/EMOJIS.md)
