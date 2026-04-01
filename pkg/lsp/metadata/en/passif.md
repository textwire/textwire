(directive)
Define a conditional `@passif` in a component body to pass the content to a named `@slot` if condition is `true`.

```textwire
@component('my-component')
    @passif(boolean, 'name')
        <p>content</p>
    @end
@end
```