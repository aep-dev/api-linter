---
rule:
  aep: 4
  name: [core, '4', path-never-optional]
  summary: Resource path fields must never be labeled with proto3_optional.
permalink: /4/path-never-optional
redirect_from:
  - /4/path-never-optional
---

# Resource path field never optional

This rule enforces that the path field of a resource message is not labeled with
proto3_optional.

## Details

This rule scans for messages with a `aep.api.resource` annotation and ensures
that the configured path field (either `path` or whichever field specified via
`path_field`) is not labeled as `optional`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  // The path field should not be labeled as optional.
  optional string path = 1;
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
// (-- api-linter: core::04::path-never-optional=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  optional string path = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep.dev/not-precedent]: https://aep.dev/not-precedent
