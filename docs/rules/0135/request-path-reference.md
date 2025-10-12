---
rule:
  aep: 135
  name: [core, '0135', request-path-reference]
  summary: |
    Delete RPCs should annotate the `path` field with `google.api.resource_reference`.
permalink: /135/request-path-reference
redirect_from:
  - /0135/request-path-reference
---

# Delete methods: Resource reference

This rule enforces that all `Delete` standard methods have
`google.api.resource_reference` on their `string path` field, as mandated in
[AEP-135][].

## Details

This rule looks at the `path` field of any message matching `Delete*Request`
and complains if it does not have a `google.api.resource_reference` annotation.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message DeleteBookRequest {
  // The `google.api.resource_reference` annotation should also be included.
  string path = 1 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

**Correct** code for this rule:

```proto
// Correct.
message DeleteBookRequest {
  string path = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message DeleteBookRequest {
  // (-- api-linter: core::0135::request-path-reference=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string path = 1 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-135]: https://aep.dev/135
[aep.dev/not-precedent]: https://aep.dev/not-precedent
