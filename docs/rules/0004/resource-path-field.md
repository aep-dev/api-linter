---
rule:
  aep: 4
  name: [core, '4', resource-path-field]
  summary: Resource messages should have a `string path` field.
permalink: /4/resource-path-field
redirect_from:
  - /4/resource-path-field
---

# Resource `path` field

This rule enforces that messages that appear to represent resources have a
`string path` field, as described in [AEP-4][].

## Details

This rule scans all messages that have a `aep.api.resource` annotation, and
complains if the `path` field is missing or if it is any type other than
singular `string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect: missing `string path` field.
message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };
}
```

```proto
// Incorrect.
message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  // Should be `string`, not `bytes`.
  bytes path = 1;
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

If you need to violate this rule, use a leading comment above the message, or
above the field if it is the wrong type.

```proto
// (-- api-linter: core::4::resource-path-field=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-4]: http://aep.dev/4
[aep.dev/not-precedent]: https://aep.dev/not-precedent
