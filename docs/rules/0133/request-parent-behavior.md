---
rule:
  aep: 133
  name: [core, '0133', request-parent-behavior]
  summary: |
    Create RPCs should annotate the `parent` field with `aep.api.field_behavior`.
permalink: /133/request-parent-behavior
redirect_from:
  - /0133/request-parent-behavior
---

# Create methods: Field behavior

This rule enforces that all `Create` standard methods have
`aep.api.field_behavior` set to `FIELD_BEHAVIOR_REQUIRED` on their `string parent` field,
as mandated in [AEP-133][].

## Details

This rule looks at any message matching `Create*Request` and complains if the
`parent` field does not have a `aep.api.field_behavior` annotation with a
value of `FIELD_BEHAVIOR_REQUIRED`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message CreateBooksRequest {
  // The `aep.api.field_behavior` annotation should also be included.
  string parent = 1 [(aep.api.field_info).resource_reference = "library.googleapis.com/Publisher"];
  Book book = 2 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

**Correct** code for this rule:

```proto
// Correct.
message CreateBooksRequest {
  string parent = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (aep.api.field_info).resource_reference = "library.googleapis.com/Publisher"
  ];
  Book book = 2 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message CreateBooksRequest {
  // (-- api-linter: core::0133::request-parent-behavior=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string parent = 1 [(aep.api.field_info).resource_reference = "library.googleapis.com/Publisher"];
  Book book = 2 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-133]: https://aep.dev/133
[aep.dev/not-precedent]: https://aep.dev/not-precedent
