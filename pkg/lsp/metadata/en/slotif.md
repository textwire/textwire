(directive)
Define a conditional default slotIf in a component directive's body to provide a content to a placeholder if condition is true.

```textwire
@slotIf(boolean)
    <p>content</p>
@end
```

Alternatively, you can define a conditional named slot.

```textwire
@slotIf(boolean, 'name')
    <p>content</p>
@end
```
