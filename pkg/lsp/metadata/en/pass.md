(directive)
Define a `@pass` in a component body to provide a content to a named `@slot`.

```textwire
@component('my-component')
    @pass('name')
        <p>content</p>
    @end
@end
```