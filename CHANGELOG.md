# Release notes for 1.x

## v1.3.0 (2024-03-02)

- Added `@breakIf` and `@continueIf` directives for the `@each` loop.

## v1.2.0 (2024-03-02)

- Added `@break` and `@continue` directives for the `@each` loop.
- Added `@break` and `@continue` directives for the `@for` loop.
- Fixed bug where you couldn't compare strings with the `==` operator. For example, `@if(n == "hello")` was returning error.

## v1.1.0 (2024-02-28)

- Added `@else` directive for the `@each` loop which will be executed if the loop has no items to iterate over.
- Added `@else` directive for the `@for` loop which will be executed if the loop has no items to iterate over.
- Code refactoring and naming improvements.

## v1.0.0 (2024-02-26)

- Initial release of the first stable version with support for all [the features](https://textwire.github.io/1.x/language-elements/) that we wanted to have in the 1.0.0 version.