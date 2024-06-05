---
rule:
  aep: 131
  name: [core, '0131', method-signature]
  summary: |
    Get RPCs should annotate a method signature of "path".
permalink: /131/method-signature
redirect_from:
  - /0131/method-signature
---

# Get methods: Method signature

This rule enforces that all `Get` standard methods have a
`google.api.method_signature` annotation with a value of `"path"`, as mandated
in [AEP-131][].

## Details

This rule looks at any method beginning with `Get`, and complains if the
`google.api.method_signature` annotation is missing, or if it is set to any
value other than `"path"`. Additional method signatures, if present, are
ignored.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc GetBook(GetBookRequest) returns (Book) {
  // A google.api.method_signature annotation should be present.
}
```

```proto
// Incorrect.
rpc GetBook(GetBookRequest) returns (Book) {
  option (google.api.method_signature) = "book";  // Should be "path".
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc GetBook(GetBookRequest) returns (Book) {
  option (google.api.method_signature) = "path";
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0131::method-signature=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
rpc GetBook(GetBookRequest) returns (Book);
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-131]: https://aep.dev/131
[aep.dev/not-precedent]: https://aep.dev/not-precedent
