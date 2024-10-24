# Release notes

* 🧑‍💻 - Improving developer experience
* ✨ - Introducing new features
* 🐛 - Fixing bugs
* ♻️ - Refactoring code
* 📝 - Adding or updating documentation

## v2.1.0 (2024-10-24)
For more detailed information about this release, read the [Textwire v2.1.0 Release Notes](https://textwire.github.io/blog/2024/10/23/textwire-v2.1.0-release-notes)

- ✨ Features
    - For array literals, added `4` built-in functions: `rand`, `reverse`, `slice`, `shuffle`
    - For integer literals, added `2` built-in functions: `abs`, `str`
    - For float literals, added `5` built-in functions: `abs`, `ceil`, `floor`, `round`, `str`
    - For string literals, added `3` built-in functions: `capitalize`, `reverse`, `contains`
    - For boolean literals, added `1` built-in function: `binary`
    - New error page while rendering a template. Instead of black screen we now get a simple error page with `Sorry! We’re having some trouble right now. Please check back shortly`
- 🧑‍💻 Improvements
    - Improve error handling for built-in function. If you pass the wrong argument type it will generate an error
    - Improve error handling for custom functions. Now, when you use function that is not defined, you'll get an error that the function x doesn't exists on type y
    - Improve error handling for division by zero cases
- 📝 Remove `CONTRIBUTING.md` file
- 🐛 Fixed bug with incorrect precedence with prefixed expressions like `{{ -1.abs() }}`. This expression would left out the `-`

## v2.0.0 (2024-10-18)
- ♻️ [BREAKING CHANGE!] Moved `textwire.Config` to a separate package `config.Config`
- ✨ Added the ability to register your own custom functions for specific types and use them in your Textwire code like built-in functions. If you are upgrading from version 1, make these changes:
    1. Change all the imports from `github.com/textwire/textwire` to `github.com/textwire/textwire/v2`
    2. Run `go mod tidy` to update the dependencies
    3. Change the package name from `textwire.Config` to `config.Config` in your code if you use configuration and import `"github.com/textwire/textwire/v2/config"`. If you already have a package named `config`, you can alias the import like `twconfig "github.com/textwire/textwire/v2/config"`

## v1.7.1 (2024-09-08)
- 🧑‍💻 Improve error handling for component slots. When you pass a slot that isn't defined, you'll get an error
- 🧑‍💻 Improve error handling for component slots. When you pass multiple slots with the same name, you'll get an error

## v1.7.0 (2024-09-05)
- ✨ Added `upper` function to strings. For example, `{{ "hello".upper() }}` will print `HELLO`
- ✨ Added `lower` function to strings. For example, `{{ "HELLO".lower() }}` will print `hello`
- 🐛 Fixed bug that was appearing if you put HTML after the `@insert` directive. For example, `@insert('content', 'nice')<h2>Text</h2>` would result in error
- 🐛 Fixed bug where you couldn't define `@component("person")` directives without the second argument
- ✨ Added `@slot` directive for components. You can define slots in components and then pass content to them when using the component

## v1.6.1 (2024-08-22)
- 🧑‍💻 Improve `join` function for arrays by adding default separator ",". If you don't provide a separator, it will use a comma as a default separator
- 📝 Added emojis for each changelog item

## v1.6.0 (2024-08-22)
- 🧑‍💻 Improve error handling for functions. If you call a function that doesn't exist, it will not only tell that function doesn't exist, but also the type of the target that you are trying to call a function on
- ✨ Added `join` function to arrays. For example, `{{ arr = [1, 2, 3]; arr.join(", ") }}` will print `1, 2, 3`

## v1.5.2 (2024-03-22)
- 🐛 Fixed a bug where you were you couldn't comment a block of Textwire code
- 🧑‍💻 Added error check when executing `@each` statement. Occasionally, if passed array was invalid, it would panic. Now, it will return an error message

## v1.5.1 (2024-03-21)
- 🧑‍💻 Removed escaping single and double quotes in strings when printing them. For example, `{{ "Hello, 'world'" }}` and `{{ 'Hello, "world"' }}` will now print `Hello, 'world'` and `Hello, "world"` respectively instead of using HTML entities to escape the quotes

## v1.5.0 (2024-03-21)
- ✨ Added trailing comma support in object and array literals. For example, `{{ obj = { key: "value", } }}` and `{{ arr = [1, 2, 3, ] }}` are now valid
- ✨ Added support for comment with `{{-- --}}` syntax. For example, `{{-- This is a comment --}}`

## v1.4.1 (2024-03-18)
- 📝 Added link to a [Textwire VSCode extension link](https://marketplace.visualstudio.com/items?itemName=SerhiiCho.textwire) in the README.md file

## v1.4.0 (2024-03-11)
- ✨ Added `@component` directive for creating reusable components
- ♻️ Simplified parsing logic for `@use`, `@insert` and `@reserve` directives
- ✨ Added ability to put `@use`, `@insert` and `@reserve` directives inside any other directive like `@if`, `@each`, `@for`, `@component` etc. Previously, you could only put them at the top level of the template
- ✨ Added ability to define objects key without quotes. For example, `{{ obj = { key: "value" } }}` is now the same as `{{ obj = { "key": "value" } }}`
- ✨ Added shorthand property notation for objects, similar to JavaScript. For example, `{{ obj = { key } }}` is now the same as `{{ obj = { "key": key } }}`
- ✨ Changed so that `@use` statement is not required to be the first in the file. Now, you can place it anywhere in the file

## v1.3.0 (2024-03-04)
- ✨ Added `@breakIf` and `@continueIf` directives for the `@each` loop
- ✨ Added `@breakIf` and `@continueIf` directives for the `@for` loop

## v1.2.0 (2024-03-02)
- ✨ Added `@break` and `@continue` directives for the `@each` loop
- ✨ Added `@break` and `@continue` directives for the `@for` loop
- 🐛 Fixed bug where you couldn't compare strings with the `==` operator. For example, `@if(n == "hello")` was returning error

## v1.1.0 (2024-02-28)
- ✨ Added `@else` directive for the `@each` loop which will be executed if the loop has no items to iterate over
- ✨ Added `@else` directive for the `@for` loop which will be executed if the loop has no items to iterate over
- ♻️ Code refactoring and naming improvements

## v1.0.0 (2024-02-26)
- Initial release of the first stable version with support for all [the features](https://textwire.github.io/1.x/language-elements/) that we wanted to have in the 1.0.0 version