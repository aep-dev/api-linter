---
rule:
  aep: 140
  name: [core, '0140', uri]
  summary: Field names should prefer `uri` to `url`.
permalink: /140/uri
redirect_from:
  - /0140/uri
---

# Field names: URIs

This rule enforces that field names use `uri` rather than `url`, as mandated in
[AEP-140][].

## Details

This rule checks every field in the proto and complains if the field name
includes `url`. (It splits on `_` to avoid false positives.)

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string url = 1;  // Should be `uri`.
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string uri = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0140::uri=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
message Book {
  string url = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-140]: https://aep.dev/140
[aep.dev/not-precedent]: https://aep.dev/not-precedent
