---
rule:
  aep: 4
  name: [core, '4', resource-name-components-alternate]
  summary: |
    Resource name components should alternate between collection and
    identifiers.
permalink: /4/resource-name-components-alternate
redirect_from:
  - /4/resource-name-components-alternate
---

# Resource name components alternate

This rule enforces that messages that have a `aep.api.resource` annotation
have `pattern` annotations that alternate between collection and identifier, as
described in [AEP-4][].

## Details

This rule scans messages with a `aep.api.resource` annotation, and validates
that each `pattern` alternated between collection and identifiers.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    // two collections next to each other.
    pattern: "publishers/books/{book}"
  };
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
// (-- api-linter: core::4::resource-name-components-alternate=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/books/{book}"
  };
  string path = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-4]: http://aep.dev/4
[aep.dev/not-precedent]: https://aep.dev/not-precedent
