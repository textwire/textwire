(directive)
Conditionally render content with additional conditions using `@elseif`.

```textwire
@if(condition1)
    <p>condition1 is true</p>
@elseif(condition2)
    <p>condition2 is true</p>
@end
```

Use the `@elseif` directive to handle additional conditional branches. If none of the conditions are met, use `@else` to provide fallback content.