---
rule:
  aep: 134
  name: [core, '0134', request-path-required]
  summary: Update RPCs must have a `path` field in the request.
permalink: /134/request-path-required
redirect_from:
  - /0134/request-path-required
---

# Update methods: Name field

This rule enforces that all `Update` standard methods have a `string path`
field in the request message, as mandated in [AEP-134][].

## Details

This rule looks at any message matching `Update*Request` and complains if
the `path` field is missing.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message UpdateBookRequest {
  // Field path should be `path`.
  string book = 1;
}
```

**Correct** code for this rule:

```proto
// Correct.
message UpdateBookRequest {
  string path = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0134::request-path-required=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
message UpdateBookRequest {
  string book = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-134]: https://aep.dev/134
[aep.dev/not-precedent]: https://aep.dev/not-precedent
