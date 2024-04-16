---
rule:
  aep: 131
  name: [core, '0131', request-name-behavior]
  summary: |
    Get RPCs should annotate the `name` field with `google.api.field_behavior`.
permalink: /131/request-name-behavior
redirect_from:
  - /0131/request-name-behavior
---

# Get methods: Field behavior

This rule enforces that all `Get` standard methods have
`google.api.field_behavior` set to `REQUIRED` on their `string name` field, as
mandated in [AEP-131][].

## Details

This rule looks at any message matching `Get*Request` and complains if the
`name` field does not have a `google.api.field_behavior` annotation with a
value of `REQUIRED`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message GetBookRequest {
  // The `google.api.field_behavior` annotation should also be included.
  string name = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Book"
  }];
}
```

**Correct** code for this rule:

```proto
// Correct.
message GetBookRequest {
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message GetBookRequest {
  // (-- api-linter: core::0131::request-name-behavior=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string name = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Book"
  }];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-131]: https://aep.dev/131
[aep.dev/not-precedent]: https://aep.dev/not-precedent
