---
rule:
  aep: 162
  name: [core, '0162', tag-revision-request-tag-field]
  summary: Tag Revision RPCs must have a `tag` field in the request.
permalink: /162/tag-revision-request-tag-field
redirect_from:
  - /0162/tag-revision-request-tag-field
---

# Tag Revision requests: Tag field

This rule enforces that all Tag Revision methods have a `string tag`
field in the request message, as mandated in [AEP-162][].

## Details

This rule looks at any message matching `Tag*RevisionRequest` and complains if
either the `tag` field is missing or it has any type other than `string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should include a `string tag` field.
message TagBookRevisionRequest {
  string name = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (aep.api.field_info).resource_reference.type = "library.googleapis.com/Book"
  ];
}
```

```proto
// Incorrect.
message TagBookRevisionRequest {
  string name = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (aep.api.field_info).resource_reference.type = "library.googleapis.com/Book"
  ];

  // Field type should be `string`.
  bytes tag = 2 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

**Correct** code for this rule:

```proto
// Correct.
message TagBookRevisionRequest {
  string name = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (aep.api.field_info).resource_reference.type = "library.googleapis.com/Book"
  ];

  string tag = 2 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message (if
the `tag` field is missing) or above the field (if it is the wrong type).
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message TagBookRevisionRequest {
  string name = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (aep.api.field_info).resource_reference.type = "library.googleapis.com/Book"
  ];

  // (-- api-linter: core::0162::tag-revision-request-tag-field=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  bytes tag = 2 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-162]: https://aep.dev/162
[aep.dev/not-precedent]: https://aep.dev/not-precedent
