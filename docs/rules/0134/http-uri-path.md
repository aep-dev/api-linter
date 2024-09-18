---
rule:
  aep: 134
  name: [core, '0134', http-uri-path]
  summary: Update methods must map the resource's path field to the URI.
permalink: /134/http-uri-path
redirect_from:
  - /0134/http-uri-path
---

# Update methods: HTTP URI path field

This rule enforces that all `Update` RPCs map the `path` field from the
resource object to the HTTP URI, as mandated in [AEP-134][].

## Details

This rule looks at any message matching beginning with `Update`, and complains
if the `path` variable from the resource (not the request message) is not
included in the URI. It _does_ check additional bindings if they are present.

Additionally, if the resource uses multiple words, it ensures that word
separation uses `snake_case`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc UpdateBookRequest(UpdateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{book.path=publishers/*/books/*}"  // Should be `path`
    body: "book"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc UpdateBookRequest(UpdateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{path=publishers/*/books/*}"
    body: "book"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0134::http-uri-path=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
rpc UpdateBookRequest(UpdateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{book.path=publishers/*/books/*}"
    body: "book"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-134]: https://aep.dev/134
[aep.dev/not-precedent]: https://aep.dev/not-precedent
