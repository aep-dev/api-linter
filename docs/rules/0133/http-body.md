---
rule:
  aep: 133
  name: [core, '0133', http-body]
  summary: Create methods must have the HTTP body set to the resource.
permalink: /133/http-body
redirect_from:
  - /0133/http-body
---

# Create methods: HTTP body

This rule enforces that all `Create` RPCs set the HTTP `body` to the resource,
as mandated in [AEP-133][].

## Details

This rule looks at any message matching beginning with `Create`, and complains
if the HTTP `body` field is not set to the resource being created.

Note that any `additional_bindings` need their own `body` field.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books"
    body: "*"  // This should be "book".
  };
}
```

```proto
// Incorrect.
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books"
    body: "book"
    additional_bindings: {
      post: "/v1/books"
      // There should be a "body" here too.
    }
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


```proto
// Correct.
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books"
    body: "book"
    additional_bindings: {
      post: "/v1/books"
      body: "book"
    }
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0133::http-body=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-133]: https://aep.dev/133
[aep.dev/not-precedent]: https://aep.dev/not-precedent
