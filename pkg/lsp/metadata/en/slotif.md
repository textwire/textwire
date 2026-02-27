(directive)
Define a conditional default slotif in a component directive's body to provide a content to a placeholder if condition is true.

```textwire
@slotif(boolean)
    <p>content</p>
@end
```

Alternatively, you can define a conditional named slot.

```textwire
@slotif(boolean, 'name')
    <p>content</p>
@end
```
