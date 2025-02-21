(directive)
Inject content into reserved placeholders defined in the layout file by providing a block of content.

```textwire
@insert('reservedName')
    <p>content</p>
@end
```

As an alternative, you can provide a second argument as content instead of using a block.

```textwire
@insert('reservedName', 'content')
```

Use this directive in your templates to inject content into reserved placeholders defined in the layout file with `@reserve` directive.
