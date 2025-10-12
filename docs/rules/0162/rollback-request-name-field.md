---
rule:
  aep: 162
  name: [core, '0162', rollback-request-name-field]
  summary: Rollback RPCs must have a `name` field in the request.
permalink: /162/rollback-request-name-field
redirect_from:
  - /0162/rollback-request-name-field
---

# Rollback requests: Name field

This rule enforces that all `Rollback` methods have a `string name`
field in the request message, as mandated in [AEP-162][].

## Details

This rule looks at any message matching `Rollback*Request` and complains if
either the `name` field is missing or it has any type other than `string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should include a `string name` field.
message RollbackBookRequest {
  string revision_id = 2 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

```proto
// Incorrect.
message RollbackBookRequest {
  // Field type should be `string`.
  bytes name = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

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

If you need to violate this rule, use a leading comment above the message (if
the `name` field is missing) or above the field (if it is the wrong type).
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message RollbackBookRequest {
  // (-- api-linter: core::0162::rollback-request-name-field=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  bytes name = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  string revision_id = 2 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-162]: https://aep.dev/162
[aep.dev/not-precedent]: https://aep.dev/not-precedent
