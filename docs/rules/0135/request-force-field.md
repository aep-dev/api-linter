---
rule:
  aep: 135
  name: [core, '0135', request-force-field]
  summary: Delete request `force` fields must have type `bool`.
permalink: /135/request-force-field
redirect_from:
  - /0135/request-force-field
---

# Delete requests: force field

This rule enforces that all `Delete` request `force` fields have type `bool`, as
mandated in [AEP-135][].

## Details

This rule looks at any message matching `Delete*Request` that contains a `force`
field and complains if the field is not a singular `bool`.

## Examples

**Incorrect** code for this rule:

```proto
message DeletePublisherRequest {
  string path = 1 [
    (google.api.resource_reference).type = "library.googleapis.com/Publisher",
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED
  ];

  int32 force = 2;  // Field type should be `bool`.
}
```

**Correct** code for this rule:

```proto
message DeletePublisherRequest {
  string path = 1 [
    (google.api.resource_reference).type = "library.googleapis.com/Publisher",
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED
  ];

  bool force = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message DeletePublisherRequest {
  string path = 1 [
    (google.api.resource_reference).type = "library.googleapis.com/Publisher",
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED
  ];

  // (-- api-linter: core::0135::request-force-field=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  int32 force = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-135]: https://aep.dev/135
[aep.dev/not-precedent]: https://aep.dev/not-precedent
