# Release notes

## v4.0.0 (2026-03-20)

- ⚠️ BREAKING CHANGE! Changed increment and decrement expressions into statements to match Go's style. `{{ x++ }}` and `{{ x-- }}`  are not returning any value now, they equivalent to `{{ x = x + 1 }}` and `{{ x = x - 1 }}`. Closes [#88](https://github.com/textwire/textwire/issues/88). [Read more about this change](https://textwire.github.io/blog/2026-03-20-textwire-v4).
- ⚠️ BREAKING CHANGE! Remove support for `@breakIf` and `@continueIf`, use lovercased versions `@breakif` and `@continueif`. Closes [#89](https://github.com/textwire/textwire/issues/89).
- ♻️ Improved internal logic for evaluator. Separated literal values from nonliteral. Closes [#90](https://github.com/textwire/textwire/issues/90).
- 🧑‍💻 All public API functions like `NewTemplate()`, `EvaluateString`, etc., now return `*fail.Error` instead of Go's `error` type.
- 🧑‍💻 Added proper position to error messages. Closes [#101](https://github.com/textwire/textwire/issues/101).
- 🧑‍💻 Added type check to `@component`, `@insert` and `@reserve`. If the first argument you provide is not a string literal, you'll get a clear error. Before, the error wasn't clear. Closes [#106](https://github.com/textwire/textwire/issues/106).
- 🧑‍💻 Added error check if you are passing empty string as a `@component` name, `@insert` name or `@reserve` name.
- 🧑‍💻 Added `Empty` AST node for simplifying evaluator's logic. Closes [#104](https://github.com/textwire/textwire/issues/104).
- 🚀 Printing a string with `{{ }}` braces are not change `"` to `&#34;` and `'` to `&#39;` to correctly escape the string. Like before, you can use `str.raw()` function to get the raw output. If you call `{{ myStr.raw() }}` the string will not be mangled. Before, we would unescape string when you call `raw()` function but now there is not need to unescape them at all.

## [Release Notes V1](.github/CHANGELOG-V1.md)

## [Release Notes V2](.github/CHANGELOG-V2.md)

## [Release Notes V3](.github/CHANGELOG-V3.md)

## [Emojis Meaning](.github/EMOJIS.md)
