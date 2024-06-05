---
rule:
  aep: 131
  name: [core, '0131', http-method]
  summary: Get methods must use the GET HTTP verb.
permalink: /131/http-method
redirect_from:
  - /0131/http-method
---

# Get methods: GET HTTP verb

This rule enforces that all `Get` RPCs use the `GET` HTTP verb, as mandated in
[AEP-131][].

## Details

This rule looks at any message matching beginning with `Get`, and complains if
the HTTP verb is anything other than `GET`. It _does_ check additional bindings
if they are present.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc GetBook(GetBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{path=publishers/*/books/*}"  // Should be `get:`.
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc GetBook(GetBookRequest) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{path=publishers/*/books/*}"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0131::http-method=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
rpc GetBook(GetBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{path=publishers/*/books/*}"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-131]: https://aep.dev/131
[aep.dev/not-precedent]: https://aep.dev/not-precedent
