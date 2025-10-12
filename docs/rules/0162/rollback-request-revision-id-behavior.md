---
rule:
  aep: 162
  name: [core, '0162', rollback-request-revision-id-behavior]
  summary: |
    Rollback requests should annotate the `revision_id` field with `aep.api.field_behavior`.
permalink: /162/rollback-request-revision-id-behavior
redirect_from:
  - /0162/rollback-request-revision-id-behavior
---

# Rollback requests: Revision ID field behavior

This rule enforces that all `Rollback` requests have
`aep.api.field_behavior` set to `REQUIRED` on their `string revision_id` field, as
mandated in [AEP-162][].

## Details

This rule looks at any message matching `Rollback*Request` and complains if the
`revision_id` field does not have a `aep.api.field_behavior` annotation with a
value of `REQUIRED`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message RollbackBookRequest {
  string name = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  // The `aep.api.field_behavior` annotation should be included.
  string revision_id = 2;
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
  string name = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  // (-- api-linter: core::0162::rollback-request-revision-id-behavior=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string revision_id = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-162]: https://aep.dev/162
[aep.dev/not-precedent]: https://aep.dev/not-precedent
