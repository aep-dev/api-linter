---
rule:
  aep: 162
  name: [core, '0162', commit-request-name-field]
  summary: Commit RPCs must have a `name` field in the request.
permalink: /162/commit-request-name-field
redirect_from:
  - /0162/commit-request-name-field
---

# Commit requests: Name field

This rule enforces that all `Commit` methods have a `string name`
field in the request message, as mandated in [AEP-162][].

## Details

This rule looks at any message matching `Commit*Request` and complains if
either the `name` field is missing or it has any type other than `string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should include a `string name` field.
message CommitBookRequest {
}
```

```proto
// Incorrect.
message CommitBookRequest {
  // Field type should be `string`.
  bytes name = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

**Correct** code for this rule:

```proto
// Correct.
message CommitBookRequest {
  string name = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message (if
the `name` field is missing) or above the field (if it is the wrong type).
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message CommitBookRequest {
  // (-- api-linter: core::0162::commit-request-name-field=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  bytes name = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-162]: https://aep.dev/162
[aep.dev/not-precedent]: https://aep.dev/not-precedent
