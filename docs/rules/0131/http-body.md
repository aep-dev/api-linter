---
rule:
  aep: 131
  name: [core, '0131', http-body]
  summary: Get methods must not have an HTTP body.
permalink: /131/http-body
redirect_from:
  - /0131/http-body
---

# Get methods: No HTTP body

This rule enforces that all `Get` RPCs omit the HTTP `body`, as mandated in
[AEP-131][].

## Details

This rule looks at any message matching beginning with `Get`, and complains if
the HTTP `body` field is set.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc GetBook(GetBookRequest) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{path=publishers/*/books/*}"
    body: "*"  // This should be absent.
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
// (-- api-linter: core::0131::http-body=disabled
//     api-linter: core::0131::http-method=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
rpc GetBook(GetBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{path=publishers/*/books/*}"
    body: "*"
  };
}
```

**Important:** HTTP `GET` requests are unable to have an HTTP body, due to the
nature of the protocol. The only valid way to include a body is to also use a
different HTTP method (as depicted above).

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-131]: https://aep.dev/131
[aep.dev/not-precedent]: https://aep.dev/not-precedent
