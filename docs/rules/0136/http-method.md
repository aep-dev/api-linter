---
rule:
  aep: 136
  name: [core, '0136', http-method]
  summary: Custom methods must use the POST or GET HTTP verb.
permalink: /136/http-method
redirect_from:
  - /0136/http-method
---

# Custom methods: HTTP method

This rule enforces that all custom methods use the `POST` or `GET` HTTP verbs,
as mandated in [AEP-136][].

## Details

This rule looks at any method that is not a standard method, and complains if
the HTTP verb is anything other than `POST` or `GET`. It _does_ check
additional bindings if they are present.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc CheckoutBook(CheckoutBookRequest) returns (CheckoutBookResponse) {
  option (google.api.http) = {
    put: "/v1/{path=publishers/*/books/*}:checkout"  // Should be `post:`.
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc CheckoutBook(CheckoutBookRequest) returns (CheckoutBookResponse) {
  option (google.api.http) = {
    post: "/v1/{path=publishers/*/books/*}:checkout"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0136::http-method=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
rpc CheckoutBook(CheckoutBookRequest) returns (CheckoutBookResponse) {
  option (google.api.http) = {
    put: "/v1/{path=publishers/*/books/*}:checkout"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-136]: https://aep.dev/136
[aep.dev/not-precedent]: https://aep.dev/not-precedent
