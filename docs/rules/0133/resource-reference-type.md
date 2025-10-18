---
rule:
  aep: 133
  name: [core, '0133', resource-reference-type]
  summary: Create should use a `child_type` reference to the created resource.
permalink: /133/resource-reference-type
redirect_from:
  - /0133/resource-reference-type
---

# Create methods: Parent field resource reference

This rule enforces that all `Create` standard methods with a `string parent`
field use a proper `(aep.api.field_info).resource_reference`, that being either a
`child_type` referring to the created resource or a `type` referring directly
to the parent resource, as mandated in [AEP-133][].

## Details

This rule looks at any message matching `Create*Request` and complains if the 
`(aep.api.field_info).resource_reference` on the `parent` field refers to the wrong
resource.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message CreateBooksRequest {
  // `child_type` should be used instead of `type` when referring to the
  // created resource on a parent field.
  string parent = 1 [(aep.api.field_info).resource_reference.type = "library.googleapis.com/Book"];
  Book book = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message CreateBooksRequest {
  string parent = 1 [(aep.api.field_info).resource_reference.child_type = "library.googleapis.com/Book"];
  Book book = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message CreateBooksRequest {
  // (-- api-linter: core::0133::resource-reference-type=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string parent = 1 [(aep.api.field_info).resource_reference.type = "library.googleapis.com/Book"];
  Book book = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-133]: https://aep.dev/133
[aep.dev/not-precedent]: https://aep.dev/not-precedent
