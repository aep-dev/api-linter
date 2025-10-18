---
rule:
  aep: 162
  name: [core, '0162', commit-request-name-reference]
  summary: |
    Commit requests should annotate the `name` field with `(aep.api.field_info).resource_reference`.
permalink: /162/commit-request-name-reference
redirect_from:
  - /0162/commit-request-name-reference
---

# Commit requests: Name resource reference

This rule enforces that all `Commit` requests have
`(aep.api.field_info).resource_reference` on their `string name` field, as mandated in
[AEP-162][].

## Details

This rule looks at the `name` field of any message matching `Commit*Request`
and complains if it does not have a `(aep.api.field_info).resource_reference` annotation.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message CommitBookRequest {
  // The `(aep.api.field_info).resource_reference` annotation should also be included.
  string name = 1 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

**Correct** code for this rule:

```proto
// Correct.
message CommitBookRequest {
  string name = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (aep.api.field_info).resource_reference.type = "library.googleapis.com/Book"
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message CommitBookRequest {
  // (-- api-linter: core::0162::commit-request-name-reference=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string name = 1 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-162]: https://aep.dev/162
[aep.dev/not-precedent]: https://aep.dev/not-precedent
