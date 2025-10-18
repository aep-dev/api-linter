---
rule:
  aep: 162
  name: [core, '0162', tag-revision-request-name-behavior]
  summary: |
    Tag Revision requests should annotate the `name` field with `aep.api.field_behavior`.
permalink: /162/tag-revision-request-name-behavior
redirect_from:
  - /0162/tag-revision-request-name-behavior
---

# Tag Revision requests: Name field behavior

This rule enforces that all Tag Revision requests have
`aep.api.field_behavior` set to `FIELD_BEHAVIOR_REQUIRED` on their `string name` field, as
mandated in [AEP-162][].

## Details

This rule looks at any message matching `Tag*RevisionRequest` and complains if the
`name` field does not have a `aep.api.field_behavior` annotation with a
value of `FIELD_BEHAVIOR_REQUIRED`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message TagBookRevisionRequest {
  // The `aep.api.field_behavior` annotation should also be included.
  string name = 1 [
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  string tag = 2 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

**Correct** code for this rule:

```proto
// Correct.
message TagBookRevisionRequest {
  string name = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  string tag = 2 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message TagBookRevisionRequest {
  // (-- api-linter: core::0162::tag-revision-request-name-behavior=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string name = 1 [
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  string tag = 2 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-162]: https://aep.dev/162
[aep.dev/not-precedent]: https://aep.dev/not-precedent
