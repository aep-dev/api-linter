---
rule:
  aep: 162
  name: [core, '0162', rollback-request-name-behavior]
  summary: |
    Rollback requests should annotate the `name` field with `google.api.field_behavior`.
permalink: /162/rollback-request-name-behavior
redirect_from:
  - /0162/rollback-request-name-behavior
---

# Rollback requests: Name field behavior

This rule enforces that all `Rollback` requests have
`google.api.field_behavior` set to `REQUIRED` on their `string name` field, as
mandated in [AEP-162][].

## Details

This rule looks at any message matching `Rollback*Request` and complains if the
`name` field does not have a `google.api.field_behavior` annotation with a
value of `REQUIRED`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message RollbackBookRequest {
  // The `google.api.field_behavior` annotation should also be included.
  string name = 1 [
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  string revision_id = 2 [(google.api.field_behavior) = REQUIRED];
}
```

**Correct** code for this rule:

```proto
// Correct.
message RollbackBookRequest {
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  string revision_id = 2 [(google.api.field_behavior) = REQUIRED];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message RollbackBookRequest {
  // (-- api-linter: core::0162::rollback-request-name-behavior=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string name = 1 [
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  string revision_id = 2 [(google.api.field_behavior) = REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-162]: https://aep.dev/162
[aep.dev/not-precedent]: https://aep.dev/not-precedent
