---
rule:
  aep: 4
  name: [core, '4', resource-variables]
  summary: Resource patterns should use consistent variable naming.
permalink: /4/resource-variables
redirect_from:
  - /4/resource-variables
---

# Resource pattern variables

This rule enforces that resource patterns use consistent variable naming
conventions, as described in [AEP-4][].

## Details

This rule scans all messages with `aep.api.resource` annotations, and
complains if variables in a `pattern` use camel case, or end in `_id`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    // Should be: publishers/{publisher}/books/{book}
    pattern: "publishers/{publisher_id}/books/{book_id}"
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
    pattern: "publishers/{publisher}/electronicBooks/{electronicBook}"
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
// (-- api-linter: core::4::resource-variables=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher_id}/books/{book_id}"
  };

  string path = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-4]: http://aep.dev/4
[aep.dev/not-precedent]: https://aep.dev/not-precedent
