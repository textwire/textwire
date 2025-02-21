(directive)
Else directive can be used with if statements, for loops and each loops. Example usage with a for loop:

```textwire
@for(i = 0; i < items.len(); i++)
    <p>{{ items[i] }}</p>
@else
    <p>No items available.</p>
@end
```

Example usage with an if statement:

```textwire
@if(items.len() > 0)
    <p>Items available.</p>
@else
    <p>No items available.</p>
@end
