(directive)
Components help to organize and structure templates by encapsulating reusable parts of your user interface.

```textwire
@component('path/to', { prop })
```

Optionally, you can include a block of content within the component through a slots.

```textwire
@component('path/to', { prop })
    @slot
        <p>header</p>
    @end

    @slot('name')
        <p>footer</p>
    @end
@end
```

Use this directive to include a component into your template.