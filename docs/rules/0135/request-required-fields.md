---
rule:
  aep: 135
  name: [core, '0135', request-required-fields]
  summary: Delete RPCs must not have unexpected required fields in the request.
permalink: /135/request-required-fields
redirect_from:
  - /0135/request-required-fields
---

# Delete methods: Required fields

This rule enforces that all `Delete` standard methods do not have unexpected
required fields, as mandated in [AEP-135][].

## Details

This rule looks at any message matching `Delete*Request` and complains if it
comes across any required fields other than:

- `string path` ([AEP-135][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message DeleteBookRequest {
  string path = 1 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
  // Non-standard required field.
  bool allow_missing = 4 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

**Correct** code for this rule:

```proto
// Correct.
message DeleteBookRequest {
  string path = 1 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
  bool allow_missing = 4 [(aep.api.field_info).field_behavior = OPTIONAL];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message DeleteBookRequest {
  string path = 1 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
  // (-- api-linter: core::0135::request-required-fields=disabled
  //     aep.dev/not-precedent: We really need this field to be required because
  // reasons. --)
  bool allow_missing = 4 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-135]: https://aep.dev/135
[aep.dev/not-precedent]: https://aep.dev/not-precedent
