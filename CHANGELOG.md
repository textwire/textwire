# Release notes

## v4.0.0 (2026-04-01)

- ⚠️ Breaking changes:
    - Changed increment and decrement expressions into statements to match Go's style. `{{ x++ }}` and `{{ x-- }}`  are not returning any value now, they equivalent to `{{ x = x + 1 }}` and `{{ x = x - 1 }}`. Closes [#88](https://github.com/textwire/textwire/issues/88). [Read more about this change](https://textwire.github.io/blog/2026-03-20-textwire-v4).
    - Remove support for `@breakIf` and `@continueIf`, use lovercased versions `@breakif` and `@continueif`. Closes [#89](https://github.com/textwire/textwire/issues/89).
    - Printing a string with `{{ }}` braces are not change `"` to `&#34;` and `'` to `&#39;` to correctly escape the string. Like before, you can use `str.raw()` function to get the raw output. If you call `{{ myStr.raw() }}` the string will not be mangled. Before, we would unescape string when you call `raw()` function but now there is not need to unescape them at all.
    - Printing a string with `{{ }}` braces are not change `"` to `&#34;` and `'` to `&#39;` to correctly escape the string. Like before, you can use `str.raw()` function to get the raw output. If you call `{{ myStr.raw() }}` the string will not be mangled. Before, we would unescape string when you call `raw()` function but now there is not need to unescape them at all.
    - Fixed issue with Go's `time.Time` type being converted to Textwire's empty object. Now, it's converted to a string format like `2006-01-02 15:04:05`. Closes [#110](https://github.com/textwire/textwire/issues/110).
    - Global functions `hasValue` and `defined` are now require at least one argument. You'll get parse error if you don't provide any arguments.
    - Changed the way components work. Complete rewrite due to incorrect implementation:
      - Instead of passing default slot like this `@slot<content>@end` you just pass content inside of a component body like that `@component('my-comp')<content>@end`.
      - Now, you should use directive `@pass('name')` instead of `@slot('name')` inside of a component body when you want to pass content from template file to component file. Example:
        ```blade
        @component('user')
            <h1>Content for default slot</h1>
            @pass('name')
                <h2>{{ user.name }}</h2>
            @end
        @end
        ```
      - Now, components always require the ending directive even if they are empty. Example: `@component('name')@end`.
    - Blocks are not trimmed from both, left and the right side. Example: `@if(true) content @end` will produce `content` instead of ` content ` like before.
- 🧑‍💻 Improvements:
    - All public API functions like `NewTemplate()`, `EvaluateString`, etc., now return `*fail.Error` instead of Go's `error` type.
    - Added proper position to error messages. Closes [#101](https://github.com/textwire/textwire/issues/101).
    - Added type check to `@component`, `@insert` and `@reserve`. If the first argument you provide is not a string literal, you'll get a clear error. Before, the error wasn't clear. Closes [#106](https://github.com/textwire/textwire/issues/106).
    - Added error check if you are passing empty string as a `@component` name, `@insert` name or `@reserve` name.
    - Added `Empty` AST node for simplifying evaluator's logic. Closes [#104](https://github.com/textwire/textwire/issues/104).
    - Added trailing commas to components. Example: `@component('name',)@end`.
- ♻️ Refactored internal logic for evaluator. Separated literal values from nonliteral. Closes [#90](https://github.com/textwire/textwire/issues/90).
- ✨ Added [formatDate](https://textwire.github.io/v4/functions/global#formatdate) global function which converts string date. Closes [#111](https://github.com/textwire/textwire/issues/111).

## [Release Notes V1](.github/CHANGELOG-V1.md)

## [Release Notes V2](.github/CHANGELOG-V2.md)

## [Release Notes V3](.github/CHANGELOG-V3.md)

## [Emojis Meaning](.github/EMOJIS.md)
