---
rule:
  aep: 135
  name: [core, '0135', request-path-behavior]
  summary: |
    Delete RPCs should annotate the `path` field with `aep.api.field_behavior`.
permalink: /135/request-path-behavior
redirect_from:
  - /0135/request-path-behavior
---

# Delete methods: Field behavior

This rule enforces that all `Delete` standard methods have
`aep.api.field_behavior` set to `FIELD_BEHAVIOR_REQUIRED` on their `string path` field, as
mandated in [AEP-135][].

## Details

This rule looks at any message matching `Delete*Request` and complains if the
`path` field does not have a `aep.api.field_behavior` annotation with a
value of `FIELD_BEHAVIOR_REQUIRED`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message DeleteBookRequest {
  // The `aep.api.field_behavior` annotation should also be included.
  string path = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Book"
  }];
}
```

**Correct** code for this rule:

```proto
// Correct.
message DeleteBookRequest {
  string path = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message DeleteBookRequest {
  // (-- api-linter: core::0135::request-path-behavior=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string path = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Book"
  }];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-135]: https://aep.dev/135
[aep.dev/not-precedent]: https://aep.dev/not-precedent
