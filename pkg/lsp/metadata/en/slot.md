(directive)
Define a default slot inside a component directive's body to provide content for placeholder.

```textwire
@slot
    <p>content</p>
@end
```

Alternatively, you can define a named slot.

```textwire
@slot('name')
    <p>content</p>
@end
```
