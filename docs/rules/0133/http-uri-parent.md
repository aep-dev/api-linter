---
rule:
  aep: 133
  name: [core, '0133', http-uri-parent]
  summary: Create methods must map the parent field to the URI.
permalink: /133/http-uri-parent
redirect_from:
  - /0133/http-uri-parent
---

# Create methods: HTTP URI parent field

This rule enforces that all `Create` RPCs map the `parent` field to the HTTP
URI, as mandated in [AEP-133][].

## Details

This rule looks at any message beginning with `Create`, and complains
if `parent` is not the only variable in the URI path. It _does_ check
additional bindings if they are present.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/publishers/*/books"  // The `parent` field should be extracted.
    body: "book"
  };
}
```

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    // The only variable should be `parent`.
    post: "/v1/{parent=publishers/*}/{book=books/*}"
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books"
    body: "book"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0133::http-uri-parent=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/publishers/*/books"
    body: "book"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-133]: https://aep.dev/133
[aep.dev/not-precedent]: https://aep.dev/not-precedent
