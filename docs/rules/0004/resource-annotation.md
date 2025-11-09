---
rule:
  aep: 4
  name: [core, '4', resource-annotation]
  summary: Resource messages should be annotated with `aep.api.resource`.
permalink: /4/resource-annotation
redirect_from:
  - /4/resource-annotation
---

# Resource annotation presence

This rule enforces that top-level messages that appear to represent resources
have a `aep.api.resource` annotation, as described in [AEP-4][].

## Details

This rule scans all top-level messages, and assumes that messages with a
 `string path` field are resources unless the message name ends with `Request`.
For messages that are resources, it complains if the `aep.api.resource`
annotation is missing.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  // A `aep.api.resource` annotation should be here.
  string path = 1;
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  string path = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.

```proto
// (-- api-linter: core::4::resource-annotation=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
message Book {
  string path = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-4]: http://aep.dev/4
[aep.dev/not-precedent]: https://aep.dev/not-precedent
