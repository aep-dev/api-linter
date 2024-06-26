---
rule:
  aep: 162
  name: [core, '0162', commit-http-method]
  summary: Commit methods must use the POST HTTP verb.
permalink: /162/commit-http-method
redirect_from:
  - /0162/commit-http-method
---

# Commit methods: POST HTTP verb

This rule enforces that all `Commit` RPCs use the `POST` HTTP verb, as
mandated in [AEP-162][].

## Details

This rule looks at any method beginning with `Commit`, and complains
if the HTTP verb is anything other than `POST`. It _does_ check additional
bindings if they are present.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc CommitBook(CommitBookRequest) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}:commit"  // Should be `post:`.
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc CommitBook(CommitBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:commit"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0162::commit-http-method=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
rpc CommitBook(CommitBookRequest) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}:commit"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-162]: https://aep.dev/162
[aep.dev/not-precedent]: https://aep.dev/not-precedent
