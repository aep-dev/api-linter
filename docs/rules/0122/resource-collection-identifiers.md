---
rule:
  aep: 122
  name: [core, '0122', resource-collection-identifiers]
  summary: Resource patterns must use lowerCamelCase for collection identifiers.
permalink: /122/resource-collection-identifiers
redirect_from:
  - /0122/resource-collection-identifiers
---

# Resource pattern collection identifiers

This rule enforces that messages that have a `aep.api.resource` annotation
have properly formatted collection identifiers in each `pattern`, as described
in [AEP-122][].

## Details

This rule scans messages with a `aep.api.resource` annotation, and validates
the format of `pattern` collection identifiers, specifically that they are in
lowerCamelCase form and must start with a lowercase letter.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    // Collection identifiers must be kebab-case.
    pattern: "Publishers/{publisher}/publishedBooks/{book}"
  };
  string path = 1;
}
```

```proto
// Incorrect.
message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    // Collection identifiers must begin with a lower-cased letter.
    pattern: "/publishers/{publisher}/Published-books/{book}"
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
    pattern: "publishers/{publisher}/published-books/{book}"
  };
  string path = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.

```proto
// (-- api-linter: core::0122::resource-collection-identifiers=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "Publishers/{publisher}/Books/{book}"
  };
  string path = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-122]: http://aep.dev/122
[aep.dev/not-precedent]: https://aep.dev/not-precedent
