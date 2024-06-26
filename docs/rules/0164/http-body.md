---
rule:
  aep: 164
  name: [core, '0164', http-body]
  summary: Undelete methods should use `*` as the HTTP body.
permalink: /164/http-body
redirect_from:
  - /0164/http-body
---

# Undelete methods: HTTP body

This rule enforces that all `Undelete` RPCs use `*` as the HTTP `body`, as mandated in
[AEP-164][].

## Details

This rule looks at any message beginning with `Undelete`, and complains
if the HTTP `body` field is anything other than `*`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc UndeleteBook(UndeleteBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:undelete"
    // body: "*" should be set.
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc UndeleteBook(UndeleteBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:undelete"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0164::http-body=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
rpc UndeleteBook(UndeleteBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:undelete"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-164]: https://aep.dev/164
[aep.dev/not-precedent]: https://aep.dev/not-precedent
