# Release notes

## v1.5.2 (2024-03-22)

- Fixed bug where you were you couldn't comment a block of Textwire code

## v1.5.1 (2024-03-21)

- Removed escaping single and double quotes in strings when printing them. For example, `{{ "Hello, 'world'" }}` and `{{ 'Hello, "world"' }}` will now print `Hello, 'world'` and `Hello, "world"` respectively instead of using HTML entities to escape the quotes

## v1.5.0 (2024-03-21)

- Added trailing comma support in object and array literals. For example, `{{ obj = { key: "value", } }}` and `{{ arr = [1, 2, 3, ] }}` are now valid
- Added support for comment with `{{-- --}}` syntax. For example, `{{-- This is a comment --}}`

## v1.4.1 (2024-03-18)

- Added link to a [Textwire VSCode extension link](https://marketplace.visualstudio.com/items?itemName=SerhiiCho.textwire) in the README.md file

## v1.4.0 (2024-03-11)

- Added `@component` directive for creating reusable components
- Simplified parsing logic for `@use`, `@insert` and `@reserve` directives
- Added ability to put `@use`, `@insert` and `@reserve` directives inside any other directive like `@if`, `@each`, `@for`, `@component` etc. Previously, you could only put them at the top level of the template
- Added ability to define objects key without quotes. For example, `{{ obj = { key: "value" } }}` is now the same as `{{ obj = { "key": "value" } }}`
- Added shorthand property notation for objects, similar to JavaScript. For example, `{{ obj = { key } }}` is now the same as `{{ obj = { "key": key } }}`
- Changed so that `@use` statement is not required to be the first in the file. Now, you can place it anywhere in the file

## v1.3.0 (2024-03-04)

- Added `@breakIf` and `@continueIf` directives for the `@each` loop
- Added `@breakIf` and `@continueIf` directives for the `@for` loop

## v1.2.0 (2024-03-02)

- Added `@break` and `@continue` directives for the `@each` loop
- Added `@break` and `@continue` directives for the `@for` loop
- Fixed bug where you couldn't compare strings with the `==` operator. For example, `@if(n == "hello")` was returning error

## v1.1.0 (2024-02-28)

- Added `@else` directive for the `@each` loop which will be executed if the loop has no items to iterate over
- Added `@else` directive for the `@for` loop which will be executed if the loop has no items to iterate over
- Code refactoring and naming improvements

## v1.0.0 (2024-02-26)

- Initial release of the first stable version with support for all [the features](https://textwire.github.io/1.x/language-elements/) that we wanted to have in the 1.0.0 version