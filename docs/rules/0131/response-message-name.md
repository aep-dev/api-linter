---
rule:
  aep: 131
  name: [core, '0131', response-message-name]
  summary: Get methods must return the resource.
permalink: /131/response-message-name
redirect_from:
  - /0131/response-message-name
---

# Get methods: Response message

This rule enforces that all `Get` RPCs have a response message of the resource,
as mandated in [AEP-131][].

## Details

This rule looks at any message matching beginning with `Get`, and complains if
the name of the corresponding output message does not match the name of the RPC
with the prefix `Get` removed.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc GetBook(GetBookRequest) returns (GetBookResponse) {  // Should be `Book`.
  option (google.api.http) = {
    get: "/v1/{path=publishers/*/books/*}"
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
// (-- api-linter: core::0131::response-message-path=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
rpc GetBook(GetBookRequest) returns (GetBookResponse) {
  option (google.api.http) = {
    get: "/v1/{path=publishers/*/books/*}"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-131]: https://aep.dev/131
[aep.dev/not-precedent]: https://aep.dev/not-precedent
