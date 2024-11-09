# Release notes

- 🧑‍💻 Improving developer experience
- ✨ Introducing new features
- 🐛 Fixing bugs
- ♻️ Refactoring code

## v2.3.0 (2024-11-08)
- 🐛 Fixed bugs with `len` function for `STRING` type that was not handling Unicode correctly
- ✨ Added `at` function for `STRING` type

## v2.2.0 (2024-11-07)
- ✨ Added `truncate` function for `STRING` type
- ✨ Added `len` function for `INTEGER` type that returns the number of digits in the integer
- ✨ Added `decimal` function for `INTEGER` and `STRING` types that converts to a string with decimal points
- ♻️ Added more tests for utility functions

## v2.1.1 (2024-10-24)
- 🐛 Fixed bug where you couldn't pass `nil` to `textwire.NewTemplate` function

## v2.1.0 (2024-10-24)
For more detailed information about this release, read the [Textwire v2.1.0 Release Notes](https://textwire.github.io/blog/2024/10/24/textwire-v2.1.0-release-notes)

- ✨ Features
    - For array literals, added `4` built-in functions: `rand`, `reverse`, `slice`, `shuffle`
    - For integer literals, added `2` built-in functions: `abs`, `str`
    - For float literals, added `5` built-in functions: `abs`, `ceil`, `floor`, `round`, `str`
    - For string literals, added `3` built-in functions: `capitalize`, `reverse`, `contains`
    - For boolean literals, added `1` built-in function: `binary`
    - New error page while rendering a template. Instead of black screen we now get a simple error page with `Sorry! We’re having some trouble right now. Please check back shortly`
- 🧑‍💻 Improvements
    - 🐛 **Fixed Bug with Prefix Expression Precedence**: Resolved an issue where prefix expressions like `{{ -1.abs() }}` were not being processed correctly. Previously, the parser evaluated the expression as `{{ (-(1.abs())) }}`, resulting in an incorrect output of `-1`. Now, the parser correctly handles the precedence, evaluating it as `{{ ((-1).abs()) }}`
    - 🧑‍💻 **Enhanced Error Handling for Built-in Functions:** Improved error messages when an incorrect argument type is passed to a built-in function. Users will now receive clear error messages indicating the type mismatch
    - 🧑‍💻 **Enhanced Error Handling for Custom Functions:** If a function is called on a type where it doesn’t exist, Textwire now provides a detailed error message specifying that the function is undefined for that type. For example, an error message might read: `[Textwire ERROR in /var/www/html/templates/home.tw.html:3]: function 'some' doesn't exist for type 'STRING'`
    - 🧑‍💻 **Enhanced Error Handling for Division by Zero:** Improved error messages for division-by-zero cases, replacing previous vague messages with more meaningful ones
- 📝 Remove `CONTRIBUTING.md` file

## v2.0.0 (2024-10-18)
- ♻️ [BREAKING CHANGE!] Moved `textwire.Config` to a separate package `config.Config`
- ✨ [suggested by @joeyjurjens](https://github.com/joeyjurjens) Added the ability to register your own custom functions for specific types and use them in your Textwire code like built-in functions. If you are upgrading from version 1, make these changes:
    1. Change all the imports from `github.com/textwire/textwire` to `github.com/textwire/textwire/v2`
    2. Run `go mod tidy` to update the dependencies
    3. Change the package name from `textwire.Config` to `config.Config` in your code if you use configuration and import `"github.com/textwire/textwire/v2/config"`. If you already have a package named `config`, you can alias the import like `twconfig "github.com/textwire/textwire/v2/config"`

## [Release Notes V1](.github/CHANGELOG-V1.md)