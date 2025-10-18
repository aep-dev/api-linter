---
rule:
  aep: 131
  name: [core, '0131', request-path-reference-type]
  summary: |
    The `(aep.api.field_info).resource_reference` on the `path` field of a Get RPC request
    message should use `type`, not `child_type`.
permalink: /131/request-path-reference-type
redirect_from:
  - /0131/request-path-reference-type
---

# Get methods: Resource reference

This rule enforces that the `(aep.api.field_info).resource_reference` on the `path` field
of a Get RPC request message uses `type`, not `child_type`, as suggested in
[AEP-131][].

## Details

This rule looks at the `(aep.api.field_info).resource_reference` annotation on the `path`
field of any message matching `Get*Request` and complains if it does not use a
direct `type` reference.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message GetBookRequest {
  // The `(aep.api.field_info).resource_reference` annotation should be a direct `type`
  // reference.
  string path = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (aep.api.field_info).resource_reference.child_type = "library.googleapis.com/Book"
  ];
}
```

**Correct** code for this rule:

```proto
// Correct.
message GetBookRequest {
  string path = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (aep.api.field_info).resource_reference.type = "library.googleapis.com/Book"
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message GetBookRequest {
  // (-- api-linter: core::0131::request-path-reference-type=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string path = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (aep.api.field_info).resource_reference.child_type = "library.googleapis.com/Book"
  ];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-131]: https://aep.dev/131
[aep.dev/not-precedent]: https://aep.dev/not-precedent
