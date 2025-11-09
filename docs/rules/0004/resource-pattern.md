---
rule:
  aep: 4
  name: [core, '4', resource-pattern]
  summary: Resource annotations should define a pattern.
permalink: /4/resource-pattern
redirect_from:
  - /4/resource-pattern
---

# Resource patterns

This rule enforces that messages that appear to represent resources have a
`pattern` defined on their `aep.api.resource` annotation, as described in
[AEP-4][].

## Details

This rule scans all messages with `aep.api.resource` annotations, and
complains if `pattern` is not provided at least once. It also complains if the
segments outside of variable names contain underscores.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    // pattern should be here
  };

  string path = 1;
}
```

```proto
// Incorrect.
message ElectronicBook {
  option (aep.api.resource) = {
    type: "library.googleapis.com/ElectronicBook"
    // Should be: publishers/{publisher}/electronicBooks/{electronic_book}
    pattern: "publishers/{publisher}/electronic_books/{electronic_book}"
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

```proto
// Correct.
message ElectronicBook {
  option (aep.api.resource) = {
    type: "library.googleapis.com/ElectronicBook"
    pattern: "publishers/{publisher}/electronicBooks/{electronic_book}"
  };

  string path = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.

```proto
// (-- api-linter: core::4::resource-pattern=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
  };

  string path = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-4]: http://aep.dev/4
[aep.dev/not-precedent]: https://aep.dev/not-precedent
