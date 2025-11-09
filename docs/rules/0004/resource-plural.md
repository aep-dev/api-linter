---
rule:
  aep: 4
  name: [core, '4', resource-plural]
  summary: Resource plural is required
permalink: /4/resource-plural
redirect_from:
  - /4/resource-plural
---

# Resource type name

This rule enforces that messages that have a `aep.api.resource` annotation,
have a properly formatted `plural`, as described in [AEP-4][].

## Details

This rule scans messages with a `aep.api.resource` annotation, and
verifies the `plural` field exists.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  // no plural annotation
  option (aep.api.resource) = {
    type: "library.googleapis.com/BookShelf"
    pattern: "publishers/{publisher}/bookShelves/{book_shelf}"
  };

  string path = 1;
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/BookShelf"
    pattern: "publishers/{publisher}/bookShelves/{book_shelf}"
    plural: "bookShelves",
  };

  string path = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.

```proto
// (-- api-linter: core::4::resource-plural=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Genre/Mystery/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  string path = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-4]: http://aep.dev/4
[aep.dev/not-precedent]: https://aep.dev/not-precedent
