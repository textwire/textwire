(directive)
Define a default slot in a component to provide a placeholder for content.

```textwire
@slot
    <p>content</p>
@end
```

Alternatily, you can define a named slot in a component.

```textwire
@slot('name')
    <p>content</p>
@end
```