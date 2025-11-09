---
rule:
  aep: 4
  name: [core, '4', duplicate-resource]
  summary: Resource types should not be defined more than once.
permalink: /4/duplicate-resource
redirect_from:
  - /4/duplicate-resource
---

# Resource annotation presence

This rule enforces that the same resource type doesn't appear in more than one
`aep.api.resource` annotation, as described in [AEP-004][].

## Details

This rule complains about messages that have the same `type` for the
`aep.api.resource` annotation, which frequently occur due to copy-paste
errors and messages spread across multiple files and/or packages. Duplicate
resource definitions can cause compilation problems in generated client code.

## Examples

**Incorrect** code for this rule:

```proto
message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  string path = 1;
}

message Author {
  option (aep.api.resource) = {
    // Incorrect: should be "library.googleapis.com/Author".
    type: "library.googleapis.com/Book"
    pattern: "authors/{author}"
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

message Author {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Author"
    pattern: "authors/{author}"
  };

  string path = 1;
}
```

## Disabling

If you need to violate this rule, use a comment at the top of the file.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0004::duplicate-resource=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
syntax = "proto3";

message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  string path = 1;
}

message Author {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "authors/{author}"
  };

  string path = 1;
}
```

[aep-004]: http://aep.dev/123
[aep.dev/not-precedent]: https://aep.dev/not-precedent
