---
rule:
  aep: 158
  name: [core, '0158', request-skip-field]
  summary: Paginated RPC `skip` fields must have type `int32`.
permalink: /158/request-skip-field
redirect_from:
  - /0158/request-skip-field
---

# Paginated methods: skip field

This rule enforces that all `List` and `Search` request `skip` fields have type `int32`, as
mandated in [AEP-158][].

## Details

This rule looks at any message matching `List*Request` or `Search*Request` that
contains a `skip` field, and complains if the field is not a singular `int32`.

## Examples

**Incorrect** code for this rule:

```proto
message ListBooksRequest {
  string parent = 1 [
    (aep.api.field_info).resource_reference.child_type = "library.googleapis.com/Book",
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED
  ];

  int32 max_page_size = 2;

  string page_token = 3;

  string skip = 4;  // Field type should be `int32`.
}
```

**Correct** code for this rule:

```proto
message ListBooksRequest {
  string parent = 1 [
    (aep.api.field_info).resource_reference.child_type = "library.googleapis.com/Book",
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED
  ];

  int32 max_page_size = 2;

  string page_token = 3;

  int32 skip = 4;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message ListBooksRequest {
  string parent = 1 [
    (aep.api.field_info).resource_reference.child_type = "library.googleapis.com/Book",
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED
  ];

  int32 max_page_size = 2;

  string page_token = 3;

  // (-- api-linter: core::0158::request-skip-field=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string skip = 4;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-158]: https://aep.dev/158
[aep.dev/not-precedent]: https://aep.dev/not-precedent
