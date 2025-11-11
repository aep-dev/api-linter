---
rule:
  aep: 124
  name: [core, '0124', reference-same-package]
  summary: Resource references should refer to resources in the same package.
permalink: /124/reference-same-package
redirect_from:
  - /0124/reference-same-package
---

# Resource reference package

This rule enforces that resource reference annotations refer resources defined
in the same package, as described in [AEP-124][].

## Details

This rule scans all fields with `(aep.api.field_info).resource_reference` annotations,
and complains if the `type` on them refers to a resource that is defined in a
different protobuf package.

Certain common resource types are exempt from this rule.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
package google.example.library.v1;

message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  string name = 1;

  // ...
}
```

```proto
// Incorrect.
package google.example.libray.v1;  // Typo: Different package.

message GetBookRequest {
  string name = 1 [(aep.api.field_info).resource_reference = "library.googleapis.com/Book"];  // Lint warning: package mismatch.
}
```

**Correct** code for this rule:

```proto
// Correct.
package google.example.library;

message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  string name = 1;

  // ...
}

message GetBookRequest {
  string name = 1 [(aep.api.field_info).resource_reference = "library.googleapis.com/Book"];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
package google.example.library.common;

message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };
}
```

```proto
package google.example.library.v1;

message GetBookRequest {
  // (-- api-linter: core::0124::reference-same-package=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string name = 1 [(aep.api.field_info).resource_reference = "library.googleapis.com/Book"];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-124]: http://aep.dev/124
[aep.dev/not-precedent]: https://aep.dev/not-precedent
