---
rule:
  aep: 164
  name: [core, '0164', request-name-behavior]
  summary: |
    Undelete RPCs should annotate the `name` field with `aep.api.field_behavior`.
permalink: /164/request-name-behavior
redirect_from:
  - /0164/request-name-behavior
---

# Undelete methods: Field behavior

This rule enforces that all `Undelete` methods have
`aep.api.field_behavior` set to `FIELD_BEHAVIOR_REQUIRED` on their `string name` field, as
mandated in [AEP-164][].

## Details

This rule looks at any message matching `Undelete*Request` and complains if the
`name` field does not have a `aep.api.field_behavior` annotation with a
value of `FIELD_BEHAVIOR_REQUIRED`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message UndeleteBookRequest {
  // The `aep.api.field_behavior` annotation should also be included.
  string name = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Book"
  }];
}
```

**Correct** code for this rule:

```proto
// Correct.
message UndeleteBookRequest {
  string name = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message UndeleteBookRequest {
  // (-- api-linter: core::0164::request-name-behavior=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string name = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Book"
  }];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-164]: https://aep.dev/164
[aep.dev/not-precedent]: https://aep.dev/not-precedent
