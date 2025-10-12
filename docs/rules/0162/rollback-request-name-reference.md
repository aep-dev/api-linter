---
rule:
  aep: 162
  name: [core, '0162', rollback-request-name-reference]
  summary: |
    Rollback requests should annotate the `name` field with `google.api.resource_reference`.
permalink: /162/rollback-request-name-reference
redirect_from:
  - /0162/rollback-request-name-reference
---

# Rollback requests: Name resource reference

This rule enforces that all `Rollback` requests have
`google.api.resource_reference` on their `string name` field, as mandated in
[AEP-162][].

## Details

This rule looks at the `name` field of any message matching `Rollback*Request`
and complains if it does not have a `google.api.resource_reference` annotation.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message RollbackBookRequest {
  // The `google.api.resource_reference` annotation should also be included.
  string name = 1 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];

  string revision_id = 2 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

**Correct** code for this rule:

```proto
// Correct.
message RollbackBookRequest {
  string name = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  string revision_id = 2 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message RollbackBookRequest {
  // (-- api-linter: core::0162::rollback-request-name-reference=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string name = 1 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];

  string revision_id = 2 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-162]: https://aep.dev/162
[aep.dev/not-precedent]: https://aep.dev/not-precedent
