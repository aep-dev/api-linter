---
rule:
  aep: 131
  name: [core, '0131', request-path-behavior]
  summary: |
    Get RPCs should annotate the `path` field with `aep.api.field_behavior`.
permalink: /131/request-path-behavior
redirect_from:
  - /0131/request-path-behavior
---

# Get methods: Field behavior

This rule enforces that all `Get` standard methods have
`aep.api.field_behavior` set to `FIELD_BEHAVIOR_REQUIRED` on their `string path` field, as
mandated in [AEP-131][].

## Details

This rule looks at any message matching `Get*Request` and complains if the
`path` field does not have a `aep.api.field_behavior` annotation with a
value of `FIELD_BEHAVIOR_REQUIRED`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message GetBookRequest {
  // The `aep.api.field_behavior` annotation should also be included.
  string path = 1 [(aep.api.field_info).resource_reference = "library.googleapis.com/Book"];
}
```

**Correct** code for this rule:

```proto
// Correct.
message GetBookRequest {
  string path = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (aep.api.field_info).resource_reference = "library.googleapis.com/Book"
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message GetBookRequest {
  // (-- api-linter: core::0131::request-path-behavior=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string path = 1 [(aep.api.field_info).resource_reference = "library.googleapis.com/Book"];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-131]: https://aep.dev/131
[aep.dev/not-precedent]: https://aep.dev/not-precedent
