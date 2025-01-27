# Release notes

## v2.5.0 (2025-01-26)
### ⚠️ Possibly Breaking Changes
- **Error when inserting unknown**: If you try to insert content that is not reserved in the layout, you will now receive an error message. Previously, the it was silently ignored.
- **Error when using duplicate inserts**: If you try to insert content multiple times with the same insert name, you will now receive an error message. Previously, the it was silently ignored.

### 🐛 Bug Fixes
- **Incorrect Filepath in `@use` Statement**: Resolved an issue where invalid file paths in the `@use()` statement did not produce proper Textwire errors. Users will now see clear error messages for incorrect file paths.
- **Duplicate Unnamed Slots**: Fixed a bug where templates with duplicate unnamed slots silently used the first slot and ignored others. Errors will now be shown to alert users of this issue.
- **EvaluateFile Function**: Fixed the output of the `EvaluateFile` function to return the evaluated content as a string. Previously, it was bugged and was returning the absolute path of the file.

### ✨ New Features
- **Layout Shortcut**: You can now use a shortcut for referencing layouts. Instead of writing `@use("layouts/main")`, simply use `@use("~main")`. The `~` symbol is an alias for the `layouts/` directory.
- **Custom Error Page**: Now you can provide a path to a custom error page in the `textwire.NewTemplate` function by passing a config with the `ErrorPagePath` field. If an error occurs while rendering a template, the custom error page will be displayed instead of the [default error page](https://textwire.github.io/docs/v2/guides/error-handling#error-output-in-templates).
- **Improve Error Page**: Now you can set a `DebugMode` field in the config to `true` to see the error messages in the browser. This is useful for debugging purposes. Default is `false`. Plus, if you use VSCode, you can now click on the file path in the error message to open the file directly in the editor.
- **Custom Error Page**: Now you can provide a path to a custom error page through the `textwire.NewTemplate` function by passing a config with the `ErrorPagePath` field. If an error occurs while rendering a template, the custom error page will be displayed instead of the default error page.

## v2.4.1 (2025-01-23)
### 🐛 Bug Fixes
- **Optional `@insert` Statement**: Fixed an issue where the `@insert` statement was mandatory even when a `@reserve` statement was defined in the layout. All `@insert` statements are now optional in such cases.

## v2.4.0 (2025-01-19)
### ✨ New Features
- **String Functions**: Added 3 new utility functions for strings:
    - `trimRight`: Removes trailing spaces.
    - `trimLeft`: Removes leading spaces.
    - `repeat`: Repeats a string a specified number of times.
- **Component Shortcut**: Simplify component usage with a new alias. Instead of writing `@component("components/post-card", { post })`, you can now write `@component("~post-card", { post })`. The `~` symbol is an alias for the `components/` directory.
- **Debugging Helper**: Added the `@dump` directive to output the value of a variable for debugging purposes.
- **Array Functions**: Added 2 new utility functions for arrays:
    - `append`: Adds an element to the end of an array.
    - `prepend`: Adds an element to the beginning of an array.

For more details, read the [blog post](https://textwire.github.io/blog/2025/01/10/textwire-v2.4.0).

## v2.3.0 (2024-12-29)
- 🐛 Fixed bugs with function `len` for `STRING` type that was not handling Unicode correctly
- ✨ Added function `at` for `STRING` type that returns the character at the specified index
- ✨ Added function `first` for `STRING` type that returns the first characters of the string
- ✨ Added function `last` for `STRING` type that returns the last characters of the string
- ✨ Added function `then` for `BOOLEAN` type that returns the passed value if the boolean is true, otherwise returns `nil`. If second argument is passed, it returns the second argument if the boolean is false
- ✨ Added function `contains` for `ARRAY` type that returns true if the array contains the specified value, otherwise returns false

## v2.2.0 (2024-11-07)
- ✨ Added `truncate` function for `STRING` type
- ✨ Added `len` function for `INTEGER` type that returns the number of digits in the integer
- ✨ Added `decimal` function for `INTEGER` and `STRING` types that converts to a string with decimal points
- ♻️ Added more tests for utility functions

## v2.1.1 (2024-10-24)
- 🐛 Fixed bug where you couldn't pass `nil` to `textwire.NewTemplate` function

## v2.1.0 (2024-10-24)
### ✨ New Features
- For array literals, added `4` built-in functions: `rand`, `reverse`, `slice`, `shuffle`
- For integer literals, added `2` built-in functions: `abs`, `str`
- For float literals, added `5` built-in functions: `abs`, `ceil`, `floor`, `round`, `str`
- For string literals, added `3` built-in functions: `capitalize`, `reverse`, `contains`
- For boolean literals, added `1` built-in function: `binary`
- New error page while rendering a template. Instead of black screen we now get a simple error page with `Sorry! We’re having some trouble right now. Please check back shortly`

### 🧑‍💻 Improvements
- **Enhanced Error Handling for Built-in Functions**: Improved error messages when an incorrect argument type is passed to a built-in function. Users will now receive clear error messages indicating the type mismatch
- **Enhanced Error Handling for Custom Functions**: If a function is called on a type where it doesn’t exist, Textwire now provides a detailed error message specifying that the function is undefined for that type. For example, an error message might read: `[Textwire ERROR in /var/www/html/templates/home.tw.html:3]: function 'some' doesn't exist for type 'STRING'`
- **Enhanced Error Handling for Division by Zero**: Improved error messages for division-by-zero cases, replacing previous vague messages with more meaningful ones
- Remove `CONTRIBUTING.md` file

### 🐛 Bug Fixes
- **Fixed Bug with Prefix Expression Precedence**. Resolved an issue where prefix expressions like `{{ -1.abs() }}` were not being processed correctly. Previously, the parser evaluated the expression as `{{ (-(1.abs())) }}`, resulting in an incorrect output of `-1`. Now, the parser correctly handles the precedence, evaluating it as `{{ ((-1).abs()) }}`

For more detailed information about this release, read the [Textwire v2.1.0 Release Notes](https://textwire.github.io/blog/2024/10/24/textwire-v2.1.0-release-notes)

## v2.0.0 (2024-10-18)
- ♻️ [BREAKING CHANGE!] Moved `textwire.Config` to a separate package `config.Config`
- ✨ [suggested by @joeyjurjens](https://github.com/joeyjurjens) Added the ability to register your own custom functions for specific types and use them in your Textwire code like built-in functions. If you are upgrading from version 1, make these changes:
    1. Change all the imports from `github.com/textwire/textwire` to `github.com/textwire/textwire/v2`
    2. Run `go mod tidy` to update the dependencies
    3. Change the package name from `textwire.Config` to `config.Config` in your code if you use configuration and import `"github.com/textwire/textwire/v2/config"`. If you already have a package named `config`, you can alias the import like `twconfig "github.com/textwire/textwire/v2/config"`

## [Release Notes V1](.github/CHANGELOG-V1.md)
## [Emojis Meaning](.github/EMOJIS.md)