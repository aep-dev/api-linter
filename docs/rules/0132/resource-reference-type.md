---
rule:
  aep: 132
  name: [core, '0132', resource-reference-type]
  summary: List should use `resource_reference_child_type` to reference the paginated resource.
permalink: /132/resource-reference-type
redirect_from:
  - /0132/resource-reference-type
---

# List methods: Parent field resource reference

This rule enforces that all `List` standard methods with a `string parent`
field use a proper `(aep.api.field_info).resource_reference_child_type` to refer to the
paginated resource, as mandated in [AEP-132][].

## Details

This rule looks at any message matching `List*Request` and complains if the
`(aep.api.field_info).resource_reference_child_type` or `(aep.api.field_info).resource_reference`
on the `parent` field refers to the wrong resource. The preferred approach is to use
`resource_reference_child_type` to reference the child resource being paginated.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message ListBooksRequest {
  // Should reference the correct child resource type.
  string parent = 1 [(aep.api.field_info).resource_reference_child_type = "library.googleapis.com/Shelf"];
  int32 page_size = 2;
  string page_token = 3;
}
```

**Correct** code for this rule:

```proto
// Correct.
message ListBooksRequest {
  string parent = 1 [(aep.api.field_info).resource_reference_child_type = "library.googleapis.com/Book"];
  int32 page_size = 2;
  string page_token = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message ListBooksRequest {
  // (-- api-linter: core::0132::resource-reference-type=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string parent = 1 [(aep.api.field_info).resource_reference_child_type = "library.googleapis.com/Shelf"];
  int32 page_size = 2;
  string page_token = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-132]: https://aep.dev/132
[aep.dev/not-precedent]: https://aep.dev/not-precedent
