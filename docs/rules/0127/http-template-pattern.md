---
rule:
  aep: 127
  name: [core, '0127', http-template-pattern]
  summary: |
    HTTP template variable patterns should match the patterns defined by their resources.
permalink: /127/http-template-pattern
redirect_from:
  - /0127/http-template-pattern
---

# HTTP Pattern Variables

This rule enforces that any HTTP annotations that reference a resource must
match one of the pattern strings defined by that resource, as mandated in
[AEP-127][].

## Details

This rule ensures that `google.api.http` path template variables that represent
a resource path match one of the resource path patterns of the resource that the
field being referenced represents.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// The template for the `path` variable in the `google.api.http` annotation
// is missing segments from the Book message's `pattern`.
rpc GetBook(GetBookRequest) returns (Book) {
    option (google.api.http) = {
        get: "v1/{path=shelves/*}"
    };
}
message GetBookRequest {
    string path = 1 [
        (aep.api.field_info).resource_reference = "library.googleapis.com/Book"
    ];
}
message Book {
    option (google.api.resource) = {
        type: "library.googleapis.com/Book"
        pattern: "shelves/{shelf}/books/{book}"
    };

    // Book resource path.
    string path = 1;
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc GetBook(GetBookRequest) returns (Book) {
    option (google.api.http) = {
        get: "v1/{path=shelves/*/books/*}"
    };
}
message GetBookRequest {
    string path = 1 [
        (aep.api.field_info).resource_reference = "library.googleapis.com/Book"
    ];
}
message Book {
    option (google.api.resource) = {
        type: "library.googleapis.com/Book"
        pattern: "shelves/{shelf}/books/{book}"
    };

    // Book resource path.
    string path = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0127::http-template-pattern=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
rpc GetBook(GetBookRequest) returns (Book) {
    option (google.api.http) = {
        get: "v1/{path=shelves/*}"
    };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-127]: https://aep.dev/127
[aep.dev/not-precedent]: https://aep.dev/not-precedent