---
rule:
  aep: 132
  name: [core, '0132', request-parent-reference]
  summary: |
    List RPCs should annotate the `parent` field with `google.api.resource_reference`.
permalink: /132/request-parent-reference
redirect_from:
  - /0132/request-parent-reference
---

# List methods: Resource reference

This rule enforces that all `List` standard methods have
`google.api.resource_reference` on their `string parent` field, as mandated in
[AEP-132][].

## Details

This rule looks at the `parent` field of any message matching `List*Request`
and complains if it does not have a `google.api.resource_reference` annotation.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message ListBooksRequest {
  // The `google.api.resource_reference` annotation should also be included.
  string parent = 1 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
  int32 page_size = 2;
  string page_token = 3;
}
```

**Correct** code for this rule:

```proto
// Correct.
message ListBooksRequest {
  string parent = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Publisher"
  ];
  int32 page_size = 2;
  string page_token = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0132::request-parent-reference=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
message ListBooksRequest {
  string parent = 1 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-132]: https://aep.dev/132
[aep.dev/not-precedent]: https://aep.dev/not-precedent
