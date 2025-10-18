---
rule:
  aep: 132
  name: [core, '0132', request-parent-valid-reference]
  summary: |
    List RPCs should reference the parent resource, not the
    listed resource.
permalink: /132/request-parent-valid-reference
redirect_from:
  - /0132/request-parent-valid-reference
---

# List methods: Resource reference

This rule enforces that all `List` standard methods reference a resource other
than the resource being listed with the `(aep.api.field_info).resource_reference` on
their `string parent` field, as mandated in [AEP-132][].

## Details

This rule looks at the `parent` field of any message matching `List*Request`
and complains if the `(aep.api.field_info).resource_reference` annotation references
the resource being listed.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message ListBooksRequest {
  // The `(aep.api.field_info).resource_reference` should not reference the resource
  // being listed.
  string parent = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (aep.api.field_info).resource_reference.type = "library.googleapis.com/Book"
  ];
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
    (aep.api.field_info).resource_reference.child_type = "library.googleapis.com/Book"
  ];
  int32 page_size = 2;
  string page_token = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0132::request-parent-valid-reference=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
message ListBooksRequest {
  string parent = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (aep.api.field_info).resource_reference.type = "library.googleapis.com/Book"
  ];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-132]: https://aep.dev/132
[aep.dev/not-precedent]: https://aep.dev/not-precedent
